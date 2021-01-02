// cmd/app/main.go

package main

import (
	"fmt"
	"github.com/amryamanah/go-boilerplate/pkg/application"
	"github.com/amryamanah/go-boilerplate/pkg/config"
	"github.com/amryamanah/go-boilerplate/pkg/exithandler"
	"github.com/amryamanah/go-boilerplate/pkg/logger"
	"log"
)

func main() {
	err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	fmt.Printf("[VIPER] Config: %+v\n", config.Config)
	app := application.NewApplication()

	app.InitStore()

	go func() {
		if err := app.Start(); err != nil {
			logger.Error.Fatal(err.Error())
		}
	}()

	exithandler.Init(func() {
		if err := app.Close(); err != nil {
			logger.Error.Println(err.Error())
		}
	})
}
