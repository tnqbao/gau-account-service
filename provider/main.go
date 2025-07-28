package provider

import (
	"github.com/tnqbao/gau-account-service/config"
)

type Provider struct {
	AuthorizationServiceProvider *AuthorizationServiceProvider
}

var provider *Provider

func InitProvider(cfg *config.EnvConfig) *Provider {
	authorizationServiceProvider := NewAuthorizationServiceProvider(cfg)

	provider = &Provider{
		AuthorizationServiceProvider: authorizationServiceProvider,
	}

	return provider
}

func GetProvider() *Provider {
	if provider == nil {
		panic("Provider not initialized")
	}
	return provider
}
