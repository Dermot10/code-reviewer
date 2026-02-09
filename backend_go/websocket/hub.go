package websocket

import (
	"context"
	"log"
	"sync"
)

// The hub maintains active clients and broadcast msgs
type Hub struct {
	// registered clients mapped to userID
	clients map[uint]map[*Client]bool

	broadcast chan Message

	// register and unregister requests from clients
	Register chan *Client

	Unregister chan *Client

	mu sync.RWMutex
}

// Internal routing wrapper
type Message struct {
	UserID uint
	Data   []byte // pre-marshaled JSON from dto
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]map[*Client]bool),
		broadcast:  make(chan Message, 256),
		Register:   make(chan *Client),
		Unregister: make(chan *Client, 64),
	}
}

// starts hub from goroutine
func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("hub shutting down")
			return
		case client := <-h.Register:
			h.mu.Lock()
			if h.clients[client.UserID] == nil {
				h.clients[client.UserID] = make(map[*Client]bool)
			}
			h.clients[client.UserID][client] = true
			h.mu.Unlock()

			log.Printf("client registered: user_id=%d, total=%d", client.UserID, len(h.clients[client.UserID]))

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.UserID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)

					if len(clients) == 0 {
						delete(h.clients, client.UserID)
					}
				}
			}
			h.mu.Unlock()

			log.Printf("client unregistered: user_id=%d", client.UserID)

		case message := <-h.broadcast:
			h.mu.RLock()
			clients := h.clients[message.UserID]
			h.mu.RUnlock()

			for client := range clients {
				select {
				case client.Send <- message.Data:
				default:
					h.Unregister <- client
				}
			}
		}
	}
}

func (h *Hub) Broadcast(msg Message) {
	h.broadcast <- msg
}
