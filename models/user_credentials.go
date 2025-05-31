package models

import "github.com/google/uuid"

type UserCredentials struct {
	UserId     uuid.UUID `gorm:"primaryKey,index" json:"user_id"`
	Permission string    `gorm:"index:idx_username_permission" json:"permission"`
	Username   *string   `gorm:"unique;index:idx_username_permission" json:"username"`
	Password   *string   `gorm:"size:255" json:"password"`
	Email      *string   `gorm:"unique;index:idx_email_permission" json:"email"`
	Phone      *string   `gorm:"size:15" json:"phone"`
}
