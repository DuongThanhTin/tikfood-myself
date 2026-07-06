package auth

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestMemoryUserRepository_CreateAndFindByEmail(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	created, err := repo.CreateUser(ctx, User{Email: "Alice@Example.com", PasswordHash: "hash", DisplayName: "Alice"})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected generated id")
	}
	if created.CreatedAt.IsZero() || created.UpdatedAt.IsZero() {
		t.Fatal("expected timestamps to be set")
	}

	// Email lookup is case-insensitive (mirrors citext).
	found, err := repo.FindByEmail(ctx, "alice@example.com")
	if err != nil {
		t.Fatalf("FindByEmail: %v", err)
	}
	if found.ID != created.ID {
		t.Fatalf("expected same user, got %q vs %q", found.ID, created.ID)
	}

	if _, err := repo.FindByID(ctx, created.ID); err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if _, err := repo.FindByEmail(ctx, "nobody@example.com"); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestMemoryUserRepository_DuplicateEmailRejected(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	if _, err := repo.CreateUser(ctx, User{Email: "dup@example.com", PasswordHash: "h"}); err != nil {
		t.Fatalf("first CreateUser: %v", err)
	}
	// Case-insensitive duplicate must be rejected.
	if _, err := repo.CreateUser(ctx, User{Email: "DUP@example.com", PasswordHash: "h2"}); !errors.Is(err, ErrEmailTaken) {
		t.Fatalf("expected ErrEmailTaken, got %v", err)
	}
}

func TestMemoryUserRepository_FindByGoogleSubAndLink(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	if _, err := repo.FindByGoogleSub(ctx, "sub-123"); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound for unknown sub, got %v", err)
	}

	created, err := repo.CreateUser(ctx, User{Email: "g@example.com", GoogleSub: "sub-123", EmailVerified: true})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	found, err := repo.FindByGoogleSub(ctx, "sub-123")
	if err != nil || found.ID != created.ID {
		t.Fatalf("FindByGoogleSub: got %+v err %v", found, err)
	}

	linkTarget, err := repo.CreateUser(ctx, User{Email: "link@example.com", PasswordHash: "h"})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	linked, err := repo.LinkGoogleSub(ctx, linkTarget.ID, "sub-999")
	if err != nil {
		t.Fatalf("LinkGoogleSub: %v", err)
	}
	if linked.GoogleSub != "sub-999" {
		t.Fatalf("expected linked google sub, got %q", linked.GoogleSub)
	}
	if again, err := repo.FindByGoogleSub(ctx, "sub-999"); err != nil || again.ID != linkTarget.ID {
		t.Fatalf("expected to find linked user, got %+v err %v", again, err)
	}
}

func TestMemoryRefreshTokenRepository_StoreFindRevoke(t *testing.T) {
	repo := NewMemoryRefreshTokenRepository()
	ctx := context.Background()
	now := time.Now()

	stored, err := repo.StoreRefreshToken(ctx, RefreshToken{
		UserID:    "user-1",
		TokenHash: "hash-1",
		ExpiresAt: now.Add(time.Hour),
		UserAgent: "test-agent",
	})
	if err != nil {
		t.Fatalf("StoreRefreshToken: %v", err)
	}
	if stored.ID == "" || stored.CreatedAt.IsZero() {
		t.Fatal("expected generated id and created_at")
	}

	found, err := repo.FindRefreshTokenByHash(ctx, "hash-1")
	if err != nil {
		t.Fatalf("FindRefreshTokenByHash: %v", err)
	}
	if !found.IsUsable(now) {
		t.Fatal("expected token to be usable before revoke/expiry")
	}

	didRevoke, err := repo.RevokeRefreshToken(ctx, "hash-1")
	if err != nil {
		t.Fatalf("RevokeRefreshToken: %v", err)
	}
	if !didRevoke {
		t.Fatal("expected RevokeRefreshToken to report a live token was revoked")
	}
	revoked, err := repo.FindRefreshTokenByHash(ctx, "hash-1")
	if err != nil {
		t.Fatalf("FindRefreshTokenByHash after revoke: %v", err)
	}
	if revoked.IsUsable(now) {
		t.Fatal("expected token to be unusable after revoke")
	}

	// Revoke is idempotent, including for unknown hashes; a second revoke reports false.
	if again, err := repo.RevokeRefreshToken(ctx, "hash-1"); err != nil || again {
		t.Fatalf("second RevokeRefreshToken: again=%v err=%v", again, err)
	}
	if again, err := repo.RevokeRefreshToken(ctx, "does-not-exist"); err != nil || again {
		t.Fatalf("revoke unknown hash should be idempotent false: again=%v err=%v", again, err)
	}

	if _, err := repo.FindRefreshTokenByHash(ctx, "missing"); !errors.Is(err, ErrRefreshTokenNotFound) {
		t.Fatalf("expected ErrRefreshTokenNotFound, got %v", err)
	}
}

