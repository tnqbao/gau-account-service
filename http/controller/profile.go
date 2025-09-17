package controller

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	entity2 "github.com/tnqbao/gau-account-service/shared/entity"
	utils2 "github.com/tnqbao/gau-account-service/shared/utils"
	"gorm.io/gorm"
)

// GetAccountBasicInfo returns only basic account information (no security data)
func (ctrl *Controller) GetAccountBasicInfo(c *gin.Context) {
	ctx := c.Request.Context()

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Basic] Get basic account info request received")

	userId := c.MustGet("user_id")
	if userId == nil {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Profile Basic] User ID is missing from context")
		utils2.JSON400(c, "User ID is required")
		return
	}

	var uuidUserId uuid.UUID
	switch v := userId.(type) {
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Basic] Invalid User ID format: %s", v)
			utils2.JSON400(c, "Invalid User ID format")
			return
		}
		uuidUserId = parsed
	case uuid.UUID:
		uuidUserId = v
	default:
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Profile Basic] Invalid User ID type: %T", v)
		utils2.JSON400(c, "Invalid User ID type")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Basic] Fetching basic info for user: %s", uuidUserId.String())

	userInfo, err := ctrl.Repository.GetUserById(uuidUserId)
	if err != nil {
		if err.Error() == "record not found" {
			ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Profile Basic] User not found: %s", uuidUserId.String())
			utils2.JSON404(c, "User not found")
		} else {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Basic] Database error while fetching user: %s", uuidUserId.String())
			utils2.JSON500(c, "Internal server error")
		}
		return
	}

	response := UserBasicInfoResponse{
		UserId:      userInfo.UserID,
		FullName:    ctrl.CheckNullString(userInfo.FullName),
		Email:       ctrl.CheckNullString(userInfo.Email),
		Phone:       ctrl.CheckNullString(userInfo.Phone),
		GithubUrl:   ctrl.CheckNullString(userInfo.GithubURL),
		FacebookUrl: ctrl.CheckNullString(userInfo.FacebookURL),
		AvatarURL:   ctrl.CheckNullString(userInfo.AvatarURL),
		Username:    ctrl.CheckNullString(userInfo.Username),
		Gender:      ctrl.CheckNullString(userInfo.Gender),
		Permission:  userInfo.Permission,
		DateOfBirth: userInfo.DateOfBirth,
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Basic] Successfully retrieved basic info for user: %s", uuidUserId.String())
	utils2.JSON200(c, gin.H{
		"user_info": response,
	})
}

// GetAccountSecurityInfo returns only verification and MFA information
func (ctrl *Controller) GetAccountSecurityInfo(c *gin.Context) {
	ctx := c.Request.Context()

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Security] Get security info request received")

	userId := c.MustGet("user_id")
	if userId == nil {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Profile Security] User ID is missing from context")
		utils2.JSON400(c, "User ID is required")
		return
	}

	var uuidUserId uuid.UUID
	switch v := userId.(type) {
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Security] Invalid User ID format: %s", v)
			utils2.JSON400(c, "Invalid User ID format")
			return
		}
		uuidUserId = parsed
	case uuid.UUID:
		uuidUserId = v
	default:
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Profile Security] Invalid User ID type: %T", v)
		utils2.JSON400(c, "Invalid User ID type")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Security] Fetching security info for user: %s", uuidUserId.String())

	// Get user verifications
	verifications, err := ctrl.Repository.GetUserVerifications(uuidUserId)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Security] Error fetching user verifications for user: %s", uuidUserId.String())
		utils2.JSON500(c, "Error fetching user verifications")
		return
	}

	// Get user MFAs
	mfas, err := ctrl.Repository.GetUserMFAs(uuidUserId)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Security] Error fetching user MFAs for user: %s", uuidUserId.String())
		utils2.JSON500(c, "Error fetching user MFAs")
		return
	}

	// Convert verifications to response format
	var verificationInfos []UserVerificationInfo
	var isEmailVerified, isPhoneVerified bool

	for _, verification := range verifications {
		verificationInfo := UserVerificationInfo{
			ID:         verification.ID,
			Method:     verification.Method,
			Value:      verification.Value,
			IsVerified: verification.IsVerified,
			VerifiedAt: verification.VerifiedAt,
		}
		verificationInfos = append(verificationInfos, verificationInfo)

		// Set backward compatibility flags
		if verification.Method == "email" && verification.IsVerified {
			isEmailVerified = true
		}
		if verification.Method == "phone" && verification.IsVerified {
			isPhoneVerified = true
		}
	}

	// Convert MFAs to response format
	var mfaInfos []UserMFAInfo
	for _, mfa := range mfas {
		mfaInfo := UserMFAInfo{
			ID:         mfa.ID,
			Type:       mfa.Type,
			Enabled:    mfa.Enabled,
			VerifiedAt: mfa.VerifiedAt,
		}
		mfaInfos = append(mfaInfos, mfaInfo)
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Security] Successfully retrieved security info for user: %s - Verifications: %d, MFAs: %d", uuidUserId.String(), len(verificationInfos), len(mfaInfos))

	response := UserSecurityInfoResponse{
		UserId:          uuidUserId,
		IsEmailVerified: isEmailVerified,
		IsPhoneVerified: isPhoneVerified,
		Verifications:   verificationInfos,
		MFAs:            mfaInfos,
	}
	utils2.JSON200(c, gin.H{
		"user_info": response,
	})
}

