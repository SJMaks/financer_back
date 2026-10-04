package app

import (
	"context"

	"github.com/SJMaks/financer_back/internal/core/auth"
	"github.com/SJMaks/financer_back/internal/core/config"
	"github.com/SJMaks/financer_back/internal/core/httpserver"
	"github.com/SJMaks/financer_back/internal/core/middleware"
	accountsRepo "github.com/SJMaks/financer_back/internal/features/accounts/repository"
	accountsService "github.com/SJMaks/financer_back/internal/features/accounts/service"
	accountsTransport "github.com/SJMaks/financer_back/internal/features/accounts/transport"
	categoriesRepo "github.com/SJMaks/financer_back/internal/features/categories/repository"
	categoriesService "github.com/SJMaks/financer_back/internal/features/categories/service"
	categoriesTransport "github.com/SJMaks/financer_back/internal/features/categories/transport"
	operationTypesRepo "github.com/SJMaks/financer_back/internal/features/operation_types/repository"
	operationTypesService "github.com/SJMaks/financer_back/internal/features/operation_types/service"
	operationTypesTransport "github.com/SJMaks/financer_back/internal/features/operation_types/transport"
	operationsRepo "github.com/SJMaks/financer_back/internal/features/operations/repository"
	operationsService "github.com/SJMaks/financer_back/internal/features/operations/service"
	operationsTransport "github.com/SJMaks/financer_back/internal/features/operations/transport"
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

	usersTransport.RegisterPublicRouter(api, usersHandler)
	usersTransport.RegisterPrivateRouter(api, usersHandler, authMiddleware.CheckAuth)

	accsRepo := accountsRepo.NewPostgres(db)
	accsSrv := accountsService.NewService(accsRepo)
	accsHandler := accountsTransport.NewHandler(log, accsSrv)

	accountsTransport.RegisterPrivateRouter(api, accsHandler, authMiddleware.CheckAuth)

	ctgRepo := categoriesRepo.NewPostgres(db)
	ctgSrv := categoriesService.NewService(ctgRepo)
	ctgHandler := categoriesTransport.NewHandler(log, ctgSrv)

	categoriesTransport.RegisterPrivateRouter(api, ctgHandler, authMiddleware.CheckAuth)

	opTypesRepo := operationTypesRepo.NewPostgres(db)
	opTypesSrv := operationTypesService.NewService(opTypesRepo)
	opTypesHandler := operationTypesTransport.NewHandler(log, opTypesSrv)

	operationTypesTransport.RegisterPrivateRouter(api, opTypesHandler, authMiddleware.CheckAuth)

	opsRepo := operationsRepo.NewPostgres(db)
	opsSrv := operationsService.NewService(opsRepo)
	opsHandler := operationsTransport.NewHandler(log, opsSrv)

	operationsTransport.RegisterPrivateRouter(api, opsHandler, authMiddleware.CheckAuth)

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
