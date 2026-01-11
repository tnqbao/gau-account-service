package infra

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/tnqbao/gau-account-service/shared/config"
)

type RedisClient struct {
	Client *redis.Client
}

func InitRedisClient(cfg *config.EnvConfig) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	log.Println("Connected to Redis:", cfg.Redis.Address)

	return &RedisClient{Client: client}
}
