package httpapi

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

const dbHealthTimeout = 2 * time.Second

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Health(w http.ResponseWriter, _ *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (h *HealthHandler) Database(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), dbHealthTimeout)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		WriteJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status":   "error",
			"database": "unavailable",
		})
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"status":   "ok",
		"database": "ok",
	})
}
