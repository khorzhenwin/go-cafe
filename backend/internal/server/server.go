package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/khorzhenwin/go-cafe/backend/internal/auth"
	"github.com/khorzhenwin/go-cafe/backend/internal/cafelisting"
	appconfig "github.com/khorzhenwin/go-cafe/backend/internal/config"
	"github.com/khorzhenwin/go-cafe/backend/internal/rating"
	"github.com/khorzhenwin/go-cafe/backend/internal/user"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

// Config holds server and API config.
type Config struct {
	BasePath    string
	Address     string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

// New builds the HTTP handler from DB connection and configs. Caller must run migrations separately.
func New(dbConn *gorm.DB, authCfg *appconfig.AuthConfig, srvCfg Config) http.Handler {
	userRepo := user.NewRepository(dbConn)
	userSvc := user.NewService(userRepo)
	cafeRepo := cafelisting.NewRepository(dbConn)
	cafeSvc := cafelisting.NewService(cafeRepo)
	ratingRepo := rating.NewRepository(dbConn)
	ratingSvc := rating.NewService(ratingRepo, cafeSvc)

	authMiddleware := auth.Middleware(authCfg)
	authHandler := &auth.Handler{AuthCfg: authCfg, Finder: userSvc, Creator: userSvc}

	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Route(srvCfg.BasePath, func(r chi.Router) {
		auth.RegisterRoutes(r, authHandler)
		user.RegisterRoutes(r, userSvc)
		cafelisting.RegisterRoutes(r, cafeSvc, authMiddleware)
		rating.RegisterRoutes(r, ratingSvc, authMiddleware)
	})
	return r
}

// NewServer returns an http.Server using the same handler (for ListenAndServe).
func NewServer(handler http.Handler, cfg Config) *http.Server {
	return &http.Server{
		Addr:         cfg.Address,
		Handler:      handler,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
	}
}
