package auth

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"
)

// Service-level validation and outcome errors. Persistence errors live in model.go.
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidEmail       = errors.New("invalid email address")
	ErrWeakPassword       = errors.New("password does not meet requirements")
	ErrRefreshInvalid     = errors.New("refresh token is invalid or expired")
)

const (
	minPasswordLength = 8
	maxPasswordBytes  = 72 // bcrypt truncates beyond this; reject rather than silently cut.
)

// dummyPasswordHash is compared against when a login targets an unknown user, so the
// response time does not reveal whether the email exists (anti-enumeration).
var dummyPasswordHash, _ = HashPassword("tikfood-timing-equalizer")

// TokenPair is what a successful auth returns: a short-lived access token, the raw
// refresh token (shown once, then only its hash is stored), and the access expiry.
type TokenPair struct {
	AccessToken     string
	RefreshTokenRaw string
	AccessExpiresAt time.Time
}

// SessionMeta is optional audit context captured when a refresh token is issued.
type SessionMeta struct {
	UserAgent string
	IP        string
}

// AuthService orchestrates the security primitives and repositories into the
// register/login/refresh/logout use-cases. It never imports gin.
type AuthService struct {
	users         UserRepository
	refreshTokens RefreshTokenRepository
	issuer        *TokenIssuer
	refreshTTL    time.Duration
}

func NewAuthService(users UserRepository, refreshTokens RefreshTokenRepository, issuer *TokenIssuer, refreshTTL time.Duration) *AuthService {
	return &AuthService{
		users:         users,
		refreshTokens: refreshTokens,
		issuer:        issuer,
		refreshTTL:    refreshTTL,
	}
}

// Register validates input, creates the user, and issues a first token pair.
func (s *AuthService) Register(ctx context.Context, email, password, displayName string, meta SessionMeta) (User, TokenPair, error) {
	email = normalizeEmail(email)
	if !validEmail(email) {
		return User{}, TokenPair{}, ErrInvalidEmail
	}
	if !validPassword(password) {
		return User{}, TokenPair{}, ErrWeakPassword
	}

	hash, err := HashPassword(password)
	if err != nil {
		return User{}, TokenPair{}, err
	}

	user, err := s.users.CreateUser(ctx, User{
		Email:        email,
		PasswordHash: hash,
		DisplayName:  strings.TrimSpace(displayName),
	})
	if err != nil {
		return User{}, TokenPair{}, err // includes ErrEmailTaken
	}

	pair, err := s.issueTokens(ctx, user, meta)
	if err != nil {
		return User{}, TokenPair{}, err
	}
	return sanitize(user), pair, nil
}

// Login verifies credentials and issues a token pair. Unknown-user and wrong-password
// return the identical ErrInvalidCredentials, and both pay a bcrypt comparison, so the
// caller cannot enumerate accounts by response or timing.
func (s *AuthService) Login(ctx context.Context, email, password string, meta SessionMeta) (User, TokenPair, error) {
	email = normalizeEmail(email)

	user, err := s.users.FindByEmail(ctx, email)
	if errors.Is(err, ErrUserNotFound) {
		_ = VerifyPassword(dummyPasswordHash, password)
		return User{}, TokenPair{}, ErrInvalidCredentials
	}
	if err != nil {
		return User{}, TokenPair{}, err
	}

	if user.PasswordHash == "" { // OAuth-only account has no password to check
		_ = VerifyPassword(dummyPasswordHash, password)
		return User{}, TokenPair{}, ErrInvalidCredentials
	}
	if err := VerifyPassword(user.PasswordHash, password); err != nil {
		return User{}, TokenPair{}, ErrInvalidCredentials
	}

	pair, err := s.issueTokens(ctx, user, meta)
	if err != nil {
		return User{}, TokenPair{}, err
	}
	return sanitize(user), pair, nil
}

// Refresh rotates a refresh token: it validates the presented token, revokes it, and
// issues a fresh pair. A replayed (already-rotated), revoked, or expired token fails.
func (s *AuthService) Refresh(ctx context.Context, rawRefresh string, meta SessionMeta) (TokenPair, error) {
	if rawRefresh == "" {
		return TokenPair{}, ErrRefreshInvalid
	}
	hash := HashRefreshToken(rawRefresh)

	stored, err := s.refreshTokens.FindRefreshTokenByHash(ctx, hash)
	if errors.Is(err, ErrRefreshTokenNotFound) {
		return TokenPair{}, ErrRefreshInvalid
	}
	if err != nil {
		return TokenPair{}, err
	}
	if !stored.IsUsable(time.Now()) {
		return TokenPair{}, ErrRefreshInvalid
	}

	user, err := s.users.FindByID(ctx, stored.UserID)
	if errors.Is(err, ErrUserNotFound) {
		return TokenPair{}, ErrRefreshInvalid
	}
	if err != nil {
		return TokenPair{}, err
	}

	// Revoke the old token before issuing the new one so it can never be reused.
	if err := s.refreshTokens.RevokeRefreshToken(ctx, hash); err != nil {
		return TokenPair{}, err
	}
	return s.issueTokens(ctx, user, meta)
}

// Logout revokes the presented refresh token. It is idempotent — logging out an
// already-revoked or unknown token is not an error.
func (s *AuthService) Logout(ctx context.Context, rawRefresh string) error {
	if rawRefresh == "" {
		return nil
	}
	return s.refreshTokens.RevokeRefreshToken(ctx, HashRefreshToken(rawRefresh))
}

// ParseAccessToken validates a Bearer access token and returns its claims. It lets
// HTTP middleware authenticate requests without reaching into the token issuer.
func (s *AuthService) ParseAccessToken(raw string) (AccessClaims, error) {
	return s.issuer.ParseAccessToken(raw)
}

// Me returns the current user by id.
func (s *AuthService) Me(ctx context.Context, userID string) (User, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return User{}, err
	}
	return sanitize(user), nil
}

// sanitize clears secrets/internal identifiers before a user leaves the service
// boundary. Defense in depth: these fields are already json:"-", but callers should
// never receive them in memory either.
func sanitize(user User) User {
	user.PasswordHash = ""
	user.GoogleSub = ""
	return user
}

func (s *AuthService) issueTokens(ctx context.Context, user User, meta SessionMeta) (TokenPair, error) {
	access, expiresAt, err := s.issuer.IssueAccessToken(user)
	if err != nil {
		return TokenPair{}, err
	}
	raw, hash, err := NewRefreshTokenValue()
	if err != nil {
		return TokenPair{}, err
	}
	if _, err := s.refreshTokens.StoreRefreshToken(ctx, RefreshToken{
		UserID:    user.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(s.refreshTTL),
		UserAgent: meta.UserAgent,
		IP:        meta.IP,
	}); err != nil {
		return TokenPair{}, err
	}
	return TokenPair{AccessToken: access, RefreshTokenRaw: raw, AccessExpiresAt: expiresAt}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validEmail(email string) bool {
	if len(email) == 0 || len(email) > 254 {
		return false
	}
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	// mail.ParseAddress accepts display names; require the bare address to match.
	return addr.Address == email
}

func validPassword(password string) bool {
	return len(password) >= minPasswordLength && len([]byte(password)) <= maxPasswordBytes
}
