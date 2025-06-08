package repositories

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-account-service/schemas"
	"gorm.io/gorm"
	"time"
)

func CreateRefreshToken(token *schemas.RefreshToken, c *gin.Context) error {
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(token).Error; err != nil {
		return err
	}
	return nil
}

func DeleteRefreshTokenByTokenAndDevice(token string, deviceID string, c *gin.Context) (int64, error) {
	db := c.MustGet("db").(*gorm.DB)
	result := db.Where("token = ? AND device_id = ?", token, deviceID).Delete(&schemas.RefreshToken{})
	return result.RowsAffected, result.Error
}

func GetUserInfoFromRefreshToken(token string, c *gin.Context) (*schemas.User, error) {
	db := c.MustGet("db").(*gorm.DB)

	var refreshToken schemas.RefreshToken
	if err := db.Where("token = ?", token).First(&refreshToken).Error; err != nil {
		return nil, err
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("refresh token expired")
	}

	var user schemas.User
	if err := db.Select("user_id, permission, fullname").
		Where("user_id = ?", refreshToken.UserID).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
