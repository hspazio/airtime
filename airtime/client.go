package airtime

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// Client struct
type Client struct {
	Host    string
	socket  *websocket.Conn
	closing chan bool
}

// Connection handled back to the user
type Connection struct {
	Inbox   chan []byte
	Publish chan []byte
	Errors  chan error
	Closed  chan bool
}

// Connect to a feed in order to establish a connection
func (c *Client) Connect(feed string) (*Connection, error) {
	conn, _, err := websocket.DefaultDialer.Dial(c.subscribeURL(feed), nil)
	if err != nil {
		return nil, fmt.Errorf("could not subscribe to feed '%s': %v", feed, err)
	}
	c.socket = conn
	c.closing = make(chan bool, 1)

	connection := &Connection{
		Inbox:   make(chan []byte),
		Publish: make(chan []byte),
		Errors:  make(chan error),
		Closed:  make(chan bool, 1),
	}

	go c.awaitClosure(connection)
	go c.readLoop(connection)
	go c.writeLoop(connection)

	return connection, nil
}

// Close the connection currently opened for the client
func (c *Client) Close() {
	c.closing <- true
}

func (c *Client) readLoop(conn *Connection) {
	for {
		_, message, err := c.socket.ReadMessage()

		if err != nil {
			conn.Errors <- err
			continue
		}

		conn.Inbox <- message
	}
}

func (c *Client) writeLoop(conn *Connection) {
	for {
		message := <-conn.Publish
		log.Printf("%s\n", message)
		c.socket.WriteMessage(websocket.TextMessage, message)
	}
}

func (c *Client) awaitClosure(conn *Connection) {
	signal := <-c.closing
	err := c.socket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		conn.Errors <- err
	}
	if err := c.socket.Close(); err != nil {
		conn.Errors <- err
	}
	conn.Closed <- signal
}

func (c *Client) subscribeURL(feed string) string {
	path := fmt.Sprintf("/feeds/%s", feed)
	u := url.URL{Scheme: "ws", Host: c.Host, Path: path}
	return u.String()
}
