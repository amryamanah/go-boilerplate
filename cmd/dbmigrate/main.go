package main

import (
	"fmt"
	"github.com/amryamanah/go-boilerplate/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func main() {
	err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	fmt.Printf("[VIPER] Config: %+v\n", config.Config)

	direction := config.Config.Migrate

	if direction != "down" && direction != "up" {
		log.Fatal("-migrate accepts [up, down] values only")
	}

	log.Printf("[MIGRATE] DBCONNSTRING: %v\n", config.Config.GetDBConnStr())
	m, err := migrate.New("file://internal/store/migration", config.Config.GetDBConnStr())
	log.Printf("[MIGRATE] RUNNING MIGRATION WITH: %+v\n", m)
	if err != nil {
		log.Fatal(err)
	}

	if direction == "up" {
		if err := m.Up(); err != nil {
			if err.Error() == "no change" {
				return
			} else {
				log.Fatal(err)
			}
		}
	}

	if direction == "down" {
		if err := m.Down(); err != nil {
			if err.Error() == "no change" {
				return
			} else {
				log.Fatal(err)
			}
		}
	}
}
