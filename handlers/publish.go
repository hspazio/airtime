package handlers

import (
	"log"
	"net/http"

	"github.com/hspazio/hermes-lite/hub"

	"github.com/gorilla/mux"
)

// TODO: move this out
var h hub.Hub

func init() {
	h = hub.NewHub()
}

// Publish a message to a feed
func Publish(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "could not establish connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	feed, _ := mux.Vars(r)["name"]

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			http.Error(w, "could not read message", http.StatusInternalServerError)
			break
		}
		log.Printf("received: %s", message)
		h.Publish(feed, message)
	}
}
