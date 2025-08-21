package provider

import (
	"github.com/tnqbao/gau-account-service/config"
)

type Provider struct {
	AuthorizationServiceProvider *AuthorizationServiceProvider
	UploadServiceProvider        *UploadServiceProvider
	LoggerProvider               *LoggerProvider
}

var provider *Provider

func InitProvider(cfg *config.EnvConfig) *Provider {
	authorizationServiceProvider := NewAuthorizationServiceProvider(cfg)
	uploadServiceProvider := NewUploadServiceProvider(cfg)
	loggerProvider := NewLoggerProvider()
	provider = &Provider{
		AuthorizationServiceProvider: authorizationServiceProvider,
		UploadServiceProvider:        uploadServiceProvider,
		LoggerProvider:               loggerProvider,
	}

	return provider
}

func GetProvider() *Provider {
	if provider == nil {
		panic("Provider not initialized")
	}
	return provider
}
