package controller

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type UserRegistryReq struct {
	Username    *string   `json:"username,omitempty"`
	Password    string    `json:"password,omitempty"`
	FullName    string    `json:"fullname,omitempty"`
	Email       *string   `json:"email,omitempty"`
	Phone       *string   `json:"phone,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	Gender      string    `json:"gender,omitempty"`
}

type UserInfoResponse struct {
	Username    *string   `json:"username,omitempty"`
	Password    *string   `json:"password,omitempty"`
	FullName    string    `json:"fullname,omitempty"`
	Email       *string   `json:"email,omitempty"`
	Phone       *string   `json:"phone,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	Gender      string    `json:"gender,omitempty"`
}

type ClaimsResponse struct {
	UserID         uuid.UUID `json:"user_id"`
	UserPermission string    `json:"permission"`
	FullName       string    `json:"fullname"`
	jwt.RegisteredClaims
}
