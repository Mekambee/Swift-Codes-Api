package main

import (
	"log"
	"os"

	"github.com/Mekambee/Swift-Codes-Api/internal/api"
	"github.com/Mekambee/Swift-Codes-Api/internal/database"
)

func main() {

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}

	err := database.ConnectAndMigrate(dbHost, dbUser, dbPass, dbName, dbPort)
	if err != nil {
		log.Fatalf("Cannot connect to DB: %v\n", err)
	}
	defer database.DB.Close()

	router := api.SetupRouter()
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v\n", err)
	}
}
