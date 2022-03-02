package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	// Config map for database
	DB_HOST = "database"
	DB_PORT = "5432"
	DB_NAME = "mygram"
	DB_USER = "postgres"
	DB_PASS = "pass"

	// Config map for server
	SRV_PORT = "3030"

	// Config for JWT
	JWT_SECRET_KEY = "bq8G]0P(_O2,zm{0A#$@ZPFMV^v<8wbR"
)

// LoadConfig loads config directly from .env file
func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("POSTGRES_HOST")
	if dbHost != "" {
		DB_HOST = dbHost
	}
	dbUser := os.Getenv("POSTGRES_USER")
	if dbUser != "" {
		DB_USER = dbUser
	}
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	if dbPass != "" {
		DB_PASS = dbPass
	}
	dbName := os.Getenv("POSTGRES_DB")
	if dbName != "" {
		DB_NAME = dbName
	}

	srvPort := os.Getenv("SERVER_PORT")
	if srvPort != "" {
		SRV_PORT = srvPort
	}
}
