// Command cleanuptokens purges refresh_tokens rows that can no longer authenticate
// (expired or revoked), bounding table growth. Run it on a schedule (cron / periodic
// job). It requires DATABASE_URL; there is nothing to clean in no-DB mode.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/config"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/storage/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})).With("service", "cleanuptokens")
	if err := run(logger); err != nil {
		logger.Error("refresh-token cleanup failed", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer func() { _ = db.Close() }()
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	repo := postgres.NewRefreshTokenRepository(db)
	removed, err := repo.PurgeExpiredRefreshTokens(context.Background(), time.Now())
	if err != nil {
		return err
	}

	logger.Info("refresh-token cleanup complete", "removed", removed)
	return nil
}
