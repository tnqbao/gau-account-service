package repositories

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/models"
	"gorm.io/gorm"
	"strings"
)

func CreateUserCredential(userCredential *models.UserCredentials, c *gin.Context) error {
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(userCredential).Error; err != nil {
		return fmt.Errorf("error creating user credential: %v", err)
	}
	return nil
}

func UpdateUserCredential(userCredential *models.UserCredentials, c *gin.Context) (*models.UserCredentials, error) {
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Save(&userCredential).Error; err != nil {
		return nil, err
	}
	return userCredential, nil
}

func DeleteUserCredentialById(id uuid.UUID, c *gin.Context) error {
	db := c.MustGet("db").(*gorm.DB)
	var userCredential models.UserCredentials
	if err := db.Where("id = ?", id).First(&userCredential).Error; err != nil {
		fmt.Errorf("error finding user credential with id %s: %v", id, err)
	}
	if err := db.Delete(&userCredential).Error; err != nil {
		fmt.Errorf("error deleting user credential with id %s: %v", id, err)
	}
	return nil
}

func CreateUserInfo(userInfo *models.UserInformation, c *gin.Context) error {
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(userInfo).Error; err != nil {
		return fmt.Errorf("error creating user info: %v", err)
	}
	return nil
}

func UpdateUserInfo(userInfo *models.UserInformation, c *gin.Context) (*models.UserInformation, error) {
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Save(&userInfo).Error; err != nil {
		return nil, err
	}
	return userInfo, nil
}

func DeleteUserInfoById(id uuid.UUID, c *gin.Context) error {
	db := c.MustGet("db").(*gorm.DB)
	var userInfo models.UserInformation
	if err := db.Where("id = ?", id).First(&userInfo).Error; err != nil {
		return fmt.Errorf("error finding user info with id %s: %v", id, err)
	}
	if err := db.Delete(&userInfo).Error; err != nil {
		return fmt.Errorf("error deleting user info with id %s: %v", id, err)
	}
	return nil
}

func GetUserCredentialById(id uuid.UUID, c *gin.Context) (*models.UserCredentials, error) {
	db := c.MustGet("db").(*gorm.DB)
	var userCredential models.UserCredentials
	if err := db.Where("id = ?", id).First(&userCredential).Error; err != nil {
		return nil, fmt.Errorf("error finding user credential with id %s: %v", id, err)
	}
	return &userCredential, nil
}

func GetUserInfoById(id uuid.UUID, c *gin.Context) (*models.UserInformation, error) {
	db := c.MustGet("db").(*gorm.DB)
	var userInfo models.UserInformation
	if err := db.Where("id = ?", id).First(&userInfo).Error; err != nil {
		return nil, fmt.Errorf("error finding user info with id %s: %v", id, err)
	}
	return &userInfo, nil
}

func UpdateUserInfoById(id uuid.UUID, userInfo *models.UserInformation, c *gin.Context) (*models.UserInformation, error) {
	db := c.MustGet("db").(*gorm.DB)
	var existingUserInfo models.UserInformation
	if err := db.Where("user_id = ?", id).First(&existingUserInfo).Error; err != nil {
		return nil, fmt.Errorf("error finding user info with id %s: %v", id, err)
	}

	if err := db.Model(&existingUserInfo).Updates(userInfo).Error; err != nil {
		return nil, fmt.Errorf("error updating user info with id %s: %v", id, err)
	}
	return &existingUserInfo, nil
}

func GetUserCredentialByEmail(email string, c *gin.Context) (*models.UserCredentials, error) {
	db := c.MustGet("db").(*gorm.DB)
	var userCredential models.UserCredentials
	if err := db.Where("email = ?", email).First(&userCredential).Error; err != nil {
		return nil, fmt.Errorf("error finding user credential with email %s: %v", email, err)
	}
	return &userCredential, nil
}

func GetUserCredentialByPhone(phone string, c *gin.Context) (*models.UserCredentials, error) {
	db := c.MustGet("db").(*gorm.DB)
	var userCredential models.UserCredentials
	if err := db.Where("phone = ?", phone).First(&userCredential).Error; err != nil {
		return nil, fmt.Errorf("error finding user credential with phone %s: %v", phone, err)
	}
	return &userCredential, nil
}

func GetUserCredentialByUsername(username string, c *gin.Context) (*models.UserCredentials, error) {
	db := c.MustGet("db").(*gorm.DB)
	var userCredential models.UserCredentials
	if err := db.Where("username = ?", username).First(&userCredential).Error; err != nil {
		return nil, fmt.Errorf("error finding user credential with username %s: %v", username, err)
	}
	return &userCredential, nil
}

func GetUserByIdentifierAndPassword(identifier, identifierType, hashedPassword string, db *gorm.DB) (*models.UserInformation, error) {
	var userInfo models.UserInformation

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

func UpdateEmailById(id uuid.UUID, newEmail string, c *gin.Context) error {
	db := c.MustGet("db").(*gorm.DB)
	var userCredential models.UserCredentials

	if err := db.Where("user_id = ?", id).First(&userCredential).Error; err != nil {
		return fmt.Errorf("error finding user credential with id %s: %v", id, err)
	}

	userCredential.Email = &newEmail
	if err := db.Save(&userCredential).Error; err != nil {
		return fmt.Errorf("error updating email for user credential with id %s: %v", id, err)
	}

	return nil
}

func UpdatePhoneById(id uuid.UUID, newPhone string, c *gin.Context) error {
	db := c.MustGet("db").(*gorm.DB)
	var userCredential models.UserCredentials

	if err := db.Where("user_id = ?", id).First(&userCredential).Error; err != nil {
		return fmt.Errorf("error finding user credential with id %s: %v", id, err)
	}

	userCredential.Phone = &newPhone
	if err := db.Save(&userCredential).Error; err != nil {
		return fmt.Errorf("error updating phone for user credential with id %s: %v", id, err)
	}

	return nil
}

func UpdateUsernameById(id uuid.UUID, newUsername string, c *gin.Context) error {
	db := c.MustGet("db").(*gorm.DB)
	var userCredential models.UserCredentials

	if err := db.Where("user_id = ?", id).First(&userCredential).Error; err != nil {
		return fmt.Errorf("error finding user credential with id %s: %v", id, err)
	}

	userCredential.Username = &newUsername
	if err := db.Save(&userCredential).Error; err != nil {
		return fmt.Errorf("error updating username for user credential with id %s: %v", id, err)
	}

	return nil
}
