package routes

import (
	"github.com/gin-gonic/gin"
	api_authed_user "github.com/tnqbao/gau_user_service/api/authed/user"
	"github.com/tnqbao/gau_user_service/api/public"
	api_public_auth "github.com/tnqbao/gau_user_service/api/public/auth"

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

				authedRoutes.GET("/:id", api_authed_user.GetUserById)
				authedRoutes.GET("/me", api_authed_user.GetMe)

				authedRoutes.DELETE("/delete", api_authed_user.DeleteUserById)
				authedRoutes.PUT("/update", api_authed_user.UpdateUserInformation)
			}
			publicRoutes := userRoutes.Group("/public")
			{
				publicRoutes.GET("/check", api_public_auth.HealthCheck)
				publicRoutes.GET("/user/:id", public.GetPublicUserById)

				publicRoutes.POST("/register", api_public_auth.Register)
				publicRoutes.POST("/login", api_public_auth.Authentication)
				publicRoutes.POST("/logout", middlewares.AuthMiddleware(), api_public_auth.Logout)

			}

		}
	}
	return r
}
