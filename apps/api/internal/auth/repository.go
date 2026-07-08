package auth

import (
	"context"
	"time"
)

// UserRepository persists users. Implementations: postgres.UserRepository and the
// in-memory MemoryUserRepository (tests / no-DB mode). Persistence only — no business rules.
type UserRepository interface {
	// CreateUser stores a new user and returns it with the generated id and timestamps.
	// Returns ErrEmailTaken if the email is already registered.
	CreateUser(ctx context.Context, user User) (User, error)
	// FindByEmail returns the user with the given (case-insensitive) email, or ErrUserNotFound.
	FindByEmail(ctx context.Context, email string) (User, error)
	// FindByID returns the user with the given id, or ErrUserNotFound.
	FindByID(ctx context.Context, id string) (User, error)
	// FindByGoogleSub returns the user linked to the given Google subject id, or ErrUserNotFound.
	FindByGoogleSub(ctx context.Context, googleSub string) (User, error)
	// LinkGoogleSub attaches a Google subject id to an existing user and returns the updated user.
	LinkGoogleSub(ctx context.Context, userID string, googleSub string) (User, error)
	// MarkEmailVerified sets email_verified=true for the user and returns the updated user,
	// or ErrUserNotFound. Idempotent: verifying an already-verified user is not an error.
	MarkEmailVerified(ctx context.Context, userID string) (User, error)
}

// EmailVerificationTokenRepository persists single-use email-verification tokens.
type EmailVerificationTokenRepository interface {
	// CreateEmailVerificationToken persists a token and returns it with the generated id
	// and created_at.
	CreateEmailVerificationToken(ctx context.Context, token EmailVerificationToken) (EmailVerificationToken, error)
	// FindEmailVerificationTokenByHash returns the token with the given hash, or
	// ErrEmailVerificationTokenNotFound.
	FindEmailVerificationTokenByHash(ctx context.Context, tokenHash string) (EmailVerificationToken, error)
	// ConsumeEmailVerificationToken atomically marks the token consumed and reports whether
	// it flipped a live token (true) versus finding it unknown/already-consumed (false).
	// The bool makes verification single-use even under concurrent submits.
	ConsumeEmailVerificationToken(ctx context.Context, tokenHash string) (bool, error)
}

// RefreshTokenRepository persists rotating refresh tokens.
type RefreshTokenRepository interface {
	// StoreRefreshToken persists a token and returns it with the generated id and created_at.
	StoreRefreshToken(ctx context.Context, token RefreshToken) (RefreshToken, error)
	// FindRefreshTokenByHash returns the token with the given hash, or ErrRefreshTokenNotFound.
	FindRefreshTokenByHash(ctx context.Context, tokenHash string) (RefreshToken, error)
	// RevokeRefreshToken marks the token with the given hash revoked and reports whether it
	// actually flipped a live token (true) versus finding it unknown/already-revoked (false).
	// Idempotent: revoking an unknown/already-revoked token is not an error. The bool lets
	// rotation detect a lost race (a concurrent refresh already consumed the token).
	RevokeRefreshToken(ctx context.Context, tokenHash string) (bool, error)
	// RotateRefreshToken atomically revokes oldHash and, only if that revoke won the race
	// (the token was live), stores newToken — both in a single transaction. It reports
	// whether the rotation happened (true) or was rejected because oldHash was already
	// consumed/unknown (false, newToken not stored). Atomicity guarantees a session is
	// never left with its old token revoked but no replacement persisted.
	RotateRefreshToken(ctx context.Context, oldHash string, newToken RefreshToken) (bool, error)
	// RevokeAllForUser revokes every active token belonging to the user.
	RevokeAllForUser(ctx context.Context, userID string) error
	// PurgeExpiredRefreshTokens deletes tokens that can no longer authenticate as of `now`
	// — expired or revoked — to bound table growth. Returns the number of rows removed.
	PurgeExpiredRefreshTokens(ctx context.Context, now time.Time) (int64, error)
}
