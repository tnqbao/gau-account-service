package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	// === Refresh Token ===

	// Lấy ID rảnh từ Redis bitmap
	refreshTokenID, err := ctrl.Infra.Redis.AllocateRefreshTokenID(c.Request.Context())
	if err != nil {
		log.Println("Failed to allocate refresh token ID:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not allocate refresh token ID"})
		return
	}

	refreshTokenPlain := ctrl.GenerateToken()
	refreshTokenHashed := ctrl.hashToken(refreshTokenPlain)
	refreshTokenExpiry := time.Now().Add(30 * 24 * time.Hour)

	refreshTokenModel := &schemas.RefreshToken{
		ID:        refreshTokenID,
		UserID:    user.UserID,
		Token:     refreshTokenHashed,
		DeviceID:  c.GetHeader("X-Device-ID"),
		ExpiresAt: refreshTokenExpiry,
	}

	if err := ctrl.Repository.CreateRefreshToken(refreshTokenModel); err != nil {
		log.Println("Failed to save refresh token:", err)
		// Nếu lỗi xảy ra, nên trả ID lại
		_ = ctrl.Infra.Redis.ReleaseID(c.Request.Context(), refreshTokenID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not store refresh token"})
		return
	}

	// === Access Token ===
	accessTokenDuration := 15 * time.Minute
	if req.KeepLogin != nil && *req.KeepLogin == "true" {
		accessTokenDuration = 7 * 24 * time.Hour
	}
	accessTokenExpiry := time.Now().Add(accessTokenDuration)

	claims := &ClaimsToken{
		JID:            refreshTokenModel.ID,
		UserID:         user.UserID,
		FullName:       *user.FullName,
		UserPermission: user.Permission,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken, err := ctrl.CreateAccessToken(*claims)
	if err != nil {
		log.Println("Failed to create access token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create access token"})
		return
	}

	// === Set Cookies  ===
	ctrl.SetAccessCookie(c, accessToken, int(accessTokenDuration.Seconds()))
	ctrl.SetRefreshCookie(c, refreshTokenPlain, int((30 * 24 * time.Hour).Seconds()))

	// === Response ===
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshTokenPlain,
	})
}
