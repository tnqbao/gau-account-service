package auth

import (
	"github.com/tnqbao/gau_user_service/providers"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func Authentication(c *gin.Context) {
	var req providers.ClientRequestLogin
	jwtKey := os.Getenv("JWT_SECRET")
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("UserRequest binding error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}
	if (req.Username == nil && req.Email == nil) || req.Password == nil || (*req.Username == "" && *req.Email == "") || *req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email/Username and Password are required"})
		return
	}

	var user providers.ServerResponseLogin
	var err error

	hashedPassword := providers.HashPassword(*req.Password)
	if req.Username != nil && *req.Username != "" {
		user, err = verifyCredentialsByUsername(c, *req.Username, hashedPassword)
	} else if req.Email != nil && *req.Email != "" {
		user, err = verifyCredentialsByEmail(c, *req.Email, hashedPassword)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email/Username and Password are required"})
		return
	}

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/email or password"})
		return
	}

	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := &providers.ClaimsResponse{
		UserID:         user.UserId,
		FullName:       user.FullName,
		UserPermission: user.Permission,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	var timeExpired int
	if req.KeepLogin != nil && *req.KeepLogin == "true" {
		timeExpired = 3600 * 24 * 30
	} else {
		timeExpired = 0
	}
	c.SetCookie("auth_token", tokenString, timeExpired, "/", os.Getenv("GLOBAL_DOMAIN"), false, true)
	c.JSON(http.StatusOK, gin.H{"token": tokenString, "user": user})
}

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
