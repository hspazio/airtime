package airtime

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

// Client struct
type Client struct {
	Host string
}

// Subscribe to a feed with a strategy function that handles the message
func (c *Client) Subscribe(feed string, strategy func(msg []byte)) error {
	conn, _, err := websocket.DefaultDialer.Dial(c.subscribeURL(feed), nil)
	if err != nil {
		return fmt.Errorf("could not subscribe to feed '%s': %v", feed, err)
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("could not read message: %v", err)
		}
		strategy(message)
	}
	// conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	// return nil
}

func (c *Client) subscribeURL(feed string) string {
	path := fmt.Sprintf("/feeds/%s/subscribe", feed)
	u := url.URL{Scheme: "ws", Host: c.Host, Path: path}
	return u.String()
}
