package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-account-service/utils"
)

func (ctrl *Controller) Logout(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		refreshToken, _ = c.Cookie("refresh_token")
	}
	if refreshToken == "" {
		utils.JSON400(c, "No refresh token provided")
		return
	}

	deviceID := c.GetHeader("X-Device-ID")
	if deviceID == "" {
		utils.JSON400(c, "Device ID is required")
		return
	}

	if err := ctrl.Provider.AuthorizationServiceProvider.RevokeToken(refreshToken, deviceID); err != nil {
		log.Println("[Logout] RevokeToken failed:", err)
		utils.JSON500(c, "Failed to revoke token")
		return
	}

	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	utils.JSON200(c, gin.H{"message": "Logout successful"})
}
