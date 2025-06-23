package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/providers/helper"
	"github.com/tnqbao/gau-account-service/schemas"
	"github.com/tnqbao/gau-account-service/utils"
	"gorm.io/gorm"
	"log"
	"time"
)

func (ctrl *Controller) LoginWithGoogle(c *gin.Context) {
	var req ClientRequestGoogleAuthentication
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "invalid request")
		return
	}

	googleUser, err := helper.GetUserInfoFromGoogle(req.Token)
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
			user = &schemas.User{
				UserID:          uuid.New(),
				Email:           googleUser.Email,
				FullName:        googleUser.FullName,
				ImageURL:        googleUser.ImageURL,
				Username:        googleUser.Username,
				IsEmailVerified: googleUser.IsEmailVerified,
				Permission:      "member",
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

	// === Refresh Token ===

	refreshTokenID, err := ctrl.Repository.AllocateRefreshTokenID(c.Request.Context())
	if err != nil {
		log.Println("Failed to allocate refresh token ID:", err)
		utils.JSON500(c, "could not allocate refresh token ID")
		return
	}

	token := ctrl.GenerateToken()
	hashedRefresh := ctrl.hashToken(token)
	refreshTokenExpiry := time.Now().Add(30 * 24 * time.Hour)

	refreshToken := &schemas.RefreshToken{
		ID:        refreshTokenID,
		UserID:    user.UserID,
		Token:     hashedRefresh,
		DeviceID:  c.GetHeader("X-Device-ID"),
		CreatedAt: time.Now(),
		ExpiresAt: refreshTokenExpiry,
	}

	if err := ctrl.Repository.CreateRefreshToken(refreshToken); err != nil {
		_ = ctrl.Repository.ReleaseID(c.Request.Context(), refreshTokenID)
		utils.JSON500(c, "failed to save refresh token")
		return
	}

	// === Access Token ===
	accessTokenDuration := 15 * time.Minute
	accessTokenExpiry := time.Now().Add(accessTokenDuration)

	accessToken, err := ctrl.CreateAccessToken(ClaimsToken{
		JID:            refreshToken.ID,
		UserID:         user.UserID,
		UserPermission: user.Permission,
		FullName:       ctrl.CheckNullString(user.FullName),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	if err != nil {
		utils.JSON500(c, "failed to create access token")
		return
	}

	ctrl.SetAccessCookie(c, accessToken, int(accessTokenDuration.Seconds()))
	ctrl.SetRefreshCookie(c, token, int((30 * 24 * time.Hour).Seconds()))

	utils.JSON200(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": token,
	})
}

func (ctrl *Controller) LoginWithFacebook(c *gin.Context) {
	utils.JSON200(c, gin.H{
		"message": "Login with Facebook is not implemented yet",
	})
}
