package middlewares

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CORSMiddleware() gin.HandlerFunc {
	domains := os.Getenv("LIST_DOMAIN")
	domainList := strings.Split(domains, "^")
	return cors.New(cors.Config{
		AllowOrigins:     domainList,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

func AuthMiddleware() gin.HandlerFunc {
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

		token, err := validateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if userIdFloat, ok := claims["user_id"].(float64); ok {
				c.Set("user_id", uint(userIdFloat))
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

func validateToken(tokenString string) (*jwt.Token, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
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
