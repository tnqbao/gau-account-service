package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/entity"
	"github.com/tnqbao/gau-account-service/provider"
	"github.com/tnqbao/gau-account-service/utils"
	"gorm.io/gorm"
	"time"
)

func (ctrl *Controller) LoginWithGoogle(c *gin.Context) {
	var req ClientRequestGoogleAuthentication
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "invalid request")
		return
	}

	googleUser, err := provider.GetUserInfoFromGoogle(req.Token)
	if err != nil {
		utils.JSON401(c, "invalid Google token")
		return
	}

	email := ctrl.CheckNullString(googleUser.Email)
	if !ctrl.IsValidEmail(email) {
		utils.JSON400(c, "invalid email from Google")
		return
	}

	user, err := ctrl.Repository.GetUserByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			userID := uuid.New()
			user = &entity.User{
				UserID:          userID,
				Email:           googleUser.Email,
				FullName:        googleUser.FullName,
				Username:        googleUser.Username,
				IsEmailVerified: googleUser.IsEmailVerified,
				Permission:      "member",
			}

			// Upload avatar image if Google provides one
			if googleUser.AvatarURL != nil && *googleUser.AvatarURL != "" {
				imageURL, err := ctrl.UploadAvatarFromURL(userID, *googleUser.AvatarURL)
				if err != nil {
					utils.JSON500(c, "Failed to upload avatar: "+err.Error())
					return
				}
				// Add CDN URL prefix
				fullImageURL := fmt.Sprintf("%s/images/%s", ctrl.Config.EnvConfig.ExternalService.CDNServiceURL, imageURL)
				user.AvatarURL = &fullImageURL
			}

			if err := ctrl.Repository.CreateUser(user); err != nil {
				utils.JSON500(c, "cannot create user")
				return
			}
		} else {
			utils.JSON500(c, "database error")
			return
		}
	}

	deviceID := c.GetHeader("X-Device-ID")
	if deviceID == "" {
		utils.JSON400(c, "X-Device-ID header is required")
		return
	}

	accessToken, refreshToken, expiresAt, err := ctrl.Provider.AuthorizationServiceProvider.CreateNewToken(
		user.UserID,
		user.Permission,
		deviceID,
	)
	if err != nil {
		utils.JSON500(c, "failed to generate token")
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

func (ctrl *Controller) LoginWithFacebook(c *gin.Context) {
	utils.JSON200(c, gin.H{
		"message": "Login with Facebook is not implemented yet",
	})
}
