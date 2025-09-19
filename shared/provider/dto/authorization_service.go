package dto

import "github.com/google/uuid"

type CreateTokenRequest struct {
	UserID     uuid.UUID `json:"user_id"`
	Permission string    `json:"permission"`
}

type CreateTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}
type RenewTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
