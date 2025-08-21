package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-account-service/utils"
)

func (ctrl *Controller) LoginWithIdentifierAndPassword(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Basic Login] Received login request")

	var req ClientRequestBasicLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Basic Login] Failed to bind request: %v", err)
		utils.JSON400(c, "Invalid request format")
		return
	}

	if !isValidLoginRequest(req) {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Basic Login] Invalid login request for identifier: %s - missing required fields", req)
		utils.JSON400(c, "Email/Username/Phone and Password are required")
		return
	}

	deviceID := c.GetHeader("X-Device-ID")
	if deviceID == "" {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Basic Login] Missing device ID")
		utils.JSON400(c, "X-Device-ID header is required")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Basic Login] Starting authentication for device: %s", deviceID)

	user, err := ctrl.AuthenticateUser(&req, c)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Basic Login] Authentication failed for, device: %s", deviceID)
		utils.JSON401(c, "Failed to authenticate user")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Basic Login] User authenticated successfully - UserID: %s, Device: %s", user.UserID, deviceID)

	accessToken, refreshToken, expiresAt, err := ctrl.Provider.AuthorizationServiceProvider.CreateNewToken(user.UserID, user.Permission, deviceID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Basic Login] Failed to create token for UserID: %s, Device: %s", user.UserID, deviceID)
		utils.JSON500(c, "Could not create token")
		return
	}

	expiresIn := int(time.Until(expiresAt).Seconds())

	ctrl.SetAccessCookie(c, accessToken, expiresIn)
	ctrl.SetRefreshCookie(c, refreshToken, 30*24*60*60)

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Basic Login] Login completed successfully - UserID: %s, Device: %s, ExpiresIn: %d", user.UserID, deviceID, expiresIn)

	utils.JSON200(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    expiresIn,
	})
}
