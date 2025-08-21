package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/entity"
	"github.com/tnqbao/gau-account-service/utils"
)

func (ctrl *Controller) RegisterWithIdentifierAndPassword(c *gin.Context) {
	var req UserBasicRegistryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("UserRequest binding error:", err)
		utils.JSON400(c, "UserRequest binding error: "+err.Error())
		return
	}
	req.Password = ctrl.HashPassword(req.Password)

	if req.Email != nil && !ctrl.IsValidEmail(*req.Email) {
		utils.JSON400(c, "Invalid email format")
		return
	}

	if req.Phone != nil && !ctrl.IsValidPhone(*req.Phone) {
		utils.JSON400(c, "Invalid phone format")
		return
	}

	if (req.FullName == "" && req.Email == nil && req.Phone == nil) || req.Password == "" {
		utils.JSON400(c, "Missing required fields: FullName/Email/Phone, or Password")
		return
	}

	// Start a database transaction
	tx := ctrl.Repository.Db.Begin()
	if tx.Error != nil {
		log.Println("Error starting transaction:", tx.Error)
		utils.JSON500(c, "Internal server error")
		return
	}
	defer func() {
		if r := recover(); r != nil {
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

	// Create the user
	if err := ctrl.Repository.CreateUserWithTransaction(tx, &user); err != nil {
		tx.Rollback()
		log.Println("Error creating user:", err)
		utils.JSON500(c, "Internal server error")
		return
	}

	// Create verification records for email and phone if provided
	if req.Email != nil && *req.Email != "" {
		emailVerification := entity.UserVerification{
			ID:         uuid.New(),
			UserID:     user.UserID,
			Method:     "email",
			Value:      *req.Email,
			IsVerified: false,
		}
		if err := ctrl.Repository.CreateUserVerificationWithTransaction(tx, &emailVerification); err != nil {
			tx.Rollback()
			log.Println("Error creating email verification:", err)
			utils.JSON500(c, "Internal server error")
			return
		}
	}

	if req.Phone != nil && *req.Phone != "" {
		phoneVerification := entity.UserVerification{
			ID:         uuid.New(),
			UserID:     user.UserID,
			Method:     "phone",
			Value:      *req.Phone,
			IsVerified: false,
		}
		if err := ctrl.Repository.CreateUserVerificationWithTransaction(tx, &phoneVerification); err != nil {
			tx.Rollback()
			log.Println("Error creating phone verification:", err)
			utils.JSON500(c, "Internal server error")
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Println("Error committing transaction:", err)
		utils.JSON500(c, "Internal server error")
		return
	}

	utils.JSON200(c, gin.H{
		"message": "User successfully created",
		"user_id": user.UserID,
	})
}
