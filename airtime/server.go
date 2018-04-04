package airtime

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Server for Airtime
type Server struct {
	port     string
	hub      Hub
	upgrader websocket.Upgrader
}

// NewServer creates a server
func NewServer(port string) *Server {
	return &Server{
		port:     port,
		hub:      NewHub(),
		upgrader: websocket.Upgrader{},
	}
}

// Run the server
func (s *Server) Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/feeds/{name}/publish", s.publish).Methods("GET")
	router.HandleFunc("/feeds/{name}/subscribe", s.subscribe).Methods("GET")

	addr := fmt.Sprintf(":%s", s.port)
	return http.ListenAndServe(addr, router)
}

func (s *Server) publish(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
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
		s.hub.Publish(feed, message)
	}
}

func (s *Server) subscribe(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "could not establish connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	feed, _ := mux.Vars(r)["name"]

	log.Printf("subscription created for feed '%s'", feed)
	sb := s.hub.Subscribe(feed)
	for message := range sb.Inbox {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			http.Error(w, "could forward message", http.StatusInternalServerError)
			break
		}
	}
}
