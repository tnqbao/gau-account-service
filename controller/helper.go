package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tnqbao/gau-account-service/providers"
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