// GetAccountCompleteInfo returns all account information (basic + security)
func (ctrl *Controller) GetAccountCompleteInfo(c *gin.Context) {
	ctx := c.Request.Context()

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Complete] Get complete account info request received")

	userId := c.MustGet("user_id")
	if userId == nil {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Profile Complete] User ID is missing from context")
		utils2.JSON400(c, "User ID is required")
		return
	}

	var uuidUserId uuid.UUID
	switch v := userId.(type) {
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Complete] Invalid User ID format: %s", v)
			utils2.JSON400(c, "Invalid User ID format")
			return
		}
		uuidUserId = parsed
	case uuid.UUID:
		uuidUserId = v
	default:
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Profile Complete] Invalid User ID type: %T", v)
		utils2.JSON400(c, "Invalid User ID type")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Complete] Fetching complete info for user: %s", uuidUserId.String())

	userInfo, err := ctrl.Repository.GetUserById(uuidUserId)
	if err != nil {
		if err.Error() == "record not found" {
			ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Profile Complete] User not found: %s", uuidUserId.String())
			utils2.JSON404(c, "User not found")
		} else {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Complete] Database error while fetching user: %s", uuidUserId.String())
			utils2.JSON500(c, "Internal server error")
		}
		return
	}

	// Get user verifications
	verifications, err := ctrl.Repository.GetUserVerifications(uuidUserId)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Complete] Error fetching user verifications for user: %s", uuidUserId.String())
		utils2.JSON500(c, "Error fetching user verifications")
		return
	}

	// Get user MFAs
	mfas, err := ctrl.Repository.GetUserMFAs(uuidUserId)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Complete] Error fetching user MFAs for user: %s", uuidUserId.String())
		utils2.JSON500(c, "Error fetching user MFAs")
		return
	}

	// Convert verifications to response format
	var verificationInfos []UserVerificationInfo
	var isEmailVerified, isPhoneVerified bool

	for _, verification := range verifications {
		verificationInfo := UserVerificationInfo{
			ID:         verification.ID,
			Method:     verification.Method,
			Value:      verification.Value,
			IsVerified: verification.IsVerified,
			VerifiedAt: verification.VerifiedAt,
		}
		verificationInfos = append(verificationInfos, verificationInfo)

		if verification.Method == "email" && verification.IsVerified {
			isEmailVerified = true
		}
		if verification.Method == "phone" && verification.IsVerified {
			isPhoneVerified = true
		}
	}

	// Convert MFAs to response format
	var mfaInfos []UserMFAInfo
	for _, mfa := range mfas {
		mfaInfo := UserMFAInfo{
			ID:         mfa.ID,
			Type:       mfa.Type,
			Enabled:    mfa.Enabled,
			VerifiedAt: mfa.VerifiedAt,
		}
		mfaInfos = append(mfaInfos, mfaInfo)
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Complete] Successfully retrieved complete info for user: %s - Verifications: %d, MFAs: %d", uuidUserId.String(), len(verificationInfos), len(mfaInfos))

	response := UserCompleteInfoResponse{
		UserId:          userInfo.UserID,
		FullName:        ctrl.CheckNullString(userInfo.FullName),
		Email:           ctrl.CheckNullString(userInfo.Email),
		Phone:           ctrl.CheckNullString(userInfo.Phone),
		GithubUrl:       ctrl.CheckNullString(userInfo.GithubURL),
		FacebookUrl:     ctrl.CheckNullString(userInfo.FacebookURL),
		AvatarURL:       ctrl.CheckNullString(userInfo.AvatarURL),
		Username:        ctrl.CheckNullString(userInfo.Username),
		Gender:          ctrl.CheckNullString(userInfo.Gender),
		Permission:      userInfo.Permission,
		IsEmailVerified: isEmailVerified,
		IsPhoneVerified: isPhoneVerified,
		DateOfBirth:     userInfo.DateOfBirth,
		Verifications:   verificationInfos,
		MFAs:            mfaInfos,
	}
	utils2.JSON200(c, gin.H{
		"user_info": response,
	})
}

