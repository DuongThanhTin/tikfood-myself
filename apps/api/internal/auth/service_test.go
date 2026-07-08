package auth

import (
	"context"
	"errors"
	"net/url"
	"testing"
	"time"
)

// captureMailer records the last verification email so tests can extract the raw token.
type captureMailer struct {
	calls   int
	lastTo  string
	lastURL string
}

func (m *captureMailer) SendVerificationEmail(_ context.Context, to string, verifyURL string) error {
	m.calls++
	m.lastTo = to
	m.lastURL = verifyURL
	return nil
}

func tokenFromURL(t *testing.T, raw string) string {
	t.Helper()
	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("parse verify url %q: %v", raw, err)
	}
	token := parsed.Query().Get("token")
	if token == "" {
		t.Fatalf("no token in verify url %q", raw)
	}
	return token
}

func newVerificationService(t *testing.T, mailer Mailer) (*AuthService, *MemoryUserRepository) {
	t.Helper()
	users := NewMemoryUserRepository()
	return NewAuthService(ServiceConfig{
		Users:              users,
		RefreshTokens:      NewMemoryRefreshTokenRepository(),
		VerificationTokens: NewMemoryEmailVerificationTokenRepository(),
		Issuer:             mustIssuer(t, 15*time.Minute),
		Mailer:             mailer,
		RefreshTTL:         time.Hour,
		VerificationTTL:    time.Hour,
		VerifyBaseURL:      "http://web.local",
	}), users
}

func newTestService(t *testing.T, refreshTTL time.Duration) *AuthService {
	t.Helper()
	// Mailer left nil: email verification no-ops here, keeping these cases focused on the
	// register/login/refresh use-cases. The verification flow has its own tests below.
	return NewAuthService(ServiceConfig{
		Users:         NewMemoryUserRepository(),
		RefreshTokens: NewMemoryRefreshTokenRepository(),
		Issuer:        mustIssuer(t, 15*time.Minute),
		RefreshTTL:    refreshTTL,
	})
}

var testMeta = SessionMeta{UserAgent: "test-agent", IP: "127.0.0.1"}

