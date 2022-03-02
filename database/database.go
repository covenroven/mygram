package database

import (
	"fmt"
	"log"
	"time"

	"github.com/covenroven/mygram/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Open connection to database
func InitDB() (*sqlx.DB, error) {
	maxRetry := 10
	var err error
	for r := 0; r <= maxRetry; r++ {
		db, err := sqlx.Connect("postgres", DBConnectionString())
		if err != nil {
			log.Println("Cannot connect to database. Waiting for 2 seconds to retry...")
			time.Sleep(2 * time.Second)
		} else {
			return db, nil
		}
	}

	return nil, err
}

// DBConnectionString returns string for database connection
func DBConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DB_HOST, config.DB_PORT, config.DB_USER, config.DB_PASS, config.DB_NAME)
}
