package repositories

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/models"
	"gorm.io/gorm"
	"strings"
)

func CreateUser(user *models.User, c *gin.Context) error {
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(user).Error; err != nil {
		return fmt.Errorf("error creating user credential: %v", err)
	}
	return nil
}

func UpdateUser(user *models.User, c *gin.Context) (*models.User, error) {
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Save(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func DeleteUser(id uuid.UUID, c *gin.Context) error {
	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("user_id = ?", id).First(&user).Error; err != nil {
		fmt.Errorf("error finding user with id %s: %v", id, err)
	}
	if err := db.Delete(&user).Error; err != nil {
		fmt.Errorf("error deleting user with id %s: %v", id, err)
	}
	return nil
}

func GetUserById(id uuid.UUID, c *gin.Context) (*models.User, error) {
	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("user_id = ?", id).First(&user).Error; err != nil {
		return nil, fmt.Errorf("error finding user with id %s: %v", id, err)
	}
	return &user, nil
}

func GetUserByIdentifierAndPassword(identifier, identifierType, hashedPassword string, c *gin.Context) (*models.User, error) {
	db := c.MustGet("db").(*gorm.DB)
	var userInfo models.User

	var queryField string
	switch strings.ToLower(identifierType) {
	case "email":
		queryField = "email"
	case "phone":
		queryField = "phone"
	case "username":
		queryField = "username"
	default:
		return nil, fmt.Errorf("invalid identifier type: %s", identifierType)
	}

	if err := db.Where(fmt.Sprintf("%s = ? AND password = ?", queryField), identifier, hashedPassword).First(&userInfo).Error; err != nil {
		return nil, fmt.Errorf("user not found with %s and password: %v", queryField, err)
	}

	return &userInfo, nil
}
