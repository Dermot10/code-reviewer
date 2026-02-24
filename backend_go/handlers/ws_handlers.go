package handlers

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"github.com/dermot10/code-reviewer/backend_go/dto"
	"github.com/dermot10/code-reviewer/backend_go/middleware"
	"github.com/dermot10/code-reviewer/backend_go/services"
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
	logger      *slog.Logger
	hub         *websocket.Hub
	fileService *services.FileService
	chatService *services.ChatService
}

func NewWSHandler(logger *slog.Logger, hub *websocket.Hub, fileService *services.FileService, chatService *services.ChatService) *WSHandler {
	return &WSHandler{
		logger:      logger,
		hub:         hub,
		fileService: fileService,
		chatService: chatService,
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
		OnMessage: func(userID uint, msg []byte) {
			h.routeEvent(userID, msg)
		},
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

func (h *WSHandler) routeEvent(userID uint, raw []byte) {
	var event dto.WSEvent
	if err := json.Unmarshal(raw, &event); err != nil {
		return
	}

	switch event.Type {

	case dto.EventFileUpload:
		h.handleFileUpload(userID, event)

	case dto.EventChatMessage:
		h.handleChatMessage(userID, event)

	default:
		log.Println("unknown event type:", event.Type)
	}
}

func (h *WSHandler) handleFileUpload(userID uint, event dto.WSEvent) {
	var payload dto.FileUpdatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		h.logger.Error("invalid file upload payload", "error", err)
		return
	}

	file, err := h.fileService.UpdateFile(userID, payload.FileID, payload.Content)
	if err != nil {
		return
	}

	responsePayload, _ := json.Marshal(dto.FileUpdatedPayload{
		FileID:  file.ID,
		Content: file.Content,
	})

	response := dto.WSEvent{
		Type:    dto.EventFileUpdated,
		Payload: responsePayload,
	}

	data, _ := json.Marshal(response)

	h.hub.Broadcast(websocket.Message{
		UserID: userID,
		Data:   data,
	})
}

func (h *WSHandler) handleChatMessage(userID uint, event dto.WSEvent) {}
