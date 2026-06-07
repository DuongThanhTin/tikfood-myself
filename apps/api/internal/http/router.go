package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RouteRegistrar interface {
	RegisterRoutes(v1 *gin.RouterGroup)
}

type RouterDependencies struct {
	Logger          *slog.Logger
	RouteRegistrars []RouteRegistrar
}

func NewRouter(deps RouterDependencies) http.Handler {
	logger := deps.Logger
	if logger == nil {
		logger = slog.Default()
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(recoveryMiddleware(logger), requestIDMiddleware(), requestLoggingMiddleware(logger))
	_ = router.SetTrustedProxies(nil)

	router.GET("/health", func(c *gin.Context) {
		respondWithData(c, map[string]bool{"ok": true})
	})

	v1 := router.Group("/api/v1")
	for _, registrar := range deps.RouteRegistrars {
		if registrar == nil {
			continue
		}
		registrar.RegisterRoutes(v1)
	}

	return router
}