func TestMemoryRefreshTokenRepository_Rotate(t *testing.T) {
	repo := NewMemoryRefreshTokenRepository()
	ctx := context.Background()
	now := time.Now()

	if _, err := repo.StoreRefreshToken(ctx, RefreshToken{UserID: "u1", TokenHash: "old", ExpiresAt: now.Add(time.Hour)}); err != nil {
		t.Fatalf("StoreRefreshToken: %v", err)
	}

	// First rotation wins: old is revoked and the replacement is stored atomically.
	rotated, err := repo.RotateRefreshToken(ctx, "old", RefreshToken{UserID: "u1", TokenHash: "new", ExpiresAt: now.Add(time.Hour)})
	if err != nil || !rotated {
		t.Fatalf("expected rotation to win: rotated=%v err=%v", rotated, err)
	}
	if old, _ := repo.FindRefreshTokenByHash(ctx, "old"); old.IsUsable(now) {
		t.Fatal("expected old token revoked after rotation")
	}
	if fresh, _ := repo.FindRefreshTokenByHash(ctx, "new"); !fresh.IsUsable(now) {
		t.Fatal("expected replacement token stored and usable")
	}

	// Replaying the already-rotated token loses the race and stores nothing.
	rotated, err = repo.RotateRefreshToken(ctx, "old", RefreshToken{UserID: "u1", TokenHash: "new2", ExpiresAt: now.Add(time.Hour)})
	if err != nil || rotated {
		t.Fatalf("expected replay to lose: rotated=%v err=%v", rotated, err)
	}
	if _, err := repo.FindRefreshTokenByHash(ctx, "new2"); !errors.Is(err, ErrRefreshTokenNotFound) {
		t.Fatal("a lost rotation must not store its replacement")
	}
}

func TestMemoryRefreshTokenRepository_PurgeExpired(t *testing.T) {
	repo := NewMemoryRefreshTokenRepository()
	ctx := context.Background()
	now := time.Now()

	// live (kept), expired (purged), and revoked (purged).
	if _, err := repo.StoreRefreshToken(ctx, RefreshToken{UserID: "u1", TokenHash: "live", ExpiresAt: now.Add(time.Hour)}); err != nil {
		t.Fatalf("store live: %v", err)
	}
	if _, err := repo.StoreRefreshToken(ctx, RefreshToken{UserID: "u1", TokenHash: "expired", ExpiresAt: now.Add(-time.Hour)}); err != nil {
		t.Fatalf("store expired: %v", err)
	}
	if _, err := repo.StoreRefreshToken(ctx, RefreshToken{UserID: "u1", TokenHash: "revoked", ExpiresAt: now.Add(time.Hour)}); err != nil {
		t.Fatalf("store revoked: %v", err)
	}
	if _, err := repo.RevokeRefreshToken(ctx, "revoked"); err != nil {
		t.Fatalf("revoke: %v", err)
	}

	removed, err := repo.PurgeExpiredRefreshTokens(ctx, now)
	if err != nil {
		t.Fatalf("PurgeExpiredRefreshTokens: %v", err)
	}
	if removed != 2 {
		t.Fatalf("expected 2 rows purged, got %d", removed)
	}
	if _, err := repo.FindRefreshTokenByHash(ctx, "live"); err != nil {
		t.Fatalf("live token should survive purge: %v", err)
	}
	if _, err := repo.FindRefreshTokenByHash(ctx, "expired"); !errors.Is(err, ErrRefreshTokenNotFound) {
		t.Fatal("expired token should be purged")
	}
	if _, err := repo.FindRefreshTokenByHash(ctx, "revoked"); !errors.Is(err, ErrRefreshTokenNotFound) {
		t.Fatal("revoked token should be purged")
	}
}

func TestMemoryRefreshTokenRepository_RevokeAllForUser(t *testing.T) {
	repo := NewMemoryRefreshTokenRepository()
	ctx := context.Background()
	now := time.Now()

	for _, hash := range []string{"a", "b"} {
		if _, err := repo.StoreRefreshToken(ctx, RefreshToken{UserID: "u1", TokenHash: hash, ExpiresAt: now.Add(time.Hour)}); err != nil {
			t.Fatalf("store %s: %v", hash, err)
		}
	}
	if _, err := repo.StoreRefreshToken(ctx, RefreshToken{UserID: "u2", TokenHash: "c", ExpiresAt: now.Add(time.Hour)}); err != nil {
		t.Fatalf("store c: %v", err)
	}

	if err := repo.RevokeAllForUser(ctx, "u1"); err != nil {
		t.Fatalf("RevokeAllForUser: %v", err)
	}
	for _, hash := range []string{"a", "b"} {
		tok, _ := repo.FindRefreshTokenByHash(ctx, hash)
		if tok.IsUsable(now) {
			t.Fatalf("expected %s revoked", hash)
		}
	}
	other, _ := repo.FindRefreshTokenByHash(ctx, "c")
	if !other.IsUsable(now) {
		t.Fatal("expected u2 token to remain usable")
	}
}
