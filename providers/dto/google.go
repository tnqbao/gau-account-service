package dto

type Google struct {
	GoogleID    string `json:"sub" binding:"required"`
	AccessToken string `json:"access_token" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Picture     string `json:"picture" binding:"required"`
}
