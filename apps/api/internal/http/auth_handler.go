package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/auth"
	"github.com/gin-gonic/gin"
)

const (
	refreshCookieName = "tikfood_refresh"
	refreshCookiePath = "/api/v1/auth"
	stateCookieName   = "tikfood_oauth_state"
	stateCookieMaxAge = 300 // seconds; the state cookie only needs to survive the round-trip
)

// AuthHandler exposes the authentication endpoints. It parses/validates requests,
// calls the auth service, sets/clears the refresh cookie, and responds with the
// standard {data,error} envelope. It holds no business logic.
type AuthHandler struct {
	auth         *auth.AuthService
	google       auth.GoogleAuthenticator // optional; nil disables Google login
	refreshTTL   time.Duration
	cookieSecure bool
	webOrigin    string
}

// AuthHandlerConfig wires the auth handler's dependencies.
type AuthHandlerConfig struct {
	Service      *auth.AuthService
	Google       auth.GoogleAuthenticator
	RefreshTTL   time.Duration
	CookieSecure bool
	WebOrigin    string
}

func NewAuthHandler(cfg AuthHandlerConfig) *AuthHandler {
	if cfg.Service == nil {
		panic("http.NewAuthHandler requires an auth service")
	}
	return &AuthHandler{
		auth:         cfg.Service,
		google:       cfg.Google,
		refreshTTL:   cfg.RefreshTTL,
		cookieSecure: cfg.CookieSecure,
		webOrigin:    cfg.WebOrigin,
	}
}

func (handler *AuthHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	group := v1.Group("/auth")
	group.POST("/register", handler.Register)
	group.POST("/login", handler.Login)
	group.POST("/refresh", handler.Refresh)
	group.POST("/logout", handler.Logout)

	if handler.google != nil {
		group.GET("/google/login", handler.GoogleLogin)
		group.GET("/google/callback", handler.GoogleCallback)
	}

	protected := group.Group("")
	protected.Use(authMiddleware(handler.auth))
	protected.GET("/me", handler.Me)
}

type registerRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (handler *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithBadRequest(c, invalidBody())
		return
	}

	user, pair, err := handler.auth.Register(c.Request.Context(), req.Email, req.Password, req.DisplayName, sessionMetaFrom(c))
	if err != nil {
		respondWithAuthError(c, err)
		return
	}

	handler.setRefreshCookie(c, pair.RefreshTokenRaw)
	respondWithCreated(c, authPayload(user, pair))
}

func (handler *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithBadRequest(c, invalidBody())
		return
	}

	user, pair, err := handler.auth.Login(c.Request.Context(), req.Email, req.Password, sessionMetaFrom(c))
	if err != nil {
		respondWithAuthError(c, err)
		return
	}

	handler.setRefreshCookie(c, pair.RefreshTokenRaw)
	respondWithData(c, authPayload(user, pair))
}

func (handler *AuthHandler) Refresh(c *gin.Context) {
	raw, err := c.Cookie(refreshCookieName)
	if err != nil || raw == "" {
		respondWithError(c, http.StatusUnauthorized, ErrorCodeUnauthorized, MessageSessionExpired)
		return
	}

	pair, err := handler.auth.Refresh(c.Request.Context(), raw, sessionMetaFrom(c))
	if err != nil {
		respondWithAuthError(c, err)
		return
	}

	handler.setRefreshCookie(c, pair.RefreshTokenRaw)
	respondWithData(c, gin.H{
		"access_token":      pair.AccessToken,
		"access_expires_at": pair.AccessExpiresAt,
	})
}

func (handler *AuthHandler) Logout(c *gin.Context) {
	raw, _ := c.Cookie(refreshCookieName)
	// Best-effort, idempotent: always clear the client cookie so the browser is logged out
	// even if server-side revocation hit a transient error (the token expires on its own).
	_ = handler.auth.Logout(c.Request.Context(), raw)
	handler.clearRefreshCookie(c)
	respondWithData(c, gin.H{"logged_out": true})
}

