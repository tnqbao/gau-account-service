package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/shared/utils"
)

// SendEmailVerification sends verification email to user
func (ctrl *Controller) SendEmailVerification(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("user_id")

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[SendEmailVerification] Request received for user: %s", userID)

	// Parse user ID
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[SendEmailVerification] Invalid user ID: %s", userID)
		utils.JSON400(c, "Invalid user ID")
		return
	}

	// Get user from database
	user, err := ctrl.Repository.GetUserById(parsedUserID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[SendEmailVerification] User not found: %s", userID)
		utils.JSON404(c, "User not found")
		return
	}

	// Check if user has email
	if user.Email == nil || *user.Email == "" {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[SendEmailVerification] User has no email: %s", userID)
		utils.JSON400(c, "User has no email address")
		return
	}

	// Generate verification token
	token, err := ctrl.Repository.GenerateVerificationToken(ctx, userID, *user.Email)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[SendEmailVerification] Failed to generate token for user: %s", userID)
		utils.JSON500(c, "Failed to generate verification token")
		return
	}

	// Generate verification link
	verificationLink := fmt.Sprintf("https://%s/api/v2/account/verify-email/%s", ctrl.Config.EnvConfig.DomainName, token)

	// Prepare email content
	recipientName := "User"
	if user.FullName != nil && *user.FullName != "" {
		recipientName = *user.FullName
	}

	content := fmt.Sprintf("Xin chào %s,\n\nVui lòng xác thực địa chỉ email của bạn bằng cách nhấp vào liên kết bên dưới.\n\nLiên kết này sẽ hết hạn sau 24 giờ.", recipientName)

	// Send email via RabbitMQ
	err = ctrl.Provider.EmailProducer.SendEmailConfirmation(
		ctx,
		*user.Email,
		recipientName,
		content,
		verificationLink,
	)

	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[SendEmailVerification] Failed to send email for user: %s", userID)
		utils.JSON500(c, "Failed to send verification email")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[SendEmailVerification] Verification email sent successfully to: %s", *user.Email)

	utils.JSON200(c, gin.H{
		"message": "Verification email sent successfully",
	})
}

// VerifyEmail verifies user email with token
func (ctrl *Controller) VerifyEmail(c *gin.Context) {
	ctx := c.Request.Context()
	token := c.Param("token")

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[VerifyEmail] Verification request received with token")

	// Validate token and get user info
	userID, email, err := ctrl.Repository.ValidateVerificationToken(ctx, token)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[VerifyEmail] Invalid or expired token")
		utils.JSON400(c, "Invalid or expired verification token")
		return
	}

	// Parse user ID
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[VerifyEmail] Invalid user ID from token: %s", userID)
		utils.JSON400(c, "Invalid user ID")
		return
	}

	// Update verification status
	err = ctrl.Repository.UpdateEmailVerificationStatus(parsedUserID, email)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[VerifyEmail] Failed to update verification status for user: %s", userID)
		utils.JSON500(c, "Failed to verify email")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[VerifyEmail] Email verified successfully for user: %s, email: %s", userID, email)

	utils.JSON200(c, gin.H{
		"message": "Email verified successfully",
	})
}
