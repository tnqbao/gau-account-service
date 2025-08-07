package provider

import (
	"github.com/tnqbao/gau-account-service/config"
)

type Provider struct {
	AuthorizationServiceProvider *AuthorizationServiceProvider
	UploadServiceProvider        *UploadServiceProvider
}

var provider *Provider

func InitProvider(cfg *config.EnvConfig) *Provider {
	authorizationServiceProvider := NewAuthorizationServiceProvider(cfg)
	uploadServiceProvider := NewUploadServiceProvider(cfg)
	provider = &Provider{
		AuthorizationServiceProvider: authorizationServiceProvider,
		UploadServiceProvider:        uploadServiceProvider,
	}

	return provider
}

func GetProvider() *Provider {
	if provider == nil {
		panic("Provider not initialized")
	}
	return provider
}
