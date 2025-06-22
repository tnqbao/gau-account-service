package controller

import (
	"github.com/tnqbao/gau-account-service/config"
	"github.com/tnqbao/gau-account-service/infra"
	"github.com/tnqbao/gau-account-service/repository"
)

type Controller struct {
	Config     *config.Config
	Infra      *infra.Infra
	Repository *repository.Repository
}

func NewController(config *config.Config, infra *infra.Infra) *Controller {

	repo := repository.InitRepository(infra)
	if repo == nil {
		panic("Failed to initialize Repository")
	}
	return &Controller{
		Config:     config,
		Infra:      infra,
		Repository: repo,
	}
}
