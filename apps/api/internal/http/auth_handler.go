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
)

// AuthHandler exposes the authentication endpoints. It parses/validates requests,
// calls the auth service, sets/clears the refresh cookie, and responds with the
// standard {data,error} envelope. It holds no business logic.
type AuthHandler struct {
	auth         *auth.AuthService
	refreshTTL   time.Duration
	cookieSecure bool
}

func NewAuthHandler(service *auth.AuthService, refreshTTL time.Duration, cookieSecure bool) *AuthHandler {
	if service == nil {
		panic("http.NewAuthHandler requires an auth service")
	}
	return &AuthHandler{auth: service, refreshTTL: refreshTTL, cookieSecure: cookieSecure}
}

func (handler *AuthHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	group := v1.Group("/auth")
	group.POST("/register", handler.Register)
	group.POST("/login", handler.Login)
	group.POST("/refresh", handler.Refresh)
	group.POST("/logout", handler.Logout)

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
	if err := handler.auth.Logout(c.Request.Context(), raw); err != nil {
		respondWithInternalServerError(c, MessageAuthFailed)
		return
	}

	handler.clearRefreshCookie(c)
	respondWithData(c, gin.H{"logged_out": true})
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
	default:
		respondWithInternalServerError(c, MessageAuthFailed)
	}
}
