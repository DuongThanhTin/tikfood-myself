package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/auth"
	"github.com/jackc/pgx/v5/pgconn"
)

// UserRepository is the Postgres-backed auth.UserRepository. Persistence only; all
// queries are parameterized and it never leaks SQL to callers.
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

var _ auth.UserRepository = (*UserRepository)(nil)

const userColumns = `id::text, email::text, coalesce(password_hash, ''), display_name, coalesce(google_sub, ''), email_verified, created_at, updated_at`

func scanUser(row interface{ Scan(...any) error }) (auth.User, error) {
	var user auth.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.DisplayName,
		&user.GoogleSub,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return user, err
}

func (repo *UserRepository) CreateUser(ctx context.Context, user auth.User) (auth.User, error) {
	const query = `
insert into users (email, password_hash, display_name, google_sub, email_verified)
values ($1, nullif($2, ''), $3, nullif($4, ''), $5)
returning ` + userColumns
	created, err := scanUser(repo.db.QueryRowContext(ctx, query,
		user.Email, user.PasswordHash, user.DisplayName, user.GoogleSub, user.EmailVerified,
	))
	if err != nil {
		switch uniqueViolation(err) {
		case "users_email_key":
			return auth.User{}, auth.ErrEmailTaken
		case "users_google_sub_key":
			return auth.User{}, auth.ErrGoogleSubTaken
		}
		return auth.User{}, fmt.Errorf("create user: %w", err)
	}
	return created, nil
}

// uniqueViolation returns the constraint name when err is a Postgres unique-violation
// (SQLSTATE 23505), or "" otherwise. It lets repositories map a specific collision to a
// clean domain error instead of leaking a generic 500.
func uniqueViolation(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return pgErr.ConstraintName
	}
	return ""
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (auth.User, error) {
	const query = `select ` + userColumns + ` from users where email = $1 limit 1`
	return repo.findOne(ctx, query, email)
}

func (repo *UserRepository) FindByID(ctx context.Context, id string) (auth.User, error) {
	const query = `select ` + userColumns + ` from users where id = $1::uuid limit 1`
	return repo.findOne(ctx, query, id)
}

func (repo *UserRepository) FindByGoogleSub(ctx context.Context, googleSub string) (auth.User, error) {
	if googleSub == "" {
		return auth.User{}, auth.ErrUserNotFound
	}
	const query = `select ` + userColumns + ` from users where google_sub = $1 limit 1`
	return repo.findOne(ctx, query, googleSub)
}

func (repo *UserRepository) LinkGoogleSub(ctx context.Context, userID string, googleSub string) (auth.User, error) {
	const query = `
update users set google_sub = nullif($2, ''), updated_at = now()
where id = $1::uuid
returning ` + userColumns
	user, err := scanUser(repo.db.QueryRowContext(ctx, query, userID, googleSub))
	if errors.Is(err, sql.ErrNoRows) {
		return auth.User{}, auth.ErrUserNotFound
	}
	if err != nil {
		if uniqueViolation(err) == "users_google_sub_key" {
			return auth.User{}, auth.ErrGoogleSubTaken
		}
		return auth.User{}, fmt.Errorf("link google sub: %w", err)
	}
	return user, nil
}

func (repo *UserRepository) findOne(ctx context.Context, query string, arg any) (auth.User, error) {
	user, err := scanUser(repo.db.QueryRowContext(ctx, query, arg))
	if errors.Is(err, sql.ErrNoRows) {
		return auth.User{}, auth.ErrUserNotFound
	}
	if err != nil {
		return auth.User{}, fmt.Errorf("find user: %w", err)
	}
	return user, nil
}
