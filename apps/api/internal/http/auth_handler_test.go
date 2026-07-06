package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/auth"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
)

func testAuthRouter(t *testing.T) http.Handler {
	t.Helper()
	issuer, err := auth.NewTokenIssuer("test-secret-value", 15*time.Minute)
	if err != nil {
		t.Fatalf("NewTokenIssuer: %v", err)
	}
	service := auth.NewAuthService(
		auth.NewMemoryUserRepository(),
		auth.NewMemoryRefreshTokenRepository(),
		issuer,
		720*time.Hour,
	)
	return NewRouter(RouterDependencies{
		RouteRegistrars: DefaultRouteRegistrars(HandlerDependencies{
			Venues:       newTestVenueService(),
			Auth:         service,
			RefreshTTL:   720 * time.Hour,
			CookieSecure: false,
		}),
	})
}

func newTestVenueService() *discovery.VenueService {
	return discovery.NewVenueService(discovery.NewFallbackVenueRepository())
}

func doJSON(router http.Handler, method, path, body string, cookie *http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if cookie != nil {
		req.AddCookie(cookie)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func refreshCookieFrom(t *testing.T, rec *httptest.ResponseRecorder) *http.Cookie {
	t.Helper()
	for _, c := range rec.Result().Cookies() {
		if c.Name == refreshCookieName {
			return c
		}
	}
	t.Fatal("expected a refresh cookie to be set")
	return nil
}

func TestRegister_201SetsCookie(t *testing.T) {
	router := testAuthRouter(t)
	rec := doJSON(router, http.MethodPost, "/api/v1/auth/register",
		`{"email":"new@example.com","password":"password123","display_name":"New"}`, nil)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d (%s)", rec.Code, rec.Body.String())
	}
	cookie := refreshCookieFrom(t, rec)
	if cookie.Value == "" || !cookie.HttpOnly || cookie.Path != refreshCookiePath {
		t.Fatalf("unexpected refresh cookie: %+v", cookie)
	}

	var body struct {
		Data struct {
			User struct {
				ID           string `json:"id"`
				Email        string `json:"email"`
				PasswordHash string `json:"password_hash"`
			} `json:"user"`
			AccessToken     string `json:"access_token"`
			AccessExpiresAt string `json:"access_expires_at"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body.Data.User.ID == "" || body.Data.User.Email != "new@example.com" {
		t.Fatalf("unexpected user: %+v", body.Data.User)
	}
	if body.Data.User.PasswordHash != "" {
		t.Fatal("password_hash must never be serialized")
	}
	if body.Data.AccessToken == "" || body.Data.AccessExpiresAt == "" {
		t.Fatal("expected access token and expiry")
	}
}

func TestRegister_DuplicateReturns422(t *testing.T) {
	router := testAuthRouter(t)
	const body = `{"email":"dup@example.com","password":"password123"}`
	if rec := doJSON(router, http.MethodPost, "/api/v1/auth/register", body, nil); rec.Code != http.StatusCreated {
		t.Fatalf("first register expected 201, got %d", rec.Code)
	}
	rec := doJSON(router, http.MethodPost, "/api/v1/auth/register", body, nil)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rec.Code)
	}
	assertErrorCode(t, rec, "domain_rejected")
}

func TestRegister_WeakPasswordReturns422(t *testing.T) {
	router := testAuthRouter(t)
	rec := doJSON(router, http.MethodPost, "/api/v1/auth/register",
		`{"email":"weak@example.com","password":"short"}`, nil)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rec.Code)
	}
	assertErrorCode(t, rec, "domain_rejected")
}

func TestRegister_MalformedBodyReturns400(t *testing.T) {
	router := testAuthRouter(t)
	rec := doJSON(router, http.MethodPost, "/api/v1/auth/register", `{not-json`, nil)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	assertErrorCode(t, rec, "invalid_request")
}

func TestLogin_200(t *testing.T) {
	router := testAuthRouter(t)
	if rec := doJSON(router, http.MethodPost, "/api/v1/auth/register",
		`{"email":"log@example.com","password":"password123"}`, nil); rec.Code != http.StatusCreated {
		t.Fatalf("register: %d", rec.Code)
	}
	rec := doJSON(router, http.MethodPost, "/api/v1/auth/login",
		`{"email":"log@example.com","password":"password123"}`, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (%s)", rec.Code, rec.Body.String())
	}
	refreshCookieFrom(t, rec) // login sets a cookie too
}

func TestLogin_401OnBadCreds(t *testing.T) {
	router := testAuthRouter(t)
	if rec := doJSON(router, http.MethodPost, "/api/v1/auth/register",
		`{"email":"bad@example.com","password":"password123"}`, nil); rec.Code != http.StatusCreated {
		t.Fatalf("register: %d", rec.Code)
	}
	rec := doJSON(router, http.MethodPost, "/api/v1/auth/login",
		`{"email":"bad@example.com","password":"wrong-password"}`, nil)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	assertErrorCode(t, rec, "unauthorized")
}

func TestRefresh_200RotatesCookie(t *testing.T) {
	router := testAuthRouter(t)
	reg := doJSON(router, http.MethodPost, "/api/v1/auth/register",
		`{"email":"rot@example.com","password":"password123"}`, nil)
	if reg.Code != http.StatusCreated {
		t.Fatalf("register: %d", reg.Code)
	}
	old := refreshCookieFrom(t, reg)

	rec := doJSON(router, http.MethodPost, "/api/v1/auth/refresh", "", old)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (%s)", rec.Code, rec.Body.String())
	}
	rotated := refreshCookieFrom(t, rec)
	if rotated.Value == old.Value {
		t.Fatal("expected the refresh cookie to be rotated")
	}

	// The old cookie must no longer refresh (rotation revoked it).
	if again := doJSON(router, http.MethodPost, "/api/v1/auth/refresh", "", old); again.Code != http.StatusUnauthorized {
		t.Fatalf("expected old cookie rejected with 401, got %d", again.Code)
	}
}

func TestRefresh_401WithoutCookie(t *testing.T) {
	router := testAuthRouter(t)
	rec := doJSON(router, http.MethodPost, "/api/v1/auth/refresh", "", nil)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	assertErrorCode(t, rec, "unauthorized")
}

func TestLogout_ClearsCookie(t *testing.T) {
	router := testAuthRouter(t)
	reg := doJSON(router, http.MethodPost, "/api/v1/auth/register",
		`{"email":"out@example.com","password":"password123"}`, nil)
	cookie := refreshCookieFrom(t, reg)

	rec := doJSON(router, http.MethodPost, "/api/v1/auth/logout", "", cookie)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	cleared := refreshCookieFrom(t, rec)
	if cleared.MaxAge >= 0 && cleared.Value != "" {
		t.Fatalf("expected cleared cookie, got %+v", cleared)
	}
	// Logout is idempotent.
	if again := doJSON(router, http.MethodPost, "/api/v1/auth/logout", "", cookie); again.Code != http.StatusOK {
		t.Fatalf("second logout expected 200, got %d", again.Code)
	}
}

func TestAuthResponsesUseEnvelope(t *testing.T) {
	router := testAuthRouter(t)
	rec := doJSON(router, http.MethodPost, "/api/v1/auth/login",
		`{"email":"ghost@example.com","password":"password123"}`, nil)

	var body map[string]json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if _, ok := body["error"]; !ok {
		t.Fatal("expected an error field in the envelope")
	}
	if _, ok := body["data"]; !ok {
		t.Fatal("expected a data field in the envelope")
	}
}

func assertErrorCode(t *testing.T, rec *httptest.ResponseRecorder, code string) {
	t.Helper()
	var body struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body.Error.Code != code {
		t.Fatalf("expected error code %q, got %q (%s)", code, body.Error.Code, rec.Body.String())
	}
}
