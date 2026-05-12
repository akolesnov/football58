package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/akolesnov/football58/backend/internal/domain"
)

type GameMemberRepository struct {
	db dbtx
}

func NewGameMemberRepository(db dbtx) *GameMemberRepository {
	return &GameMemberRepository{db: db}
}

func (r *GameMemberRepository) Create(ctx context.Context, member domain.GameMember) (domain.GameMember, error) {
	const query = `
INSERT INTO game_members (
    game_id,
    user_id,
    added_by_user_id,
    name,
    status,
    source,
    position_number
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, game_id, user_id, added_by_user_id, name, status, source,
    position_number, created_at, cancelled_at;`

	created, err := scanGameMember(r.db.QueryRowContext(
		ctx,
		query,
		member.GameID,
		member.UserID,
		member.AddedByUserID,
		member.Name,
		member.Status,
		member.Source,
		member.PositionNumber,
	))
	if err != nil {
		return domain.GameMember{}, fmt.Errorf("create game member: %w", err)
	}

	return created, nil
}

func (r *GameMemberRepository) GetByID(ctx context.Context, id int64) (domain.GameMember, error) {
	const query = `
SELECT id, game_id, user_id, added_by_user_id, name, status, source,
    position_number, created_at, cancelled_at
FROM game_members
WHERE id = $1;`

	member, err := scanGameMember(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.GameMember{}, domain.ErrNotFound
		}

		return domain.GameMember{}, fmt.Errorf("get game member by id: %w", err)
	}

	return member, nil
}

func (r *GameMemberRepository) GetActiveByGameIDAndUserID(ctx context.Context, gameID, userID int64) (domain.GameMember, error) {
	const query = `
SELECT id, game_id, user_id, added_by_user_id, name, status, source,
    position_number, created_at, cancelled_at
FROM game_members
WHERE game_id = $1
  AND user_id = $2
  AND status IN ($3, $4)
ORDER BY position_number ASC NULLS LAST, created_at ASC
LIMIT 1;`

	member, err := scanGameMember(
		r.db.QueryRowContext(
			ctx,
			query,
			gameID,
			userID,
			domain.GameMemberStatusActive,
			domain.GameMemberStatusWaitlist,
		),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.GameMember{}, domain.ErrNotFound
		}

		return domain.GameMember{}, fmt.Errorf("get active game member by user id: %w", err)
	}

	return member, nil
}

func (r *GameMemberRepository) ListByGameID(ctx context.Context, gameID int64) ([]domain.GameMember, error) {
	const query = `
SELECT id, game_id, user_id, added_by_user_id, name, status, source,
    position_number, created_at, cancelled_at
FROM game_members
WHERE game_id = $1
ORDER BY position_number ASC NULLS LAST, created_at ASC;`

	rows, err := r.db.QueryContext(ctx, query, gameID)
	if err != nil {
		return nil, fmt.Errorf("list game members: %w", err)
	}
	defer rows.Close()

	members := make([]domain.GameMember, 0)
	for rows.Next() {
		member, err := scanGameMember(rows)
		if err != nil {
			return nil, fmt.Errorf("scan game member: %w", err)
		}

		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate game members: %w", err)
	}

	return members, nil
}

func (r *GameMemberRepository) CountActiveByGameID(ctx context.Context, gameID int64) (int, error) {
	const query = `
SELECT count(*)
FROM game_members
WHERE game_id = $1 AND status = $2;`

	var count int
	if err := r.db.QueryRowContext(ctx, query, gameID, domain.GameMemberStatusActive).Scan(&count); err != nil {
		return 0, fmt.Errorf("count active game members: %w", err)
	}

	return count, nil
}

func (r *GameMemberRepository) MaxPositionNumberByGameID(ctx context.Context, gameID int64) (int, error) {
	const query = `
SELECT coalesce(max(position_number), 0)
FROM game_members
WHERE game_id = $1;`

	var maxPosition int
	if err := r.db.QueryRowContext(ctx, query, gameID).Scan(&maxPosition); err != nil {
		return 0, fmt.Errorf("get max game member position: %w", err)
	}

	return maxPosition, nil
}

func (r *GameMemberRepository) Cancel(ctx context.Context, id int64) (domain.GameMember, error) {
	const query = `
UPDATE game_members
SET status = $2,
    cancelled_at = now()
WHERE id = $1
RETURNING id, game_id, user_id, added_by_user_id, name, status, source,
    position_number, created_at, cancelled_at;`

	member, err := scanGameMember(r.db.QueryRowContext(ctx, query, id, domain.GameMemberStatusCancelled))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.GameMember{}, domain.ErrNotFound
		}

		return domain.GameMember{}, fmt.Errorf("cancel game member: %w", err)
	}

	return member, nil
}

func (r *GameMemberRepository) GetNextWaitlistByGameID(ctx context.Context, gameID int64) (domain.GameMember, error) {
	const query = `
SELECT id, game_id, user_id, added_by_user_id, name, status, source,
    position_number, created_at, cancelled_at
FROM game_members
WHERE game_id = $1 AND status = $2
ORDER BY position_number ASC NULLS LAST, created_at ASC
LIMIT 1;`

	member, err := scanGameMember(r.db.QueryRowContext(ctx, query, gameID, domain.GameMemberStatusWaitlist))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.GameMember{}, domain.ErrNotFound
		}

		return domain.GameMember{}, fmt.Errorf("get next waitlist game member: %w", err)
	}

	return member, nil
}

func (r *GameMemberRepository) PromoteToActive(ctx context.Context, id int64) (domain.GameMember, error) {
	const query = `
UPDATE game_members
SET status = $2
WHERE id = $1
RETURNING id, game_id, user_id, added_by_user_id, name, status, source,
    position_number, created_at, cancelled_at;`

	member, err := scanGameMember(r.db.QueryRowContext(ctx, query, id, domain.GameMemberStatusActive))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.GameMember{}, domain.ErrNotFound
		}

		return domain.GameMember{}, fmt.Errorf("promote game member to active: %w", err)
	}

	return member, nil
}

type gameMemberScanner interface {
	Scan(dest ...any) error
}

func scanGameMember(scanner gameMemberScanner) (domain.GameMember, error) {
	var member domain.GameMember
	var nullableUserID sql.NullInt64
	var nullableAddedByUserID sql.NullInt64
	var nullablePositionNumber sql.NullInt64
	var nullableCancelledAt sql.NullTime

	if err := scanner.Scan(
		&member.ID,
		&member.GameID,
		&nullableUserID,
		&nullableAddedByUserID,
		&member.Name,
		&member.Status,
		&member.Source,
		&nullablePositionNumber,
		&member.CreatedAt,
		&nullableCancelledAt,
	); err != nil {
		return domain.GameMember{}, err
	}

	member.UserID = nullableInt64Ptr(nullableUserID)
	member.AddedByUserID = nullableInt64Ptr(nullableAddedByUserID)
	member.PositionNumber = nullableIntPtr(nullablePositionNumber)
	member.CancelledAt = nullableTimePtr(nullableCancelledAt)
	return member, nil
}

func nullableIntPtr(value sql.NullInt64) *int {
	if !value.Valid {
		return nil
	}

	number := int(value.Int64)
	return &number
}

func nullableTimePtr(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}

	return &value.Time
}
