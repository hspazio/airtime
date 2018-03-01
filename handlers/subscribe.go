package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Subscribe to messages coming from a specific feed
func Subscribe(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "could not establish connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	feed, _ := mux.Vars(r)["name"]

	log.Printf("subscription created for feed '%s'", feed)
	s := h.Subscribe(feed)
	for message := range s.Inbox {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			http.Error(w, "could forward message", http.StatusInternalServerError)
			break
		}
	}
}
