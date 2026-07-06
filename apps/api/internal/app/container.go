package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/auth"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/config"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/discovery"
	apihttp "github.com/DuongThanhTin/tikfood-myself/apps/api/internal/http"
	"github.com/DuongThanhTin/tikfood-myself/apps/api/internal/storage/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Container struct {
	Logger          *slog.Logger
	VenueRepository discovery.VenueRepository
	VenueService    *discovery.VenueService
	AuthService     *auth.AuthService
	RouteRegistrars []apihttp.RouteRegistrar
	close           func() error
}

func NewContainer(cfg config.Config) (*Container, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})).With("service", "api")

	db, closeDB, err := openDatabase(cfg, logger)
	if err != nil {
		return nil, err
	}

	venueRepo := buildVenueRepository(cfg, logger, db)
	venueService := discovery.NewVenueService(venueRepo)

	authService, err := buildAuthService(cfg, logger, db)
	if err != nil {
		if closeDB != nil {
			_ = closeDB()
		}
		return nil, err
	}

	// nil when Google credentials are not configured; that disables the Google routes.
	googleAuth := auth.NewGoogleOAuth(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL)
	if googleAuth == nil {
		logger.Info("google oauth disabled (credentials not configured)")
	}

	routeRegistrars := apihttp.DefaultRouteRegistrars(apihttp.HandlerDependencies{
		Venues:       venueService,
		Auth:         authService,
		Google:       googleAuthOrNil(googleAuth),
		RefreshTTL:   cfg.RefreshTokenTTL,
		CookieSecure: cfg.CookieSecure,
		WebOrigin:    cfg.WebOrigin,
	})

	return &Container{
		Logger:          logger,
		VenueRepository: venueRepo,
		VenueService:    venueService,
		AuthService:     authService,
		RouteRegistrars: routeRegistrars,
		close:           closeDB,
	}, nil
}

func (container *Container) Close() error {
	if container.close == nil {
		return nil
	}
	return container.close()
}

// openDatabase opens and pings the database when DATABASE_URL is set, returning the
// shared handle and its closer. When unset, it returns (nil, nil, nil) and the app
// runs on in-memory repositories.
func openDatabase(cfg config.Config, logger *slog.Logger) (*sql.DB, func() error, error) {
	if cfg.DatabaseURL == "" {
		return nil, nil, nil
	}

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, nil, fmt.Errorf("ping database: %w", err)
	}
	logger.Info("using postgres storage")
	return db, db.Close, nil
}

func buildVenueRepository(cfg config.Config, logger *slog.Logger, db *sql.DB) discovery.VenueRepository {
	if db == nil {
		logger.Info("using in-memory discovery storage")
		repo := discovery.NewFallbackVenueRepository()
		if district := strings.TrimSpace(os.Getenv("INGEST_ON_START")); district != "" {
			limit := envInt("INGEST_LIMIT", 12)
			if err := seedFallbackFromOSM(repo, cfg.OverpassEndpoint, district, limit, logger); err != nil {
				logger.Warn("OSM preview seed failed; serving base seed only", "error", err)
			}
		}
		return repo
	}
	return postgres.NewDiscoveryRepository(db)
}

func buildAuthService(cfg config.Config, logger *slog.Logger, db *sql.DB) (*auth.AuthService, error) {
	secret := cfg.JWTSecret
	if secret == "" {
		if db != nil {
			// A real (non-dev) deployment must never run with an empty JWT secret.
			return nil, fmt.Errorf("JWT_SECRET is required when DATABASE_URL is set")
		}
		generated, err := auth.RandomHexToken()
		if err != nil {
			// Fail closed: never fall back to a predictable secret. A boot without usable
			// randomness must abort rather than silently sign tokens with a guessable key.
			return nil, fmt.Errorf("generate ephemeral dev jwt secret: %w", err)
		}
		secret = generated
		logger.Warn("JWT_SECRET not set; using an ephemeral development secret (access tokens do not survive a restart)")
	}

	issuer, err := auth.NewTokenIssuer(secret, cfg.AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("build token issuer: %w", err)
	}

	var userRepo auth.UserRepository
	var refreshRepo auth.RefreshTokenRepository
	if db == nil {
		logger.Info("using in-memory auth storage")
		userRepo = auth.NewMemoryUserRepository()
		refreshRepo = auth.NewMemoryRefreshTokenRepository()
	} else {
		userRepo = postgres.NewUserRepository(db)
		refreshRepo = postgres.NewRefreshTokenRepository(db)
	}

	return auth.NewAuthService(userRepo, refreshRepo, issuer, cfg.RefreshTokenTTL), nil
}

// googleAuthOrNil converts a possibly-nil *auth.GoogleOAuth into a genuinely nil
// interface value, avoiding the "typed nil in an interface is non-nil" trap that would
// make the handler register Google routes over a nil authenticator.
func googleAuthOrNil(g *auth.GoogleOAuth) auth.GoogleAuthenticator {
	if g == nil {
		return nil
	}
	return g
}

func envInt(key string, fallback int) int {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return fallback
}
