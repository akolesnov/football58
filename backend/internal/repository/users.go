package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/akolesnov/football58/backend/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, name string, telegramID *int64, telegramUsername *string) (domain.User, error) {
	const query = `
INSERT INTO users (name, telegram_id, telegram_username)
VALUES ($1, $2, $3)
RETURNING id, name, telegram_id, telegram_username, created_at;`

	var user domain.User
	var nullableTelegramID sql.NullInt64
	var nullableTelegramUsername sql.NullString

	if err := r.db.QueryRowContext(ctx, query, name, telegramID, telegramUsername).Scan(
		&user.ID,
		&user.Name,
		&nullableTelegramID,
		&nullableTelegramUsername,
		&user.CreatedAt,
	); err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	user.TelegramID = nullableInt64Ptr(nullableTelegramID)
	user.TelegramUsername = nullableStringPtr(nullableTelegramUsername)
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (domain.User, error) {
	const query = `
SELECT id, name, telegram_id, telegram_username, created_at
FROM users
WHERE id = $1;`

	return r.getOne(ctx, query, id)
}

func (r *UserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (domain.User, error) {
	const query = `
SELECT id, name, telegram_id, telegram_username, created_at
FROM users
WHERE telegram_id = $1;`

	return r.getOne(ctx, query, telegramID)
}

func (r *UserRepository) UpsertTelegram(ctx context.Context, name string, telegramID int64, telegramUsername *string) (domain.User, error) {
	const query = `
INSERT INTO users (name, telegram_id, telegram_username)
VALUES ($1, $2, $3)
ON CONFLICT (telegram_id) DO UPDATE
SET name = EXCLUDED.name,
    telegram_username = EXCLUDED.telegram_username
RETURNING id, name, telegram_id, telegram_username, created_at;`

	var user domain.User
	var nullableTelegramID sql.NullInt64
	var nullableTelegramUsername sql.NullString

	if err := r.db.QueryRowContext(ctx, query, name, telegramID, telegramUsername).Scan(
		&user.ID,
		&user.Name,
		&nullableTelegramID,
		&nullableTelegramUsername,
		&user.CreatedAt,
	); err != nil {
		return domain.User{}, fmt.Errorf("upsert telegram user: %w", err)
	}

	user.TelegramID = nullableInt64Ptr(nullableTelegramID)
	user.TelegramUsername = nullableStringPtr(nullableTelegramUsername)
	return user, nil
}

func (r *UserRepository) getOne(ctx context.Context, query string, args ...any) (domain.User, error) {
	var user domain.User
	var nullableTelegramID sql.NullInt64
	var nullableTelegramUsername sql.NullString

	if err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&nullableTelegramID,
		&nullableTelegramUsername,
		&user.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}

		return domain.User{}, fmt.Errorf("get user: %w", err)
	}

	user.TelegramID = nullableInt64Ptr(nullableTelegramID)
	user.TelegramUsername = nullableStringPtr(nullableTelegramUsername)
	return user, nil
}

func nullableInt64Ptr(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}

	return &value.Int64
}
