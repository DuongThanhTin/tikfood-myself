package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/auth"
)

// RefreshTokenRepository is the Postgres-backed auth.RefreshTokenRepository.
type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

var _ auth.RefreshTokenRepository = (*RefreshTokenRepository)(nil)

func (repo *RefreshTokenRepository) StoreRefreshToken(ctx context.Context, token auth.RefreshToken) (auth.RefreshToken, error) {
	const query = `
insert into refresh_tokens (user_id, token_hash, expires_at, user_agent, ip)
values ($1::uuid, $2, $3, $4, nullif($5, '')::inet)
returning id::text, created_at`
	err := repo.db.QueryRowContext(ctx, query,
		token.UserID, token.TokenHash, token.ExpiresAt, token.UserAgent, token.IP,
	).Scan(&token.ID, &token.CreatedAt)
	if err != nil {
		return auth.RefreshToken{}, fmt.Errorf("store refresh token: %w", err)
	}
	return token, nil
}

func (repo *RefreshTokenRepository) FindRefreshTokenByHash(ctx context.Context, tokenHash string) (auth.RefreshToken, error) {
	const query = `
select id::text, user_id::text, token_hash, expires_at, revoked_at, created_at,
       coalesce(user_agent, ''), coalesce(host(ip), '')
from refresh_tokens where token_hash = $1 limit 1`
	var token auth.RefreshToken
	var revokedAt sql.NullTime
	err := repo.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &revokedAt,
		&token.CreatedAt, &token.UserAgent, &token.IP,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return auth.RefreshToken{}, auth.ErrRefreshTokenNotFound
	}
	if err != nil {
		return auth.RefreshToken{}, fmt.Errorf("find refresh token: %w", err)
	}
	if revokedAt.Valid {
		token.RevokedAt = &revokedAt.Time
	}
	return token, nil
}

func (repo *RefreshTokenRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) (bool, error) {
	const query = `update refresh_tokens set revoked_at = now() where token_hash = $1 and revoked_at is null`
	result, err := repo.db.ExecContext(ctx, query, tokenHash)
	if err != nil {
		return false, fmt.Errorf("revoke refresh token: %w", err)
	}
	// The `revoked_at is null` predicate makes this an atomic compare-and-revoke: exactly
	// one concurrent refresh flips the row (1 affected), the loser sees 0.
	affected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("revoke refresh token: %w", err)
	}
	return affected > 0, nil
}

func (repo *RefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID string) error {
	const query = `update refresh_tokens set revoked_at = now() where user_id = $1::uuid and revoked_at is null`
	if _, err := repo.db.ExecContext(ctx, query, userID); err != nil {
		return fmt.Errorf("revoke all refresh tokens: %w", err)
	}
	return nil
}
