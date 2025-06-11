package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/providers/helper"
	"github.com/tnqbao/gau-account-service/schemas"
	"gorm.io/gorm"
	"log"
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

	user, err := ctrl.repository.GetUserByEmail(email)
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
			if err := ctrl.repository.CreateUser(user); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create user"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}
	}

	// === Refresh Token ===

	// Get a free ID from Redis bitmap
	refreshTokenID, err := ctrl.service.Redis.AllocateRefreshTokenID(c.Request.Context())
	if err != nil {
		log.Println("Failed to allocate refresh token ID:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not allocate refresh token ID"})
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

	if err := ctrl.repository.CreateRefreshToken(refreshToken); err != nil {
		// If saving the refresh token fails, release the ID back to Redis
		_ = ctrl.service.Redis.ReleaseID(c.Request.Context(), refreshTokenID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save refresh token"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create access token"})
		return
	}

	ctrl.SetAccessCookie(c, accessToken, int(accessTokenDuration.Seconds()))
	ctrl.SetRefreshCookie(c, token, int((30 * 24 * time.Hour).Seconds()))

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": token,
	})
}

func (ctrl *Controller) LoginWithFacebook(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Login with Facebook is not implemented yet",
	})
}
