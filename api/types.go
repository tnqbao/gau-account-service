package api_user

import "github.com/golang-jwt/jwt/v5"

type ClientReq struct {
	Username    *string `json:"username"`
	Password    *string `json:"password"`
	Fullname    *string `json:"fullname"`
	Email       *string `json:"email"`
	Phone       *string `json:"phone"`
	DateOfBirth *string `json:"date_of_birth"`
}

type ServerResp struct {
	UserId     uint   `json:"user_id"`
	Fullame    string `json:"fullname"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	DateBirth  string `json:"date_of_birth"`
	Permission string `json:"permission"`
}

type ClaimsResponse struct {
	UserID         uint   `json:"user_id"`
	UserPermission string `json:"permission"`
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
	Password  *string `json:"password"`
	KeepLogin *string `json:"keepMeLogin"`
}
