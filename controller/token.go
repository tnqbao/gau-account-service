package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-account-service/repositories"
)

func (ctrl *Controller) RenewAccessToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Refresh token not found in cookie",
		})
		return
	}

	refreshTokenHeader := c.GetHeader("X-Refresh-Token")
	if refreshTokenHeader != "" {
		refreshToken = refreshTokenHeader
	}

	if refreshToken == "" {
		c.JSON(400, gin.H{
			"error": "Refresh token is required",
		})
		return
	}

	hashedRefreshToken := ctrl.hashToken(refreshToken)

	user, err := repositories.GetUserInfoFromRefreshToken(hashedRefreshToken, c)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(404, gin.H{
				"error": "Refresh token not found",
			})
			return
		} else if err.Error() == "refresh token expired" {
			c.JSON(401, gin.H{
				"error": "Refresh token expired",
			})
			return
		}
		c.JSON(500, gin.H{
			"error": "Internal server error",
		})
		return
	}
	claims := ClaimsToken{
		UserID:         user.UserID,
		UserPermission: user.Permission,
		FullName:       *user.FullName,
	}

	accessToken, err := ctrl.CreateAccessToken(claims)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Could not create access token",
		})
		return
	}

	c.SetCookie("auth_token", accessToken, 3600, "/", "", false, true)

	c.JSON(200, gin.H{
		"access_token": accessToken,
	})
}
