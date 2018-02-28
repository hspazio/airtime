package handlers

import (
	"fmt"
	"net/http"
)

// Subscribe to messages coming from a specific feed
func Subscribe(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "subscribing to messages")
}
