package app

import (
	"context"

	"github.com/SJMaks/financer_back/internal/core/config"
	"github.com/SJMaks/financer_back/internal/core/httpserver"
	usersRepo "github.com/SJMaks/financer_back/internal/features/users/repository"
	usersService "github.com/SJMaks/financer_back/internal/features/users/service"
	usersTransport "github.com/SJMaks/financer_back/internal/features/users/transport"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type App struct {
	server *httpserver.Server
}

func NewApp(ctx context.Context, log *zap.SugaredLogger, db *pgxpool.Pool, cfg *config.Config) (*App, error) {
	engine := gin.New()
	engine.Use(gin.Recovery())

	api := engine.Group("/api/v1")

	usersRepo := usersRepo.NewPostgres(db)
	usersSrv := usersService.NewService(usersRepo)
	usersHandler := usersTransport.NewHandler(usersSrv, log)

	usersTransport.RegisterRouter(api, usersHandler)

	return &App{
		server: httpserver.New(engine, cfg.HTTPAddr),
	}, nil
}

func (a *App) Run() error {
	return a.server.Start()
}

func (a *App) Stop(ctx context.Context) error {
	return a.server.Stop(ctx)
}
