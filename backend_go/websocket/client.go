package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/dermot10/code-reviewer/backend_go/dto"
	gorilla_ws "github.com/gorilla/websocket"
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
	Conn   *gorilla_ws.Conn
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
			if gorilla_ws.IsUnexpectedCloseError(err, gorilla_ws.CloseGoingAway, gorilla_ws.CloseAbnormalClosure) {
				log.Printf("websocket error: %v", err)
			}
			break
		}

		log.Printf("received message from user %d: %s", c.UserID, message)

		var event dto.WSEvent
		if err := json.Unmarshal(message, &event); err != nil {
			log.Println("invalid ws payload:", err)
			continue
		}

		// return event
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
				c.Conn.WriteMessage(gorilla_ws.CloseMessage, []byte{})
				return
			}

			// one msg = one ws frame (clean json event for parsing)
			// alternative to streaming batched msgs and flushing unneeded frames under high-throughput, less batching efficiency
			if err := c.Conn.WriteMessage(gorilla_ws.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			// ping, keep alive conn
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(gorilla_ws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
