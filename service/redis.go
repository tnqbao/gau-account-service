package service

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/tnqbao/gau-account-service/config"
)

type RedisService struct {
	client *redis.Client
}

func InitRedisService(cfg *config.EnvConfig) *RedisService {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	log.Println("Connected to Redis:", cfg.Redis.Address)

	return &RedisService{client: client}
}

func (r *RedisService) Set(key string, value string) error {
	return r.client.Set(context.Background(), key, value, 0).Err()
}

func (r *RedisService) Get(key string) (string, error) {
	return r.client.Get(context.Background(), key).Result()
}
