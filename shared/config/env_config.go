package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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
		Password  string
		Database  int
		RedisHost string
		RedisPort string
	}
	RabbitMQ struct {
		Host     string
		Port     string
		Username string
		Password string
	}
	ExternalService struct {
		AuthorizationServiceURL string
		UploadServiceURL        string
		CDNServiceURL           string
	}
	Grafana struct {
		OTLPEndpoint string
		ServiceName  string
	}
	PrivateKey string

	Environment struct {
		Mode  string
		Group string
	}
	DomainName string
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
		config.JWT.Expire = 3600 * 24 * 7
	}

	config.CORS.AllowDomains = os.Getenv("ALLOWED_DOMAINS")
	config.CORS.GlobalDomain = os.Getenv("GLOBAL_DOMAIN")

	config.Redis.Password = os.Getenv("REDIS_PASSWORD")
	config.Redis.Database, _ = strconv.Atoi(os.Getenv("REDIS_DB"))
	if config.Redis.Database == 0 {
		config.Redis.Database = 0
	}
	config.Redis.RedisHost = os.Getenv("REDIS_HOST")
	config.Redis.RedisPort = os.Getenv("REDIS_PORT")

	// RabbitMQ
	config.RabbitMQ.Host = os.Getenv("RABBITMQ_HOST")
	if config.RabbitMQ.Host == "" {
		config.RabbitMQ.Host = "localhost"
	}
	config.RabbitMQ.Port = os.Getenv("RABBITMQ_PORT")
	if config.RabbitMQ.Port == "" {
		config.RabbitMQ.Port = "5672"
	}
	config.RabbitMQ.Username = os.Getenv("RABBITMQ_USER")
	if config.RabbitMQ.Username == "" {
		config.RabbitMQ.Username = "guest"
	}
	config.RabbitMQ.Password = os.Getenv("RABBITMQ_PASSWORD")
	if config.RabbitMQ.Password == "" {
		config.RabbitMQ.Password = "guest"
	}

	config.PrivateKey = os.Getenv("PRIVATE_KEY")

	config.ExternalService.AuthorizationServiceURL = os.Getenv("AUTHORIZATION_SERVICE_URL")
	if config.ExternalService.AuthorizationServiceURL == "" {
		config.ExternalService.AuthorizationServiceURL = "http://localhost:8080"
	}
	config.ExternalService.UploadServiceURL = os.Getenv("UPLOAD_SERVICE_URL")
	if config.ExternalService.UploadServiceURL == "" {
		config.ExternalService.UploadServiceURL = "http://localhost:8081"
	}
	config.ExternalService.CDNServiceURL = os.Getenv("CDN_SERVICE_URL")
	if config.ExternalService.CDNServiceURL == "" {
		config.ExternalService.CDNServiceURL = "http://localhost:8082"
	}

	// Grafana/OpenTelemetry
	grafanaEndpoint := os.Getenv("GRAFANA_OTLP_ENDPOINT")
	if grafanaEndpoint == "" {
		grafanaEndpoint = "https://grafana.gauas.online"
	}
	// Remove protocol for OpenTelemetry client to avoid duplicate protocols
	if strings.HasPrefix(grafanaEndpoint, "https://") {
		config.Grafana.OTLPEndpoint = strings.TrimPrefix(grafanaEndpoint, "https://")
	} else if strings.HasPrefix(grafanaEndpoint, "http://") {
		config.Grafana.OTLPEndpoint = strings.TrimPrefix(grafanaEndpoint, "http://")
	} else {
		config.Grafana.OTLPEndpoint = grafanaEndpoint
	}
	config.Grafana.ServiceName = os.Getenv("SERVICE_NAME")
	if config.Grafana.ServiceName == "" {
		config.Grafana.ServiceName = "gau-account-service"
	}

	config.Environment.Mode = os.Getenv("DEPLOY_ENV")
	if config.Environment.Mode == "" {
		config.Environment.Mode = "development"
	}

	config.Environment.Group = os.Getenv("GROUP_NAME")
	if config.Environment.Group == "" {
		config.Environment.Group = "local"
	}

	config.DomainName = os.Getenv("DOMAIN_NAME")
	if config.DomainName == "" {
		config.DomainName = "localhost:8080"
	}

	return &config
}
