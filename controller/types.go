package controller

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

// Basic login request structure
type ClientRequestBasicLogin struct {
	Username  *string `json:"username"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	Password  *string `json:"password"`
	KeepLogin *string `json:"keepMeLogin"`
}

// Basic registration request structure
type UserBasicRegistryReq struct {
	Username    *string   `json:"username,omitempty"`
	Password    string    `json:"password,omitempty"`
	FullName    string    `json:"fullname,omitempty"`
	Email       *string   `json:"email,omitempty"`
	Phone       *string   `json:"phone,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	Gender      string    `json:"gender,omitempty"`
}

// User information response structure
type UserInfoResponse struct {
	UserId          uuid.UUID  `json:"user_id"`
	FullName        string     `json:"fullname,omitempty"`
	Email           string     `json:"email,omitempty"`
	Phone           string     `json:"phone,omitempty"`
	DateOfBirth     *time.Time `json:"date_of_birth,omitempty"`
	GithubUrl       string     `json:"github_url,omitempty"`
	FacebookUrl     string     `json:"facebook_url,omitempty"`
	IsEmailVerified bool       `json:"is_email_verified"`
	IsPhoneVerified bool       `json:"is_phone_verified"`
}

// User information update request structure
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

// Client request structure for google register

// Client request structure for google login
type ClientRequestGoogleAuthentication struct {
	Token string `json:"token" binding:"required"`
}

// Client access token response structure
type ClaimsToken struct {
	UserID         uuid.UUID `json:"user_id"`
	UserPermission string    `json:"permission"`
	FullName       string    `json:"fullname"`
	jwt.RegisteredClaims
}
