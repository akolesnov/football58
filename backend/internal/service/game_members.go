package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/akolesnov/football58/backend/internal/domain"
	"github.com/akolesnov/football58/backend/internal/repository"
)

var ErrGameNotOpen = errors.New("game is not open")
var ErrGameRequired = errors.New("game is required")
var ErrGameMemberRequired = errors.New("game member is required")
var ErrGameMemberGameMismatch = errors.New("game member belongs to another game")
var ErrGameMemberNameRequired = errors.New("game member name is required")
var ErrGameMemberSourceInvalid = errors.New("game member source is invalid")

type GameMemberService struct {
	db *sql.DB
}

type JoinGameInput struct {
	GameID        int64
	UserID        *int64
	AddedByUserID *int64
	Name          string
	Source        string
}

type JoinGameResult struct {
	Game   domain.Game
	Member domain.GameMember
}

type CancelMemberResult struct {
	Game     domain.Game
	Member   domain.GameMember
	Promoted *domain.GameMember
}

func NewGameMemberService(db *sql.DB) *GameMemberService {
	return &GameMemberService{db: db}
}

func (s *GameMemberService) JoinGame(ctx context.Context, input JoinGameInput) (JoinGameResult, error) {
	if input.GameID <= 0 {
		return JoinGameResult{}, ErrGameRequired
	}

	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		return JoinGameResult{}, ErrGameMemberNameRequired
	}

	input.Source = strings.TrimSpace(input.Source)
	if !validGameMemberSource(input.Source) {
		return JoinGameResult{}, ErrGameMemberSourceInvalid
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return JoinGameResult{}, err
	}
	defer tx.Rollback()

	games := repository.NewGameRepository(tx)
	members := repository.NewGameMemberRepository(tx)

	game, err := games.GetByIDForUpdate(ctx, input.GameID)
	if err != nil {
		return JoinGameResult{}, err
	}
	if game.Status != domain.GameStatusOpen {
		return JoinGameResult{}, ErrGameNotOpen
	}

	activeCount, err := members.CountActiveByGameID(ctx, game.ID)
	if err != nil {
		return JoinGameResult{}, err
	}

	status := domain.GameMemberStatusWaitlist
	if activeCount < game.MaxPlayers {
		status = domain.GameMemberStatusActive
	}

	maxPosition, err := members.MaxPositionNumberByGameID(ctx, game.ID)
	if err != nil {
		return JoinGameResult{}, err
	}
	position := maxPosition + 1

	member, err := members.Create(ctx, domain.GameMember{
		GameID:         game.ID,
		UserID:         input.UserID,
		AddedByUserID:  input.AddedByUserID,
		Name:           input.Name,
		Status:         status,
		Source:         input.Source,
		PositionNumber: &position,
	})
	if err != nil {
		return JoinGameResult{}, err
	}

	game, err = games.IncrementVersion(ctx, game.ID)
	if err != nil {
		return JoinGameResult{}, err
	}

	if err := tx.Commit(); err != nil {
		return JoinGameResult{}, err
	}

	return JoinGameResult{
		Game:   game,
		Member: member,
	}, nil
}

func (s *GameMemberService) CancelMember(ctx context.Context, gameID, memberID int64) (CancelMemberResult, error) {
	if gameID <= 0 {
		return CancelMemberResult{}, ErrGameRequired
	}
	if memberID <= 0 {
		return CancelMemberResult{}, ErrGameMemberRequired
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return CancelMemberResult{}, err
	}
	defer tx.Rollback()

	games := repository.NewGameRepository(tx)
	members := repository.NewGameMemberRepository(tx)

	game, err := games.GetByIDForUpdate(ctx, gameID)
	if err != nil {
		return CancelMemberResult{}, err
	}

	member, err := members.GetByID(ctx, memberID)
	if err != nil {
		return CancelMemberResult{}, err
	}
	if member.GameID != game.ID {
		return CancelMemberResult{}, ErrGameMemberGameMismatch
	}

	result, err := cancelMemberInTx(ctx, games, members, game, member)
	if err != nil {
		return CancelMemberResult{}, err
	}

	if err := tx.Commit(); err != nil {
		return CancelMemberResult{}, err
	}

	return result, nil
}

func (s *GameMemberService) CancelOwnTelegramMembership(ctx context.Context, gameID, telegramID int64) (CancelMemberResult, error) {
	if gameID <= 0 {
		return CancelMemberResult{}, ErrGameRequired
	}
	if telegramID == 0 {
		return CancelMemberResult{}, ErrTelegramIDRequired
	}

	users := repository.NewUserRepository(s.db)
	user, err := users.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return CancelMemberResult{}, err
	}

	if user.TelegramID == nil {
		return CancelMemberResult{}, domain.ErrNotFound
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return CancelMemberResult{}, err
	}
	defer tx.Rollback()

	games := repository.NewGameRepository(tx)
	members := repository.NewGameMemberRepository(tx)

	game, err := games.GetByIDForUpdate(ctx, gameID)
	if err != nil {
		return CancelMemberResult{}, err
	}

	member, err := members.GetActiveByGameIDAndUserID(ctx, game.ID, user.ID)
	if err != nil {
		return CancelMemberResult{}, err
	}

	result, err := cancelMemberInTx(ctx, games, members, game, member)
	if err != nil {
		return CancelMemberResult{}, err
	}

	if err := tx.Commit(); err != nil {
		return CancelMemberResult{}, err
	}

	return result, nil
}

func cancelMemberInTx(
	ctx context.Context,
	games *repository.GameRepository,
	members *repository.GameMemberRepository,
	game domain.Game,
	member domain.GameMember,
) (CancelMemberResult, error) {
	if member.GameID != game.ID {
		return CancelMemberResult{}, ErrGameMemberGameMismatch
	}

	if member.Status == domain.GameMemberStatusCancelled {
		return CancelMemberResult{
			Game:   game,
			Member: member,
		}, nil
	}

	cancelledMember, err := members.Cancel(ctx, member.ID)
	if err != nil {
		return CancelMemberResult{}, err
	}

	var promoted *domain.GameMember
	if member.Status == domain.GameMemberStatusActive {
		nextWaitlistMember, err := members.GetNextWaitlistByGameID(ctx, game.ID)
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			return CancelMemberResult{}, err
		}

		if err == nil {
			promotedMember, err := members.PromoteToActive(ctx, nextWaitlistMember.ID)
			if err != nil {
				return CancelMemberResult{}, err
			}

			promoted = &promotedMember
		}
	}

	game, err = games.IncrementVersion(ctx, game.ID)
	if err != nil {
		return CancelMemberResult{}, err
	}

	return CancelMemberResult{
		Game:     game,
		Member:   cancelledMember,
		Promoted: promoted,
	}, nil
}

func validGameMemberSource(source string) bool {
	switch source {
	case domain.GameMemberSourceTelegram, domain.GameMemberSourceWeb, domain.GameMemberSourceAdmin:
		return true
	default:
		return false
	}
}
