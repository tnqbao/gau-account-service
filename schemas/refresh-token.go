package schemas

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        int64     `gorm:"primaryKey;autoIncrement:false"`
	Token     string    `gorm:"uniqueIndex"`
	DeviceID  string    `gorm:"uniqueIndex"`
	ExpiresAt time.Time `gorm:"index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UserID    uuid.UUID `gorm:"index;not null"`
}
