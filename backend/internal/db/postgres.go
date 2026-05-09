package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const connectTimeout = 5 * time.Second

func OpenPostgres(ctx context.Context, databaseURL string) (*sql.DB, error) {
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	pool, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	if err := pool.PingContext(pingCtx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