func (ctrl *Controller) GetAccountInfo(c *gin.Context) {
	ctrl.GetAccountBasicInfo(c)
}

func (ctrl *Controller) UpdateAccountInfo(c *gin.Context) {
	ctx := c.Request.Context()

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Update] Update account info request received")

	userIdRaw := c.MustGet("user_id")
	if userIdRaw == nil {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Profile Update] User ID is missing from context")
		utils2.JSON400(c, "User ID is required")
		return
	}

	var userID uuid.UUID
	switch v := userIdRaw.(type) {
	case string:
		id, err := uuid.Parse(v)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Update] Invalid User ID format: %s", v)
			utils2.JSON400(c, "Invalid User ID format")
			return
		}
		userID = id
	case uuid.UUID:
		userID = v
	default:
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Profile Update] Invalid User ID type: %T", v)
		utils2.JSON400(c, "Invalid User ID type")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Update] Starting account update for user: %s", userID.String())

	user, err := ctrl.Repository.GetUserById(userID)
	if err != nil {
		if err.Error() == "record not found" {
			ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Profile Update] User not found: %s", userID.String())
			utils2.JSON404(c, "User not found")
		} else {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Update] Database error while fetching user: %s", userID.String())
			utils2.JSON500(c, "Internal server error")
		}
		return
	}

	var req UserInformationUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Update] Failed to bind JSON request for user: %s", userID.String())
		utils2.JSON400(c, "Invalid request format: "+err.Error())
		return
	}

	// Validate email and phone format if provided
	if req.Email != nil && !ctrl.IsValidEmail(*req.Email) {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Profile Update] Invalid email format provided for user: %s", userID.String())
		utils2.JSON400(c, "Invalid email format")
		return
	}

	if req.Phone != nil && !ctrl.IsValidPhone(*req.Phone) {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Profile Update] Invalid phone format provided for user: %s", userID.String())
		utils2.JSON400(c, "Invalid phone format")
		return
	}

	// Start a database transaction
	tx := ctrl.Repository.Db.Begin()
	if tx.Error != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, tx.Error, "[Profile Update] Failed to start transaction for user: %s", userID.String())
		utils2.JSON500(c, "Internal server error")
		return
	}
	defer func() {
		if r := recover(); r != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Profile Update] Transaction panicked for user: %s - %v", userID.String(), r)
			tx.Rollback()
		}
	}()

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Update] Updating user information for user: %s", userID.String())

	updateData := &entity2.User{
		UserID:      user.UserID,
		Username:    utils2.Coalesce(req.Username, user.Username),
		FullName:    utils2.Coalesce(req.FullName, user.FullName),
		Email:       utils2.Coalesce(req.Email, user.Email),
		Phone:       utils2.Coalesce(req.Phone, user.Phone),
		DateOfBirth: utils2.Coalesce(req.DateOfBirth, user.DateOfBirth),
		Gender:      utils2.Coalesce(req.Gender, user.Gender),
		FacebookURL: utils2.Coalesce(req.FacebookURL, user.FacebookURL),
		GithubURL:   utils2.Coalesce(req.GitHubURL, user.GithubURL),
		AvatarURL:   user.AvatarURL,
		Permission:  user.Permission,
	}

	// Update user information
	updatedUser, err := ctrl.Repository.UpdateUserWithTransaction(tx, updateData)
	if err != nil {
		tx.Rollback()
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Update] Failed to update user information for user: %s", userID.String())
		utils2.JSON500(c, "Internal server error")
		return
	}

	// Handle email verification if email is being updated
	if req.Email != nil && (user.Email == nil || *req.Email != *user.Email) {
		ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Update] Email change detected for user: %s - creating verification record", userID.String())

		// Check if verification record already exists for this email
		existingVerification, err := ctrl.Repository.GetUserVerificationByMethodAndValue(userID, "email", *req.Email)
		if err != nil {
			tx.Rollback()
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Update] Error checking email verification for user: %s", userID.String())
			utils2.JSON500(c, "Error checking email verification")
			return
		}

		if existingVerification == nil {
			// Create new email verification record
			emailVerification := entity2.UserVerification{
				ID:         uuid.New(),
				UserID:     userID,
				Method:     "email",
				Value:      *req.Email,
				IsVerified: false,
			}
			if err := ctrl.Repository.CreateUserVerificationWithTransaction(tx, &emailVerification); err != nil {
				tx.Rollback()
				ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Update] Error creating email verification for user: %s", userID.String())
				utils2.JSON500(c, "Error creating email verification")
				return
			}
			ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Update] Email verification record created for user: %s", userID.String())
		} else {
			ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Update] Email verification record already exists for user: %s", userID.String())
		}
	}

	// Handle phone verification if phone is being updated
	if req.Phone != nil && (user.Phone == nil || *req.Phone != *user.Phone) {
		ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Update] Phone change detected for user: %s - creating verification record", userID.String())

		// Check if verification record already exists for this phone
		existingVerification, err := ctrl.Repository.GetUserVerificationByMethodAndValue(userID, "phone", *req.Phone)
		if err != nil {
			tx.Rollback()
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Update] Error checking phone verification for user: %s", userID.String())
			utils2.JSON500(c, "Error checking phone verification")
			return
		}

		if existingVerification == nil {
			// Create new phone verification record
			phoneVerification := entity2.UserVerification{
				ID:         uuid.New(),
				UserID:     userID,
				Method:     "phone",
				Value:      *req.Phone,
				IsVerified: false,
			}
			if err := ctrl.Repository.CreateUserVerificationWithTransaction(tx, &phoneVerification); err != nil {
				tx.Rollback()
				ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Update] Error creating phone verification for user: %s", userID.String())
				utils2.JSON500(c, "Error creating phone verification")
				return
			}
			ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Update] Phone verification record created for user: %s", userID.String())
		} else {
			ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Update] Phone verification record already exists for user: %s", userID.String())
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Update] Failed to commit transaction for user: %s", userID.String())
		utils2.JSON500(c, "Internal server error")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Update] Successfully updated account info for user: %s", userID.String())

	utils2.JSON200(c, gin.H{
		"message":   "User information updated successfully",
		"user_info": updatedUser,
	})
}

