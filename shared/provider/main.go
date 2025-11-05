package provider

import (
	"github.com/tnqbao/gau-account-service/shared/config"
	"github.com/tnqbao/gau-account-service/shared/infra"
)

type Provider struct {
	AuthorizationServiceProvider *AuthorizationServiceProvider
	UploadServiceProvider        *UploadServiceProvider
	LoggerProvider               *LoggerProvider
	EmailProducer                *EmailProducer
}

var provider *Provider

func InitProvider(cfg *config.EnvConfig, inf *infra.Infra) *Provider {
	authorizationServiceProvider := NewAuthorizationServiceProvider(cfg)
	uploadServiceProvider := NewUploadServiceProvider(cfg)
	loggerProvider := NewLoggerProvider()
	emailProducer := NewEmailProducer(inf.RabbitMQ)
	provider = &Provider{
		AuthorizationServiceProvider: authorizationServiceProvider,
		UploadServiceProvider:        uploadServiceProvider,
		LoggerProvider:               loggerProvider,
		EmailProducer:                emailProducer,
	}

	return provider
}

func GetProvider() *Provider {
	if provider == nil {
		panic("Provider not initialized")
	}
	return provider
}
