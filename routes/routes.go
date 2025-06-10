package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-account-service/config"
	"github.com/tnqbao/gau-account-service/controller"
	"github.com/tnqbao/gau-account-service/middlewares"
	"github.com/tnqbao/gau-account-service/service"
)

func SetupRouter(config *config.Config) *gin.Engine {
	svc := service.InitServices(config)
	ctrl := controller.NewController(config, svc)
	r := gin.Default()

	useMiddlewares, err := middlewares.NewMiddlewares(config.EnvConfig, svc)
	if err != nil {
		panic(err)
	}

	r.Use(useMiddlewares.CORSMiddleware)
	apiRoutes := r.Group("/api/account/v2")
	{
		identifierRoutes := apiRoutes.Group("/basic")
		{
			identifierRoutes.POST("/register", ctrl.RegisterWithIdentifierAndPassword)
			identifierRoutes.POST("/login", ctrl.LoginWithIdentifierAndPassword)
		}

		profileRoutes := apiRoutes.Group("/profile")
		{
			profileRoutes.Use(useMiddlewares.AuthMiddleware)
			profileRoutes.GET("/", ctrl.GetAccountInfo)
			profileRoutes.PUT("/", ctrl.UpdateAccountInfo)
		}

		apiRoutes.GET("/token", ctrl.RenewAccessToken, useMiddlewares.AuthMiddleware)
		apiRoutes.POST("/logout", ctrl.Logout, useMiddlewares.AuthMiddleware)

		ssoRoutes := apiRoutes.Group("/sso")
		{
			ssoRoutes.POST("/google", ctrl.LoginWithGoogle)
		}
	}
	return r
}
