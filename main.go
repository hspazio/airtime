package main

import (
	"log"
	"net/http"

	"github.com/hspazio/hermes-lite/handlers"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {
	router := handlers.NewRouter()
	log.Printf("serving port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
