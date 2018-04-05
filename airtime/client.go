package airtime

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

// Client struct
type Client struct {
	Host   string
	socket *websocket.Conn
}

// Connection handled back to the user
type Connection struct {
	Inbox  chan []byte
	Errors chan error
	Closed chan bool
}

// Connect to a feed in order to establish a connection
func (c *Client) Connect(feed string) (*Connection, error) {
	conn, _, err := websocket.DefaultDialer.Dial(c.subscribeURL(feed), nil)
	if err != nil {
		return nil, fmt.Errorf("could not subscribe to feed '%s': %v", feed, err)
	}
	c.socket = conn

	connection := &Connection{
		Inbox:  make(chan []byte),
		Errors: make(chan error),
		Closed: make(chan bool, 1),
	}

	go c.loop(connection)

	return connection, nil
}

// Close the connection currently opened for the client
func (c *Client) Close() error {
	return c.socket.Close()
}

func (c *Client) loop(conn *Connection) {
	for {
		_, message, err := c.socket.ReadMessage()

		if err != nil {
			conn.Errors <- err
			continue
		}

		conn.Inbox <- message
	}
}

// conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
// return nil

func (c *Client) subscribeURL(feed string) string {
	path := fmt.Sprintf("/feeds/%s/subscribe", feed)
	u := url.URL{Scheme: "ws", Host: c.Host, Path: path}
	return u.String()
}
