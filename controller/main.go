package controller

import (
	"github.com/tnqbao/gau-account-service/config"
	"github.com/tnqbao/gau-account-service/repository"
	"github.com/tnqbao/gau-account-service/service"
)

type Controller struct {
	config     *config.Config
	service    *service.Service
	repository *repository.Repository
}

func NewController(config *config.Config) *Controller {
	svc := service.InitServices(config)
	if svc == nil {
		panic("Failed to initialize service")
	}
	repo := repository.InitRepository(svc.Postgres.DB)
	if repo == nil {
		panic("Failed to initialize repository")
	}
	return &Controller{
		config:     config,
		service:    svc,
		repository: repo,
	}
}
