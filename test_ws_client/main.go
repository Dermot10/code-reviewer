package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	header := http.Header{}
	header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImRlcm1vdC10ZXN0QGV4YW1wbGUuY29tIiwiZXhwIjoxNzcwNjUxMjAwLCJ1c2VyX2lkIjoxfQ.ufq-GokUAxmzPBIDz-Xahsy0pprZTfvCPgLrUc8J4YE") // <-- add a valid JWT here

	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", header)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Send a test message
	c.WriteMessage(websocket.TextMessage, []byte(`{"type":"ping"}`))

	// Read response
	_, msg, err := c.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("received:", string(msg))
}