// GoogleLogin begins the Authorization Code flow: it sets a short-lived anti-CSRF
// state cookie and redirects the browser to Google's consent screen.
func (handler *AuthHandler) GoogleLogin(c *gin.Context) {
	state, err := auth.RandomHexToken()
	if err != nil {
		respondWithInternalServerError(c, MessageAuthFailed)
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(stateCookieName, state, stateCookieMaxAge, refreshCookiePath, "", handler.cookieSecure, true)
	c.Redirect(http.StatusFound, handler.google.AuthCodeURL(state))
}

// GoogleCallback completes the flow: it verifies the state cookie (CSRF), exchanges
// the code, provisions/links the user, sets the refresh cookie, and redirects to the
// web app. The access token is NOT placed in the URL — the SPA obtains it via a silent
// refresh using the httpOnly cookie, so no token lands in browser history.
func (handler *AuthHandler) GoogleCallback(c *gin.Context) {
	stateCookie, _ := c.Cookie(stateCookieName)
	state := c.Query("state")
	if stateCookie == "" || state == "" || stateCookie != state {
		respondWithError(c, http.StatusUnauthorized, ErrorCodeUnauthorized, MessageAuthFailed)
		return
	}
	handler.clearStateCookie(c)

	profile, err := handler.google.ExchangeAndFetchProfile(c.Request.Context(), c.Query("code"))
	if err != nil {
		respondWithError(c, http.StatusUnauthorized, ErrorCodeUnauthorized, MessageAuthFailed)
		return
	}

	_, pair, err := handler.auth.LoginWithGoogle(c.Request.Context(), profile, sessionMetaFrom(c))
	if errors.Is(err, auth.ErrEmailNotVerified) {
		respondWithError(c, http.StatusForbidden, ErrorCodeForbidden, MessageEmailNotVerified)
		return
	}
	if err != nil {
		respondWithAuthError(c, err)
		return
	}

	handler.setRefreshCookie(c, pair.RefreshTokenRaw)
	c.Redirect(http.StatusFound, handler.webOrigin+"/auth/google/callback")
}

func (handler *AuthHandler) Me(c *gin.Context) {
	userID := currentUserID(c)
	if userID == "" {
		respondWithError(c, http.StatusUnauthorized, ErrorCodeUnauthorized, MessageSessionExpired)
		return
	}

	user, err := handler.auth.Me(c.Request.Context(), userID)
	if errors.Is(err, auth.ErrUserNotFound) {
		respondWithError(c, http.StatusUnauthorized, ErrorCodeUnauthorized, MessageSessionExpired)
		return
	}
	if err != nil {
		respondWithInternalServerError(c, MessageAuthFailed)
		return
	}

	respondWithData(c, gin.H{"user": user})
}

// authPayload builds the shared success body for register/login.
func authPayload(user auth.User, pair auth.TokenPair) gin.H {
	return gin.H{
		"user":              user,
		"access_token":      pair.AccessToken,
		"access_expires_at": pair.AccessExpiresAt,
	}
}

func sessionMetaFrom(c *gin.Context) auth.SessionMeta {
	return auth.SessionMeta{
		UserAgent: c.Request.UserAgent(),
		IP:        c.ClientIP(),
	}
}

func (handler *AuthHandler) setRefreshCookie(c *gin.Context, raw string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(refreshCookieName, raw, int(handler.refreshTTL.Seconds()), refreshCookiePath, "", handler.cookieSecure, true)
}

func (handler *AuthHandler) clearRefreshCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(refreshCookieName, "", -1, refreshCookiePath, "", handler.cookieSecure, true)
}

func (handler *AuthHandler) clearStateCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(stateCookieName, "", -1, refreshCookiePath, "", handler.cookieSecure, true)
}

func invalidBody() *errorResponse {
	return &errorResponse{Code: ErrorCodeInvalidRequest, Message: MessageInvalidRequestBody}
}

// respondWithAuthError maps auth service errors to HTTP status codes and safe
// messages. It never leaks the raw error to the client.
func respondWithAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, auth.ErrInvalidCredentials):
		respondWithError(c, http.StatusUnauthorized, ErrorCodeUnauthorized, MessageInvalidCredentials)
	case errors.Is(err, auth.ErrRefreshInvalid):
		respondWithError(c, http.StatusUnauthorized, ErrorCodeUnauthorized, MessageSessionExpired)
	case errors.Is(err, auth.ErrInvalidEmail):
		respondWithError(c, http.StatusUnprocessableEntity, ErrorCodeDomainRejected, MessageInvalidEmail)
	case errors.Is(err, auth.ErrWeakPassword):
		respondWithError(c, http.StatusUnprocessableEntity, ErrorCodeDomainRejected, MessageWeakPassword)
	case errors.Is(err, auth.ErrEmailTaken):
		respondWithError(c, http.StatusUnprocessableEntity, ErrorCodeDomainRejected, MessageEmailTaken)
	case errors.Is(err, auth.ErrEmailNotVerified):
		respondWithError(c, http.StatusForbidden, ErrorCodeForbidden, MessageEmailNotVerified)
	case errors.Is(err, auth.ErrAccountExistsUsePassword):
		respondWithError(c, http.StatusConflict, ErrorCodeConflict, MessageAccountExistsUsePassword)
	case errors.Is(err, auth.ErrGoogleSubTaken):
		respondWithError(c, http.StatusConflict, ErrorCodeConflict, MessageGoogleAccountLinked)
	default:
		respondWithInternalServerError(c, MessageAuthFailed)
	}
}
