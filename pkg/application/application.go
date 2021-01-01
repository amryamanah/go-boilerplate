package application

import (
	"errors"
	"fmt"
	"github.com/amryamanah/go-boilerplate/pkg/config"
	"github.com/amryamanah/go-boilerplate/pkg/db"
	"log"
)

type AppEnv string

const (
	DEV AppEnv = "dev"
	TEST = "test"
	PROD = "prod"
)

type Application struct {
	DBConn *db.DB
	Cfg    *config.Config
}

func Get(kind AppEnv) (*Application, error) {
	log.Println("Restart")
	cfg := config.Get()
	var dbConnStr string
	switch kind {
	case DEV:
		dbConnStr = cfg.GetDBConnStr()
	case TEST:
		dbConnStr = cfg.GetTestDBConnStr()
	default:
		return nil, errors.New(fmt.Sprintf("unsupported application environment: %s", kind))
	}
	dbConn, err := db.Get(dbConnStr)

	if err != nil {
		return nil, err
	}

	return &Application{
		DBConn: dbConn,
		Cfg:    cfg,
	}, nil
}