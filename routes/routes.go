package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-account-service/config"
	"github.com/tnqbao/gau-account-service/controller"
	"github.com/tnqbao/gau-account-service/infra"
	"github.com/tnqbao/gau-account-service/middlewares"
)

func SetupRouter(config *config.Config) *gin.Engine {
	inf := infra.InitInfra(config)
	ctrl := controller.NewController(config, inf)

	r := gin.Default()
	useMiddlewares, err := middlewares.NewMiddlewares(ctrl)
	if err != nil {
		panic(err)
	}

	r.Use(useMiddlewares.CORSMiddleware)
	apiRoutes := r.Group("/api/v2/account/")
	{
		identifierRoutes := apiRoutes.Group("/basic")
		{
			identifierRoutes.POST("/register", ctrl.RegisterWithIdentifierAndPassword)
			identifierRoutes.POST("/login", ctrl.LoginWithIdentifierAndPassword)
		}

		profileRoutes := apiRoutes.Group("/profile")
		{
			profileRoutes.Use(useMiddlewares.AuthMiddleware)
			// Basic profile info (no security data)
			profileRoutes.GET("/basic", ctrl.GetAccountBasicInfo)
			// Security info only (verifications & MFA)
			profileRoutes.GET("/security", ctrl.GetAccountSecurityInfo)
			// Complete info (basic + security)
			profileRoutes.GET("/complete", ctrl.GetAccountCompleteInfo)

			// Legacy endpoints for backward compatibility
			profileRoutes.GET("/", ctrl.GetAccountInfo) // -> maps to basic

			// Update and avatar endpoints
			profileRoutes.PUT("/", ctrl.UpdateAccountInfo)
			profileRoutes.PATCH("/avatar", ctrl.UpdateAvatarImage)
		}

		apiRoutes.POST("/logout", useMiddlewares.AuthMiddleware, ctrl.Logout)

		ssoRoutes := apiRoutes.Group("/sso")
		{
			ssoRoutes.POST("/google", ctrl.LoginWithGoogle)
		}
		apiRoutes.GET("/", ctrl.CheckHealth)
	}
	return r
}
