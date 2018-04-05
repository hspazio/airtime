package main

import (
	"log"

	"github.com/hspazio/hermes-go/airtime"
)

func main() {
	client := &airtime.Client{Host: "localhost:8080"}

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
		}
	}
}
