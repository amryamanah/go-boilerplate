// cmd/api/router/router.go

package router

import (
	"github.com/amryamanah/go-boilerplate/internal/api/controllers"
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

