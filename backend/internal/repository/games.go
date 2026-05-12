package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/akolesnov/football58/backend/internal/domain"
)

type GameRepository struct {
	db dbtx
}

func NewGameRepository(db dbtx) *GameRepository {
	return &GameRepository{db: db}
}

func (r *GameRepository) Create(ctx context.Context, game domain.Game) (domain.Game, error) {
	const query = `
INSERT INTO games (
    venue_id,
    starts_at,
    duration_minutes,
    min_players,
    max_players,
    price_rub,
    notes,
    created_by_user_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, venue_id, starts_at, duration_minutes, min_players, max_players,
    price_rub, notes, status, created_by_user_id, telegram_chat_id,
    telegram_message_id, version, created_at;`

	created, err := scanGame(r.db.QueryRowContext(
		ctx,
		query,
		game.VenueID,
		game.StartsAt,
		game.DurationMinutes,
		game.MinPlayers,
		game.MaxPlayers,
		game.PriceRub,
		game.Notes,
		game.CreatedByUserID,
	))
	if err != nil {
		return domain.Game{}, fmt.Errorf("create game: %w", err)
	}

	return created, nil
}

func (r *GameRepository) GetByID(ctx context.Context, id int64) (domain.Game, error) {
	const query = `
SELECT id, venue_id, starts_at, duration_minutes, min_players, max_players,
    price_rub, notes, status, created_by_user_id, telegram_chat_id,
    telegram_message_id, version, created_at
FROM games
WHERE id = $1;`

	game, err := scanGame(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Game{}, domain.ErrNotFound
		}

		return domain.Game{}, fmt.Errorf("get game by id: %w", err)
	}

	return game, nil
}

func (r *GameRepository) GetByIDForUpdate(ctx context.Context, id int64) (domain.Game, error) {
	const query = `
SELECT id, venue_id, starts_at, duration_minutes, min_players, max_players,
    price_rub, notes, status, created_by_user_id, telegram_chat_id,
    telegram_message_id, version, created_at
FROM games
WHERE id = $1
FOR UPDATE;`

	game, err := scanGame(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Game{}, domain.ErrNotFound
		}

		return domain.Game{}, fmt.Errorf("get game by id for update: %w", err)
	}

	return game, nil
}

func (r *GameRepository) List(ctx context.Context) ([]domain.Game, error) {
	const query = `
SELECT id, venue_id, starts_at, duration_minutes, min_players, max_players,
    price_rub, notes, status, created_by_user_id, telegram_chat_id,
    telegram_message_id, version, created_at
FROM games
ORDER BY starts_at DESC;`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list games: %w", err)
	}
	defer rows.Close()

	games := make([]domain.Game, 0)
	for rows.Next() {
		game, err := scanGame(rows)
		if err != nil {
			return nil, fmt.Errorf("scan game: %w", err)
		}

		games = append(games, game)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate games: %w", err)
	}

	return games, nil
}

func (r *GameRepository) Update(ctx context.Context, game domain.Game) (domain.Game, error) {
	const query = `
UPDATE games
SET venue_id = $2,
    starts_at = $3,
    duration_minutes = $4,
    min_players = $5,
    max_players = $6,
    price_rub = $7,
    notes = $8,
    version = version + 1
WHERE id = $1
RETURNING id, venue_id, starts_at, duration_minutes, min_players, max_players,
    price_rub, notes, status, created_by_user_id, telegram_chat_id,
    telegram_message_id, version, created_at;`

	updated, err := scanGame(r.db.QueryRowContext(
		ctx,
		query,
		game.ID,
		game.VenueID,
		game.StartsAt,
		game.DurationMinutes,
		game.MinPlayers,
		game.MaxPlayers,
		game.PriceRub,
		game.Notes,
	))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Game{}, domain.ErrNotFound
		}

		return domain.Game{}, fmt.Errorf("update game: %w", err)
	}

	return updated, nil
}

func (r *GameRepository) UpdateTelegramMessage(ctx context.Context, id int64, chatID, messageID int64) (domain.Game, error) {
	const query = `
UPDATE games
SET telegram_chat_id = $2,
    telegram_message_id = $3,
    version = version + 1
WHERE id = $1
RETURNING id, venue_id, starts_at, duration_minutes, min_players, max_players,
    price_rub, notes, status, created_by_user_id, telegram_chat_id,
    telegram_message_id, version, created_at;`

	game, err := scanGame(r.db.QueryRowContext(ctx, query, id, chatID, messageID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Game{}, domain.ErrNotFound
		}

		return domain.Game{}, fmt.Errorf("update telegram message: %w", err)
	}

	return game, nil
}

func (r *GameRepository) SetStatus(ctx context.Context, id int64, status string) (domain.Game, error) {
	const query = `
UPDATE games
SET status = $2,
    version = version + 1
WHERE id = $1
RETURNING id, venue_id, starts_at, duration_minutes, min_players, max_players,
    price_rub, notes, status, created_by_user_id, telegram_chat_id,
    telegram_message_id, version, created_at;`

	game, err := scanGame(r.db.QueryRowContext(ctx, query, id, status))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Game{}, domain.ErrNotFound
		}

		return domain.Game{}, fmt.Errorf("set game status: %w", err)
	}

	return game, nil
}

func (r *GameRepository) IncrementVersion(ctx context.Context, id int64) (domain.Game, error) {
	const query = `
UPDATE games
SET version = version + 1
WHERE id = $1
RETURNING id, venue_id, starts_at, duration_minutes, min_players, max_players,
    price_rub, notes, status, created_by_user_id, telegram_chat_id,
    telegram_message_id, version, created_at;`

	game, err := scanGame(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Game{}, domain.ErrNotFound
		}

		return domain.Game{}, fmt.Errorf("increment game version: %w", err)
	}

	return game, nil
}

type gameScanner interface {
	Scan(dest ...any) error
}

func scanGame(scanner gameScanner) (domain.Game, error) {
	var game domain.Game
	var nullableNotes sql.NullString
	var nullableCreatedByUserID sql.NullInt64
	var nullableTelegramChatID sql.NullInt64
	var nullableTelegramMessageID sql.NullInt64

	if err := scanner.Scan(
		&game.ID,
		&game.VenueID,
		&game.StartsAt,
		&game.DurationMinutes,
		&game.MinPlayers,
		&game.MaxPlayers,
		&game.PriceRub,
		&nullableNotes,
		&game.Status,
		&nullableCreatedByUserID,
		&nullableTelegramChatID,
		&nullableTelegramMessageID,
		&game.Version,
		&game.CreatedAt,
	); err != nil {
		return domain.Game{}, err
	}

	game.Notes = nullableStringPtr(nullableNotes)
	game.CreatedByUserID = nullableInt64Ptr(nullableCreatedByUserID)
	game.TelegramChatID = nullableInt64Ptr(nullableTelegramChatID)
	game.TelegramMessageID = nullableInt64Ptr(nullableTelegramMessageID)
	return game, nil
}
