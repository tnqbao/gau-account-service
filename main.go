// build
package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tnqbao/gau-account-service/config"
	"github.com/tnqbao/gau-account-service/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cf := config.LoadEnvConfig()
	db := config.InitDB()

	router := routes.SetupRouter(db, cf)
	router.Run(":8080")
}
