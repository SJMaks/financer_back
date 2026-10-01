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
)

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
		err = closeFile()
	}()

	db, err := database.NewPostgres(ctx, cfg)
	if err != nil {
		log.Error(err)
		return err
	}

	app, err := app.NewApp(ctx, log, db, cfg)
	if err != nil {
		log.Error(err)
		return err
	}

	errChan := make(chan error, 1)

	go func() {
		log.Info("http server starting on addr: ", cfg.HTTPAddr)
		if err := app.Run(); err != nil && errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("start http server: %w", err)
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		log.Info("shutdown signal recieved")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return app.Stop(shutdownCtx)
}
