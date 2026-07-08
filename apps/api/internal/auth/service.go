package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/mail"
	"net/url"
	"strings"
	"time"
)

// Service-level validation and outcome errors. Persistence errors live in model.go.
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidEmail       = errors.New("invalid email address")
	ErrWeakPassword       = errors.New("password does not meet requirements")
	ErrRefreshInvalid     = errors.New("refresh token is invalid or expired")
	ErrEmailNotVerified   = errors.New("google email is not verified")
	// ErrAccountExistsUsePassword is returned when a Google login matches an existing
	// account by email that has never proven ownership of that email (email_verified is
	// false — the case for every password account today). Auto-linking there would let an
	// attacker who pre-registered the victim's email absorb the victim's Google identity,
	// so we refuse and steer the user to their existing password login instead.
	ErrAccountExistsUsePassword = errors.New("an account with this email already exists; sign in with your password")
	// ErrEmailVerificationInvalid covers a missing, unknown, expired, or already-consumed
	// email-verification token.
	ErrEmailVerificationInvalid = errors.New("email verification token is invalid or expired")
	// ErrEmailAlreadyVerified is returned when a resend is requested for an already-verified account.
	ErrEmailAlreadyVerified = errors.New("email is already verified")
)

const (
	// MinPasswordLength is the single source of truth for the password floor; the HTTP
	// layer derives its user-facing message from it so the two cannot drift.
	MinPasswordLength = 8
	maxPasswordBytes  = 72 // bcrypt truncates beyond this; reject rather than silently cut.
	// defaultVerificationTTL bounds how long an email-verification link stays valid.
	defaultVerificationTTL = 24 * time.Hour
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
// register/login/refresh/logout/verify use-cases. It never imports gin.
type AuthService struct {
	users              UserRepository
	refreshTokens      RefreshTokenRepository
	verificationTokens EmailVerificationTokenRepository
	issuer             *TokenIssuer
	mailer             Mailer
	logger             *slog.Logger
	refreshTTL         time.Duration
	verificationTTL    time.Duration
	verifyBaseURL      string // link base, e.g. WEB_ORIGIN; link = verifyBaseURL + "/verify-email?token=..."
}

// ServiceConfig wires the auth service's dependencies. VerificationTokens/Mailer may be
// nil (email verification then simply does nothing); Logger defaults to slog.Default and
// VerificationTTL to defaultVerificationTTL when unset.
type ServiceConfig struct {
	Users              UserRepository
	RefreshTokens      RefreshTokenRepository
	VerificationTokens EmailVerificationTokenRepository
	Issuer             *TokenIssuer
	Mailer             Mailer
	Logger             *slog.Logger
	RefreshTTL         time.Duration
	VerificationTTL    time.Duration
	VerifyBaseURL      string
}

func NewAuthService(cfg ServiceConfig) *AuthService {
	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}
	verificationTTL := cfg.VerificationTTL
	if verificationTTL <= 0 {
		verificationTTL = defaultVerificationTTL
	}
	return &AuthService{
		users:              cfg.Users,
		refreshTokens:      cfg.RefreshTokens,
		verificationTokens: cfg.VerificationTokens,
		issuer:             cfg.Issuer,
		mailer:             cfg.Mailer,
		logger:             logger,
		refreshTTL:         cfg.RefreshTTL,
		verificationTTL:    verificationTTL,
		verifyBaseURL:      cfg.VerifyBaseURL,
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

	// Best-effort: a failure to send the verification email must not fail registration —
	// the account exists and the user can request a resend.
	s.sendVerificationEmail(ctx, user)

	pair, err := s.issueTokens(ctx, user, meta)
	if err != nil {
		return User{}, TokenPair{}, err
	}
	return sanitize(user), pair, nil
}

// VerifyEmail consumes a single-use verification token and marks its user email-verified.
// A missing, unknown, expired, or already-consumed token returns ErrEmailVerificationInvalid.
func (s *AuthService) VerifyEmail(ctx context.Context, rawToken string) (User, error) {
	if rawToken == "" || s.verificationTokens == nil {
		return User{}, ErrEmailVerificationInvalid
	}
	hash := HashRefreshToken(rawToken) // shared opaque-token SHA-256 hex

	stored, err := s.verificationTokens.FindEmailVerificationTokenByHash(ctx, hash)
	if errors.Is(err, ErrEmailVerificationTokenNotFound) {
		return User{}, ErrEmailVerificationInvalid
	}
	if err != nil {
		return User{}, err
	}
	if !stored.IsUsable(time.Now()) {
		return User{}, ErrEmailVerificationInvalid
	}

	// Atomic compare-and-consume: only the first submit flips the token; a replay loses.
	consumed, err := s.verificationTokens.ConsumeEmailVerificationToken(ctx, hash)
	if err != nil {
		return User{}, err
	}
	if !consumed {
		return User{}, ErrEmailVerificationInvalid
	}

	user, err := s.users.MarkEmailVerified(ctx, stored.UserID)
	if err != nil {
		return User{}, err
	}
	return sanitize(user), nil
}

// ResendVerification re-issues a verification email for the given (authenticated) user.
// It returns ErrEmailAlreadyVerified if there is nothing to verify.
func (s *AuthService) ResendVerification(ctx context.Context, userID string) error {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.EmailVerified {
		return ErrEmailAlreadyVerified
	}
	s.sendVerificationEmail(ctx, user)
	return nil
}

// sendVerificationEmail mints a verification token and hands the link to the mailer. It is
// best-effort: it no-ops when verification isn't configured and logs (never returns) on
// failure, so callers can proceed regardless.
func (s *AuthService) sendVerificationEmail(ctx context.Context, user User) {
	if s.verificationTokens == nil || s.mailer == nil {
		return
	}
	raw, hash, err := NewEmailVerificationTokenValue()
	if err != nil {
		s.logger.Warn("email verification token generation failed", "error", err)
		return
	}
	if _, err := s.verificationTokens.CreateEmailVerificationToken(ctx, EmailVerificationToken{
		UserID:    user.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(s.verificationTTL),
	}); err != nil {
		s.logger.Warn("email verification token store failed", "error", err)
		return
	}
	verifyURL := s.verifyBaseURL + "/verify-email?token=" + url.QueryEscape(raw)
	if err := s.mailer.SendVerificationEmail(ctx, user.Email, verifyURL); err != nil {
		s.logger.Warn("email verification send failed", "error", err)
	}
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

	return s.rotateTokens(ctx, user, hash, meta)
}

// Logout revokes the presented refresh token. It is idempotent — logging out an
// already-revoked or unknown token is not an error.
func (s *AuthService) Logout(ctx context.Context, rawRefresh string) error {
	if rawRefresh == "" {
		return nil
	}
	_, err := s.refreshTokens.RevokeRefreshToken(ctx, HashRefreshToken(rawRefresh))
	return err
}

// LoginWithGoogle provisions or links a user from a Google profile and issues a token
// pair. Resolution order: (1) existing google_sub, (2) existing email — auto-linked only
// when Google reports the email verified AND the existing account is itself already
// email-verified (guarding against account takeover), otherwise rejected, (3) otherwise a
// new OAuth-only account (no password).
func (s *AuthService) LoginWithGoogle(ctx context.Context, profile GoogleProfile, meta SessionMeta) (User, TokenPair, error) {
	if profile.Sub == "" {
		return User{}, TokenPair{}, ErrRefreshInvalid
	}

	if user, err := s.users.FindByGoogleSub(ctx, profile.Sub); err == nil {
		return s.finishGoogleLogin(ctx, user, meta)
	} else if !errors.Is(err, ErrUserNotFound) {
		return User{}, TokenPair{}, err
	}

	email := normalizeEmail(profile.Email)
	if !validEmail(email) {
		return User{}, TokenPair{}, ErrInvalidEmail
	}

	if existing, err := s.users.FindByEmail(ctx, email); err == nil {
		if !profile.EmailVerified {
			return User{}, TokenPair{}, ErrEmailNotVerified
		}
		// Only auto-link when the existing account has itself proven ownership of this
		// email. A password account with email_verified=false may have been seeded by an
		// attacker under the victim's address; linking the victim's verified Google
		// identity onto it would silently merge the two. Refuse and send the user to their
		// existing password login (from where an explicit, authenticated link is safe).
		if !existing.EmailVerified {
			return User{}, TokenPair{}, ErrAccountExistsUsePassword
		}
		linked, err := s.users.LinkGoogleSub(ctx, existing.ID, profile.Sub)
		if err != nil {
			return User{}, TokenPair{}, err
		}
		return s.finishGoogleLogin(ctx, linked, meta)
	} else if !errors.Is(err, ErrUserNotFound) {
		return User{}, TokenPair{}, err
	}

	created, err := s.users.CreateUser(ctx, User{
		Email:         email,
		DisplayName:   strings.TrimSpace(profile.Name),
		GoogleSub:     profile.Sub,
		EmailVerified: profile.EmailVerified,
	})
	if err != nil {
		return User{}, TokenPair{}, err
	}
	return s.finishGoogleLogin(ctx, created, meta)
}

func (s *AuthService) finishGoogleLogin(ctx context.Context, user User, meta SessionMeta) (User, TokenPair, error) {
	pair, err := s.issueTokens(ctx, user, meta)
	if err != nil {
		return User{}, TokenPair{}, err
	}
	return sanitize(user), pair, nil
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

// rotateTokens issues a new pair and swaps it for the presented refresh token in a
// single atomic step: the store revokes oldHash and persists the replacement in one DB
// transaction (compare-and-revoke inside the tx). Only the first concurrent refresh wins
// (revoked=true); a loser or a replay of an already-rotated token sees revoked=false and
// is rejected, so a single token never yields two live families — and a crash between
// revoke and store can never leave a session with no valid refresh token.
func (s *AuthService) rotateTokens(ctx context.Context, user User, oldHash string, meta SessionMeta) (TokenPair, error) {
	access, expiresAt, err := s.issuer.IssueAccessToken(user)
	if err != nil {
		return TokenPair{}, err
	}
	raw, hash, err := NewRefreshTokenValue()
	if err != nil {
		return TokenPair{}, err
	}
	revoked, err := s.refreshTokens.RotateRefreshToken(ctx, oldHash, RefreshToken{
		UserID:    user.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(s.refreshTTL),
		UserAgent: meta.UserAgent,
		IP:        meta.IP,
	})
	if err != nil {
		return TokenPair{}, err
	}
	if !revoked {
		return TokenPair{}, ErrRefreshInvalid
	}
	return TokenPair{AccessToken: access, RefreshTokenRaw: raw, AccessExpiresAt: expiresAt}, nil
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
	return len(password) >= MinPasswordLength && len([]byte(password)) <= maxPasswordBytes
}
