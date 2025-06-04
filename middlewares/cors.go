package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-account-service/config"
	"strings"
	"time"
)

func CORSMiddleware(config *config.EnvConfig) gin.HandlerFunc {
	domains := config.CORS.AllowDomains
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
