package domain

import "time"

const (
	GameStatusOpen      = "open"
	GameStatusClosed    = "closed"
	GameStatusCancelled = "cancelled"
	GameStatusFinished  = "finished"
)

type Game struct {
	ID                int64
	VenueID           int64
	StartsAt          time.Time
	DurationMinutes   int
	MinPlayers        int
	MaxPlayers        int
	PriceRub          int
	Notes             *string
	Status            string
	CreatedByUserID   *int64
	TelegramChatID    *int64
	TelegramMessageID *int64
	Version           int
	CreatedAt         time.Time
}
