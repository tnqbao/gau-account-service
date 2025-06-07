package schemas

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	UserID          uuid.UUID  `gorm:"primaryKey,index" json:"user_id,omitempty"`
	Permission      string     `gorm:"index:idx_username_permission" json:"permission,omitempty"`
	Username        *string    `gorm:"unique;index:idx_username_permission" json:"username,omitempty"`
	Password        *string    `gorm:"size:255" json:"password,omitempty"`
	Email           *string    `gorm:"unique;index:idx_email_permission" json:"email,omitempty"`
	Phone           *string    `gorm:"size:15" json:"phone,omitempty"`
	FullName        *string    `json:"fullname,omitempty"`
	Gender          *string    `json:"gender,omitempty"`
	DateOfBirth     *time.Time `json:"date_of_birth,omitempty"`
	IsPhoneVerified bool       `gorm:"default:false" json:"is_phone_verified,omitempty"`
	IsEmailVerified bool       `gorm:"default:false" json:"is_email_verified,omitempty"`
	FacebookURL     *string    `gorm:"unique" json:"facebook_url,omitempty"`
	GithubURL       *string    `gorm:"unique" json:"github_url,omitempty"`
}
