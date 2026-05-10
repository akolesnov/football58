package domain

import "time"

type User struct {
	ID               int64
	Name             string
	TelegramID       *int64
	TelegramUsername *string
	CreatedAt        time.Time
}
