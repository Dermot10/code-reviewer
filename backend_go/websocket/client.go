package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 64 * 1024 // 64kb
)

// single websocket conn
type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte // buffered chan for outbound msgs
	UserID uint
}

// reads msgs from ws connection
// run in own goroutine
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	// heartbeat mechanism, periodically checks
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket error: %v", err)
			}
			break
		}

		log.Printf("received message from user %d: %s", c.UserID, message)

		c.Send <- message
	}
}

// sends msgs to ws connection
// own goroutine
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// hub closed the chan, shutdown mechanism
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// one msg = one ws frame (clean json event for parsing)
			// alternative to streaming batched msgs and flushing unneeded frames under high-throughput, less batching efficiency
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			// ping, keep alive conn
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
