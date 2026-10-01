package config

import (
	"os"

	"github.com/SJMaks/financer_back/internal/core/errors"
)

type Config struct {
	HTTPAddr string

	DBConnString string
}

func LoadConfig() (*Config, error) {
	httpAddr, ok := os.LookupEnv("HTTP_ADDRESS")
	if !ok {
		return nil, errors.EnvNotFound
	}

	dbConnStr, ok := os.LookupEnv("DB_URL")
	if !ok {
		return nil, errors.EnvNotFound
	}

	return &Config{
		HTTPAddr:     httpAddr,
		DBConnString: dbConnStr,
	}, nil
}
