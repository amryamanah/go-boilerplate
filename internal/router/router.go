// cmd/api/router/router.go

package router

import (
	"github.com/amryamanah/go-boilerplate/internal/handlers/ping_controller"
	"github.com/gorilla/mux"
	"net/http"
)

func Get() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ping", ping_controller.PingController.Ping).Methods(http.MethodGet)
	return router
}

