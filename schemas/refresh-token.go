package schemas

import (
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	ID        string    `gorm:"primaryKey"`
	Token     string    `gorm:"uniqueIndex"`
	DeviceID  string    `gorm:"uniqueIndex"`
	ExpiresAt time.Time `gorm:"index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UserID    uuid.UUID `gorm:"index;not null"`
}
