package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/SJMaks/financer_back/internal/core/apperrors"
)

type Config struct {
	HTTPAddr string

	JWTSecret string

	DBConnString string
}

func LoadConfig() (*Config, error) {
	httpAddr, err := loadEnv("HTTP_ADDRESS")
	if err != nil {
		return nil, err
	}

	dbHost, err := loadEnv("DB_HOST")
	if err != nil {
		return nil, err
	}

	dbPort, err := loadEnv("DB_PORT")
	if err != nil {
		return nil, err
	}

	dbUser, err := loadEnv("DB_USER")
	if err != nil {
		return nil, err
	}

	dbPass, err := loadEnv("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	dbName, err := loadEnv("DB_NAME")
	if err != nil {
		return nil, err
	}

	jwtSecret, err := loadEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}

	dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	return &Config{
		HTTPAddr:     httpAddr,
		JWTSecret:    jwtSecret,
		DBConnString: dbConnStr,
	}, nil
}

func loadEnv(name string) (string, error) {
	env, ok := os.LookupEnv(name)
	if !ok {
		return "", fmt.Errorf("%w: %s", apperrors.EnvNotFound, name)
	}
	if len(strings.TrimSpace(env)) < 1 {
		return "", fmt.Errorf("%w: %s", apperrors.EmptyEnv, name)
	}
	return env, nil
}
