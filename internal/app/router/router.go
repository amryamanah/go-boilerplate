// cmd/app/router/router.go

package router

import (
	"github.com/amryamanah/go-boilerplate/internal/app/controllers"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func Get() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/ping", controllers.PingController.Ping).Methods(http.MethodGet)
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	return loggedRouter
}

