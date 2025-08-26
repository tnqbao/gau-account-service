package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/entity"
	"github.com/tnqbao/gau-account-service/utils"
)

func (ctrl *Controller) RegisterWithIdentifierAndPassword(c *gin.Context) {
	ctx := c.Request.Context()

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Register] Registration request received")

	var req UserBasicRegistryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Register] Failed to bind JSON request")
		utils.JSON400(c, "UserRequest binding error: "+err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Register] Processing registration for user with email: %v, phone: %v", req.Email != nil, req.Phone != nil)

	req.Password = ctrl.HashPassword(req.Password)

	if req.Email != nil && !ctrl.IsValidEmail(*req.Email) {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Register] Invalid email format provided: %s", *req.Email)
		utils.JSON400(c, "Invalid email format")
		return
	}

	if req.Phone != nil && !ctrl.IsValidPhone(*req.Phone) {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Register] Invalid phone format provided: %s", *req.Phone)
		utils.JSON400(c, "Invalid phone format")
		return
	}

	if (req.FullName == "" && req.Email == nil && req.Phone == nil) || req.Password == "" {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Register] Missing required fields - FullName: %s, Email: %v, Phone: %v, Password present: %v",
			req.FullName, req.Email != nil, req.Phone != nil, req.Password != "")
		utils.JSON400(c, "Missing required fields: FullName/Email/Phone, or Password")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Register] Starting user creation transaction")

	// Start a database transaction
	tx := ctrl.Repository.Db.Begin()
	if tx.Error != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, tx.Error, "[Register] Failed to start database transaction")
		utils.JSON500(c, "Internal server error")
		return
	}
	defer func() {
		if r := recover(); r != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Register] Transaction panicked - %v", r)
			tx.Rollback()
		}
	}()

	user := entity.User{
		UserID:      uuid.New(),
		Username:    req.Username,
		Password:    &req.Password,
		Email:       req.Email,
		Phone:       req.Phone,
		Permission:  "member",
		DateOfBirth: &req.DateOfBirth,
		FullName:    &req.FullName,
		Gender:      &req.Gender,
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Register] Creating user with ID: %s", user.UserID.String())

	// Create the user
	if err := ctrl.Repository.CreateUserWithTransaction(tx, &user); err != nil {
		tx.Rollback()
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Register] Failed to create user: %s", user.UserID.String())
		utils.JSON500(c, "Internal server error")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Register] User created successfully: %s", user.UserID.String())

	// Create verification records for email and phone if provided
	if req.Email != nil && *req.Email != "" {
		ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Register] Creating email verification record for user: %s", user.UserID.String())

		emailVerification := entity.UserVerification{
			ID:         uuid.New(),
			UserID:     user.UserID,
			Method:     "email",
			Value:      *req.Email,
			IsVerified: false,
		}
		if err := ctrl.Repository.CreateUserVerificationWithTransaction(tx, &emailVerification); err != nil {
			tx.Rollback()
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Register] Failed to create email verification for user: %s", user.UserID.String())
			utils.JSON500(c, "Internal server error")
			return
		}
		ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Register] Email verification record created for user: %s", user.UserID.String())
	}

	if req.Phone != nil && *req.Phone != "" {
		ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Register] Creating phone verification record for user: %s", user.UserID.String())

		phoneVerification := entity.UserVerification{
			ID:         uuid.New(),
			UserID:     user.UserID,
			Method:     "phone",
			Value:      *req.Phone,
			IsVerified: false,
		}
		if err := ctrl.Repository.CreateUserVerificationWithTransaction(tx, &phoneVerification); err != nil {
			tx.Rollback()
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Register] Failed to create phone verification for user: %s", user.UserID.String())
			utils.JSON500(c, "Internal server error")
			return
		}
		ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Register] Phone verification record created for user: %s", user.UserID.String())
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Register] Failed to commit registration transaction for user: %s", user.UserID.String())
		utils.JSON500(c, "Internal server error")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Register] Registration completed successfully for user: %s", user.UserID.String())

	utils.JSON200(c, gin.H{
		"message": "Registration successful",
		"user_id": user.UserID,
	})
}
