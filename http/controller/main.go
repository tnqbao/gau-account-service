package controller

import (
	"github.com/tnqbao/gau-account-service/shared/config"
	"github.com/tnqbao/gau-account-service/shared/infra"
	"github.com/tnqbao/gau-account-service/shared/provider"
	"github.com/tnqbao/gau-account-service/shared/repository"
)

type Controller struct {
	Config     *config.Config
	Infra      *infra.Infra
	Repository *repository.Repository
	Provider   *provider.Provider
}

func NewController(config *config.Config, infra *infra.Infra) *Controller {

	repo := repository.InitRepository(infra)
	provide := provider.InitProvider(config.EnvConfig, infra)
	if repo == nil {
		panic("Failed to initialize Repository")
	}
	return &Controller{
		Config:     config,
		Infra:      infra,
		Repository: repo,
		Provider:   provide,
	}
}
