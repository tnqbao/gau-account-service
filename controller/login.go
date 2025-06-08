package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/repositories"
	"github.com/tnqbao/gau-account-service/schemas"
	"log"
	"net/http"
	"time"
)

func (ctrl *Controller) LoginWithIdentifierAndPassword(c *gin.Context) {
	var req ClientRequestBasicLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Binding error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	if !isValidLoginRequest(req) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email/Username/Phone and Password are required"})
		return
	}

	user, err := ctrl.AuthenticateUser(&req, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// === Access Token ===
	accessTokenDuration := 15 * time.Minute
	if req.KeepLogin != nil && *req.KeepLogin == "true" {
		accessTokenDuration = 7 * 24 * time.Hour // 7 days
	}
	accessTokenExpiry := time.Now().Add(accessTokenDuration)

	claims := &ClaimsToken{
		UserID:         user.UserID,
		FullName:       *user.FullName,
		UserPermission: user.Permission,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken, err := ctrl.CreateAuthToken(*claims)
	if err != nil {
		log.Println("Failed to create access token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create access token"})
		return
	}

	// === Refresh Token ===
	refreshTokenPlain := ctrl.GenerateRefreshToken()
	refreshTokenHashed := ctrl.hashToken(refreshTokenPlain)
	refreshTokenExpiry := time.Now().Add(30 * 24 * time.Hour)

	refreshTokenModel := &schemas.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.UserID,
		Token:     refreshTokenHashed,
		DeviceID:  c.GetHeader("X-Device-ID"),
		ExpiresAt: refreshTokenExpiry,
	}

	if err := repositories.CreateRefreshToken(refreshTokenModel, c); err != nil {
		log.Println("Failed to save refresh token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not store refresh token"})
		return
	}

	// === Set Cookies  ===
	ctrl.SetAuthCookie(c, accessToken, int(accessTokenDuration.Seconds()))
	ctrl.SetRefreshCookie(c, refreshTokenPlain, int((30 * 24 * time.Hour).Seconds()))

	// === Response ===
	c.JSON(http.StatusOK, gin.H{
		"token":         accessToken,
		"refresh_token": refreshTokenPlain,
	})
}
