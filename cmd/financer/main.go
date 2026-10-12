package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SJMaks/financer_back/internal/core/app"
	"github.com/SJMaks/financer_back/internal/core/config"
	"github.com/SJMaks/financer_back/internal/core/database"
	"github.com/SJMaks/financer_back/internal/core/logger"

	_ "github.com/SJMaks/financer_back/docs"
)

// @title       Financer API
// @version     1.0
// @description Finance Application REST-API scheme
// @host        api.financerapp.ru
// @BasePath    /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Type "Bearer" followed by a space and the JWT token.

func main() {
	if err := run(); err != nil {
		fmt.Printf("error while starting the app: %s\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	log, closeFile, err := logger.NewZap()
	if err != nil {
		return err
	}
	defer func() {
		if cerr := closeFile(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	db, err := database.NewPostgres(ctx, cfg)
	if err != nil {
		log.Error(err)
		return err
	}
	defer db.Close()

	app := app.NewApp(log, db, cfg)

	errChan := make(chan error, 1)

	go func() {
		log.Infow("http server starting", "addr", cfg.HTTPAddr)
		if err := app.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("start http server: %w", err)
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		log.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return app.Stop(shutdownCtx)
}
