package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-account-service/http/controller"
	"github.com/tnqbao/gau-account-service/http/middlewares"
	"github.com/tnqbao/gau-account-service/shared/config"
	"github.com/tnqbao/gau-account-service/shared/infra"
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
			profileRoutes.PUT("/basic", ctrl.UpdateAccountBasicInfo)

			// Security info only (verifications & MFA)
			profileRoutes.GET("/security", ctrl.GetAccountSecurityInfo)
			profileRoutes.PUT("/security", ctrl.UpdateAccountSecurityInfo)

			// Complete info (basic + security)
			profileRoutes.GET("/complete", ctrl.GetAccountCompleteInfo)
			profileRoutes.PUT("/complete", ctrl.UpdateAccountCompleteInfo)

			// Legacy endpoints for backward compatibility
			profileRoutes.GET("/", ctrl.GetAccountInfo)    // -> maps to basic
			profileRoutes.PUT("/", ctrl.UpdateAccountInfo) // -> maps to complete

			// Avatar upload endpoint
			profileRoutes.PATCH("/avatar", ctrl.UpdateAvatarImage)
		}

		mfaRoutes := apiRoutes.Group("/mfa")
		{
			mfaRoutes.Use(useMiddlewares.AuthMiddleware)
			// TOTP endpoints
			mfaRoutes.GET("/totp/qr", ctrl.GenerateTOTPQR)
			mfaRoutes.POST("/totp/enable", ctrl.EnableTOTP)
			mfaRoutes.POST("/totp/verify", ctrl.VerifyTOTP)
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
