package app

import (
	"log/slog"
	"net/http"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/config"
	apihttp "github.com/DuongThanhTin/tikfood-myself/apps/api/internal/http"
)

type App struct {
	Config config.Config
	Logger *slog.Logger
	Server *http.Server
	close  func() error
}

func New(cfg config.Config) (*App, error) {
	container, err := NewContainer(cfg)
	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr: ":" + cfg.Port,
		Handler: apihttp.NewRouter(apihttp.RouterDependencies{
			Logger:          container.Logger,
			RouteRegistrars: container.RouteRegistrars,
		}),
	}

	return &App{
		Config: cfg,
		Logger: container.Logger,
		Server: server,
		close:  container.Close,
	}, nil
}

func (app *App) Close() error {
	if app.close == nil {
		return nil
	}
	return app.close()
}
