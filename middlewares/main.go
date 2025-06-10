package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-account-service/config"
	"github.com/tnqbao/gau-account-service/service"
)

type Middlewares struct {
	CORSMiddleware gin.HandlerFunc
	AuthMiddleware gin.HandlerFunc
}

func NewMiddlewares(config *config.EnvConfig, service *service.Service) (*Middlewares, error) {
	cors := CORSMiddleware(config)
	auth := AuthMiddleware(config, service)

	return &Middlewares{
		CORSMiddleware: cors,
		AuthMiddleware: auth,
	}, nil
}
