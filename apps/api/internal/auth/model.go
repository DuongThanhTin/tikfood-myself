// Package auth is the authentication domain: users, refresh tokens, and the
// register/login/refresh/logout use-cases. It mirrors internal/discovery — a
// Service holds business logic and calls Repository interfaces, and an in-memory
// fallback repository allows tests to run without a database. It never imports gin.
package auth

import (
	"errors"
	"time"
)

// Repository / persistence errors. Service-level validation errors live in service.go.
var (
	ErrUserNotFound         = errors.New("user not found")
	ErrEmailTaken           = errors.New("email already registered")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

// User is an authenticated account. password_hash and google_sub are secrets/internal
// identifiers and are never serialized to clients.
type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"` // empty means no password (OAuth-only account)
	DisplayName   string    `json:"display_name"`
	GoogleSub     string    `json:"-"` // empty means not linked to Google
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// RefreshToken is a persisted, revocable session credential. Only the hash of the raw
// token value is stored; the raw value is returned to the client once at issue time.
type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
	UserAgent string
	IP        string
}

// IsUsable reports whether the token can still authenticate a refresh, at time now.
func (t RefreshToken) IsUsable(now time.Time) bool {
	return t.RevokedAt == nil && t.ExpiresAt.After(now)
}
