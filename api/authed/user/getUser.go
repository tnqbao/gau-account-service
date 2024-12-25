package api_user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	provider "github.com/tnqbao/gau_services/api"
	"github.com/tnqbao/gau_services/models"
	"gorm.io/gorm"
)

func GetUserById(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	tokenId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
		return
	}

	tokenIdUint, ok := tokenId.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user_id format"})
		return
	}

	permission, exists := c.Get("permission")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "permission not found"})
		return
	}

	if uint64(tokenIdUint) != id && permission != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden, "error": "You don't have permission to access this resource!"})
		return
	}
	var user models.UserAuthentication
	var userInfo models.UserInformation
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, "user_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.First(&userInfo, "user_id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	response := provider.ServerResp{
		UserId:     user.UserId,
		Fullame:    provider.ToString(userInfo.FullName),
		Email:      provider.ToString(userInfo.Email),
		Phone:      provider.ToString(userInfo.Phone),
		DateBirth:  provider.FormatDateToString(userInfo.DateOfBirth),
		Permission: user.Permission,
	}
	c.JSON(http.StatusOK, response)
}

func GetMe(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	tokenId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
		return
	}

	var user models.UserAuthentication
	var userInfo models.UserInformation
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, "user_id = ?", tokenId).Error; err != nil {
			return err
		}
		if err := tx.First(&userInfo, "user_id = ?", tokenId).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	response := provider.ServerResp{
		UserId:     user.UserId,
		Fullame:    provider.ToString(userInfo.FullName),
		Email:      provider.ToString(userInfo.Email),
		Phone:      provider.ToString(userInfo.Phone),
		DateBirth:  provider.FormatDateToString(userInfo.DateOfBirth),
		Permission: user.Permission,
	}
	c.JSON(http.StatusOK, response)
}
