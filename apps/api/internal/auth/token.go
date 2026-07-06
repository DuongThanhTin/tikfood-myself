package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenIssuerName = "tikfood"
	// refreshTokenBytes is the entropy of a raw refresh token before hashing.
	refreshTokenBytes = 32
)

// ErrEmptySecret is returned when a TokenIssuer is constructed without a signing secret.
var ErrEmptySecret = errors.New("jwt secret must not be empty")

// accessClaims is the JWT payload for an access token. Subject holds the user id.
type accessClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// AccessClaims is the validated, decoded form returned to callers.
type AccessClaims struct {
	UserID    string
	Email     string
	ExpiresAt time.Time
}

// TokenIssuer signs and verifies access tokens with a single HMAC secret. It is
// safe for concurrent use.
type TokenIssuer struct {
	secret    []byte
	accessTTL time.Duration
}

// NewTokenIssuer builds an issuer. The secret must be non-empty; the service refuses
// to start without one (guarding against an accidentally blank JWT_SECRET).
func NewTokenIssuer(secret string, accessTTL time.Duration) (*TokenIssuer, error) {
	if secret == "" {
		return nil, ErrEmptySecret
	}
	return &TokenIssuer{secret: []byte(secret), accessTTL: accessTTL}, nil
}

// IssueAccessToken returns a signed HS256 access token for the user and its absolute
// expiry time.
func (issuer *TokenIssuer) IssueAccessToken(user User) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(issuer.accessTTL)
	claims := accessClaims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			Issuer:    tokenIssuerName,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(issuer.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiresAt, nil
}

// ParseAccessToken validates the token's signature, algorithm (HS256 only), issuer,
// and expiry, returning the decoded claims. It rejects the alg-confusion "none" and
// asymmetric-algorithm attacks by pinning the accepted method.
func (issuer *TokenIssuer) ParseAccessToken(raw string) (AccessClaims, error) {
	var claims accessClaims
	_, err := jwt.ParseWithClaims(raw, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return issuer.secret, nil
	},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(tokenIssuerName),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return AccessClaims{}, err
	}

	result := AccessClaims{UserID: claims.Subject, Email: claims.Email}
	if claims.ExpiresAt != nil {
		result.ExpiresAt = claims.ExpiresAt.Time
	}
	return result, nil
}

// NewRefreshTokenValue generates a high-entropy raw refresh token and its storage
// hash. The raw value is returned to the client once; only the hash is persisted.
func NewRefreshTokenValue() (raw string, hash string, err error) {
	buf := make([]byte, refreshTokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}
	raw = base64.RawURLEncoding.EncodeToString(buf)
	return raw, HashRefreshToken(raw), nil
}

// HashRefreshToken returns the deterministic SHA-256 hash (hex) of a raw refresh
// token, used both at issue time and on lookup.
func HashRefreshToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// RandomHexToken returns 32 bytes of cryptographic randomness, hex-encoded. It is the
// shared source for opaque, non-guessable strings that are not refresh tokens — the
// OAuth anti-CSRF state and the ephemeral dev-only JWT secret. It fails closed: a
// crypto/rand failure is returned, never masked with a predictable value.
func RandomHexToken() (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}
