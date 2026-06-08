package main

import (
	"net/http"
	"os"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/app"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/config"
)

func main() {
	cfg := config.Load()
	application, err := app.New(cfg)
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
	defer application.Close()

	application.Logger.Info("api listening", "addr", application.Server.Addr)
	if err := application.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		application.Logger.Error("api server stopped", "error", err)
		os.Exit(1)
	}
}
