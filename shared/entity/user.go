package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID      uuid.UUID  `gorm:"type:uuid;primaryKey" json:"user_id,omitempty"`
	Permission  string     `gorm:"index:idx_username_permission" json:"permission,omitempty"`
	Username    *string    `gorm:"unique;index:idx_username_permission" json:"username,omitempty"`
	Password    *string    `gorm:"size:255" json:"password,omitempty"`
	Email       *string    `gorm:"unique" json:"email,omitempty"`
	Phone       *string    `gorm:"size:15" json:"phone,omitempty"`
	FullName    *string    `json:"fullname,omitempty"`
	Gender      *string    `json:"gender,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	FacebookURL *string    `gorm:"unique" json:"facebook_url,omitempty"`
	GithubURL   *string    `gorm:"unique" json:"github_url,omitempty"`
	AvatarURL   *string    `json:"avatar_url,omitempty"`

	Verifications []UserVerification `gorm:"foreignKey:UserID"`
	MFAs          []UserMFA          `gorm:"foreignKey:UserID"`
}
