# Airtime

Airtime is a PubSub Go library that uses WebSockets. It includes the server implementation as well as the client.

Once the server is running, the client can connect to a feed by using the `Connect("feedname")` function. This returns a connection object that allows the client to:

1. subscribe to any messages sent through the feed
2. publish meessages to the feed
3. receive any errors from the feed
4. close the connection

## Installation

```bash
go get github.com/hspazio/airtime
```

## Server

```go
package main 

import (
  "log"

  "github.com/hspazio/airtime"
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

	"github.com/hspazio/airtime"
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

## TODO

This is library is still work-in-progress and it's still missing a lot of basic functionalities: 

* authentication
* make client more resilient to network failures
* persist and replay undelivered messages when client reconnects
* read/write permissions management
* more...
