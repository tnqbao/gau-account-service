package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_user_service/api"
	api_user2 "github.com/tnqbao/gau_user_service/api/auth"
	api_user3 "github.com/tnqbao/gau_user_service/api/user"
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
			authedRoutes := userRoutes.Group("/authed")
			{
				authedRoutes.Use(middlewares.AuthMiddleware())

				authedRoutes.GET("/:id", api_user3.GetUserById)
				authedRoutes.GET("/me", api_user3.GetMe)

				authedRoutes.DELETE("/delete", api_user3.DeleteUserById)
				authedRoutes.PUT("/update", api_user3.UpdateUserInformation)
			}
			publicRoutes := userRoutes.Group("/public")
			{
				publicRoutes.GET("/check", api_user.HealthCheck)
				publicRoutes.GET("/user/:id", api_user3.GetPublicUserById)

				publicRoutes.POST("/register", api_user2.Register)
				publicRoutes.POST("/login", api_user2.Authentication)
				publicRoutes.POST("/logout", middlewares.AuthMiddleware(), api_user2.Logout)

			}

		}
	}
	return r
}
