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
			publicRoutes := userRoutes.Group("/")
			{
				publicRoutes.GET("/check", api_user.HealthCheck)
				publicRoutes.GET("/user/:id", user.GetPublicUserById)

				publicRoutes.POST("/register", auth.Register)
				publicRoutes.POST("/login", auth.Authentication)
				publicRoutes.POST("/logout", middlewares.AuthMiddleware(), auth.Logout)
			}
		}
	}
	return r
}
