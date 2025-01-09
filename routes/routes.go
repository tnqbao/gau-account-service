package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_user_service/api"
	"github.com/tnqbao/gau_user_service/api/auth"
	"github.com/tnqbao/gau_user_service/api/user"
	"github.com/tnqbao/gau_user_service/middlewares"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.CORSMiddleware())
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})
	apiRoutes := r.Group("/api")
	{
		userRoutes := apiRoutes.Group("/user")
		{
			authedRoutes := userRoutes.Group("/")
			{
				authedRoutes.Use(middlewares.AuthMiddleware())

				authedRoutes.GET("/:id", user.GetUserById)
				authedRoutes.GET("/me", user.GetMe)

				authedRoutes.DELETE("/delete", user.DeleteUserById)
				authedRoutes.PUT("/update", user.UpdateUserInformation)
			}

			authRotues := userRoutes.Group("/auth")
			{
				authRotues.POST("/login", auth.Authentication)
				authRotues.PUT("/register", auth.Register)
				authRotues.POST("/logout", auth.Logout, middlewares.AuthMiddleware())
			}

			publicRoutes := userRoutes.Group("/public")
			{
				publicRoutes.GET("/check", api_user.HealthCheck)
				publicRoutes.GET("/:id", user.GetPublicUserById)

			}
		}
	}
	return r
}
