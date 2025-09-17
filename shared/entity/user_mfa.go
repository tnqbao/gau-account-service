package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserMFA struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id,omitempty"`
	UserID     uuid.UUID  `gorm:"type:uuid;index" json:"user_id,omitempty"`
	Type       string     `gorm:"size:30;index" json:"type,omitempty"` // "totp" | "email_otp" | "sms_otp"
	Secret     *string    `gorm:"size:255" json:"secret,omitempty"`    // secret TOTP (nếu có)
	Enabled    bool       `gorm:"default:false" json:"enabled,omitempty"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at,omitempty"`
	UpdatedAt  time.Time  `json:"updated_at,omitempty"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
