package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/providers/helper"
	"github.com/tnqbao/gau-account-service/repositories"
	"github.com/tnqbao/gau-account-service/schemas"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func (ctrl *Controller) LoginWithGoogle(c *gin.Context) {
	var req ClientRequestGoogleAuthentication
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	googleUser, err := helper.GetUserInfoFromGoogle(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Google token"})
		return
	}

	email := ctrl.CheckNullString(googleUser.Email)
	if !ctrl.IsValidEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email from Google"})
		return
	}

	user, err := repositories.GetUserByEmail(email, c)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}
	}

	if user == nil {
		user = &schemas.User{
			UserID:          uuid.New(),
			Email:           googleUser.Email,
			FullName:        googleUser.FullName,
			ImageURL:        googleUser.ImageURL,
			Username:        googleUser.Username,
			IsEmailVerified: googleUser.IsEmailVerified,
			Permission:      "member",
		}
		if err := repositories.CreateUser(user, c); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create user"})
			return
		}
	}

	accessToken, err := ctrl.CreateAccessToken(ClaimsToken{
		UserID:         user.UserID,
		UserPermission: user.Permission,
		FullName:       ctrl.CheckNullString(user.FullName),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create access token"})
		return
	}

	refreshToken := ctrl.GenerateRefreshToken()
	hashedRefresh := ctrl.hashToken(refreshToken)

	rt := &schemas.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.UserID,
		Token:     hashedRefresh,
		DeviceID:  c.GetHeader("X-Device-ID"),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := repositories.CreateRefreshToken(rt, c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save refresh token"})
		return
	}

	ctrl.SetAuthCookie(c, accessToken, 15)
	ctrl.SetRefreshCookie(c, refreshToken, ctrl.config.JWT.Expire)

	c.JSON(http.StatusOK, gin.H{
		"token": accessToken,
	})
}

func (ctrl *Controller) LoginWithFacebook(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Login with Facebook is not implemented yet",
	})
}
