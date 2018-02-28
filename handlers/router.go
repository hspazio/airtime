package handlers

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// NewRouter generates a router with all the mapped handlers
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", Home).Methods("GET")
	router.HandleFunc("/feeds/{name}/publish", Publish).Methods("GET")
	router.HandleFunc("/feeds/{name}/subscribe", Subscribe).Methods("GET")
	return router
}
