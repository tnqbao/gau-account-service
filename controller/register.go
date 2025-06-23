package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/schemas"
	"github.com/tnqbao/gau-account-service/utils"
	"log"
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

	user := schemas.User{
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

	if err := ctrl.Repository.CreateUser(&user); err != nil {
		log.Println("Error creating user :", err)
		utils.JSON500(c, "Internal server error")
		return
	}

	utils.JSON200(c, gin.H{
		"message": "User successfully created",
	})
}
