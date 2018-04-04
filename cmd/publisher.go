package main

import (
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func main() {

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/feeds/test/publish"}
	log.Printf("connecting to %s", u.String())

	conn, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("publisher could not connect to server: %s", err)
	}
	log.Printf("response code %d", resp.StatusCode)
	defer conn.Close()

	conn.WriteMessage(websocket.TextMessage, []byte("msg 1"))
	time.Sleep(time.Second)
	conn.WriteMessage(websocket.TextMessage, []byte("msg 2"))
	time.Sleep(time.Second)
	conn.WriteMessage(websocket.TextMessage, []byte("msg 3"))
	time.Sleep(time.Second)
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}
