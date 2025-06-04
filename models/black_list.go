package models

import "time"

type BlackList struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Token     string    `json:"token" gorm:"type:varchar(255);not null;uniqueIndex"`
	ExpiresAt time.Time `json:"expires_at" gorm:"type:datetime;not null"`
}
