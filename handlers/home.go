package handlers

import (
	"fmt"
	"net/http"
)

// Home shows the default message
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is the home")
}
