package controller

import "github.com/tnqbao/gau-account-service/config"

type Controller struct {
	config *config.EnvConfig
}

func NewController(config *config.EnvConfig) *Controller {
	return &Controller{
		config: config,
	}
}
