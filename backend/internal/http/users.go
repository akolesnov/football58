package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/akolesnov/football58/backend/internal/domain"
	"github.com/akolesnov/football58/backend/internal/service"
)

type UserHandler struct {
	users *service.UserService
}

type upsertTelegramUserRequest struct {
	Name             string  `json:"name"`
	TelegramID       int64   `json:"telegram_id"`
	TelegramUsername *string `json:"telegram_username"`
}

type userResponse struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	TelegramID       *int64    `json:"telegram_id,omitempty"`
	TelegramUsername *string   `json:"telegram_username,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

func NewUserHandler(users *service.UserService) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) UpsertTelegram(w http.ResponseWriter, r *http.Request) {
	var request upsertTelegramUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid_json", "некорректный JSON")
		return
	}

	user, err := h.users.UpsertTelegram(r.Context(), request.Name, request.TelegramID, request.TelegramUsername)
	if err != nil {
		if errors.Is(err, service.ErrUserNameRequired) {
			WriteError(w, http.StatusBadRequest, "user_name_required", "имя пользователя обязательно")
			return
		}
		if errors.Is(err, service.ErrTelegramIDRequired) {
			WriteError(w, http.StatusBadRequest, "telegram_id_required", "telegram_id обязателен")
			return
		}

		WriteError(w, http.StatusInternalServerError, "upsert_telegram_user_failed", "не удалось сохранить пользователя")
		return
	}

	WriteJSON(w, http.StatusOK, userToResponse(user))
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		WriteError(w, http.StatusBadRequest, "invalid_user_id", "некорректный id пользователя")
		return
	}

	user, err := h.users.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "user_not_found", "пользователь не найден")
			return
		}

		WriteError(w, http.StatusInternalServerError, "get_user_failed", "не удалось получить пользователя")
		return
	}

	WriteJSON(w, http.StatusOK, userToResponse(user))
}

func (h *UserHandler) GetByTelegramID(w http.ResponseWriter, r *http.Request) {
	telegramID, err := strconv.ParseInt(r.PathValue("telegram_id"), 10, 64)
	if err != nil || telegramID == 0 {
		WriteError(w, http.StatusBadRequest, "invalid_telegram_id", "некорректный telegram_id")
		return
	}

	user, err := h.users.GetByTelegramID(r.Context(), telegramID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "user_not_found", "пользователь не найден")
			return
		}

		WriteError(w, http.StatusInternalServerError, "get_user_failed", "не удалось получить пользователя")
		return
	}

	WriteJSON(w, http.StatusOK, userToResponse(user))
}

func userToResponse(user domain.User) userResponse {
	return userResponse{
		ID:               user.ID,
		Name:             user.Name,
		TelegramID:       user.TelegramID,
		TelegramUsername: user.TelegramUsername,
		CreatedAt:        user.CreatedAt,
	}
}
