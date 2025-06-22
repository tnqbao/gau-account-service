package infra

import (
	"fmt"
	"github.com/tnqbao/gau-account-service/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type PostgresClient struct {
	DB *gorm.DB
}

func InitPostgresClient(cfg *config.EnvConfig) *PostgresClient {
	pgUser := cfg.Postgres.Username
	pgPassword := cfg.Postgres.Password
	pgHost := cfg.Postgres.HOST
	pgDB := cfg.Postgres.Database
	pgPort := cfg.Postgres.Port

	if pgUser == "" || pgPassword == "" || pgHost == "" || pgDB == "" || pgPort == "" {
		log.Fatal("One or more required Postgres secrets are missing")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
		pgHost, pgUser, pgPassword, pgDB, pgPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("PostgreSQL connected at", pgHost)

	return &PostgresClient{DB: db}
}
