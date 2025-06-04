package controller

import (
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/models"
	"github.com/tnqbao/gau-account-service/providers"
	"github.com/tnqbao/gau-account-service/repositories"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) RegisterWithIdentifierAndPassword(c *gin.Context) {
	var req UserRegistryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("UserRequest binding error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserRequest binding error: " + err.Error()})
		return
	}
	req.Password = providers.HashPassword(req.Password)

	if req.Email != nil && !providers.IsValidEmail(*req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	if req.Phone != nil && !providers.IsValidPhone(*req.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone format"})
		return
	}

	if (req.FullName == "" && req.Email == nil && req.Phone == nil) || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields: FullName/Email/Phone, or Password"})
		return
	}

	user := models.User{
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

	if err := repositories.CreateUser(&user, c); err != nil {
		log.Println("Error creating user :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User successfully created",
	})
}
