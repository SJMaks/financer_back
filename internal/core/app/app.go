package app

import (
	"context"
	"net/http"

	"github.com/SJMaks/financer_back/internal/core/auth"
	"github.com/SJMaks/financer_back/internal/core/config"
	"github.com/SJMaks/financer_back/internal/core/httpserver"
	"github.com/SJMaks/financer_back/internal/core/middleware"
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

func NewApp(log *zap.SugaredLogger, db *pgxpool.Pool, cfg *config.Config) *App {
	engine := gin.New()
	engine.Use(gin.Recovery())

	jwt := auth.NewJWT(cfg)

	authMiddleware := middleware.NewAuthMiddleware(jwt)

	api := engine.Group("/api/v1")

	usersRepo := usersRepo.NewPostgres(db)
	usersSrv := usersService.NewService(usersRepo, jwt)
	usersHandler := usersTransport.NewHandler(usersSrv, log)

	g := api.Group("/ping", authMiddleware.CheckAuth)
	{
		g.GET("/", func(c *gin.Context) {
			p, ok := middleware.UserFromContext(c)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"id": p.ID, "login": p.Login})
		})
	}

	usersTransport.RegisterRouter(api, usersHandler)

	return &App{
		server: httpserver.New(engine, cfg.HTTPAddr),
	}
}

func (a *App) Run() error {
	return a.server.Start()
}

func (a *App) Stop(ctx context.Context) error {
	return a.server.Stop(ctx)
}
