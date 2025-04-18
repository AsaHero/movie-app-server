package app

import (
	"context"
	"net/http"
	"time"

	"github.com/AsaHero/movie-app-server/delivery/api"
	"github.com/AsaHero/movie-app-server/delivery/api/handlers"
	"github.com/AsaHero/movie-app-server/delivery/api/validation"
	genres_repo "github.com/AsaHero/movie-app-server/internal/repository/genres"
	"github.com/AsaHero/movie-app-server/internal/repository/movie_genres"
	movies_repo "github.com/AsaHero/movie-app-server/internal/repository/movies"
	users_repo "github.com/AsaHero/movie-app-server/internal/repository/users"
	"github.com/AsaHero/movie-app-server/internal/service/auth"
	"github.com/AsaHero/movie-app-server/internal/service/genres"
	"github.com/AsaHero/movie-app-server/internal/service/movies"
	"github.com/AsaHero/movie-app-server/internal/service/users"
	"github.com/AsaHero/movie-app-server/pkg/config"
	"github.com/AsaHero/movie-app-server/pkg/database/postgres"
	"github.com/AsaHero/movie-app-server/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func Run() {
	x := fx.New(
		fx.Provide(
			config.New,
			func(cfg *config.Config) string { return cfg.APP + ".log" },
			logger.Init,
			postgres.New,
			genres_repo.New,
			movie_genres.New,
			users_repo.New,
			movies_repo.New,
			// timeout provider
			func(cfg *config.Config) time.Duration {
				d, err := time.ParseDuration(cfg.Context.Timeout)
				if err != nil {
					panic(err)
				}
				return d
			},
			auth.New,
			users.New,
			genres.New,
			movies.New,
			validation.NewValidator,
			func(
				cfg *config.Config,
				validator *validation.Validator,
				authSvc auth.Service,
				userSvc users.Service,
				movieSvc movies.Service,
				genresSvc genres.Service,
			) *handlers.HandlerOptions {
				return &handlers.HandlerOptions{
					Config:        cfg,
					Validator:     validator,
					AuthService:   authSvc,
					UsersService:  userSvc,
					MoviesSerive:  movieSvc,
					GenresService: genresSvc,
				}
			},
			api.NewRouter,
			func(router *gin.Engine) http.Handler { return router },
			api.NewServer,
		),
		fx.Invoke(registerHooks),
	)

	x.Run()
}

func registerHooks(
	lc fx.Lifecycle,
	cfg *config.Config,
	log *logrus.Logger,
	server *http.Server,
	db *gorm.DB,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info(cfg.APP, "starting server...")
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Error("server error:", err.Error())
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info(cfg.APP, "shutting down...")
			if err := server.Shutdown(ctx); err != nil {
				return err
			}
			sqlDB, _ := db.DB()
			return sqlDB.Close()
		},
	})
}
