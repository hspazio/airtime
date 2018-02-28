package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlers(t *testing.T) {
	router := NewRouter()

	t.Run("home", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatalf("cannot create request: %s", err)
		}
		router.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("expected code %d, got %d", http.StatusOK, w.Code)
		}
		expected := "this is the home"
		actual := w.Body.String()
		if expected != actual {
			t.Fatalf("expected body '%s', got '%s'", expected, actual)
		}
	})

	t.Run("publish feed successfully", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("POST", "/feeds/test", strings.NewReader("the-body"))
		if err != nil {
			t.Fatalf("cannot create request: %s", err)
		}
		router.ServeHTTP(w, r)

		if w.Code != http.StatusAccepted {
			t.Fatalf("expected code %d, got %d", http.StatusAccepted, w.Code)
		}
		expected := "publishing message to feed test: the-body"
		actual := w.Body.String()
		if expected != actual {
			t.Fatalf("expected body '%s', got '%s'", expected, actual)
		}
	})
}
