package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	var err error
	pg_user := os.Getenv("PGPOOL_USER")
	pg_password := os.Getenv("PGPOOL_PASSWORD")
	pg_host := os.Getenv("PGPOOL_HOST")
	database_name := os.Getenv("PGPOOL_DB")
	pg_port := os.Getenv("PGPOOL_PORT")

	if pg_user == "" || pg_password == "" || pg_host == "" || database_name == "" || pg_port == "" {
		log.Fatal("One or more required secrets are missing")
	}

	fmt.Printf("DB connect status: %s:%s@tcp(%s:%s)/%s\n", pg_user, pg_password, pg_host, pg_port, database_name)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh", pg_host, pg_user, pg_password, database_name, pg_port)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected")

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	return DB
}
