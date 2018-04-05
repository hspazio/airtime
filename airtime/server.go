package airtime

import (
	"fmt"
	"log"
	"net/http"
	"sync"

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
	router.HandleFunc("/feeds/{name}", s.connect).Methods("GET")

	addr := fmt.Sprintf(":%s", s.port)
	return http.ListenAndServe(addr, router)
}

func (s *Server) connect(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "could not establish connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	feed, _ := mux.Vars(r)["name"]
	log.Printf("connection established for %s on feed '%s'", conn.RemoteAddr().String(), feed)

	var wg sync.WaitGroup
	wg.Add(2)
	go s.publishLoop(feed, conn, &wg)
	go s.subscribeLoop(feed, conn, &wg)
	wg.Wait()
}

func (s *Server) publishLoop(feed string, conn *websocket.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			switch err.(type) {
			case *websocket.CloseError:
				log.Printf("from: %s, to: %s, connection closed", conn.RemoteAddr().String(), feed)
			default:
				log.Printf("error in publishLoop: %v (endpoint: %s, feed: %s)", err, conn.RemoteAddr().String(), feed)
			}
			break
		}

		log.Printf("from: %s, to: %s, message: %s", conn.RemoteAddr().String(), feed, message)
		s.hub.Publish(feed, message)
	}
}

func (s *Server) subscribeLoop(feed string, conn *websocket.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	sb := s.hub.Subscribe(feed)

	for message := range sb.Inbox {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("error in subscribeLoop: %v", err)
			s.hub.Unsubscribe(feed, sb)
			break
		}
	}
}
