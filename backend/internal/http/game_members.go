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

type GameMemberHandler struct {
	members *service.GameMemberService
}

type gameMemberResponse struct {
	ID             int64      `json:"id"`
	GameID         int64      `json:"game_id"`
	UserID         *int64     `json:"user_id,omitempty"`
	AddedByUserID  *int64     `json:"added_by_user_id,omitempty"`
	Name           string     `json:"name"`
	Status         string     `json:"status"`
	Source         string     `json:"source"`
	PositionNumber *int       `json:"position_number,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	CancelledAt    *time.Time `json:"cancelled_at,omitempty"`
}

type cancelGameMemberResponse struct {
	Game     gameResponse        `json:"game"`
	Member   gameMemberResponse  `json:"member"`
	Promoted *gameMemberResponse `json:"promoted,omitempty"`
}

type cancelOwnTelegramMembershipRequest struct {
	TelegramID int64 `json:"telegram_id"`
}

func NewGameMemberHandler(members *service.GameMemberService) *GameMemberHandler {
	return &GameMemberHandler{members: members}
}

func (h *GameMemberHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	gameID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || gameID <= 0 {
		WriteError(w, http.StatusBadRequest, "invalid_game_id", "некорректный id игры")
		return
	}

	memberID, err := strconv.ParseInt(r.PathValue("member_id"), 10, 64)
	if err != nil || memberID <= 0 {
		WriteError(w, http.StatusBadRequest, "invalid_member_id", "некорректный id участника")
		return
	}

	result, err := h.members.CancelMember(r.Context(), gameID, memberID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "not_found", "игра или участник не найдены")
			return
		}
		if errors.Is(err, service.ErrGameMemberGameMismatch) {
			WriteError(w, http.StatusBadRequest, "game_member_mismatch", "участник не относится к этой игре")
			return
		}

		WriteError(w, http.StatusInternalServerError, "cancel_game_member_failed", "не удалось отменить участие")
		return
	}

	writeCancelGameMemberResponse(w, result)
}

func (h *GameMemberHandler) CancelOwnTelegramMembership(w http.ResponseWriter, r *http.Request) {
	gameID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || gameID <= 0 {
		WriteError(w, http.StatusBadRequest, "invalid_game_id", "некорректный id игры")
		return
	}

	var request cancelOwnTelegramMembershipRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid_json", "некорректный JSON")
		return
	}

	result, err := h.members.CancelOwnTelegramMembership(r.Context(), gameID, request.TelegramID)
	if err != nil {
		if errors.Is(err, service.ErrTelegramIDRequired) {
			WriteError(w, http.StatusBadRequest, "telegram_id_required", "telegram_id обязателен")
			return
		}
		if errors.Is(err, domain.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "membership_not_found", "запись пользователя на игру не найдена")
			return
		}

		WriteError(w, http.StatusInternalServerError, "cancel_own_membership_failed", "не удалось отменить участие")
		return
	}

	writeCancelGameMemberResponse(w, result)
}

func writeCancelGameMemberResponse(w http.ResponseWriter, result service.CancelMemberResult) {
	response := cancelGameMemberResponse{
		Game:   gameToResponse(result.Game),
		Member: gameMemberToResponse(result.Member),
	}
	if result.Promoted != nil {
		promoted := gameMemberToResponse(*result.Promoted)
		response.Promoted = &promoted
	}

	WriteJSON(w, http.StatusOK, response)
}

func gameMemberToResponse(member domain.GameMember) gameMemberResponse {
	return gameMemberResponse{
		ID:             member.ID,
		GameID:         member.GameID,
		UserID:         member.UserID,
		AddedByUserID:  member.AddedByUserID,
		Name:           member.Name,
		Status:         member.Status,
		Source:         member.Source,
		PositionNumber: member.PositionNumber,
		CreatedAt:      member.CreatedAt,
		CancelledAt:    member.CancelledAt,
	}
}
