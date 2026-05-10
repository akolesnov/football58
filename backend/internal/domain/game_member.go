package domain

import "time"

const (
	GameMemberStatusActive    = "active"
	GameMemberStatusWaitlist  = "waitlist"
	GameMemberStatusCancelled = "cancelled"

	GameMemberSourceTelegram = "telegram"
	GameMemberSourceWeb      = "web"
	GameMemberSourceAdmin    = "admin"
)

type GameMember struct {
	ID             int64
	GameID         int64
	UserID         *int64
	AddedByUserID  *int64
	Name           string
	Status         string
	Source         string
	PositionNumber *int
	CreatedAt      time.Time
	CancelledAt    *time.Time
}
