package schemas

import (
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	ID        string `gorm:"primaryKey"`
	Token     string // hashed token
	DeviceID  string
	ExpiresAt time.Time
	CreatedAt time.Time
	UserID    uuid.UUID
}
