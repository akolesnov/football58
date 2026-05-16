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
var ErrTelegramChatIDRequired = errors.New("telegram chat id is required")
var ErrTelegramMessageIDRequired = errors.New("telegram message id is required")

type GameService struct {
	games   *repository.GameRepository
	members *repository.GameMemberRepository
}

type GameDetails struct {
	Game          domain.Game
	Members       []domain.GameMember
	ActiveCount   int
	WaitlistCount int
}

func NewGameService(games *repository.GameRepository, members *repository.GameMemberRepository) *GameService {
	return &GameService{
		games:   games,
		members: members,
	}
}

func (s *GameService) GetByID(ctx context.Context, id int64) (domain.Game, error) {
	return s.games.GetByID(ctx, id)
}

func (s *GameService) GetDetailsByID(ctx context.Context, id int64) (GameDetails, error) {
	game, err := s.games.GetByID(ctx, id)
	if err != nil {
		return GameDetails{}, err
	}

	members, err := s.members.ListByGameID(ctx, game.ID)
	if err != nil {
		return GameDetails{}, err
	}

	return buildGameDetails(game, members), nil
}

func (s *GameService) List(ctx context.Context) ([]domain.Game, error) {
	return s.games.List(ctx)
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

func (s *GameService) UpdateTelegramMessage(ctx context.Context, id, chatID, messageID int64) (domain.Game, error) {
	if id <= 0 {
		return domain.Game{}, ErrGameRequired
	}
	if chatID == 0 {
		return domain.Game{}, ErrTelegramChatIDRequired
	}
	if messageID <= 0 {
		return domain.Game{}, ErrTelegramMessageIDRequired
	}

	return s.games.UpdateTelegramMessage(ctx, id, chatID, messageID)
}

func (s *GameService) Close(ctx context.Context, id int64) (domain.Game, error) {
	return s.setStatus(ctx, id, domain.GameStatusClosed)
}

func (s *GameService) Cancel(ctx context.Context, id int64) (domain.Game, error) {
	return s.setStatus(ctx, id, domain.GameStatusCancelled)
}

func (s *GameService) Finish(ctx context.Context, id int64) (domain.Game, error) {
	return s.setStatus(ctx, id, domain.GameStatusFinished)
}

func (s *GameService) setStatus(ctx context.Context, id int64, status string) (domain.Game, error) {
	if id <= 0 {
		return domain.Game{}, ErrGameRequired
	}

	return s.games.SetStatus(ctx, id, status)
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

func buildGameDetails(game domain.Game, members []domain.GameMember) GameDetails {
	var activeCount int
	var waitlistCount int

	for _, member := range members {
		switch member.Status {
		case domain.GameMemberStatusActive:
			activeCount++
		case domain.GameMemberStatusWaitlist:
			waitlistCount++
		}
	}

	return GameDetails{
		Game:          game,
		Members:       members,
		ActiveCount:   activeCount,
		WaitlistCount: waitlistCount,
	}
}
