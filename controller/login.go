package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-account-service/utils"
	"log"
	"time"
)

func (ctrl *Controller) LoginWithIdentifierAndPassword(c *gin.Context) {
	var req ClientRequestBasicLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("[Login] Binding error:", err)
		utils.JSON400(c, "Invalid request format: "+err.Error())
		return
	}

	if !isValidLoginRequest(req) {
		utils.JSON400(c, "Email/Username/Phone and Password are required")
		return
	}

	deviceID := c.GetHeader("X-Device-ID")
	if deviceID == "" {
		utils.JSON400(c, "X-Device-ID header is required")
		return
	}

	user, err := ctrl.AuthenticateUser(&req, c)
	if err != nil {
		utils.JSON401(c, err.Error())
		return
	}

	accessToken, refreshToken, expiresAt, err := ctrl.Provider.AuthorizationServiceProvider.CreateNewToken(user.UserID, user.Permission, deviceID)
	if err != nil {
		log.Println("[Login] Failed to call authorization service:", err)
		utils.JSON500(c, "Could not create token")
		return
	}

	expiresIn := int(time.Until(expiresAt).Seconds())

	ctrl.SetAccessCookie(c, accessToken, expiresIn)
	ctrl.SetRefreshCookie(c, refreshToken, 30*24*60*60)

	utils.JSON200(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    expiresIn,
	})
}
