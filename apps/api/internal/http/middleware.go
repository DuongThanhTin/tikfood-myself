package http

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/auth"
	"github.com/gin-gonic/gin"
)

const (
	requestIDHeader   = "X-Request-ID"
	contextUserIDKey  = "user_id"
	authHeaderName    = "Authorization"
	bearerSchemeLower = "bearer"
)

// accessTokenParser validates an access token and returns its claims. *auth.AuthService
// satisfies it, keeping the JWT secret encapsulated in the auth domain.
type accessTokenParser interface {
	ParseAccessToken(raw string) (auth.AccessClaims, error)
}

// authMiddleware requires a valid Bearer access token, injecting the user id into the
// context for downstream handlers. It aborts with a 401 envelope otherwise.
func authMiddleware(parser accessTokenParser) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, ok := bearerToken(c.GetHeader(authHeaderName))
		if !ok {
			abortUnauthorized(c)
			return
		}
		claims, err := parser.ParseAccessToken(raw)
		if err != nil {
			abortUnauthorized(c)
			return
		}
		c.Set(contextUserIDKey, claims.UserID)
		c.Next()
	}
}

// corsMiddleware allows credentialed cross-origin requests only from an explicit
// allow-list. It never emits `Access-Control-Allow-Origin: *` alongside credentials.
func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		if origin != "" {
			allowed[origin] = struct{}{}
		}
	}
	return func(c *gin.Context) {
		// Set Vary: Origin on every response the middleware touches (not just the allowed
		// branch) so a shared cache never serves one origin's ACAO header to another.
		c.Header("Vary", "Origin")
		origin := c.GetHeader("Origin")
		if origin != "" {
			if _, ok := allowed[origin]; ok {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Credentials", "true")
				c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
			}
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// currentUserID returns the authenticated user id set by authMiddleware, or "".
func currentUserID(c *gin.Context) string {
	if value, ok := c.Get(contextUserIDKey); ok {
		if id, ok := value.(string); ok {
			return id
		}
	}
	return ""
}

func bearerToken(header string) (string, bool) {
	header = strings.TrimSpace(header)
	if header == "" {
		return "", false
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != bearerSchemeLower {
		return "", false
	}
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", false
	}
	return token, true
}

func abortUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, response{
		Data:  nil,
		Error: &errorResponse{Code: ErrorCodeUnauthorized, Message: MessageSessionExpired},
	})
}

func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(requestIDHeader)
		if !validRequestID(requestID) {
			requestID = newRequestID()
		}

		c.Set("request_id", requestID)
		c.Header(requestIDHeader, requestID)
		c.Next()
	}
}

func requestLoggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		c.Next()

		requestID, _ := c.Get("request_id")
		logger.InfoContext(
			c.Request.Context(),
			"http_request",
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.FullPath(),
			"status", c.Writer.Status(),
			"duration_ms", time.Since(startedAt).Milliseconds(),
		)
	}
}

func recoveryMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				requestID, _ := c.Get("request_id")
				logger.ErrorContext(
					c.Request.Context(),
					"panic_recovered",
					"request_id", requestID,
					"method", c.Request.Method,
					"path", c.FullPath(),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, response{
					Data: nil,
					Error: &errorResponse{
						Code:    "internal_error",
						Message: "Unexpected server error.",
					},
				})
			}
		}()
		c.Next()
	}
}

func validRequestID(value string) bool {
	if len(value) < 8 || len(value) > 128 {
		return false
	}
	return !strings.ContainsAny(value, "\r\n\t ")
}

func newRequestID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "req_" + hex.EncodeToString([]byte(time.Now().UTC().Format("20060102150405.000000000")))
	}
	return "req_" + hex.EncodeToString(bytes[:])
}
