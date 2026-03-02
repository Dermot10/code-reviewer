package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, reading from system environment")
	}
	header := http.Header{}
	token := strings.TrimSpace(os.Getenv("JWT_TOKEN"))
	if token == "" {
		log.Fatal("JWT_TOKEN not set")
	}

	header.Add("Authorization", "Bearer "+token)

	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", header)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	log.Println("connected to websocket")

	prompt := `{
		"type": "assistant.prompt",
		"payload": {
			"conversation_id": 1,
			"prompt": "Say hello in Go"
		}
	}`

	err = c.WriteMessage(websocket.TextMessage, []byte(prompt))
	if err != nil {
		log.Fatal(err)
	}

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("received:", string(msg))

		if strings.Contains(string(msg), `"done":true`) {
			log.Println("stream completed")
			break
		}
	}
}
