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
	var req providers.UserRegistryCredentialReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("UserRequest binding error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserRequest binding error: " + err.Error()})
		return
	}
	*req.Password = providers.HashPassword(*req.Password)

	if req.Email != nil && !providers.IsValidEmail(*req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	if req.Phone != nil && !providers.IsValidPhone(*req.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone format"})
		return
	}

	if (req.FullName == nil && req.Email == nil && req.Phone == nil) || req.Password == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email/Username and Password must be provided"})
		return
	}

	log.Println("Parsed Request:", req)
	userCredentials := models.UserCredentials{
		UserId:     uuid.New(),
		Username:   &req.Username,
		Password:   req.Password,
		Email:      req.Email,
		Phone:      req.Phone,
		Permission: "member",
	}

	if err := repositories.CreateUserCredential(&userCredentials, c); err != nil {
		log.Println("Error creating user credentials:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create user credentials: " + err.Error()})
		return
	}

	userInfo := models.UserInformation{
		UserId:      userCredentials.UserId,
		FullName:    req.FullName,
		Email:       providers.CheckNullString(req.Email),
		Phone:       providers.CheckNullString(req.Phone),
		DateOfBirth: req.DateOfBirth,
	}

	if err := repositories.CreateUserInfo(&userInfo, c); err != nil {
		log.Println("Error creating user information:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create user information: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User successfully created",
	})
}
