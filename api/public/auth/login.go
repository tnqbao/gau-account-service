package api_user

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	provider "github.com/tnqbao/gau_services/api"
	"gorm.io/gorm"
)

func Authentication(c *gin.Context) {
	var req provider.ClientRequestLogin
	jwtKey := os.Getenv("JWT_SECRET")
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("UserRequest binding error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}
	if req.Username == nil || req.Password == nil || *req.Username == "" || *req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and Password are required"})
		return
	}
	hashedPassword := provider.HashPassword(*req.Password)
	user, err := verifyCredentials(c, *req.Username, hashedPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username or password invalid!"})
		return
	}

	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := &provider.ClaimsResponse{
		UserID:         user.UserId,
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

func verifyCredentials(c *gin.Context, username, password string) (provider.ServerResponseLogin, error) {
	var user provider.ServerResponseLogin
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Table("user_authentications").
		Select("user_authentications.user_id, user_authentications.permission , user_informations.full_name").
		Joins("INNER JOIN user_informations ON user_informations.user_id = user_authentications.user_id").
		Where("user_authentications.username = ? AND user_authentications.password = ?", username, password).
		First(&user).Error; err != nil {
		return provider.ServerResponseLogin{}, err
	}
	return user, nil
}
