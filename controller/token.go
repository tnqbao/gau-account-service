package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

func (ctrl *Controller) RenewAccessToken(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		refreshToken, _ = c.Cookie("refresh_token")
	}

	deviceID := c.GetHeader("X-Device-ID")

	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device ID is required"})
		return
	}

	hashedRefreshToken := ctrl.hashToken(refreshToken)

	refreshTokenModel, err := ctrl.Repository.GetRefreshTokenByTokenAndDevice(hashedRefreshToken, deviceID)
	if err != nil || refreshTokenModel == nil {
		handleTokenError(c, err)
		return
	}

	user, err := ctrl.Repository.GetUserInfoFromRefreshToken(hashedRefreshToken)
	if err != nil {
		handleTokenError(c, err)
		return
	}

	// === Access Token ===
	accessTokenDuration := 15 * time.Minute
	accessTokenExpiry := time.Now().Add(accessTokenDuration)

	claims := ClaimsToken{
		JID:            refreshTokenModel.ID,
		UserID:         user.UserID,
		UserPermission: user.Permission,
		FullName:       ctrl.CheckNullString(user.FullName),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken, err := ctrl.CreateAccessToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create access token"})
		return
	}

	ctrl.SetAccessCookie(c, accessToken, int(accessTokenDuration.Seconds()))

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}
