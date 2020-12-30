package main

import (
	"github.com/amryamanah/go-boilerplate/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	godotenv.Load()
	cfg := config.Get()

	direction := cfg.GetMigration()

	if direction != "down" && direction != "up" {
		log.Fatal("-migrate accepts [up, down] values only")
	}

	m, err := migrate.New("file://db/migrations", cfg.GetDBConnStr())
	if err != nil {
		log.Fatal(err)
	}

	if direction == "up" {
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	}

	if direction == "down" {
		if err := m.Down(); err != nil {
			log.Fatal(err)
		}
	}
}
