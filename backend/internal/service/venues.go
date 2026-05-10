package service

import (
	"context"
	"errors"
	"strings"

	"github.com/akolesnov/football58/backend/internal/domain"
	"github.com/akolesnov/football58/backend/internal/repository"
)

var ErrVenueNameRequired = errors.New("venue name is required")

type VenueService struct {
	venues *repository.VenueRepository
}

func NewVenueService(venues *repository.VenueRepository) *VenueService {
	return &VenueService{venues: venues}
}

func (s *VenueService) Create(ctx context.Context, name string, address *string) (domain.Venue, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.Venue{}, ErrVenueNameRequired
	}

	return s.venues.Create(ctx, name, address)
}

func (s *VenueService) List(ctx context.Context) ([]domain.Venue, error) {
	return s.venues.List(ctx)
}

func (s *VenueService) GetByID(ctx context.Context, id int64) (domain.Venue, error) {
	return s.venues.GetByID(ctx, id)
}
