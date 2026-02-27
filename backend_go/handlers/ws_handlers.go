package handlers

import (
	"encoding/json"
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

	h.hub.Register <- client

	h.logger.Info("websocket connected", "user_id", userID)

	go client.WritePump()
	go client.ReadPump()
}

func (h *WSHandler) routeEvent(userID uint, raw []byte) {
	var event dto.WSEvent
	if err := json.Unmarshal(raw, &event); err != nil {
		h.logger.Error("invalid ws event", "event", err)
		return
	}

	handlers := map[dto.WSEventType]func(uint, dto.WSEvent){
		dto.EventFileUpload:          h.FileUpload,
		dto.EventFileUpdated:         h.FileUpdate,
		dto.EventMessageSend:         h.MessageSend,
		dto.EventConversationCreate:  h.ConversationCreate,
		dto.EventConversationArchive: h.ConversationArchive,
		dto.EventConversationRename:  h.ConversationRename,
		dto.EventConvrsationDelete:   h.ConversationDelete,
	}

	if handler, ok := handlers[event.Type]; ok {
		handler(userID, event)
	} else {
		h.logger.Warn("unknown event type", "type", event.Type)
	}
}

func (h *WSHandler) FileUpload(userID uint, event dto.WSEvent) {
	var payload dto.FileUpdatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		h.logger.Error("invalid file upload payload", "error", err)
		return
	}

	file, err := h.fileService.UpdateFile(userID, payload.FileID, payload.Content)
	if err != nil {
		return
	}

	responsePayload, err := json.Marshal(dto.FileUpdatedPayload{
		FileID:  file.ID,
		Content: file.Content,
	})
	if err != nil {
		return
	}

	response := dto.WSEvent{
		Type:    dto.EventFileUpdated,
		Payload: responsePayload,
	}

	data, err := json.Marshal(response)
	if err != nil {
		return
	}

	h.hub.Broadcast(websocket.Message{
		UserID: userID,
		Data:   data,
	})
}

func (h *WSHandler) FileUpdate(userID uint, event dto.WSEvent) {
	var payload dto.FileUpdatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		h.logger.Error("invalid file update payload", "error", err)
		return
	}

	updatedFile, err := h.fileService.UpdateFile(userID, payload.FileID, payload.Content)
	if err != nil {
		h.logger.Error("failed to update file", "error", err)
		return
	}

	responsePayload, err := json.Marshal(dto.FileUpdatedPayload{
		FileID:  updatedFile.ID,
		Content: updatedFile.Content,
	})
	if err != nil {
		return
	}

	response := dto.WSEvent{
		Type:    dto.EventFileUpdated,
		Payload: responsePayload,
	}

	data, err := json.Marshal(response)
	if err != nil {
		return
	}

	h.hub.Broadcast(websocket.Message{
		UserID: userID,
		Data:   data,
	})
}

func (h *WSHandler) MessageSend(userID uint, event dto.WSEvent) {
	var payload dto.MessageSendPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		h.logger.Error("invalid chat message payload", "error", err)
		return
	}

	msg, err := h.chatService.CreateMessage(userID, payload.ConversationID, "user", payload.Content)
	if err != nil {
		h.logger.Error("failed to create chat message", "error", err)
		return
	}

	responsePayload, err := json.Marshal(dto.MessageCreatedPayload{
		ID:             msg.ID,
		ConversationID: msg.ConversationID,
		Role:           msg.Role,
		Content:        msg.Content,
		CreatedAt:      msg.CreatedAt,
	})
	if err != nil {
		return
	}

	response := dto.WSEvent{
		Type:    dto.EventMessageCreated,
		Payload: responsePayload,
	}

	data, err := json.Marshal(response)
	if err != nil {
		return
	}

	h.hub.Broadcast(websocket.Message{
		UserID: userID,
		Data:   data,
	})
}

func (h *WSHandler) ConversationCreate(userID uint, event dto.WSEvent) {
	var payload dto.ConversationCreatePayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		h.logger.Error("invalid conversartion create payload", "error", err)
		return
	}

	conv, err := h.chatService.CreateConversation(userID, payload.Title)
	if err != nil {
		h.logger.Error("failed to create conversation", "error", err)
		return
	}
	responsePayload, err := json.Marshal(dto.ConversationCreatedPayload{
		ConversationID: conv.ID,
		Title:          conv.Title,
	})
	if err != nil {
		return
	}

	response := dto.WSEvent{
		Type:    dto.EventConversationCreated,
		Payload: responsePayload,
	}

	data, err := json.Marshal(response)
	if err != nil {
		return
	}

	h.hub.Broadcast(websocket.Message{
		UserID: userID,
		Data:   data,
	})

}

func (h *WSHandler) ConversationArchive(userID uint, event dto.WSEvent) {
	var payload dto.ConversationArhivePayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		h.logger.Error("invalid conversation update payload", "error", err)
		return
	}

	err := h.chatService.ArchiveConversation(userID, payload.ConversationID)
	if err != nil {
		h.logger.Error("failed to update conversation", "error", err)
		return
	}

	responsePayload, err := json.Marshal(dto.ConversationArchivedPayload{
		ConversationID: payload.ConversationID,
	})
	if err != nil {
		return
	}

	response := dto.WSEvent{
		Type:    dto.EventConversationArchived,
		Payload: responsePayload,
	}

	data, err := json.Marshal(response)
	if err != nil {
		return
	}

	h.hub.Broadcast(websocket.Message{
		UserID: userID,
		Data:   data,
	})
}

func (h *WSHandler) ConversationRename(userID uint, event dto.WSEvent) {
	var payload dto.ConversationRenamePayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		h.logger.Error("invalid conversation rename payload", "error", err)
		return
	}

	err := h.chatService.RenameConversation(userID, payload.ConversationID, payload.Title)
	if err != nil {
		h.logger.Error("failed to rename conversation", "error", err)
		return
	}

	responsePayload, err := json.Marshal(dto.ConversationRenamedPayload{
		ConversationID: payload.ConversationID,
		Title:          payload.Title,
	})
	if err != nil {
		return
	}

	response := dto.WSEvent{
		Type:    dto.EventConversationRenamed,
		Payload: responsePayload,
	}

	data, err := json.Marshal(response)
	if err != nil {
		return
	}

	h.hub.Broadcast(websocket.Message{
		UserID: userID,
		Data:   data,
	})
}

func (h *WSHandler) ConversationDelete(userID uint, event dto.WSEvent) {
	var payload dto.ConversationDeletePayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		h.logger.Error("invalid conversation delete payload", "error", err)
		return
	}

	err := h.chatService.DeleteConversation(userID, payload.ConversationID)
	if err != nil {
		h.logger.Error("failed to delete conversation", "error", err)
		return
	}
}
