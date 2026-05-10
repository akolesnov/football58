package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
)

const (
	envHTTPAddr         = "HTTP_ADDR"
	envDatabaseURL      = "DATABASE_URL"
	envPostgresUser     = "POSTGRES_USER"
	envPostgresPassword = "POSTGRES_PASSWORD"
	envPostgresHost     = "POSTGRES_HOST"
	envPostgresPort     = "POSTGRES_PORT"
	envPostgresDB       = "POSTGRES_DB"
	envPostgresSSLMode  = "POSTGRES_SSLMODE"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
}

func Load() (Config, error) {
	httpAddr, err := requiredEnv(envHTTPAddr)
	if err != nil {
		return Config{}, err
	}

	databaseURL, err := databaseURL()
	if err != nil {
		return Config{}, err
	}

	return Config{
		HTTPAddr:    httpAddr,
		DatabaseURL: databaseURL,
	}, nil
}

func databaseURL() (string, error) {
	value := os.Getenv(envDatabaseURL)
	if value != "" {
		return value, nil
	}

	user, err := requiredEnv(envPostgresUser)
	if err != nil {
		return "", err
	}
	password, err := requiredEnv(envPostgresPassword)
	if err != nil {
		return "", err
	}
	host, err := requiredEnv(envPostgresHost)
	if err != nil {
		return "", err
	}
	port, err := requiredEnv(envPostgresPort)
	if err != nil {
		return "", err
	}
	name, err := requiredEnv(envPostgresDB)
	if err != nil {
		return "", err
	}
	sslMode, err := requiredEnv(envPostgresSSLMode)
	if err != nil {
		return "", err
	}

	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(host, port),
		Path:   name,
	}

	query := dsn.Query()
	query.Set("sslmode", sslMode)
	dsn.RawQuery = query.Encode()

	return dsn.String(), nil
}

func requiredEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("не задана обязательная переменная окружения: %s", key)
	}

	return value, nil
}
