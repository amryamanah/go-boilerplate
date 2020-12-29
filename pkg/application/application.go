package application

import (
	"github.com/amryamanah/go-boilerplate/pkg/config"
	"github.com/amryamanah/go-boilerplate/pkg/db"
	"log"
)

type Application struct {
	DB *db.DB
	Cfg *config.Config
}

func Get() (*Application, error) {
	log.Println("Restart")
	cfg := config.Get()
	db, err := db.Get(cfg.GetDBConnStr())

	if err != nil {
		return nil, err
	}

	return &Application{
		DB: db,
		Cfg: cfg,
	}, nil
}