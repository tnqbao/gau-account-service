package api_user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_services/models"
	"gorm.io/gorm"
)

func UpdateUserInformation(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

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

	userUpate := models.UserInformation{}
	if err := c.ShouldBindJSON(&userUpate); err != nil {
		log.Println("UserRequest binding error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserRequest binding error: " + err.Error()})
		return
	}

	var userInfor models.UserInformation

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&userInfor, "user_id = ?", tokenIdUint).Error; err != nil {
			return err
		}
		db.Model(&userInfor).Updates(userUpate)
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
	c.JSON(http.StatusOK, gin.H{"message": "Update successful"})
}
