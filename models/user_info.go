package models

import (
	"github.com/google/uuid"
	"time"
)

type UserInformation struct {
	UserId          uuid.UUID  `gorm:"primaryKey;index" json:"user_id"`
	FullName        *string    `json:"fullname"`
	Phone           *string    `gorm:"unique" json:"phone"`
	IsPhoneVerified bool       `gorm:"default:false" json:"is_phone_verified"`
	Email           *string    `gorm:"unique" json:"email"`
	IsEmailVerified bool       `gorm:"default:false" json:"is_email_verified"`
	DateOfBirth     *time.Time `json:"date_of_birth"`
	FacebookURL     *string    `gorm:"unique" json:"facebook_url"`
	GithubURL       *string    `gorm:"unique" json:"github_url"`
}