// UpdateAccountBasicInfo updates only basic account information (no email/phone changes)
func (ctrl *Controller) UpdateAccountBasicInfo(c *gin.Context) {

	userIdRaw := c.MustGet("user_id")
	if userIdRaw == nil {
		utils2.JSON400(c, "User ID is required")
		return
	}

	var userID uuid.UUID
	switch v := userIdRaw.(type) {
	case string:
		id, err := uuid.Parse(v)
		if err != nil {
			utils2.JSON400(c, "Invalid User ID format")
			return
		}
		userID = id
	case uuid.UUID:
		userID = v
	default:
		utils2.JSON400(c, "Invalid User ID type")
		return
	}

	user, err := ctrl.Repository.GetUserById(userID)
	if err != nil {
		if err.Error() == "record not found" {
			utils2.JSON404(c, "User not found")
		} else {
			utils2.JSON500(c, "Internal server error")
		}
		return
	}

	var req UserBasicInfoUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils2.JSON400(c, "Invalid request format: "+err.Error())
		return
	}

	updateData := &entity2.User{
		UserID:      user.UserID,
		Username:    utils2.Coalesce(req.Username, user.Username),
		FullName:    utils2.Coalesce(req.FullName, user.FullName),
		Email:       user.Email, // Keep existing email
		Phone:       user.Phone, // Keep existing phone
		DateOfBirth: utils2.Coalesce(req.DateOfBirth, user.DateOfBirth),
		Gender:      utils2.Coalesce(req.Gender, user.Gender),
		FacebookURL: utils2.Coalesce(req.FacebookURL, user.FacebookURL),
		GithubURL:   utils2.Coalesce(req.GitHubURL, user.GithubURL),
		AvatarURL:   user.AvatarURL,
		Permission:  user.Permission,
		Password:    user.Password,
	}

	updatedUser, err := ctrl.Repository.UpdateUser(updateData)
	if err != nil {
		utils2.JSON500(c, "Internal server error")
		return
	}

	utils2.JSON200(c, gin.H{
		"message":   "Basic user information updated successfully",
		"user_info": updatedUser,
	})
}

