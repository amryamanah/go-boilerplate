package application

import (
	"fmt"
	store "github.com/amryamanah/go-boilerplate/internal/store/sqlc"
	"github.com/amryamanah/go-boilerplate/pkg/config"
	"github.com/amryamanah/go-boilerplate/pkg/logger"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Application struct {
	Store  store.Store
	router *gin.Engine
}

func NewApplication(store store.Store) *Application {
	app := &Application{Store: store}
	router := gin.Default()

	router.POST("/accounts", app.CreateAccount)
	router.GET("/accounts/:id", app.GetAccount)
	router.GET("/accounts", app.ListAccount)

	app.router = router
	return app
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (a *Application) Start() error {
	address := fmt.Sprintf("0.0.0.0:%s", config.Config.GetApiPort())
	logger.Info.Printf("starting server at %s", address)
	return a.router.Run(address)
}

func (a *Application) Close() error {
	return nil
}

