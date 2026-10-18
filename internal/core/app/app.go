package app

import (
	"context"
	"net/http"

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
	statisticsRepo "github.com/SJMaks/financer_back/internal/features/statistics/repository"
	statisticsService "github.com/SJMaks/financer_back/internal/features/statistics/service"
	statisticsTransport "github.com/SJMaks/financer_back/internal/features/statistics/transport"
	usersRepo "github.com/SJMaks/financer_back/internal/features/users/repository"
	usersService "github.com/SJMaks/financer_back/internal/features/users/service"
	usersTransport "github.com/SJMaks/financer_back/internal/features/users/transport"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type App struct {
	server *httpserver.Server
}

func NewApp(log *zap.SugaredLogger, db *pgxpool.Pool, cfg *config.Config, rds *redis.Client) *App {
	engine := gin.New()
	engine.Use(gin.Recovery())

	jwt := auth.NewJWT(cfg)

	commonMiddleware := middleware.NewCommonMiddleware(log, cfg)
	authMiddleware := middleware.NewAuthMiddleware(log, jwt)

	engine.Use(commonMiddleware.CORS)
	api := engine.Group("/api/v1", commonMiddleware.CORS)

	usersRepo := usersRepo.NewPostgres(db)
	usersSrv := usersService.NewService(usersRepo, jwt)
	usersHandler := usersTransport.NewHandler(usersSrv, log, authMiddleware)

	usersTransport.RegisterPublicRouter(api, usersHandler)
	usersTransport.RegisterPrivateRouter(api, usersHandler, authMiddleware.CheckAuth)

	accsRepo := accountsRepo.NewPostgres(db)
	accsSrv := accountsService.NewService(accsRepo)
	accsHandler := accountsTransport.NewHandler(log, accsSrv, authMiddleware)

	accountsTransport.RegisterPrivateRouter(api, accsHandler, authMiddleware.CheckAuth)

	ctgRepo := categoriesRepo.NewPostgres(db)
	ctgSrv := categoriesService.NewService(ctgRepo)
	ctgHandler := categoriesTransport.NewHandler(log, ctgSrv, authMiddleware)

	categoriesTransport.RegisterPrivateRouter(api, ctgHandler, authMiddleware.CheckAuth)

	opTypesRepo := operationTypesRepo.NewPostgres(db)
	opTypesSrv := operationTypesService.NewService(opTypesRepo)
	opTypesHandler := operationTypesTransport.NewHandler(log, opTypesSrv)

	operationTypesTransport.RegisterPrivateRouter(api, opTypesHandler, authMiddleware.CheckAuth)

	opsRepo := operationsRepo.NewPostgres(db)
	opsSrv := operationsService.NewService(opsRepo)
	opsHandler := operationsTransport.NewHandler(log, opsSrv, authMiddleware)

	operationsTransport.RegisterPrivateRouter(api, opsHandler, authMiddleware.CheckAuth)

	stcsRepo := statisticsRepo.NewPostgres(db)
	stcsSrv := statisticsService.NewService(stcsRepo, rds, log)
	stcsHandler := statisticsTransport.NewHandler(log, stcsSrv, authMiddleware)

	statisticsTransport.RegisterPrivateRouter(api, stcsHandler, authMiddleware.CheckAuth)

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler), commonMiddleware.CORS)
	engine.GET("/healthy", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Server is running and healthy")
	})

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
