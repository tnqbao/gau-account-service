package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/tnqbao/gau-account-service/shared/config"
)

func main() {
	_ = godotenv.Load()
	cfg := config.NewConfig()
	log.Printf("Consumer service starting with config: %+v", cfg.EnvConfig.Environment)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("Shutting down consumer service gracefully...")
}
