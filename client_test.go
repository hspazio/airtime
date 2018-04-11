package airtime

import (
	"testing"

	"github.com/hspazio/airtime"
)

func TestClient(t *testing.T) {
	t.Run("should send and receive from a feed", func(t *testing.T) {
		client := airtime.Client{Host: "localhost:8080"}

		conn, err := client.Connect("test")
		if err != nil {
			t.Fatalf("could not connect client: %v", err)
		}

		go func() {
			for {
				select {
				case msg := <-conn.Inbox:
					t.Logf("received msg: %s", msg)
					client.Close()
					return
				case err := <-conn.Errors:
					t.Errorf("unexpected error received via connection struct: %v", err)
				}
			}
		}()
		conn.Publish <- []byte("ping")
		<-conn.Closed
	})
}
