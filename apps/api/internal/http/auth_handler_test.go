package http

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/auth"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
)

func testAuthRouter(t *testing.T) http.Handler {
	return testAuthRouterWithGoogle(t, nil)
}

func testAuthRouterWithGoogle(t *testing.T, google auth.GoogleAuthenticator) http.Handler {
	t.Helper()
	issuer, err := auth.NewTokenIssuer(testJWTSecret, 15*time.Minute)
	if err != nil {
		t.Fatalf("NewTokenIssuer: %v", err)
	}
	service := auth.NewAuthService(auth.ServiceConfig{
		Users:              auth.NewMemoryUserRepository(),
		RefreshTokens:      auth.NewMemoryRefreshTokenRepository(),
		VerificationTokens: auth.NewMemoryEmailVerificationTokenRepository(),
		Issuer:             issuer,
		Mailer:             auth.NewLogMailer(slog.New(slog.NewTextHandler(io.Discard, nil))),
		RefreshTTL:         720 * time.Hour,
		VerificationTTL:    24 * time.Hour,
		VerifyBaseURL:      "http://localhost:3000",
	})
	return NewRouter(RouterDependencies{
		AllowedOrigins: []string{"http://localhost:3000"},
		RouteRegistrars: DefaultRouteRegistrars(HandlerDependencies{
			Venues:       newTestVenueService(),
			Auth:         service,
			Google:       google,
			RefreshTTL:   720 * time.Hour,
			CookieSecure: false,
			WebOrigin:    "http://localhost:3000",
		}),
	})
}

// fakeGoogleAuthenticator injects a canned profile so tests never call Google.
type fakeGoogleAuthenticator struct {
	authURL string
	profile auth.GoogleProfile
	err     error
}

func (f fakeGoogleAuthenticator) AuthCodeURL(state string) string {
	return f.authURL + "?state=" + state
}

func (f fakeGoogleAuthenticator) ExchangeAndFetchProfile(_ context.Context, _ string) (auth.GoogleProfile, error) {
	return f.profile, f.err
}

const testJWTSecret = "test-secret-value"

// accessTokenFrom extracts the access_token from a register/login response body.
func accessTokenFrom(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body.Data.AccessToken == "" {
		t.Fatal("expected an access token in the response")
	}
	return body.Data.AccessToken
}

