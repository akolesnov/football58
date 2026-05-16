package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/akolesnov/football58/backend/internal/domain"
	"github.com/akolesnov/football58/backend/internal/service"
)

type GameHandler struct {
	games *service.GameService
}

type createGameRequest struct {
	VenueID         int64   `json:"venue_id"`
	StartsAt        string  `json:"starts_at"`
	DurationMinutes int     `json:"duration_minutes"`
	MinPlayers      int     `json:"min_players"`
	MaxPlayers      int     `json:"max_players"`
	PriceRub        int     `json:"price_rub"`
	Notes           *string `json:"notes"`
	CreatedByUserID *int64  `json:"created_by_user_id"`
}

type updateTelegramMessageRequest struct {
	TelegramChatID    int64 `json:"telegram_chat_id"`
	TelegramMessageID int64 `json:"telegram_message_id"`
}

type gameResponse struct {
	ID                int64     `json:"id"`
	VenueID           int64     `json:"venue_id"`
	StartsAt          time.Time `json:"starts_at"`
	DurationMinutes   int       `json:"duration_minutes"`
	MinPlayers        int       `json:"min_players"`
	MaxPlayers        int       `json:"max_players"`
	PriceRub          int       `json:"price_rub"`
	Notes             *string   `json:"notes,omitempty"`
	Status            string    `json:"status"`
	CreatedByUserID   *int64    `json:"created_by_user_id,omitempty"`
	TelegramChatID    *int64    `json:"telegram_chat_id,omitempty"`
	TelegramMessageID *int64    `json:"telegram_message_id,omitempty"`
	Version           int       `json:"version"`
	CreatedAt         time.Time `json:"created_at"`
}

type gameDetailsResponse struct {
	gameResponse
	ActiveCount   int                  `json:"active_count"`
	WaitlistCount int                  `json:"waitlist_count"`
	Members       []gameMemberResponse `json:"members"`
}

func NewGameHandler(games *service.GameService) *GameHandler {
	return &GameHandler{games: games}
}

func (h *GameHandler) List(w http.ResponseWriter, r *http.Request) {
	games, err := h.games.List(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "list_games_failed", "не удалось получить игры")
		return
	}

	response := make([]gameResponse, 0, len(games))
	for _, game := range games {
		response = append(response, gameToResponse(game))
	}

	WriteJSON(w, http.StatusOK, response)
}

func (h *GameHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		WriteError(w, http.StatusBadRequest, "invalid_game_id", "некорректный id игры")
		return
	}

	details, err := h.games.GetDetailsByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "game_not_found", "игра не найдена")
			return
		}

		WriteError(w, http.StatusInternalServerError, "get_game_failed", "не удалось получить игру")
		return
	}

	WriteJSON(w, http.StatusOK, gameDetailsToResponse(details))
}

func (h *GameHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request createGameRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid_json", "некорректный JSON")
		return
	}

	startsAt, err := service.ParseGameStartsAt(request.StartsAt)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid_starts_at", "starts_at должен быть в формате RFC3339")
		return
	}

	game, err := h.games.Create(r.Context(), domain.Game{
		VenueID:         request.VenueID,
		StartsAt:        startsAt,
		DurationMinutes: request.DurationMinutes,
		MinPlayers:      request.MinPlayers,
		MaxPlayers:      request.MaxPlayers,
		PriceRub:        request.PriceRub,
		Notes:           request.Notes,
		CreatedByUserID: request.CreatedByUserID,
	})
	if err != nil {
		if errors.Is(err, service.ErrGameVenueRequired) {
			WriteError(w, http.StatusBadRequest, "game_venue_required", "площадка обязательна")
			return
		}
		if errors.Is(err, service.ErrGameStartsAtRequired) {
			WriteError(w, http.StatusBadRequest, "game_starts_at_required", "время начала обязательно")
			return
		}
		if errors.Is(err, service.ErrGameDurationInvalid) {
			WriteError(w, http.StatusBadRequest, "game_duration_invalid", "некорректная длительность игры")
			return
		}
		if errors.Is(err, service.ErrGamePlayersLimitInvalid) {
			WriteError(w, http.StatusBadRequest, "game_players_limit_invalid", "некорректный лимит игроков")
			return
		}
		if errors.Is(err, service.ErrGamePriceInvalid) {
			WriteError(w, http.StatusBadRequest, "game_price_invalid", "некорректная стоимость игры")
			return
		}

		WriteError(w, http.StatusInternalServerError, "create_game_failed", "не удалось создать игру")
		return
	}

	WriteJSON(w, http.StatusCreated, gameToResponse(game))
}

