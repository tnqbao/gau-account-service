package controller

import (
	"crypto/rand"
	"encoding/base32"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"github.com/tnqbao/gau-account-service/entity"
	"github.com/tnqbao/gau-account-service/utils"
)

// GenerateTOTPQR generates QR code for TOTP setup
func (ctrl *Controller) GenerateTOTPQR(c *gin.Context) {
	userID := c.MustGet("user_id")
	if userID == nil {
		utils.JSON400(c, "User ID is required")
		return
	}

	var uuidUserID uuid.UUID
	switch v := userID.(type) {
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			utils.JSON400(c, "Invalid User ID format")
			return
		}
		uuidUserID = parsed
	case uuid.UUID:
		uuidUserID = v
	default:
		utils.JSON400(c, "Invalid User ID type")
		return
	}

	// Get user info
	user, err := ctrl.Repository.GetUserById(uuidUserID)
	if err != nil {
		utils.JSON404(c, "User not found")
		return
	}

	// Check if user already has a TOTP MFA enabled
	mfas, err := ctrl.Repository.GetUserMFAs(uuidUserID)
	if err != nil {
		utils.JSON500(c, "Error checking user MFA status")
		return
	}

	for _, mfa := range mfas {
		if mfa.Type == "totp" && mfa.Enabled {
			utils.JSON400(c, "TOTP is already enabled for this user")
			return
		}
	}

	// Generate a random secret for TOTP
	secret := make([]byte, 20)
	_, err = rand.Read(secret)
	if err != nil {
		utils.JSON500(c, "Failed to generate secret")
		return
	}

	secretString := base32.StdEncoding.EncodeToString(secret)

	// Create account name from user email or username
	accountName := ctrl.CheckNullString(user.Email)
	if accountName == "" {
		accountName = ctrl.CheckNullString(user.Username)
	}
	if accountName == "" {
		accountName = user.UserID.String()
	}

	// Generate TOTP key
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Gauas Account Service",
		AccountName: accountName,
		Secret:      secret,
	})
	if err != nil {
		utils.JSON500(c, "Failed to generate TOTP key")
		return
	}

	// Create or update MFA record (not enabled yet)
	mfaRecord := entity.UserMFA{
		ID:      uuid.New(),
		UserID:  uuidUserID,
		Type:    "totp",
		Secret:  &secretString,
		Enabled: false,
	}

	// Check if there's an existing disabled TOTP record and update it
	var existingMFA *entity.UserMFA
	for _, mfa := range mfas {
		if mfa.Type == "totp" && !mfa.Enabled {
			existingMFA = &mfa
			break
		}
	}

	if existingMFA != nil {
		existingMFA.Secret = &secretString
		if err := ctrl.Repository.UpdateUserMFA(existingMFA); err != nil {
			utils.JSON500(c, "Failed to update MFA record")
			return
		}
	} else {
		if err := ctrl.Repository.CreateUserMFA(&mfaRecord); err != nil {
			utils.JSON500(c, "Failed to create MFA record")
			return
		}
	}

	utils.JSON200(c, gin.H{
		"qr_code": key.URL(),
		"secret":  secretString,
		"account": accountName,
		"issuer":  "Gauas Account Service",
		"message": "Scan this QR code with your authenticator app",
	})
}

// EnableTOTP verifies the first OTP and enables TOTP for the user
func (ctrl *Controller) EnableTOTP(c *gin.Context) {
	userID := c.MustGet("user_id")
	if userID == nil {
		utils.JSON400(c, "User ID is required")
		return
	}

	var uuidUserID uuid.UUID
	switch v := userID.(type) {
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			utils.JSON400(c, "Invalid User ID format")
			return
		}
		uuidUserID = parsed
	case uuid.UUID:
		uuidUserID = v
	default:
		utils.JSON400(c, "Invalid User ID type")
		return
	}

	var req struct {
		OTPCode string `json:"otp_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request format: "+err.Error())
		return
	}

	// Get user's TOTP MFA record
	mfas, err := ctrl.Repository.GetUserMFAs(uuidUserID)
	if err != nil {
		utils.JSON500(c, "Error getting user MFA records")
		return
	}

	var totpMFA *entity.UserMFA
	for _, mfa := range mfas {
		if mfa.Type == "totp" {
			totpMFA = &mfa
			break
		}
	}

	if totpMFA == nil || totpMFA.Secret == nil {
		utils.JSON400(c, "No TOTP setup found. Please generate QR code first")
		return
	}

	if totpMFA.Enabled {
		utils.JSON400(c, "TOTP is already enabled for this user")
		return
	}

	// Verify the OTP code
	valid := totp.Validate(req.OTPCode, *totpMFA.Secret)
	if !valid {
		utils.JSON400(c, "Invalid OTP code")
		return
	}

	// Enable TOTP
	totpMFA.Enabled = true
	now := time.Now()
	totpMFA.VerifiedAt = &now

	if err := ctrl.Repository.UpdateUserMFA(totpMFA); err != nil {
		utils.JSON500(c, "Failed to enable TOTP")
		return
	}

	utils.JSON200(c, gin.H{
		"message": "TOTP has been successfully enabled",
		"enabled": true,
	})
}

// VerifyTOTP verifies TOTP code during authentication
func (ctrl *Controller) VerifyTOTP(c *gin.Context) {
	userID := c.MustGet("user_id")
	if userID == nil {
		utils.JSON400(c, "User ID is required")
		return
	}

	var uuidUserID uuid.UUID
	switch v := userID.(type) {
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			utils.JSON400(c, "Invalid User ID format")
			return
		}
		uuidUserID = parsed
	case uuid.UUID:
		uuidUserID = v
	default:
		utils.JSON400(c, "Invalid User ID type")
		return
	}

	var req struct {
		OTPCode  string `json:"otp_code" binding:"required"`
		DeviceID string `json:"device_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request format: "+err.Error())
		return
	}

	// Get user info
	user, err := ctrl.Repository.GetUserById(uuidUserID)
	if err != nil {
		utils.JSON404(c, "User not found")
		return
	}

	// Get user's TOTP MFA record
	mfas, err := ctrl.Repository.GetUserMFAs(uuidUserID)
	if err != nil {
		utils.JSON500(c, "Error getting user MFA records")
		return
	}

	var totpMFA *entity.UserMFA
	for _, mfa := range mfas {
		if mfa.Type == "totp" && mfa.Enabled {
			totpMFA = &mfa
			break
		}
	}

	if totpMFA == nil || totpMFA.Secret == nil {
		utils.JSON400(c, "TOTP is not enabled for this user")
		return
	}

	// Verify the OTP code
	valid := totp.Validate(req.OTPCode, *totpMFA.Secret)
	if !valid {
		utils.JSON400(c, "Invalid OTP code")
		return
	}

	// Generate new JWT tokens after successful TOTP verification
	accessToken, refreshToken, expiresAt, err := ctrl.Provider.AuthorizationServiceProvider.CreateNewToken(
		user.UserID,
		user.Permission,
		req.DeviceID,
	)
	if err != nil {
		utils.JSON500(c, "Failed to generate tokens")
		return
	}

	expiresIn := int(time.Until(expiresAt).Seconds())

	ctrl.SetAccessCookie(c, accessToken, expiresIn)
	ctrl.SetRefreshCookie(c, refreshToken, 30*24*60*60)

	utils.JSON200(c, gin.H{
		"message":       "TOTP verification successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    expiresIn,
	})
}
