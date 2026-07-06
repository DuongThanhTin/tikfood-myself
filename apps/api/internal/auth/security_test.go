package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestHashPassword_VerifyRoundTrip(t *testing.T) {
	hash, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if hash == "" || hash == "correct horse battery staple" {
		t.Fatal("expected a non-empty hash distinct from the plaintext")
	}
	if err := VerifyPassword(hash, "correct horse battery staple"); err != nil {
		t.Fatalf("VerifyPassword should accept the right password: %v", err)
	}
}

func TestVerifyPassword_WrongPasswordFails(t *testing.T) {
	hash, err := HashPassword("s3cret-password")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if err := VerifyPassword(hash, "wrong-password"); err == nil {
		t.Fatal("expected verification to fail for a wrong password")
	}
}

func TestNewTokenIssuer_EmptySecretFails(t *testing.T) {
	if _, err := NewTokenIssuer("", time.Minute); err == nil {
		t.Fatal("expected an error when the secret is empty")
	}
}

func TestIssueAndParseAccessToken_RoundTrip(t *testing.T) {
	issuer := mustIssuer(t, 15*time.Minute)
	user := User{ID: "user-1", Email: "alice@example.com"}

	token, expiresAt, err := issuer.IssueAccessToken(user)
	if err != nil {
		t.Fatalf("IssueAccessToken: %v", err)
	}
	if token == "" {
		t.Fatal("expected a token")
	}
	if d := time.Until(expiresAt); d <= 0 || d > 15*time.Minute+time.Second {
		t.Fatalf("unexpected expiry window: %s", d)
	}

	claims, err := issuer.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("ParseAccessToken: %v", err)
	}
	if claims.UserID != "user-1" || claims.Email != "alice@example.com" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestParseAccessToken_ExpiredFails(t *testing.T) {
	issuer := mustIssuer(t, -time.Minute) // already expired
	token, _, err := issuer.IssueAccessToken(User{ID: "user-1", Email: "a@b.com"})
	if err != nil {
		t.Fatalf("IssueAccessToken: %v", err)
	}
	if _, err := issuer.ParseAccessToken(token); err == nil {
		t.Fatal("expected an expired token to be rejected")
	}
}

func TestParseAccessToken_TamperedFails(t *testing.T) {
	issuer := mustIssuer(t, 15*time.Minute)
	token, _, err := issuer.IssueAccessToken(User{ID: "user-1", Email: "a@b.com"})
	if err != nil {
		t.Fatalf("IssueAccessToken: %v", err)
	}
	// Flip the last character of the signature.
	tampered := token[:len(token)-1]
	if token[len(token)-1] == 'a' {
		tampered += "b"
	} else {
		tampered += "a"
	}
	if _, err := issuer.ParseAccessToken(tampered); err == nil {
		t.Fatal("expected a tampered token to be rejected")
	}
}

func TestParseAccessToken_RejectsNonHS256(t *testing.T) {
	issuer := mustIssuer(t, 15*time.Minute)
	// Forge a token with the "none" algorithm — the classic alg-confusion attack.
	forged := jwt.NewWithClaims(jwt.SigningMethodNone, accessClaims{
		Email: "attacker@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "attacker",
			Issuer:    tokenIssuerName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})
	raw, err := forged.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("sign none token: %v", err)
	}
	if _, err := issuer.ParseAccessToken(raw); err == nil {
		t.Fatal("expected a non-HS256 token to be rejected")
	}
}

func TestNewRefreshTokenValue_HashDeterministic(t *testing.T) {
	raw, hash, err := NewRefreshTokenValue()
	if err != nil {
		t.Fatalf("NewRefreshTokenValue: %v", err)
	}
	if raw == "" || hash == "" || raw == hash {
		t.Fatal("expected distinct non-empty raw value and hash")
	}
	if HashRefreshToken(raw) != hash {
		t.Fatal("HashRefreshToken must be deterministic and match the issued hash")
	}
	// Different raw values must not collide.
	raw2, _, _ := NewRefreshTokenValue()
	if raw2 == raw {
		t.Fatal("expected unique refresh token values")
	}
}

func mustIssuer(t *testing.T, ttl time.Duration) *TokenIssuer {
	t.Helper()
	issuer, err := NewTokenIssuer("a-sufficiently-long-test-secret", ttl)
	if err != nil {
		t.Fatalf("NewTokenIssuer: %v", err)
	}
	return issuer
}
