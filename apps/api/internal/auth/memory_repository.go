package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
	"sync"
	"time"
)

// newID returns a random UUIDv4 string. The Postgres tables generate ids DB-side
// via gen_random_uuid(); the in-memory repositories generate them here so tests and
// no-DB mode behave the same.
func newID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand should never fail; fall back to a time-derived value.
		return fmt.Sprintf("00000000-0000-4000-8000-%012x", time.Now().UnixNano())
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 10
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// MemoryUserRepository is an in-memory UserRepository for tests and no-DB mode.
type MemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]User // keyed by id
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{users: make(map[string]User)}
}

var (
	_ UserRepository         = (*MemoryUserRepository)(nil)
	_ RefreshTokenRepository = (*MemoryRefreshTokenRepository)(nil)
)

func (r *MemoryUserRepository) CreateUser(_ context.Context, user User) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.users {
		if strings.EqualFold(existing.Email, user.Email) {
			return User{}, ErrEmailTaken
		}
	}

	now := time.Now()
	user.ID = newID()
	user.CreatedAt = now
	user.UpdatedAt = now
	r.users[user.ID] = user
	return user, nil
}

func (r *MemoryUserRepository) FindByEmail(_ context.Context, email string) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if strings.EqualFold(user.Email, email) {
			return user, nil
		}
	}
	return User{}, ErrUserNotFound
}

func (r *MemoryUserRepository) FindByID(_ context.Context, id string) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if user, ok := r.users[id]; ok {
		return user, nil
	}
	return User{}, ErrUserNotFound
}

func (r *MemoryUserRepository) FindByGoogleSub(_ context.Context, googleSub string) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if googleSub == "" {
		return User{}, ErrUserNotFound
	}
	for _, user := range r.users {
		if user.GoogleSub == googleSub {
			return user, nil
		}
	}
	return User{}, ErrUserNotFound
}

func (r *MemoryUserRepository) LinkGoogleSub(_ context.Context, userID string, googleSub string) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[userID]
	if !ok {
		return User{}, ErrUserNotFound
	}
	user.GoogleSub = googleSub
	user.UpdatedAt = time.Now()
	r.users[userID] = user
	return user, nil
}

// MemoryRefreshTokenRepository is an in-memory RefreshTokenRepository for tests and no-DB mode.
type MemoryRefreshTokenRepository struct {
	mu     sync.RWMutex
	tokens map[string]RefreshToken // keyed by token hash
}

func NewMemoryRefreshTokenRepository() *MemoryRefreshTokenRepository {
	return &MemoryRefreshTokenRepository{tokens: make(map[string]RefreshToken)}
}

func (r *MemoryRefreshTokenRepository) StoreRefreshToken(_ context.Context, token RefreshToken) (RefreshToken, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	token.ID = newID()
	token.CreatedAt = time.Now()
	r.tokens[token.TokenHash] = token
	return token, nil
}

func (r *MemoryRefreshTokenRepository) FindRefreshTokenByHash(_ context.Context, tokenHash string) (RefreshToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if token, ok := r.tokens[tokenHash]; ok {
		return token, nil
	}
	return RefreshToken{}, ErrRefreshTokenNotFound
}

func (r *MemoryRefreshTokenRepository) RevokeRefreshToken(_ context.Context, tokenHash string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	token, ok := r.tokens[tokenHash]
	if !ok || token.RevokedAt != nil {
		return false, nil // idempotent: nothing live to revoke
	}
	now := time.Now()
	token.RevokedAt = &now
	r.tokens[tokenHash] = token
	return true, nil
}

func (r *MemoryRefreshTokenRepository) RotateRefreshToken(_ context.Context, oldHash string, newToken RefreshToken) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Compare-and-revoke under the same lock that guards the store, mirroring the
	// single-transaction guarantee of the Postgres implementation.
	old, ok := r.tokens[oldHash]
	if !ok || old.RevokedAt != nil {
		return false, nil // lost the race or replay: nothing live to rotate, store nothing
	}
	now := time.Now()
	old.RevokedAt = &now
	r.tokens[oldHash] = old

	newToken.ID = newID()
	newToken.CreatedAt = now
	r.tokens[newToken.TokenHash] = newToken
	return true, nil
}

func (r *MemoryRefreshTokenRepository) RevokeAllForUser(_ context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for hash, token := range r.tokens {
		if token.UserID == userID && token.RevokedAt == nil {
			token.RevokedAt = &now
			r.tokens[hash] = token
		}
	}
	return nil
}
