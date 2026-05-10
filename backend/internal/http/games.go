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

	game, err := h.games.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "game_not_found", "игра не найдена")
			return
		}

		WriteError(w, http.StatusInternalServerError, "get_game_failed", "не удалось получить игру")
		return
	}

	WriteJSON(w, http.StatusOK, gameToResponse(game))
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
