package service

import (
	"github.com/tnqbao/gau-account-service/config"
)

type Service struct {
	Redis    *RedisService
	Postgres *PostgresService
}

var serviceInstance *Service

func InitServices(cfg *config.Config) *Service {
	if serviceInstance != nil {
		return serviceInstance
	}

	redis := InitRedisService(cfg.EnvConfig)
	if redis == nil {
		panic("Failed to initialize Redis service")
	}

	postgres := InitPostgresService(cfg.EnvConfig)
	if postgres == nil {
		panic("Failed to initialize Postgres service")
	}

	serviceInstance = &Service{
		Redis:    redis,
		Postgres: postgres,
	}

	return serviceInstance
}

func GetService() *Service {
	if serviceInstance == nil {
		panic("Service not initialized. Call InitServices() first.")
	}
	return serviceInstance
}
