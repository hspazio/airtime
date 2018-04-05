package airtime

import (
	"log"
	"sync"
	"time"
)

// Hub is the pub-sub server that manages the subscriptions
type Hub struct {
	mx            *sync.RWMutex
	subscriptions map[string]map[*Subscription]struct{}
	broadcast     chan envelope
}

// Subscription contains an inbox where the subscriber can pull messages from
type Subscription struct {
	Inbox chan []byte
}

type envelope struct {
	feed    string
	message []byte
}

// NewHub creates a new Hub and spawns it in the background listening for messages to broadcast
func NewHub() Hub {
	h := Hub{
		mx:            &sync.RWMutex{},
		subscriptions: make(map[string]map[*Subscription]struct{}),
		broadcast:     make(chan envelope),
	}

	go func() {
		log.Println("hub started")
		for envelope := range h.broadcast {
			h.mx.RLock()
			for sub := range h.subscriptions[envelope.feed] {
				select {
				case sub.Inbox <- envelope.message:
				case <-time.After(1 * time.Second):
					log.Printf("closing connection for unresponsive endpoint")
					h.Unsubscribe(envelope.feed, sub)
				}
			}
			h.mx.RUnlock()
		}
	}()

	return h
}

// Subscribe creates a subscription to a feed using the provided connection
func (h Hub) Subscribe(feed string) *Subscription {
	h.mx.Lock()
	defer h.mx.Unlock()
	sub := &Subscription{make(chan []byte)}
	if _, ok := h.subscriptions[feed]; !ok {
		h.subscriptions[feed] = make(map[*Subscription]struct{})
	}
	h.subscriptions[feed][sub] = struct{}{}

	return sub
}

// Unsubscribe removes a subscription from the feed
func (h Hub) Unsubscribe(feed string, sub *Subscription) {
	h.mx.Lock()
	defer h.mx.Unlock()
	if _, ok := h.subscriptions[feed]; ok {
		subs := h.subscriptions[feed]
		if _, ok := subs[sub]; ok {
			delete(subs, sub)
			close(sub.Inbox)
		}
	}
}

// Publish broadcast a message to all the subscribers of the feed
func (h Hub) Publish(feed string, message []byte) {
	h.broadcast <- envelope{feed, message}
}
