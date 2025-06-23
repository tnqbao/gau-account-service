package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tnqbao/gau-account-service/schemas"
	"github.com/tnqbao/gau-account-service/utils"
	"log"
	"time"
)

func (ctrl *Controller) LoginWithIdentifierAndPassword(c *gin.Context) {
	var req ClientRequestBasicLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Binding error:", err)
		utils.JSON400(c, "Invalid request format: "+err.Error())
		return
	}

	if !isValidLoginRequest(req) {
		utils.JSON400(c, "Email/Username/Phone and Password are required")
		return
	}

	user, err := ctrl.AuthenticateUser(&req, c)
	if err != nil {
		utils.JSON401(c, err.Error())
		return
	}

	// === Refresh Token ===

	refreshTokenID, err := ctrl.Repository.AllocateRefreshTokenID(c.Request.Context())
	if err != nil {
		log.Println("Failed to allocate refresh token ID:", err)
		utils.JSON500(c, "Could not allocate refresh token ID")
		return
	}

	refreshTokenPlain := ctrl.GenerateToken()
	refreshTokenHashed := ctrl.hashToken(refreshTokenPlain)
	refreshTokenExpiry := time.Now().Add(30 * 24 * time.Hour)

	refreshTokenModel := &schemas.RefreshToken{
		ID:        refreshTokenID,
		UserID:    user.UserID,
		Token:     refreshTokenHashed,
		DeviceID:  c.GetHeader("X-Device-ID"),
		ExpiresAt: refreshTokenExpiry,
	}

	if err := ctrl.Repository.CreateRefreshToken(refreshTokenModel); err != nil {
		log.Println("Failed to save refresh token:", err)
		_ = ctrl.Repository.ReleaseID(c.Request.Context(), refreshTokenID)
		utils.JSON500(c, "Could not store refresh token")
		return
	}

	// === Access Token ===
	accessTokenDuration := 15 * time.Minute
	if req.KeepLogin != nil && *req.KeepLogin == "true" {
		accessTokenDuration = 7 * 24 * time.Hour
	}
	accessTokenExpiry := time.Now().Add(accessTokenDuration)

	claims := &ClaimsToken{
		JID:            refreshTokenModel.ID,
		UserID:         user.UserID,
		FullName:       *user.FullName,
		UserPermission: user.Permission,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken, err := ctrl.CreateAccessToken(*claims)
	if err != nil {
		log.Println("Failed to create access token:", err)
		utils.JSON500(c, "Could not create access token")
		return
	}

	// === Set Cookies  ===
	ctrl.SetAccessCookie(c, accessToken, int(accessTokenDuration.Seconds()))
	ctrl.SetRefreshCookie(c, refreshTokenPlain, int((30 * 24 * time.Hour).Seconds()))

	// === Response ===
	utils.JSON200(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshTokenPlain,
	})
}
