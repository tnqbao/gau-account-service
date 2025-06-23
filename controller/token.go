package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tnqbao/gau-account-service/utils"
	"time"
)

func (ctrl *Controller) RenewAccessToken(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		refreshToken, _ = c.Cookie("refresh_token")
	}

	deviceID := c.GetHeader("X-Device-ID")

	if refreshToken == "" {
		utils.JSON400(c, "Refresh token is required")
		return
	}

	if deviceID == "" {
		utils.JSON400(c, "Device ID is required")
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
		utils.JSON500(c, "Could not create access token")
		return
	}

	ctrl.SetAccessCookie(c, accessToken, int(accessTokenDuration.Seconds()))

	utils.JSON200(c, gin.H{
		"access_token": accessToken,
	})
}

func (ctrl *Controller) CheckAccessToken(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		utils.JSON400(c, "Access token is required")
		return
	}

	claims, err := utils.ValidateToken(c.Request.Context(), token, ctrl.Config.EnvConfig, ctrl.Repository)
	if err != nil {
		utils.JSON401(c, err.Error())
		return
	}

	if claims == nil {
		utils.JSON401(c, "Invalid access token")
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Access token is valid",
	})
}