// UpdateAccountSecurityInfo updates only security-related information (email/phone)
func (ctrl *Controller) UpdateAccountSecurityInfo(c *gin.Context) {
	ctx := c.Request.Context()
	userIdRaw := c.MustGet("user_id")
	if userIdRaw == nil {
		utils2.JSON400(c, "User ID is required")
		return
	}

	var userID uuid.UUID
	switch v := userIdRaw.(type) {
	case string:
		id, err := uuid.Parse(v)
		if err != nil {
			utils2.JSON400(c, "Invalid User ID format")
			return
		}
		userID = id
	case uuid.UUID:
		userID = v
	default:
		utils2.JSON400(c, "Invalid User ID type")
		return
	}

	user, err := ctrl.Repository.GetUserById(userID)
	if err != nil {
		if err.Error() == "record not found" {
			utils2.JSON404(c, "User not found")
		} else {
			utils2.JSON500(c, "Internal server error")
		}
		return
	}

	var req UserSecurityInfoUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils2.JSON400(c, "Invalid request format: "+err.Error())
		return
	}

	// Validate email and phone format if provided
	if req.Email != nil && !ctrl.IsValidEmail(*req.Email) {
		utils2.JSON400(c, "Invalid email format")
		return
	}

	if req.Phone != nil && !ctrl.IsValidPhone(*req.Phone) {
		utils2.JSON400(c, "Invalid phone format")
		return
	}

	// Start a database transaction
	tx := ctrl.Repository.Db.Begin()
	if tx.Error != nil {
		utils2.JSON500(c, "Internal server error")
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	updateData := &entity2.User{
		UserID:      user.UserID,
		Username:    user.Username, // Keep existing
		FullName:    user.FullName, // Keep existing
		Email:       utils2.Coalesce(req.Email, user.Email),
		Phone:       utils2.Coalesce(req.Phone, user.Phone),
		DateOfBirth: user.DateOfBirth, // Keep existing
		Gender:      user.Gender,      // Keep existing
		FacebookURL: user.FacebookURL, // Keep existing
		GithubURL:   user.GithubURL,   // Keep existing
		AvatarURL:   user.AvatarURL,
		Permission:  user.Permission,
		Password:    user.Password,
	}

	// Update user information
	updatedUser, err := ctrl.Repository.UpdateUserWithTransaction(tx, updateData)
	if err != nil {
		tx.Rollback()
		utils2.JSON500(c, "Internal server error")
		return
	}

	// Handle email verification if email is being updated
	if req.Email != nil && (user.Email == nil || *req.Email != *user.Email) {
		ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Update] Email change detected for user: %s - creating verification record", userID.String())

		// Check if verification record already exists for this email
		existingVerification, err := ctrl.Repository.GetUserVerificationByMethodAndValue(userID, "email", *req.Email)
		if err != nil {
			tx.Rollback()
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Update] Error checking email verification for user: %s", userID.String())
			utils2.JSON500(c, "Error checking email verification")
			return
		}

		if existingVerification == nil {
			// Create new email verification record
			emailVerification := entity2.UserVerification{
				ID:         uuid.New(),
				UserID:     userID,
				Method:     "email",
				Value:      *req.Email,
				IsVerified: false,
			}
			if err := ctrl.Repository.CreateUserVerificationWithTransaction(tx, &emailVerification); err != nil {
				tx.Rollback()
				ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Update] Error creating email verification for user: %s", userID.String())
				utils2.JSON500(c, "Error creating email verification")
				return
			}
			ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Update] Email verification record created for user: %s", userID.String())
		} else {
			ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Update] Email verification record already exists for user: %s", userID.String())
		}
	}

	// Handle phone verification if phone is being updated
	if req.Phone != nil && (user.Phone == nil || *req.Phone != *user.Phone) {
		ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Update] Phone change detected for user: %s - creating verification record", userID.String())

		// Check if verification record already exists for this phone
		existingVerification, err := ctrl.Repository.GetUserVerificationByMethodAndValue(userID, "phone", *req.Phone)
		if err != nil {
			tx.Rollback()
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Update] Error checking phone verification for user: %s", userID.String())
			utils2.JSON500(c, "Error checking phone verification")
			return
		}

		if existingVerification == nil {
			// Create new phone verification record
			phoneVerification := entity2.UserVerification{
				ID:         uuid.New(),
				UserID:     userID,
				Method:     "phone",
				Value:      *req.Phone,
				IsVerified: false,
			}
			if err := ctrl.Repository.CreateUserVerificationWithTransaction(tx, &phoneVerification); err != nil {
				tx.Rollback()
				ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Profile Update] Error creating phone verification for user: %s", userID.String())
				utils2.JSON500(c, "Error creating phone verification")
				return
			}
			ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Update] Phone verification record created for user: %s", userID.String())
		} else {
			ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Profile Update] Phone verification record already exists for user: %s", userID.String())
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		utils2.JSON500(c, "Internal server error")
		return
	}

	utils2.JSON200(c, gin.H{
		"message":   "Security information updated successfully",
		"user_info": updatedUser,
	})
}