func (h *GameHandler) UpdateTelegramMessage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		WriteError(w, http.StatusBadRequest, "invalid_game_id", "некорректный id игры")
		return
	}

	var request updateTelegramMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid_json", "некорректный JSON")
		return
	}

	game, err := h.games.UpdateTelegramMessage(r.Context(), id, request.TelegramChatID, request.TelegramMessageID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "game_not_found", "игра не найдена")
			return
		}
		if errors.Is(err, service.ErrTelegramChatIDRequired) {
			WriteError(w, http.StatusBadRequest, "telegram_chat_id_required", "telegram_chat_id обязателен")
			return
		}
		if errors.Is(err, service.ErrTelegramMessageIDRequired) {
			WriteError(w, http.StatusBadRequest, "telegram_message_id_required", "telegram_message_id обязателен")
			return
		}

		WriteError(w, http.StatusInternalServerError, "update_telegram_message_failed", "не удалось сохранить Telegram-сообщение")
		return
	}

	WriteJSON(w, http.StatusOK, gameToResponse(game))
}

func (h *GameHandler) Close(w http.ResponseWriter, r *http.Request) {
	h.setStatus(w, r, h.games.Close, "close_game_failed")
}

func (h *GameHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	h.setStatus(w, r, h.games.Cancel, "cancel_game_failed")
}

func (h *GameHandler) Finish(w http.ResponseWriter, r *http.Request) {
	h.setStatus(w, r, h.games.Finish, "finish_game_failed")
}

func (h *GameHandler) setStatus(
	w http.ResponseWriter,
	r *http.Request,
	update func(context.Context, int64) (domain.Game, error),
	internalCode string,
) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		WriteError(w, http.StatusBadRequest, "invalid_game_id", "некорректный id игры")
		return
	}

	game, err := update(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "game_not_found", "игра не найдена")
			return
		}
		if errors.Is(err, service.ErrGameRequired) {
			WriteError(w, http.StatusBadRequest, "invalid_game_id", "некорректный id игры")
			return
		}

		WriteError(w, http.StatusInternalServerError, internalCode, "не удалось изменить статус игры")
		return
	}

	WriteJSON(w, http.StatusOK, gameToResponse(game))
}

func gameToResponse(game domain.Game) gameResponse {
	return gameResponse{
		ID:                game.ID,
		VenueID:           game.VenueID,
		StartsAt:          game.StartsAt,
		DurationMinutes:   game.DurationMinutes,
		MinPlayers:        game.MinPlayers,
		MaxPlayers:        game.MaxPlayers,
		PriceRub:          game.PriceRub,
		Notes:             game.Notes,
		Status:            game.Status,
		CreatedByUserID:   game.CreatedByUserID,
		TelegramChatID:    game.TelegramChatID,
		TelegramMessageID: game.TelegramMessageID,
		Version:           game.Version,
		CreatedAt:         game.CreatedAt,
	}
}

func gameDetailsToResponse(details service.GameDetails) gameDetailsResponse {
	members := make([]gameMemberResponse, 0, len(details.Members))
	for _, member := range details.Members {
		members = append(members, gameMemberToResponse(member))
	}

	return gameDetailsResponse{
		gameResponse:  gameToResponse(details.Game),
		ActiveCount:   details.ActiveCount,
		WaitlistCount: details.WaitlistCount,
		Members:       members,
	}
}
