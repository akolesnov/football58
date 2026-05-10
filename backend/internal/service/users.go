package service

import (
	"context"
	"errors"
	"strings"

	"github.com/akolesnov/football58/backend/internal/domain"
	"github.com/akolesnov/football58/backend/internal/repository"
)

var ErrUserNameRequired = errors.New("user name is required")
var ErrTelegramIDRequired = errors.New("telegram id is required")

type UserService struct {
	users *repository.UserRepository
}

func NewUserService(users *repository.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) GetByID(ctx context.Context, id int64) (domain.User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *UserService) GetByTelegramID(ctx context.Context, telegramID int64) (domain.User, error) {
	return s.users.GetByTelegramID(ctx, telegramID)
}

func (s *UserService) UpsertTelegram(ctx context.Context, name string, telegramID int64, telegramUsername *string) (domain.User, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.User{}, ErrUserNameRequired
	}
	if telegramID == 0 {
		return domain.User{}, ErrTelegramIDRequired
	}

	if telegramUsername != nil {
		username := strings.TrimSpace(*telegramUsername)
		if username == "" {
			telegramUsername = nil
		} else {
			telegramUsername = &username
		}
	}

	return s.users.UpsertTelegram(ctx, name, telegramID, telegramUsername)
}
