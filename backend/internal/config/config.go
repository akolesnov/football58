package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
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
	envTelegramBotToken = "TELEGRAM_BOT_TOKEN"
	envAPIBaseURL       = "API_BASE_URL"
	envTelegramAdminIDs = "TELEGRAM_ADMIN_IDS"
)

type APIConfig struct {
	HTTPAddr    string
	DatabaseURL string
}

type BotConfig struct {
	TelegramBotToken string
	APIBaseURL       string
	TelegramAdminIDs []int64
}

func LoadAPI() (APIConfig, error) {
	httpAddr, err := requiredEnv(envHTTPAddr)
	if err != nil {
		return APIConfig{}, err
	}

	databaseURL, err := databaseURL()
	if err != nil {
		return APIConfig{}, err
	}

	return APIConfig{
		HTTPAddr:    httpAddr,
		DatabaseURL: databaseURL,
	}, nil
}

func LoadBot() (BotConfig, error) {
	token, err := requiredEnv(envTelegramBotToken)
	if err != nil {
		return BotConfig{}, err
	}

	apiBaseURL, err := requiredEnv(envAPIBaseURL)
	if err != nil {
		return BotConfig{}, err
	}

	adminIDs, err := parseInt64List(os.Getenv(envTelegramAdminIDs), envTelegramAdminIDs)
	if err != nil {
		return BotConfig{}, err
	}

	return BotConfig{
		TelegramBotToken: token,
		APIBaseURL:       apiBaseURL,
		TelegramAdminIDs: adminIDs,
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

func parseInt64List(value string, envName string) ([]int64, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}

	parts := strings.Split(value, ",")
	result := make([]int64, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		number, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("некорректное значение переменной окружения %s: %s", envName, part)
		}

		result = append(result, number)
	}

	return result, nil
}
