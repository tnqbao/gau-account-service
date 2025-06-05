package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-account-service/models"
	"gorm.io/gorm"
)

func CreateRefreshToken(token *models.RefreshToken, c *gin.Context) error {
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(token).Error; err != nil {
		return err
	}
	return nil
}

func DeleteRefreshTokenByTokenAndDevice(token string, deviceID string, c *gin.Context) (int64, error) {
	db := c.MustGet("db").(*gorm.DB)
	result := db.Where("token = ? AND device_id = ?", token, deviceID).Delete(&models.RefreshToken{})
	return result.RowsAffected, result.Error
}
