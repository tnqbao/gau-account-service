package repository

import (
	"fmt"
	"github.com/tnqbao/gau-account-service/schemas"
	"time"
)

func (r *Repository) CreateRefreshToken(token *schemas.RefreshToken) error {
	if err := r.db.Create(token).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteRefreshTokenByTokenAndDevice(token string, deviceID string) (int64, error) {
	result := r.db.Where("token = ? AND device_id = ?", token, deviceID).Delete(&schemas.RefreshToken{})
	return result.RowsAffected, result.Error
}

func (r *Repository) GetUserInfoFromRefreshToken(token string) (*schemas.User, error) {
	var refreshToken schemas.RefreshToken
	if err := r.db.Where("token = ?", token).First(&refreshToken).Error; err != nil {
		return nil, err
	}
	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("refresh token expired")
	}

	var user schemas.User
	if err := r.db.Select("user_id, permission, full_name").
		Where("user_id = ?", refreshToken.UserID).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetRefreshTokenByTokenAndDevice(token string, deviceID string) (*schemas.RefreshToken, error) {
	var refreshToken schemas.RefreshToken
	if err := r.db.Where("token = ? AND device_id = ?", token, deviceID).First(&refreshToken).Error; err != nil {
		return nil, err
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		if err := r.db.Delete(&refreshToken).Error; err != nil {
			return nil, fmt.Errorf("failed to delete expired refresh token: %w", err)
		}
		return nil, fmt.Errorf("refresh token expired")
	}

	return &refreshToken, nil
}
