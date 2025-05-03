package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_user_service/controller"
	"github.com/tnqbao/gau_user_service/controller/auth"
	"github.com/tnqbao/gau_user_service/controller/public"
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

				authedRoutes.GET("/:id", controller.GetUserById)
				authedRoutes.GET("/me", controller.GetMe)

				authedRoutes.DELETE("/delete", controller.DeleteUserById)
				authedRoutes.PUT("/update", controller.UpdateUserInformation)
			}

			authRotues := userRoutes.Group("/auth")
			{
				authRotues.POST("/login", auth.Authentication)
				authRotues.PUT("/register", auth.Register)
				authRotues.POST("/logout", auth.Logout, middlewares.AuthMiddleware())
			}

			publicRoutes := userRoutes.Group("/public")
			{
				publicRoutes.GET("/check", public.HealthCheck)
				publicRoutes.GET("/:id", public.GetPublicUserByID)
				publicRoutes.POST("/list", public.GetListUserPublicByIDs)
			}
		}
	}
	return r
}
