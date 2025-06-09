package config

import (
	"fmt"
	"os"
	"strconv"
)

type EnvConfig struct {
	Postgres struct {
		HOST     string
		Database string
		Username string
		Password string
		Port     string
	}
	JWT struct {
		SecretKey string
		Algorithm string
		Expire    int
	}
	CORS struct {
		AllowDomains string
		GlobalDomain string
	}
	Redis struct {
		Address  string
		Password string
		Database int
	}
}

func LoadEnvConfig() *EnvConfig {
	var config EnvConfig

	// Postgres
	config.Postgres.HOST = os.Getenv("PGPOOL_HOST")
	config.Postgres.Database = os.Getenv("PGPOOL_DB")
	config.Postgres.Username = os.Getenv("PGPOOL_USER")
	config.Postgres.Password = os.Getenv("PGPOOL_PASSWORD")
	config.Postgres.Port = os.Getenv("PGPOOL_PORT")

	// JWT
	config.JWT.SecretKey = os.Getenv("JWT_SECRET_KEY")
	config.JWT.Algorithm = os.Getenv("JWT_ALGORITHM")

	if val := os.Getenv("JWT_EXPIRE"); val != "" {
		fmt.Sscanf(val, "%d", &config.JWT.Expire)
	} else {
		config.JWT.Expire = 3600 * 24 * 7 // Default to 7 days
	}

	// CORS
	config.CORS.AllowDomains = os.Getenv("ALLOWED_DOMAINS")
	config.CORS.GlobalDomain = os.Getenv("GLOBAL_DOMAIN")

	// Redis
	config.Redis.Address = os.Getenv("REDIS_ADDRESS")
	config.Redis.Password = os.Getenv("REDIS_PASSWORD")
	config.Redis.Database, _ = strconv.Atoi(os.Getenv("REDIS_DB"))
	if config.Redis.Database == 0 {
		config.Redis.Database = 0 // Default to 0 if not set
	}

	return &config
}
