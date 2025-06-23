package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-account-service/config"
	"github.com/tnqbao/gau-account-service/repository"
	"github.com/tnqbao/gau-account-service/utils"
	"net/http"
)

func AuthMiddleware(config *config.EnvConfig, repository *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string

		// 1. Try Authorization header: Bearer <token>
		tokenStr = utils.ExtractToken(c)

		// 2. If empty, try query param ?access_token=...
		if tokenStr == "" {
			tokenStr = c.Query("access_token")
		}

		// 3. If still empty, try route param /access_token/:token
		if tokenStr == "" {
			tokenStr = c.Param("token")
		}

		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort()
			return
		}

		// Validate token
		claims, err := utils.ValidateToken(c.Request.Context(), tokenStr, config, repository)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Inject claims into context
		if err := utils.InjectClaimsToContext(c, claims); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}
