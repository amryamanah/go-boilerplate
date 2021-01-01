// cmd/app/main.go

package main

import (
	"github.com/amryamanah/go-boilerplate/internal/app/router"
	"github.com/amryamanah/go-boilerplate/pkg/application"
	"github.com/amryamanah/go-boilerplate/pkg/exithandler"
	"github.com/amryamanah/go-boilerplate/pkg/logger"
	"github.com/amryamanah/go-boilerplate/pkg/server"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("failed to load env vars")
	}

	app, err := application.Get(application.DEV)
	if err != nil {
		log.Fatal(err.Error())
	}

	srv := server.
			Get().
			WithAddr(app.Cfg.GetApiPort()).
			WithRouter(router.Get()).
			WithErrLogger(logger.Error)

	go func() {
		logger.Info.Printf("starting server at %s", app.Cfg.GetApiPort())
		if err := srv.Start(); err != nil {
			logger.Error.Fatal(err.Error())
		}
	}()

	exithandler.Init(func() {
		if err := srv.Close(); err != nil {
			logger.Error.Println(err.Error())
		}
		app.DBConn.Close()
	})
}
