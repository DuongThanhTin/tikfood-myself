package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/auth"
)

// EmailVerificationTokenRepository is the Postgres-backed
// auth.EmailVerificationTokenRepository. Persistence only; only token hashes are stored.
type EmailVerificationTokenRepository struct {
	db *sql.DB
}

func NewEmailVerificationTokenRepository(db *sql.DB) *EmailVerificationTokenRepository {
	return &EmailVerificationTokenRepository{db: db}
}

var _ auth.EmailVerificationTokenRepository = (*EmailVerificationTokenRepository)(nil)

func (repo *EmailVerificationTokenRepository) CreateEmailVerificationToken(ctx context.Context, token auth.EmailVerificationToken) (auth.EmailVerificationToken, error) {
	const query = `
insert into email_verification_tokens (user_id, token_hash, expires_at)
values ($1::uuid, $2, $3)
returning id::text, created_at`
	err := repo.db.QueryRowContext(ctx, query, token.UserID, token.TokenHash, token.ExpiresAt).
		Scan(&token.ID, &token.CreatedAt)
	if err != nil {
		return auth.EmailVerificationToken{}, fmt.Errorf("create email verification token: %w", err)
	}
	return token, nil
}

func (repo *EmailVerificationTokenRepository) FindEmailVerificationTokenByHash(ctx context.Context, tokenHash string) (auth.EmailVerificationToken, error) {
	const query = `
select id::text, user_id::text, token_hash, expires_at, consumed_at, created_at
from email_verification_tokens where token_hash = $1 limit 1`
	var token auth.EmailVerificationToken
	var consumedAt sql.NullTime
	err := repo.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &consumedAt, &token.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return auth.EmailVerificationToken{}, auth.ErrEmailVerificationTokenNotFound
	}
	if err != nil {
		return auth.EmailVerificationToken{}, fmt.Errorf("find email verification token: %w", err)
	}
	if consumedAt.Valid {
		token.ConsumedAt = &consumedAt.Time
	}
	return token, nil
}

func (repo *EmailVerificationTokenRepository) ConsumeEmailVerificationToken(ctx context.Context, tokenHash string) (bool, error) {
	const query = `update email_verification_tokens set consumed_at = now() where token_hash = $1 and consumed_at is null`
	result, err := repo.db.ExecContext(ctx, query, tokenHash)
	if err != nil {
		return false, fmt.Errorf("consume email verification token: %w", err)
	}
	// The `consumed_at is null` predicate makes this an atomic compare-and-consume: exactly
	// one concurrent submit flips the row (1 affected), a replay sees 0.
	affected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("consume email verification token: %w", err)
	}
	return affected > 0, nil
}
