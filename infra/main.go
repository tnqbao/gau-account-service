package infra

import (
	"github.com/tnqbao/gau-account-service/config"
)

type Infra struct {
	Redis    *RedisClient
	Postgres *PostgresClient
}

var infraInstance *Infra

func InitInfra(cfg *config.Config) *Infra {
	if infraInstance != nil {
		return infraInstance
	}

	redis := InitRedisClient(cfg.EnvConfig)
	if redis == nil {
		panic("Failed to initialize Redis service")
	}

	postgres := InitPostgresClient(cfg.EnvConfig)
	if postgres == nil {
		panic("Failed to initialize Postgres service")
	}

	infraInstance = &Infra{
		Redis:    redis,
		Postgres: postgres,
	}

	return infraInstance
}

func GetClient() *Infra {
	if infraInstance == nil {
		panic("Infra not initialized. Call InitInfra() first.")
	}
	return infraInstance
}
