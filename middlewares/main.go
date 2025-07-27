package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-account-service/controller"
)

type Middlewares struct {
	CORSMiddleware gin.HandlerFunc
	AuthMiddleware gin.HandlerFunc
}

func NewMiddlewares(ctrl *controller.Controller) (*Middlewares, error) {
	cors := CORSMiddleware(ctrl.Config.EnvConfig)
	auth := AuthMiddleware(ctrl.Provider.AuthorizationServiceProvider, ctrl.Config.EnvConfig)

	return &Middlewares{
		CORSMiddleware: cors,
		AuthMiddleware: auth,
	}, nil
}