// UpdateAccountCompleteInfo updates all account information (basic + security)
func (ctrl *Controller) UpdateAccountCompleteInfo(c *gin.Context) {
	userIdRaw := c.MustGet("user_id")
	if userIdRaw == nil {
		utils2.JSON400(c, "User ID is required")
		return
	}

	var userID uuid.UUID
	switch v := userIdRaw.(type) {
	case string:
		id, err := uuid.Parse(v)
		if err != nil {
			utils2.JSON400(c, "Invalid User ID format")
			return
		}
		userID = id
	case uuid.UUID:
		userID = v
	default:
		utils2.JSON400(c, "Invalid User ID type")
		return
	}

	user, err := ctrl.Repository.GetUserById(userID)
	if err != nil {
		if err.Error() == "record not found" {
			utils2.JSON404(c, "User not found")
		} else {
			utils2.JSON500(c, "Internal server error")
		}
		return
	}

	var req UserCompleteInfoUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils2.JSON400(c, "Invalid request format: "+err.Error())
		return
	}

	// Validate email and phone format if provided
	if req.Email != nil && !ctrl.IsValidEmail(*req.Email) {
		utils2.JSON400(c, "Invalid email format")
		return
	}

	if req.Phone != nil && !ctrl.IsValidPhone(*req.Phone) {
		utils2.JSON400(c, "Invalid phone format")
		return
	}

	// Start a database transaction
	tx := ctrl.Repository.Db.Begin()
	if tx.Error != nil {
		utils2.JSON500(c, "Internal server error")
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	updateData := &entity2.User{
		UserID:      user.UserID,
		Username:    utils2.Coalesce(req.Username, user.Username),
		FullName:    utils2.Coalesce(req.FullName, user.FullName),
		Email:       utils2.Coalesce(req.Email, user.Email),
		Phone:       utils2.Coalesce(req.Phone, user.Phone),
		DateOfBirth: utils2.Coalesce(req.DateOfBirth, user.DateOfBirth),
		Gender:      utils2.Coalesce(req.Gender, user.Gender),
		FacebookURL: utils2.Coalesce(req.FacebookURL, user.FacebookURL),
		GithubURL:   utils2.Coalesce(req.GitHubURL, user.GithubURL),
		AvatarURL:   user.AvatarURL,
		Permission:  user.Permission,
		Password:    user.Password,
	}

	// Update user information
	updatedUser, err := ctrl.Repository.UpdateUserWithTransaction(tx, updateData)
	if err != nil {
		tx.Rollback()
		utils2.JSON500(c, "Internal server error")
		return
	}

	// Handle email verification if email is being updated
	if req.Email != nil && (user.Email == nil || *req.Email != *user.Email) {
		// Check if verification record already exists for this email
		existingVerification, err := ctrl.Repository.GetUserVerificationByMethodAndValue(userID, "email", *req.Email)
		if err != nil {
			tx.Rollback()
			utils2.JSON500(c, "Error checking email verification")
			return
		}

		if existingVerification == nil {
			// Create new email verification record
			emailVerification := entity2.UserVerification{
				ID:         uuid.New(),
				UserID:     userID,
				Method:     "email",
				Value:      *req.Email,
				IsVerified: false,
			}
			if err := ctrl.Repository.CreateUserVerificationWithTransaction(tx, &emailVerification); err != nil {
				tx.Rollback()
				utils2.JSON500(c, "Error creating email verification")
				return
			}
		}
	}

	// Handle phone verification if phone is being updated
	if req.Phone != nil && (user.Phone == nil || *req.Phone != *user.Phone) {
		// Check if verification record already exists for this phone
		existingVerification, err := ctrl.Repository.GetUserVerificationByMethodAndValue(userID, "phone", *req.Phone)
		if err != nil {
			tx.Rollback()
			utils2.JSON500(c, "Error checking phone verification")
			return
		}

		if existingVerification == nil {
			// Create new phone verification record
			phoneVerification := entity2.UserVerification{
				ID:         uuid.New(),
				UserID:     userID,
				Method:     "phone",
				Value:      *req.Phone,
				IsVerified: false,
			}
			if err := ctrl.Repository.CreateUserVerificationWithTransaction(tx, &phoneVerification); err != nil {
				tx.Rollback()
				utils2.JSON500(c, "Error creating phone verification")
				return
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		utils2.JSON500(c, "Internal server error")
		return
	}

	utils2.JSON200(c, gin.H{
		"message":   "Complete user information updated successfully",
		"user_info": updatedUser,
	})
}

func (ctrl *Controller) UpdateAvatarImage(c *gin.Context) {
	userIdRaw := c.MustGet("user_id")
	if userIdRaw == nil {
		utils2.JSON400(c, "User ID is required")
		return
	}

	var userID uuid.UUID
	switch v := userIdRaw.(type) {
	case string:
		id, err := uuid.Parse(v)
		if err != nil {
			utils2.JSON400(c, "Invalid User ID format")
			return
		}
		userID = id
	case uuid.UUID:
		userID = v
	default:
		utils2.JSON400(c, "Invalid User ID type")
		return
	}

	// Change from "avatar_image" to "file" to match the upload service expectation
	file, err := c.FormFile("file")
	err = c.Request.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		utils2.JSON400(c, "Invalid form data")
		return
	}

	file = c.Request.MultipartForm.File["file"][0]

	if file.Size > 50*1024*1024 {
		utils2.JSON400(c, "File size exceeds the limit of 50MB")
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		utils2.JSON500(c, "Failed to open uploaded file: "+err.Error())
		return
	}
	defer openedFile.Close()

	fileBytes, err := io.ReadAll(openedFile)
	if err != nil {
		utils2.JSON500(c, "Failed to read uploaded file: "+err.Error())
		return
	}

	// Get content type from the uploaded file and detect if needed
	contentType := file.Header.Get("Content-Type")
	if contentType == "" || contentType == "application/octet-stream" {
		contentType = http.DetectContentType(fileBytes)
	}

	// Validate image type using helper function
	if !ctrl.ValidateImageContentType(contentType) {
		utils2.JSON400(c, "Invalid file type. Only JPEG, JPG, PNG, GIF, WEBP, SVG, and ICO are allowed")
		return
	}

	// Use GORM's Transaction method for avatar upload and database update
	var fullImageURL string
	err = ctrl.ExecuteInTransaction(func(tx *gorm.DB) error {
		// Upload avatar image
		imageURL, err := ctrl.UploadAvatarFromFile(userID, fileBytes, contentType)
		if err != nil {
			return fmt.Errorf("failed to upload image: %w", err)
		}

		// Add CDN URL prefix
		fullImageURL = fmt.Sprintf("%s/images/%s", ctrl.Config.EnvConfig.ExternalService.CDNServiceURL, imageURL)

		// Get user and update avatar URL
		user, err := ctrl.Repository.GetUserById(userID)
		if err != nil {
			if err.Error() == "record not found" {
				return fmt.Errorf("user not found")
			}
			return fmt.Errorf("failed to get user: %w", err)
		}

		user.AvatarURL = &fullImageURL

		// Update user within transaction using repository method
		if _, err := ctrl.Repository.UpdateUserWithTransaction(tx, user); err != nil {
			return fmt.Errorf("failed to update user information: %w", err)
		}

		return nil
	})

	if err != nil {
		utils2.JSON500(c, err.Error())
		return
	}

	utils2.JSON200(c, gin.H{
		"message":    "Avatar image updated successfully",
		"avatar_url": fullImageURL,
	})
}
