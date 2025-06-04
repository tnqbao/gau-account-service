package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-account-service/config"
)

type Middlewares struct {
	CORSMiddleware gin.HandlerFunc
	AuthMiddleware gin.HandlerFunc
}

func NewMiddlewares(config *config.EnvConfig) (*Middlewares, error) {
	cors := CORSMiddleware(config)
	auth := AuthMiddleware(config)

	return &Middlewares{
		CORSMiddleware: cors,
		AuthMiddleware: auth,
	}, nil
}
