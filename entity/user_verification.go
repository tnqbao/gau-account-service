package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserVerification struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id,omitempty"`
	UserID     uuid.UUID  `gorm:"type:uuid;index" json:"user_id,omitempty"`
	Method     string     `gorm:"size:20;index" json:"method,omitempty"` // "email" | "phone"
	Value      string     `gorm:"size:255" json:"value,omitempty"`       // email hoặc số điện thoại
	IsVerified bool       `gorm:"default:false" json:"is_verified,omitempty"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
