package handlers

import (
	"log/slog"
	"net/http"

	"github.com/dermot10/code-reviewer/backend_go/middleware"
	"github.com/dermot10/code-reviewer/backend_go/websocket"
	gorilla_ws "github.com/gorilla/websocket"
)

var upgrader = gorilla_ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	logger *slog.Logger
	hub    *websocket.Hub
}

func NewWSHandler(logger *slog.Logger, hub *websocket.Hub) *WSHandler {
	return &WSHandler{
		logger: logger,
		hub:    hub,
	}
}

func (h *WSHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		h.logger.Error("userID missing or invalid in context")
		http.Error(w, "unathorized", http.StatusUnauthorized)
		return
	}

	//upgrade http to ws
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("websocket upgrade failed", "error", err)
		return
	}

	client := &websocket.Client{
		Hub:    h.hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		UserID: userID,
	}

	// go func() {
	// 	for {
	// 		_, msg, err := client.Conn.ReadMessage()
	// 		if err != nil {
	// 			log.Println("read error:", err)
	// 			return
	// 		}
	// 		log.Println("received from client:", string(msg))
	// 		client.Conn.WriteMessage(gorilla_ws.TextMessage, msg)
	// 	}
	// }()

	h.hub.Register <- client

	h.logger.Info("websocket connected", "user_id", userID)

	go client.WritePump()
	go client.ReadPump()
}
