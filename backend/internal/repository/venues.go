package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/akolesnov/football58/backend/internal/domain"
)

type VenueRepository struct {
	db *sql.DB
}

func NewVenueRepository(db *sql.DB) *VenueRepository {
	return &VenueRepository{db: db}
}

func (r *VenueRepository) Create(ctx context.Context, name string, address *string) (domain.Venue, error) {
	const query = `
INSERT INTO venues (name, address)
VALUES ($1, $2)
RETURNING id, name, address;`

	var venue domain.Venue
	var nullableAddress sql.NullString

	if err := r.db.QueryRowContext(ctx, query, name, address).Scan(
		&venue.ID,
		&venue.Name,
		&nullableAddress,
	); err != nil {
		return domain.Venue{}, fmt.Errorf("create venue: %w", err)
	}

	venue.Address = nullableStringPtr(nullableAddress)
	return venue, nil
}

func (r *VenueRepository) List(ctx context.Context) ([]domain.Venue, error) {
	const query = `
SELECT id, name, address
FROM venues
ORDER BY name;`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list venues: %w", err)
	}
	defer rows.Close()

	venues := make([]domain.Venue, 0)
	for rows.Next() {
		var venue domain.Venue
		var nullableAddress sql.NullString

		if err := rows.Scan(&venue.ID, &venue.Name, &nullableAddress); err != nil {
			return nil, fmt.Errorf("scan venue: %w", err)
		}

		venue.Address = nullableStringPtr(nullableAddress)
		venues = append(venues, venue)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate venues: %w", err)
	}

	return venues, nil
}

func (r *VenueRepository) GetByID(ctx context.Context, id int64) (domain.Venue, error) {
	const query = `
SELECT id, name, address
FROM venues
WHERE id = $1;`

	var venue domain.Venue
	var nullableAddress sql.NullString

	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&venue.ID,
		&venue.Name,
		&nullableAddress,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Venue{}, domain.ErrNotFound
		}

		return domain.Venue{}, fmt.Errorf("get venue by id: %w", err)
	}

	venue.Address = nullableStringPtr(nullableAddress)
	return venue, nil
}

func nullableStringPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}

	return &value.String
}
