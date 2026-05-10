package service

import (
	"context"
	"errors"
	"time"

	"github.com/akolesnov/football58/backend/internal/domain"
	"github.com/akolesnov/football58/backend/internal/repository"
)

const (
	defaultGameDurationMinutes = 120
	defaultGameMinPlayers      = 10
	defaultGameMaxPlayers      = 15
	defaultGamePriceRub        = 220
)

var ErrGameVenueRequired = errors.New("game venue is required")
var ErrGameStartsAtRequired = errors.New("game starts_at is required")
var ErrGameDurationInvalid = errors.New("game duration is invalid")
var ErrGamePlayersLimitInvalid = errors.New("game players limit is invalid")
var ErrGamePriceInvalid = errors.New("game price is invalid")

type GameService struct {
	games *repository.GameRepository
}

func NewGameService(games *repository.GameRepository) *GameService {
	return &GameService{games: games}
}

func (s *GameService) Create(ctx context.Context, game domain.Game) (domain.Game, error) {
	if game.VenueID <= 0 {
		return domain.Game{}, ErrGameVenueRequired
	}
	if game.StartsAt.IsZero() {
		return domain.Game{}, ErrGameStartsAtRequired
	}

	applyGameDefaults(&game)

	if game.DurationMinutes <= 0 {
		return domain.Game{}, ErrGameDurationInvalid
	}
	if game.MinPlayers <= 0 || game.MaxPlayers <= 0 || game.MinPlayers > game.MaxPlayers {
		return domain.Game{}, ErrGamePlayersLimitInvalid
	}
	if game.PriceRub < 0 {
		return domain.Game{}, ErrGamePriceInvalid
	}

	return s.games.Create(ctx, game)
}

func applyGameDefaults(game *domain.Game) {
	if game.DurationMinutes == 0 {
		game.DurationMinutes = defaultGameDurationMinutes
	}
	if game.MinPlayers == 0 {
		game.MinPlayers = defaultGameMinPlayers
	}
	if game.MaxPlayers == 0 {
		game.MaxPlayers = defaultGameMaxPlayers
	}
	if game.PriceRub == 0 {
		game.PriceRub = defaultGamePriceRub
	}
}

func ParseGameStartsAt(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}
