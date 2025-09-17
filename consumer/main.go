package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tnqbao/gau-account-service/shared/config"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, continuing with environment variables")
	}

	cfg := config.NewConfig()

	// TODO: Initialize message queue consumers here
	// This is a placeholder consumer service
	log.Printf("Consumer service starting with config: %+v", cfg.EnvConfig.Environment)

	// Keep the service running
	select {}
}
