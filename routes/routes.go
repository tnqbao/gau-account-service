package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-account-service/config"
	"github.com/tnqbao/gau-account-service/controller"
	"github.com/tnqbao/gau-account-service/middlewares"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, config *config.EnvConfig) *gin.Engine {
	ctrl := controller.NewController(config)
	r := gin.Default()
	useMiddlewares, err := middlewares.NewMiddlewares(config)
	if err != nil {
		panic(err)
	}

	r.Use(useMiddlewares.CORSMiddleware)
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})
	apiRoutes := r.Group("/api/account/v2/")
	{
		identifierRoutes := apiRoutes.Group("/basic")
		{
			identifierRoutes.POST("/register", ctrl.RegisterWithIdentifierAndPassword)
			identifierRoutes.POST("/login", ctrl.LoginWithIdentifierAndPassword)
		}

		//	authedRoutes := userRoutes.Group("/")
		//	{
		//		authedRoutes.Use(useMiddlewares.AuthMiddleware)
		//
		//		authedRoutes.GET("/:id", controller.GetUserById)
		//		authedRoutes.GET("/me", controller.GetMe)
		//
		//		authedRoutes.DELETE("/delete", controller.DeleteUserById)
		//		authedRoutes.PUT("/update", controller.UpdateUserInformation)
		//	}
		//
		//	authRotues := userRoutes.Group("/auth")
		//	{
		//		authRotues.POST("/login", controller.Authentication)
		//		authRotues.PUT("/register", controller.Register)
		//		authRotues.POST("/logout", controller.Logout, useMiddlewares.AuthMiddleware)
		//	}
		//
		//	publicRoutes := userRoutes.Group("/public")
		//	{
		//		publicRoutes.GET("/check", public.HealthCheck)
		//		publicRoutes.GET("/:id", public.GetPublicUserByID)
		//		publicRoutes.POST("/list", public.GetListUserPublicByIDs)
		//	}
		//
		//	userRoutes.GET("/check-deploy", controller.TestDeployment)
	}
	return r
}
