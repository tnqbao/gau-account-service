package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-account-service/config"
	"github.com/tnqbao/gau-account-service/controller"
	"github.com/tnqbao/gau-account-service/infra"
	"github.com/tnqbao/gau-account-service/middlewares"
)

func SetupRouter(config *config.Config) *gin.Engine {
	svc := infra.InitInfra(config)
	ctrl := controller.NewController(config, svc)

	r := gin.Default()
	useMiddlewares, err := middlewares.NewMiddlewares(ctrl)
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

		tokenRoutes := apiRoutes.Group("/token")
		{
			tokenRoutes.GET("/new-access-token", ctrl.RenewAccessToken)
			tokenRoutes.GET("/check/:token", ctrl.CheckAccessToken)
		}

		apiRoutes.POST("/logout", ctrl.Logout, useMiddlewares.AuthMiddleware)

		ssoRoutes := apiRoutes.Group("/sso")
		{
			ssoRoutes.POST("/google", ctrl.LoginWithGoogle)
		}
		checkRoutes := apiRoutes.Group("/check")
		{
			//checkRoutes.GET("/deployment", ctrl.TestDeployment)
			checkRoutes.GET("/", ctrl.CheckHealth)
		}
	}
	return r
}
