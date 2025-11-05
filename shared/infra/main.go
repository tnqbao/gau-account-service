package infra

import (
	"github.com/tnqbao/gau-account-service/shared/config"
)

type Infra struct {
	Redis    *RedisClient
	Postgres *PostgresClient
	Logger   *LoggerClient
	RabbitMQ *RabbitMQClient
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

	logger := InitLoggerClient(cfg.EnvConfig)
	if logger == nil {
		panic("Failed to initialize Logger service")
	}

	rabbitMQ := InitRabbitMQClient(cfg.EnvConfig)
	if rabbitMQ == nil {
		panic("Failed to initialize RabbitMQ service")
	}

	infraInstance = &Infra{
		Redis:    redis,
		Postgres: postgres,
		Logger:   logger,
		RabbitMQ: rabbitMQ,
	}

	return infraInstance
}

func GetClient() *Infra {
	if infraInstance == nil {
		panic("Infra not initialized. Call InitInfra() first.")
	}
	return infraInstance
}
