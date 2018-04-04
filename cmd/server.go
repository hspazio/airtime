package main

import (
	"log"

	"github.com/hspazio/hermes-go/airtime"
)

func main() {
	port := "8080"
	log.Printf("serving port %s", port)
	svr := airtime.NewServer(port)
	log.Fatal(svr.Run())
}
