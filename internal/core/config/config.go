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

	AllowedOrigins map[string]struct{}
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

	allowedOrigins, err := loadAllowedOrigins()
	if err != nil {
		return nil, err
	}

	dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	return &Config{
		HTTPAddr:       httpAddr,
		JWTSecret:      jwtSecret,
		DBConnString:   dbConnStr,
		AllowedOrigins: allowedOrigins,
	}, nil
}

func loadAllowedOrigins() (map[string]struct{}, error) {
	raw, err := loadEnv("CORS_ALLOWED_ORIGINS")
	if err != nil {
		return nil, err
	}

	origins := make(map[string]struct{})
	for _, o := range strings.Split(raw, ",") {
		o = strings.TrimSpace(o)
		if o == "" {
			continue
		}
		origins[o] = struct{}{}
	}

	if len(origins) == 0 {
		return nil, fmt.Errorf("%w: CORS_ALLOWED_ORIGINS", apperrors.EmptyEnv)
	}

	return origins, nil
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
