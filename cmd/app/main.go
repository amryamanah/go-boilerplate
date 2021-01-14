// cmd/app/main.go

package main

import (
	"database/sql"
	"fmt"
	"log"

	store "github.com/amryamanah/go-boilerplate/internal/store/sqlc"
	"github.com/amryamanah/go-boilerplate/pkg/application"
	"github.com/amryamanah/go-boilerplate/pkg/client"
	"github.com/amryamanah/go-boilerplate/pkg/config"
	"github.com/amryamanah/go-boilerplate/pkg/exithandler"
	"github.com/amryamanah/go-boilerplate/pkg/logger"
	_ "github.com/lib/pq"
)

func main() {
	err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	fmt.Printf("[VIPER] Config: %+v\n", config.Config)

	client.InitRedis()

	conn, err := sql.Open("postgres", config.Config.GetDBConnStr())
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	storeInst := store.NewStore(conn)

	app := application.NewApplication(storeInst)

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