func doGet(router http.Handler, path string, header map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	for k, v := range header {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
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

func TestAuthMiddleware_ValidTokenAllows(t *testing.T) {
	router := testAuthRouter(t)
	reg := doJSON(router, http.MethodPost, "/api/v1/auth/register",
		`{"email":"me@example.com","password":"password123","display_name":"Me"}`, nil)
	if reg.Code != http.StatusCreated {
		t.Fatalf("register: %d", reg.Code)
	}
	token := accessTokenFrom(t, reg)

	rec := doGet(router, "/api/v1/auth/me", map[string]string{"Authorization": "Bearer " + token})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (%s)", rec.Code, rec.Body.String())
	}
	var body struct {
		Data struct {
			User struct {
				Email string `json:"email"`
			} `json:"user"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body.Data.User.Email != "me@example.com" {
		t.Fatalf("unexpected /me user: %+v", body.Data.User)
	}
}

func TestAuthMiddleware_MissingTokenReturns401(t *testing.T) {
	router := testAuthRouter(t)
	rec := doGet(router, "/api/v1/auth/me", nil)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	assertErrorCode(t, rec, "unauthorized")
}

func TestAuthMiddleware_ExpiredReturns401(t *testing.T) {
	router := testAuthRouter(t)
	// Mint a token that is already expired, signed with the same secret the router uses.
	expiredIssuer, err := auth.NewTokenIssuer(testJWTSecret, -time.Minute)
	if err != nil {
		t.Fatalf("NewTokenIssuer: %v", err)
	}
	token, _, err := expiredIssuer.IssueAccessToken(auth.User{ID: "u1", Email: "e@x.com"})
	if err != nil {
		t.Fatalf("IssueAccessToken: %v", err)
	}
	rec := doGet(router, "/api/v1/auth/me", map[string]string{"Authorization": "Bearer " + token})
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for expired token, got %d", rec.Code)
	}
}

func TestCORS_AllowsConfiguredOrigin(t *testing.T) {
	router := testAuthRouter(t)
	rec := doGet(router, "/health", map[string]string{"Origin": "http://localhost:3000"})
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Fatalf("expected ACAO to echo the origin, got %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Fatalf("expected credentials allowed, got %q", got)
	}
}

func TestCORS_RejectsUnknownOrigin(t *testing.T) {
	router := testAuthRouter(t)
	rec := doGet(router, "/health", map[string]string{"Origin": "http://evil.example.com"})
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no ACAO for an unknown origin, got %q", got)
	}
}

func TestDiscoveryRoutesStayPublic(t *testing.T) {
	router := testAuthRouter(t)
	rec := doGet(router, "/api/v1/map/venues?district=District%201", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("discovery must stay public, got %d", rec.Code)
	}
}

func TestGoogleLogin_RedirectsToGoogle(t *testing.T) {
	fake := fakeGoogleAuthenticator{authURL: "https://accounts.google.com/o/oauth2/auth"}
	router := testAuthRouterWithGoogle(t, fake)

	rec := doGet(router, "/api/v1/auth/google/login", nil)
	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if !strings.HasPrefix(loc, "https://accounts.google.com/o/oauth2/auth") {
		t.Fatalf("expected redirect to Google, got %q", loc)
	}
	// A state cookie must be set and echoed into the redirect (CSRF binding).
	var state string
	for _, c := range rec.Result().Cookies() {
		if c.Name == stateCookieName {
			state = c.Value
		}
	}
	if state == "" || !strings.Contains(loc, "state="+state) {
		t.Fatalf("expected state cookie bound to redirect; state=%q loc=%q", state, loc)
	}
}

func TestGoogleCallback_StateMismatchRejected(t *testing.T) {
	fake := fakeGoogleAuthenticator{profile: auth.GoogleProfile{Sub: "s", Email: "a@b.com", EmailVerified: true}}
	router := testAuthRouterWithGoogle(t, fake)

	// state query does not match the state cookie -> reject.
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/google/callback?code=x&state=attacker", nil)
	req.AddCookie(&http.Cookie{Name: stateCookieName, Value: "legit"})
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 on state mismatch, got %d", rec.Code)
	}
	assertErrorCode(t, rec, "unauthorized")
}

func TestGoogleCallback_SuccessSetsCookieAndRedirects(t *testing.T) {
	fake := fakeGoogleAuthenticator{profile: auth.GoogleProfile{Sub: "sub-1", Email: "gcb@example.com", EmailVerified: true}}
	router := testAuthRouterWithGoogle(t, fake)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/google/callback?code=good&state=match", nil)
	req.AddCookie(&http.Cookie{Name: stateCookieName, Value: "match"})
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d (%s)", rec.Code, rec.Body.String())
	}
	if loc := rec.Header().Get("Location"); loc != "http://localhost:3000/auth/google/callback" {
		t.Fatalf("unexpected redirect: %q", loc)
	}
	refreshCookieFrom(t, rec) // the session is handed off via the httpOnly refresh cookie, not the URL
}

func TestGoogleCallback_UnverifiedEmailForbidden(t *testing.T) {
	// A password account exists; an unverified Google email for the same address must
	// not be allowed to link (account-takeover guard) -> 403.
	fake := fakeGoogleAuthenticator{profile: auth.GoogleProfile{Sub: "sub-x", Email: "vic@example.com", EmailVerified: false}}
	gr := testAuthRouterWithGoogle(t, fake)
	if rec := doJSON(gr, http.MethodPost, "/api/v1/auth/register",
		`{"email":"vic@example.com","password":"password123"}`, nil); rec.Code != http.StatusCreated {
		t.Fatalf("register on google router: %d", rec.Code)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/google/callback?code=good&state=match", nil)
	req.AddCookie(&http.Cookie{Name: stateCookieName, Value: "match"})
	rec := httptest.NewRecorder()
	gr.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
	assertErrorCode(t, rec, "forbidden")
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
