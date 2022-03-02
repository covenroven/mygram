package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/covenroven/mygram/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("No argument found")
		return
	}

	config.LoadConfig()

	m, err := func() (*migrate.Migrate, error) {
		maxRetry := 10
		var err error
		for r := 0; r <= maxRetry; r++ {
			m, err := migrate.New(
				"file://cmd/migrate/migrations",
				fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.DB_USER, config.DB_PASS, config.DB_HOST, config.DB_PORT, config.DB_NAME),
			)
			if err != nil {
				log.Println(err)
				log.Println("Cannot connect to database. Waiting for 2 seconds to retry...")
				time.Sleep(2 * time.Second)
			} else {
				return m, nil
			}
		}
		return nil, err
	}()
	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[1] {
	case "down":
		if err := m.Down(); err != nil {
			log.Fatal(err)
		}
	case "up":
		fallthrough
	default:
		if err := m.Up(); err != nil {
			if err != migrate.ErrNoChange {
				log.Fatal(err)
			}
		}
		log.Println("Successfully applied all migrations")
	}
}
