package api_user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_services/models"
	"gorm.io/gorm"
)

func DeleteUserById(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	tokenId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
		return
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if result := tx.Delete(&models.UserAuthentication{}, tokenId); result.Error != nil {
			return result.Error
		} else if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		if err := tx.Delete(&models.UserInformation{}, tokenId).Error; err != nil {
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
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func DeleteListUser(c *gin.Context) {

}