func TestRegister_HappyPath(t *testing.T) {
	svc := newTestService(t, time.Hour)
	user, pair, err := svc.Register(context.Background(), "alice@example.com", "password123", "Alice", testMeta)
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if user.ID == "" || user.Email != "alice@example.com" {
		t.Fatalf("unexpected user: %+v", user)
	}
	if user.PasswordHash != "" {
		t.Fatal("password hash must not be exposed on the returned user")
	}
	if pair.AccessToken == "" || pair.RefreshTokenRaw == "" || pair.AccessExpiresAt.IsZero() {
		t.Fatalf("expected a full token pair, got %+v", pair)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	svc := newTestService(t, time.Hour)
	ctx := context.Background()
	if _, _, err := svc.Register(ctx, "dup@example.com", "password123", "", testMeta); err != nil {
		t.Fatalf("first register: %v", err)
	}
	if _, _, err := svc.Register(ctx, "DUP@example.com", "password123", "", testMeta); !errors.Is(err, ErrEmailTaken) {
		t.Fatalf("expected ErrEmailTaken, got %v", err)
	}
}

func TestRegister_WeakPassword(t *testing.T) {
	svc := newTestService(t, time.Hour)
	if _, _, err := svc.Register(context.Background(), "weak@example.com", "short", "", testMeta); !errors.Is(err, ErrWeakPassword) {
		t.Fatalf("expected ErrWeakPassword, got %v", err)
	}
}

func TestRegister_InvalidEmail(t *testing.T) {
	svc := newTestService(t, time.Hour)
	if _, _, err := svc.Register(context.Background(), "not-an-email", "password123", "", testMeta); !errors.Is(err, ErrInvalidEmail) {
		t.Fatalf("expected ErrInvalidEmail, got %v", err)
	}
}

func TestLogin_HappyPath(t *testing.T) {
	svc := newTestService(t, time.Hour)
	ctx := context.Background()
	if _, _, err := svc.Register(ctx, "log@example.com", "password123", "Log", testMeta); err != nil {
		t.Fatalf("register: %v", err)
	}
	user, pair, err := svc.Login(ctx, "LOG@example.com", "password123", testMeta)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if user.Email != "log@example.com" || pair.AccessToken == "" || pair.RefreshTokenRaw == "" {
		t.Fatalf("unexpected login result: user %+v pair %+v", user, pair)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	svc := newTestService(t, time.Hour)
	ctx := context.Background()
	if _, _, err := svc.Register(ctx, "wp@example.com", "password123", "", testMeta); err != nil {
		t.Fatalf("register: %v", err)
	}
	if _, _, err := svc.Login(ctx, "wp@example.com", "wrong-password", testMeta); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_UnknownUserSameError(t *testing.T) {
	svc := newTestService(t, time.Hour)
	// Unknown user and wrong password must return the identical error (no enumeration).
	if _, _, err := svc.Login(context.Background(), "ghost@example.com", "whatever12", testMeta); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials for unknown user, got %v", err)
	}
}

func TestRefresh_RotatesAndOldRevoked(t *testing.T) {
	svc := newTestService(t, time.Hour)
	ctx := context.Background()
	_, pair, err := svc.Register(ctx, "rot@example.com", "password123", "", testMeta)
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	rotated, err := svc.Refresh(ctx, pair.RefreshTokenRaw, testMeta)
	if err != nil {
		t.Fatalf("Refresh: %v", err)
	}
	if rotated.RefreshTokenRaw == pair.RefreshTokenRaw {
		t.Fatal("expected a rotated (new) refresh token")
	}
	if rotated.AccessToken == "" {
		t.Fatal("expected a new access token")
	}

	// The old token must no longer work (rotation revoked it).
	if _, err := svc.Refresh(ctx, pair.RefreshTokenRaw, testMeta); !errors.Is(err, ErrRefreshInvalid) {
		t.Fatalf("expected old token to be revoked, got %v", err)
	}
	// The new token should still work.
	if _, err := svc.Refresh(ctx, rotated.RefreshTokenRaw, testMeta); err != nil {
		t.Fatalf("expected new token to work: %v", err)
	}
}

func TestRefresh_ExpiredRejected(t *testing.T) {
	svc := newTestService(t, -time.Minute) // refresh tokens issue already-expired
	ctx := context.Background()
	_, pair, err := svc.Register(ctx, "exp@example.com", "password123", "", testMeta)
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if _, err := svc.Refresh(ctx, pair.RefreshTokenRaw, testMeta); !errors.Is(err, ErrRefreshInvalid) {
		t.Fatalf("expected ErrRefreshInvalid for expired token, got %v", err)
	}
}

func TestRefresh_RevokedRejected(t *testing.T) {
	svc := newTestService(t, time.Hour)
	ctx := context.Background()
	_, pair, err := svc.Register(ctx, "rev@example.com", "password123", "", testMeta)
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := svc.Logout(ctx, pair.RefreshTokenRaw); err != nil {
		t.Fatalf("Logout: %v", err)
	}
	if _, err := svc.Refresh(ctx, pair.RefreshTokenRaw, testMeta); !errors.Is(err, ErrRefreshInvalid) {
		t.Fatalf("expected ErrRefreshInvalid for revoked token, got %v", err)
	}
}

func TestLogout_Idempotent(t *testing.T) {
	svc := newTestService(t, time.Hour)
	ctx := context.Background()
	_, pair, err := svc.Register(ctx, "out@example.com", "password123", "", testMeta)
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := svc.Logout(ctx, pair.RefreshTokenRaw); err != nil {
		t.Fatalf("first Logout: %v", err)
	}
	if err := svc.Logout(ctx, pair.RefreshTokenRaw); err != nil {
		t.Fatalf("second Logout should be idempotent: %v", err)
	}
	if err := svc.Logout(ctx, "never-issued"); err != nil {
		t.Fatalf("logout of unknown token should be idempotent: %v", err)
	}
}

func TestLoginWithGoogle_NewUser(t *testing.T) {
	svc := newTestService(t, time.Hour)
	profile := GoogleProfile{Sub: "google-1", Email: "gnew@example.com", Name: "G New", EmailVerified: true}
	user, pair, err := svc.LoginWithGoogle(context.Background(), profile, testMeta)
	if err != nil {
		t.Fatalf("LoginWithGoogle: %v", err)
	}
	if user.Email != "gnew@example.com" || pair.AccessToken == "" || pair.RefreshTokenRaw == "" {
		t.Fatalf("unexpected result: user %+v pair %+v", user, pair)
	}
}

func TestLoginWithGoogle_ExistingByGoogleSub(t *testing.T) {
	svc := newTestService(t, time.Hour)
	ctx := context.Background()
	profile := GoogleProfile{Sub: "google-2", Email: "gsub@example.com", EmailVerified: true}
	first, _, err := svc.LoginWithGoogle(ctx, profile, testMeta)
	if err != nil {
		t.Fatalf("first login: %v", err)
	}
	second, _, err := svc.LoginWithGoogle(ctx, profile, testMeta)
	if err != nil {
		t.Fatalf("second login: %v", err)
	}
	if first.ID != second.ID {
		t.Fatalf("expected the same user on repeat Google login, got %q vs %q", first.ID, second.ID)
	}
}

func TestLoginWithGoogle_ExistingUnverifiedAccountRejected(t *testing.T) {
	svc := newTestService(t, time.Hour)
	ctx := context.Background()
	// A password account exists but has never verified its email (the only kind today).
	if _, _, err := svc.Register(ctx, "link@example.com", "password123", "Link", testMeta); err != nil {
		t.Fatalf("register: %v", err)
	}
	// Even with a verified Google email, auto-linking onto an unverified account is refused
	// (account-takeover guard): the user is steered to their existing password login.
	if _, _, err := svc.LoginWithGoogle(ctx, GoogleProfile{Sub: "google-3", Email: "link@example.com", EmailVerified: true}, testMeta); !errors.Is(err, ErrAccountExistsUsePassword) {
		t.Fatalf("expected ErrAccountExistsUsePassword, got %v", err)
	}
}

func TestLoginWithGoogle_LinksVerifiedAccount(t *testing.T) {
	users := NewMemoryUserRepository()
	svc := NewAuthService(ServiceConfig{
		Users:         users,
		RefreshTokens: NewMemoryRefreshTokenRepository(),
		Issuer:        mustIssuer(t, 15*time.Minute),
		RefreshTTL:    time.Hour,
	})
	ctx := context.Background()
	// Seed an already email-verified account (as a real verification flow eventually would).
	seeded, err := users.CreateUser(ctx, User{Email: "link@example.com", PasswordHash: "x", DisplayName: "Link", EmailVerified: true})
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	// Google login with the same verified email links to that verified account, not a new one.
	linked, _, err := svc.LoginWithGoogle(ctx, GoogleProfile{Sub: "google-3", Email: "link@example.com", EmailVerified: true}, testMeta)
	if err != nil {
		t.Fatalf("LoginWithGoogle: %v", err)
	}
	if linked.ID != seeded.ID {
		t.Fatalf("expected link to existing account %q, got %q", seeded.ID, linked.ID)
	}
}

func TestLoginWithGoogle_UnverifiedEmailRejected(t *testing.T) {
	svc := newTestService(t, time.Hour)
	ctx := context.Background()
	if _, _, err := svc.Register(ctx, "takeover@example.com", "password123", "", testMeta); err != nil {
		t.Fatalf("register: %v", err)
	}
	// Unverified Google email must NOT be allowed to link to an existing account.
	if _, _, err := svc.LoginWithGoogle(ctx, GoogleProfile{Sub: "google-4", Email: "takeover@example.com", EmailVerified: false}, testMeta); !errors.Is(err, ErrEmailNotVerified) {
		t.Fatalf("expected ErrEmailNotVerified, got %v", err)
	}
}

func TestRegisterSendsVerificationAndVerifyEmailFlow(t *testing.T) {
	mailer := &captureMailer{}
	svc, _ := newVerificationService(t, mailer)
	ctx := context.Background()

	user, _, err := svc.Register(ctx, "verify@example.com", "password123", "V", testMeta)
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if user.EmailVerified {
		t.Fatal("new account should start unverified")
	}
	if mailer.calls != 1 || mailer.lastTo != "verify@example.com" {
		t.Fatalf("expected one verification email to the user, got calls=%d to=%q", mailer.calls, mailer.lastTo)
	}

	raw := tokenFromURL(t, mailer.lastURL)
	verified, err := svc.VerifyEmail(ctx, raw)
	if err != nil {
		t.Fatalf("VerifyEmail: %v", err)
	}
	if !verified.EmailVerified {
		t.Fatal("expected email to be verified after VerifyEmail")
	}

	// Single-use: replaying the same token fails.
	if _, err := svc.VerifyEmail(ctx, raw); !errors.Is(err, ErrEmailVerificationInvalid) {
		t.Fatalf("expected replay to be rejected, got %v", err)
	}
}

func TestVerifyEmail_InvalidToken(t *testing.T) {
	svc, _ := newVerificationService(t, &captureMailer{})
	if _, err := svc.VerifyEmail(context.Background(), "not-a-real-token"); !errors.Is(err, ErrEmailVerificationInvalid) {
		t.Fatalf("expected ErrEmailVerificationInvalid, got %v", err)
	}
	if _, err := svc.VerifyEmail(context.Background(), ""); !errors.Is(err, ErrEmailVerificationInvalid) {
		t.Fatalf("expected ErrEmailVerificationInvalid for empty token, got %v", err)
	}
}

func TestResendVerification(t *testing.T) {
	mailer := &captureMailer{}
	svc, _ := newVerificationService(t, mailer)
	ctx := context.Background()

	user, _, err := svc.Register(ctx, "resend@example.com", "password123", "R", testMeta)
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if err := svc.ResendVerification(ctx, user.ID); err != nil {
		t.Fatalf("ResendVerification: %v", err)
	}
	if mailer.calls != 2 { // one from register, one from resend
		t.Fatalf("expected 2 verification emails, got %d", mailer.calls)
	}

	// Once verified, a resend reports ErrEmailAlreadyVerified.
	if _, err := svc.VerifyEmail(ctx, tokenFromURL(t, mailer.lastURL)); err != nil {
		t.Fatalf("VerifyEmail: %v", err)
	}
	if err := svc.ResendVerification(ctx, user.ID); !errors.Is(err, ErrEmailAlreadyVerified) {
		t.Fatalf("expected ErrEmailAlreadyVerified, got %v", err)
	}
}

func TestMe_ReturnsUser(t *testing.T) {
	svc := newTestService(t, time.Hour)
	ctx := context.Background()
	user, _, err := svc.Register(ctx, "me@example.com", "password123", "Me", testMeta)
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	got, err := svc.Me(ctx, user.ID)
	if err != nil {
		t.Fatalf("Me: %v", err)
	}
	if got.ID != user.ID {
		t.Fatalf("expected same user, got %+v", got)
	}
	if _, err := svc.Me(ctx, "no-such-id"); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
