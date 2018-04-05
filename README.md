# Airtime

PubSub client and server using WebSockets in Go

## Server

```go
package main 

import (
  "log"

  "github.com/hspazio/airtime/airtime"
)

func main() {
	svr := airtime.NewServer("8080")
	log.Fatal(svr.Run())
}
```

## Client

```go
package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/hspazio/airtime/airtime"
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
    // Read messages published to the feed
		case msg := <-conn.Inbox:
			log.Printf("received msg: %s", msg)

    // Publish messages to the feed using the Publish channel
		case <-time.After(1 * time.Second):
			conn.Publish <- []byte("ping")

    // Errors are pushed to the client through a dedicated channel
		case err := <-conn.Errors:
			log.Printf("error: %v", err)

    // Close a connection using the Close() method
		case <-interrupt:
			client.Close()

    // When a connection is fully closed the client is notified 
		case <-conn.Closed:
			log.Println("connection closed")
			return
		}
	}
}
```
