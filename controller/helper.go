package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/providers"
	"github.com/tnqbao/gau-account-service/repositories"
	"gorm.io/gorm"
)

func verifyCredentialsByUsername(c *gin.Context, username, password string) (providers.ServerResponseLogin, error) {
	var user providers.ServerResponseLogin
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Table("user_authentications").
		Select("user_authentications.user_id, user_authentications.permission , user_informations.full_name").
		Joins("INNER JOIN user_informations ON user_informations.user_id = user_authentications.user_id").
		Where("user_authentications.username = ? AND user_authentications.password = ?", username, password).
		First(&user).Error; err != nil {
		return providers.ServerResponseLogin{}, err
	}
	return user, nil
}

func verifyCredentialsByEmail(c *gin.Context, email, password string) (providers.ServerResponseLogin, error) {
	var user providers.ServerResponseLogin
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Table("user_authentications").
		Select("user_authentications.user_id, user_authentications.permission , user_informations.full_name").
		Joins("INNER JOIN user_informations ON user_informations.user_id = user_authentications.user_id").
		Where("user_authentications.email = ? AND user_authentications.password = ?", email, password).
		First(&user).Error; err != nil {
		return providers.ServerResponseLogin{}, err
	}

	return user, nil
}

func verifyCredentialsByPhone(c *gin.Context, phone, password string) (providers.ServerResponseLogin, error) {
	var user providers.ServerResponseLogin
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Table("user_authentications").
		Select("user_authentications.user_id, user_authentications.permission , user_informations.full_name").
		Joins("INNER JOIN user_informations ON user_informations.user_id = user_authentications.user_id").
		Where("user_authentications.phone = ? AND user_authentications.password = ?", phone, password).
		First(&user).Error; err != nil {
		return providers.ServerResponseLogin{}, err
	}

	return user, nil
}

func (ctrl *Controller) CreateAuthToken(claims providers.ClaimsResponse) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    claims.UserID,
		"permission": claims.UserPermission,
		"fullname":   claims.FullName,
		"exp":        claims.ExpiresAt.Unix(),
		"iat":        claims.IssuedAt.Unix(),
	})

	secretKey := ctrl.config.JWT.SecretKey
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (ctrl *Controller) SetAuthCookie(c *gin.Context, token string, timeExpired int) {
	globalDomain := ctrl.config.CORS.GlobalDomain
	c.SetCookie("auth_token", token, timeExpired, "/", globalDomain, false, true)
	c.Next()
}

func (ctrl *Controller) updateEmail(userId uuid.UUID, email *string, c *gin.Context) error {
	if email == nil {
		return nil
	}

	if !providers.IsValidEmail(*email) {
		c.JSON(400, gin.H{"error": "Invalid email format"})
		return errors.New("invalid email")
	}

	existingUser, err := repositories.GetUserCredentialByEmail(*email, c)
	if err == nil && existingUser.UserId != userId {
		c.JSON(400, gin.H{"error": "Email already in use by another user"})
		return errors.New("duplicate email")
	}

	if err := repositories.UpdateEmailById(userId, *email, c); err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return err
	}

	return nil
}

func (ctrl *Controller) updatePhone(userId uuid.UUID, phone *string, c *gin.Context) error {
	if phone == nil {
		return nil
	}

	if !providers.IsValidPhone(*phone) {
		c.JSON(400, gin.H{"error": "Invalid phone format"})
		return errors.New("invalid phone")
	}

	existingUser, err := repositories.GetUserCredentialByPhone(*phone, c)
	if err == nil && existingUser.UserId != userId {
		c.JSON(400, gin.H{"error": "Phone number already in use by another user"})
		return errors.New("duplicate phone")
	}

	if err := repositories.UpdatePhoneById(userId, *phone, c); err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return err
	}

	return nil
}
