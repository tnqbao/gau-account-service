package providers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type UserInfoResponse struct {
	UserId      uuid.UUID  `json:"user_id"`
	FullName    string     `json:"fullname"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	DateOfBirth *time.Time `json:"date_of_birth"`
}

type ClaimsResponse struct {
	UserID         uint   `json:"user_id"`
	UserPermission string `json:"permission"`
	FullName       string `json:"fullname"`
	jwt.RegisteredClaims
}

// login
type ServerResponseLogin struct {
	UserId     uint   `json:"user_id"`
	Permission string `json:"permission"`
	FullName   string `json:"fullname"`
}

type ClientRequestLogin struct {
	Username  *string `json:"username"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	Password  *string `json:"password"`
	KeepLogin *string `json:"keepMeLogin"`
}

type UserPublic struct {
	UserId   uint   `json:"user_id"`
	Fullname string `json:"fullname"`
}

type UserInformationUpdateReq struct {
	FullName    *string    `json:"fullname,omitempty"`
	Username    *string    `json:"username,omitempty"`
	Email       *string    `json:"email,omitempty"`
	Phone       *string    `json:"phone,omitempty"`
	Gender      *string    `json:"gender,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	FacebookURL *string    `json:"facebook_url,omitempty"`
	GitHubURL   *string    `json:"github_url,omitempty"`
}
type UserInformationUpdateRes struct {
	UserId      uuid.UUID  `json:"user_id"`
	FullName    string     `json:"fullname"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	DateOfBirth *time.Time `json:"date_of_birth"`
}
