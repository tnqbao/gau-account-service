package controller

import (
	"crypto/rand"
	"encoding/base32"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"github.com/tnqbao/gau-account-service/shared/entity"
	"github.com/tnqbao/gau-account-service/shared/utils"
)

// GenerateTOTPQR generates QR code for TOTP setup
func (ctrl *Controller) GenerateTOTPQR(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[MFA] Generate TOTP QR request received")

	userID := c.MustGet("user_id")
	if userID == nil {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[MFA] User ID is missing from context")
		utils.JSON400(c, "User ID is required")
		return
	}

	var uuidUserID uuid.UUID
	switch v := userID.(type) {
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[MFA] Invalid User ID format: %s", v)
			utils.JSON400(c, "Invalid User ID format")
			return
		}
		uuidUserID = parsed
	case uuid.UUID:
		uuidUserID = v
	default:
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[MFA] Invalid User ID type: %T", v)
		utils.JSON400(c, "Invalid User ID type")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[MFA] Generating TOTP QR for user: %s", uuidUserID.String())

	// Get user info
	user, err := ctrl.Repository.GetUserById(uuidUserID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[MFA] User not found: %s", uuidUserID.String())
		utils.JSON404(c, "User not found")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[MFA] Checking existing MFA status for user: %s", uuidUserID.String())

	// Check if user already has a TOTP MFA enabled
	mfas, err := ctrl.Repository.GetUserMFAs(uuidUserID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[MFA] Error checking user MFA status for user: %s", uuidUserID.String())
		utils.JSON500(c, "Error checking user MFA status")
		return
	}

	for _, mfa := range mfas {
		if mfa.Type == "totp" && mfa.Enabled {
			ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[MFA] TOTP already enabled for user: %s", uuidUserID.String())
			utils.JSON400(c, "TOTP is already enabled for this user")
			return
		}
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[MFA] Generating TOTP secret for user: %s", uuidUserID.String())

	// Generate a random secret for TOTP
	secret := make([]byte, 20)
	_, err = rand.Read(secret)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[MFA] Failed to generate secret for user: %s", uuidUserID.String())
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

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[MFA] Creating TOTP key for user: %s, account: %s", uuidUserID.String(), accountName)

	// Generate TOTP key
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Gauas Account Service",
		AccountName: accountName,
		Secret:      secret,
	})
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[MFA] Failed to generate TOTP key for user: %s", uuidUserID.String())
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
		ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[MFA] Updating existing disabled TOTP record for user: %s", uuidUserID.String())
		existingMFA.Secret = &secretString
		if err := ctrl.Repository.UpdateUserMFA(existingMFA); err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[MFA] Failed to update MFA record for user: %s", uuidUserID.String())
			utils.JSON500(c, "Failed to update MFA record")
			return
		}
		ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[MFA] MFA record updated successfully for user: %s", uuidUserID.String())
	} else {
		ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[MFA] Creating new TOTP record for user: %s", uuidUserID.String())
		if err := ctrl.Repository.CreateUserMFA(&mfaRecord); err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[MFA] Failed to create MFA record for user: %s", uuidUserID.String())
			utils.JSON500(c, "Failed to create MFA record")
			return
		}
		ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[MFA] MFA record created successfully for user: %s", uuidUserID.String())
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[MFA] TOTP QR generation completed successfully for user: %s", uuidUserID.String())

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
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[MFA] Enable TOTP request received")

	userID := c.MustGet("user_id")
	if userID == nil {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[MFA] User ID is missing from context")
		utils.JSON400(c, "User ID is required")
		return
	}

	var uuidUserID uuid.UUID
	switch v := userID.(type) {
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[MFA] Invalid User ID format: %s", v)
			utils.JSON400(c, "Invalid User ID format")
			return
		}
		uuidUserID = parsed
	case uuid.UUID:
		uuidUserID = v
	default:
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[MFA] Invalid User ID type: %T", v)
		utils.JSON400(c, "Invalid User ID type")
		return
	}

	var req TOTPEnableRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[MFA] Failed to bind JSON request for user: %s", uuidUserID.String())
		utils.JSON400(c, "Invalid request format: "+err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[MFA] Enabling TOTP for user: %s", uuidUserID.String())

	// Get user's TOTP MFA record
	mfas, err := ctrl.Repository.GetUserMFAs(uuidUserID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[MFA] Error getting user MFA records for user: %s", uuidUserID.String())
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
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[MFA] No TOTP setup found for user: %s", uuidUserID.String())
		utils.JSON400(c, "No TOTP setup found. Please generate QR code first")
		return
	}

	if totpMFA.Enabled {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[MFA] TOTP already enabled for user: %s", uuidUserID.String())
		utils.JSON400(c, "TOTP is already enabled for this user")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[MFA] Verifying OTP code for user: %s", uuidUserID.String())

	// Verify the OTP code
	valid := totp.Validate(req.OTPCode, *totpMFA.Secret)
	if !valid {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[MFA] Invalid OTP code provided for user: %s", uuidUserID.String())
		utils.JSON400(c, "Invalid OTP code")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[MFA] OTP code verified successfully, enabling TOTP for user: %s", uuidUserID.String())

	// Enable TOTP
	totpMFA.Enabled = true
	now := time.Now()
	totpMFA.VerifiedAt = &now

	if err := ctrl.Repository.UpdateUserMFA(totpMFA); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[MFA] Failed to enable TOTP for user: %s", uuidUserID.String())
		utils.JSON500(c, "Failed to enable TOTP")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[MFA] TOTP enabled successfully for user: %s", uuidUserID.String())

	utils.JSON200(c, gin.H{
		"message": "TOTP has been successfully enabled",
		"enabled": true,
	})
}

// VerifyTOTP verifies TOTP code during authentication
func (ctrl *Controller) VerifyTOTP(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[MFA] TOTP verification request received")

	userID := c.MustGet("user_id")
	if userID == nil {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[MFA] User ID is missing from context")
		utils.JSON400(c, "User ID is required")
		return
	}

	var uuidUserID uuid.UUID
	switch v := userID.(type) {
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[MFA] Invalid User ID format: %s", v)
			utils.JSON400(c, "Invalid User ID format")
			return
		}
		uuidUserID = parsed
	case uuid.UUID:
		uuidUserID = v
	default:
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[MFA] Invalid User ID type: %T", v)
		utils.JSON400(c, "Invalid User ID type")
		return
	}

	var req TOTPVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[MFA] Failed to bind JSON request for user: %s", uuidUserID.String())
		utils.JSON400(c, "Invalid request format: "+err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[MFA] Verifying TOTP for user: %s, device: %s", uuidUserID.String(), req.DeviceID)

	// Get user info
	user, err := ctrl.Repository.GetUserById(uuidUserID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[MFA] User not found: %s", uuidUserID.String())
		utils.JSON404(c, "User not found")
		return
	}

	// Get user's TOTP MFA record
	mfas, err := ctrl.Repository.GetUserMFAs(uuidUserID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[MFA] Error getting user MFA records for user: %s", uuidUserID.String())
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
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[MFA] TOTP is not enabled for user: %s", uuidUserID.String())
		utils.JSON400(c, "TOTP is not enabled for this user")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[MFA] Validating OTP code for user: %s", uuidUserID.String())

	// Verify the OTP code
	valid := totp.Validate(req.OTPCode, *totpMFA.Secret)
	if !valid {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[MFA] Invalid OTP code provided for user: %s", uuidUserID.String())
		utils.JSON400(c, "Invalid OTP code")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[MFA] OTP code verified successfully, generating tokens for user: %s", uuidUserID.String())

	// Generate new JWT tokens after successful TOTP verification
	accessToken, refreshToken, expiresAt, err := ctrl.Provider.AuthorizationServiceProvider.CreateNewToken(
		user.UserID,
		user.Permission,
		req.DeviceID,
	)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[MFA] Failed to generate tokens for user: %s", uuidUserID.String())
		utils.JSON500(c, "Failed to generate tokens")
		return
	}

	expiresIn := int(time.Until(expiresAt).Seconds())

	ctrl.SetAccessCookie(c, accessToken, expiresIn)
	ctrl.SetRefreshCookie(c, refreshToken, 30*24*60*60)

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[MFA] TOTP verification completed successfully for user: %s, device: %s, expires_in: %d",
		uuidUserID.String(), req.DeviceID, expiresIn)

	utils.JSON200(c, gin.H{
		"message":       "TOTP verification successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    expiresIn,
	})
}
