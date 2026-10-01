package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
	context, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
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

	_, err = database.NewPostgres(context, cfg)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
