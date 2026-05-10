package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/akolesnov/football58/backend/internal/domain"
	"github.com/akolesnov/football58/backend/internal/service"
)

type VenueHandler struct {
	venues *service.VenueService
}

type createVenueRequest struct {
	Name    string  `json:"name"`
	Address *string `json:"address"`
}

type venueResponse struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Address *string `json:"address,omitempty"`
}

func NewVenueHandler(venues *service.VenueService) *VenueHandler {
	return &VenueHandler{venues: venues}
}

func (h *VenueHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request createVenueRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid_json", "некорректный JSON")
		return
	}

	venue, err := h.venues.Create(r.Context(), request.Name, request.Address)
	if err != nil {
		if errors.Is(err, service.ErrVenueNameRequired) {
			WriteError(w, http.StatusBadRequest, "venue_name_required", "название площадки обязательно")
			return
		}

		WriteError(w, http.StatusInternalServerError, "create_venue_failed", "не удалось создать площадку")
		return
	}

	WriteJSON(w, http.StatusCreated, venueToResponse(venue))
}

func (h *VenueHandler) List(w http.ResponseWriter, r *http.Request) {
	venues, err := h.venues.List(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "list_venues_failed", "не удалось получить площадки")
		return
	}

	response := make([]venueResponse, 0, len(venues))
	for _, venue := range venues {
		response = append(response, venueToResponse(venue))
	}

	WriteJSON(w, http.StatusOK, response)
}

func (h *VenueHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		WriteError(w, http.StatusBadRequest, "invalid_venue_id", "некорректный id площадки")
		return
	}

	venue, err := h.venues.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "venue_not_found", "площадка не найдена")
			return
		}

		WriteError(w, http.StatusInternalServerError, "get_venue_failed", "не удалось получить площадку")
		return
	}

	WriteJSON(w, http.StatusOK, venueToResponse(venue))
}

func venueToResponse(venue domain.Venue) venueResponse {
	return venueResponse{
		ID:      venue.ID,
		Name:    venue.Name,
		Address: venue.Address,
	}
}
