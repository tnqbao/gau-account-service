package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-account-service/shared/utils"
)

func (ctrl *Controller) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Logout] Received logout request")

	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		refreshToken, _ = c.Cookie("refresh_token")
	}
	if refreshToken == "" {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Logout] No refresh token provided")
		utils.JSON400(c, "No refresh token provided")
		return
	}

	deviceID := c.GetHeader("X-Device-ID")
	if deviceID == "" {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Logout] Device ID is required")
		utils.JSON400(c, "Device ID is required")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Logout] Starting token revocation for device: %s", deviceID)

	if err := ctrl.Provider.AuthorizationServiceProvider.RevokeToken(refreshToken, deviceID); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Logout] RevokeToken failed for device: %s", deviceID)
		utils.JSON500(c, "Failed to revoke token")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Logout] Token revoked successfully for device: %s", deviceID)

	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Logout] Logout completed successfully for device: %s", deviceID)

	utils.JSON200(c, gin.H{"message": "Logout successful"})
}
