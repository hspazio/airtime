package main

import (
	"log"

	"github.com/hspazio/hermes-go/airtime"
)

func main() {
	client := &airtime.Client{Host: "localhost:8080"}

	// TODO: use the following design instead
	// conn, err := client.Connect("test")
	// if err != nil {
	//     log.Fatal(err)
	// }
	// for {
	//     select {
	//     case msg := <-conn.Inbox:
	//         log.Printf("received message: %s", m)
	//     case err := <-conn.Errors:
	//         log.Printf("error: %v", err)
	//     case <-conn.Closed:
	//         log.Println("connection closed")
	//     }
	// }

	err := client.Subscribe("test", func(m []byte) {
		log.Printf("received message: %s", m)
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connection closed")
}
