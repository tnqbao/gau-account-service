package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/config"
	"net/http"
)

func AuthMiddleware(config *config.EnvConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("auth_token")
		if err != nil {
			tokenString = c.GetHeader("Authorization")
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization cookie is required"})
			c.Abort()
			return
		}

		token, err := validateToken(tokenString, config)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if userId, ok := claims["user_id"].(uuid.UUID); ok {
				c.Set("user_id", userId)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user_id format"})
				c.Abort()
				return
			}

			if permission, ok := claims["permission"].(string); ok {
				c.Set("permission", permission)
			} else {
				c.Set("permission", "")
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func validateToken(tokenString string, config *config.EnvConfig) (*jwt.Token, error) {
	jwtSecret := []byte(config.JWT.SecretKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
