package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/hspazio/hermes-go/airtime"
)

func main() {
	client := &airtime.Client{Host: "localhost:8080"}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	conn, err := client.Connect("test")
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case msg := <-conn.Inbox:
			log.Printf("received msg: %s", msg)
		case err := <-conn.Errors:
			log.Printf("error: %v", err)
		case <-conn.Closed:
			log.Println("connection closed")
			return
		case <-time.After(1 * time.Second):
			conn.Publish <- []byte("ping")
		case <-interrupt:
			client.Close()
		}
	}
}
