package config

import (
	"os"

	"github.com/SJMaks/financer_back/internal/core/errors"
)

type Config struct {
	DBConnString string
}

func LoadConfig() (*Config, error) {
	dbConnStr, ok := os.LookupEnv("DB_URL")
	if !ok {
		return nil, errors.EnvNotFound
	}

	return &Config{
		DBConnString: dbConnStr,
	}, nil
}
