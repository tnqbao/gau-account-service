package models

import "time"

type UserInformation struct {
	UserId          uint       `gorm:"primaryKey" json:"user_id"`
	FullName        *string    `json:"fullname"`
	Phone           *string    `gorm:"unique" json:"phone"`
	IsPhoneVerified bool       `gorm:"default:false" json:"is_phone_verified"`
	Email           *string    `gorm:"unique" json:"email"`
	IsEmailVerified bool       `gorm:"default:false" json:"is_email_verified"`
	DateOfBirth     *time.Time `json:"date_of_birth"`
}
