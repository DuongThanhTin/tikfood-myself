package auth

import "context"

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
}

// RefreshTokenRepository persists rotating refresh tokens.
type RefreshTokenRepository interface {
	// StoreRefreshToken persists a token and returns it with the generated id and created_at.
	StoreRefreshToken(ctx context.Context, token RefreshToken) (RefreshToken, error)
	// FindRefreshTokenByHash returns the token with the given hash, or ErrRefreshTokenNotFound.
	FindRefreshTokenByHash(ctx context.Context, tokenHash string) (RefreshToken, error)
	// RevokeRefreshToken marks the token with the given hash revoked. Idempotent: it is not
	// an error if the token is unknown or already revoked.
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	// RevokeAllForUser revokes every active token belonging to the user.
	RevokeAllForUser(ctx context.Context, userID string) error
}
