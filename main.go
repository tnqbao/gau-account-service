// build
package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tnqbao/gau_user_service/config"
	"github.com/tnqbao/gau_user_service/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	db := config.InitDB()
	router := routes.SetupRouter(db)
	router.Run(":8080")
}
