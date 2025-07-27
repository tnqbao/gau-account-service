package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tnqbao/gau-account-service/config"
	"github.com/tnqbao/gau-account-service/provider"
	"github.com/tnqbao/gau-account-service/utils"
	"net/http"
)

func AuthMiddleware(authProvider *provider.AuthorizationServiceProvider, config *config.EnvConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string

		tokenStr = utils.ExtractToken(c)

		if tokenStr == "" {
			tokenStr = c.Query("access_token")
		}

		if tokenStr == "" {
			tokenStr = c.Param("token")
		}

		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort()
			return
		}

		if err := authProvider.CheckAccessToken(tokenStr); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		parsedToken, err := utils.ParseToken(tokenStr, config)
		if err != nil || !parsedToken.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
			if err := utils.InjectClaimsToContext(c, claims); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Next()
	}
}
