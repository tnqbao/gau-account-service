package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-account-service/repositories"
	"log"
	"net/http"
)

func (ctrl *Controller) Logout(c *gin.Context) {
	// Lấy refresh token từ header hoặc cookie
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		refreshToken, _ = c.Cookie("refresh_token")
	}

	if refreshToken != "" {
		hashedToken := ctrl.hashToken(refreshToken)
		deviceID := c.GetHeader("X-Device-ID")

		if err := repositories.DeleteRefreshTokenByTokenAndDevice(hashedToken, deviceID, c); err != nil {
			log.Println("Error deleting refresh token:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not log out, please try again later"})
		}
	}

	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
