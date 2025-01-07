package api_user

import (
	provider "github.com/tnqbao/gau_user_service/providers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_user_service/models"
	"gorm.io/gorm"
)

func GetPublicUserById(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}
	var userInfo models.UserInformation
	err = db.Transaction(func(tx *gorm.DB) error {
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
	c.JSON(http.StatusOK, gin.H{"id": userInfo.UserId, "fullname": provider.ToString(userInfo.FullName)})
}
