package services

import (
	"encoding/json"
	"log"
	"log/slog"

	"github.com/dermot10/code-reviewer/backend_go/dto"
	"github.com/dermot10/code-reviewer/backend_go/models"
	"github.com/dermot10/code-reviewer/backend_go/redis"
	"github.com/dermot10/code-reviewer/backend_go/websocket"
	"gorm.io/gorm"
)

type WsService struct {
	db     *gorm.DB
	redis  *redis.RedisClient
	logger *slog.Logger
	Hub    *websocket.Hub
}

func NewWsService(db *gorm.DB, redis *redis.RedisClient, logger *slog.Logger, wsHub *websocket.Hub) *WsService {
	return &WsService{db: db, redis: redis, logger: logger, Hub: wsHub}
}

func (s *WsService) HandleEvent(userID uint, event dto.WSEvent) {
	switch event.Type {

	case dto.EventFileUpload:
		s.handleFileUpload(userID, event)

	case dto.EventChatMessage:
		s.handleChatMessage(userID, event)

	default:
		log.Println("unknown event type:", event.Type)
	}
}

// Client sends update
// → Service validates & persists
// → Server constructs authoritative event
// → Hub broadcasts to user's active sessions

func (s *WsService) handleFileUpload(userID uint, event dto.WSEvent) {
	if userID == 0 {
		return
	}

	var payload dto.FileUpdatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		s.logger.Error("invalid file upload payload", "error", err)
		return
	}

	result := s.db.Model(&models.File{}).
		Where("id = ? AND user_id = ?", payload.FileID, userID).
		Updates(map[string]interface{}{
			"content": payload.Content,
		})

	if result.RowsAffected == 0 {
		s.logger.Warn("file update unauthorized or not found",
			"user_id", userID,
			"file_id", payload.FileID,
		)
		return
	}

	payloadData, err := json.Marshal(dto.FileUpdatedPayload{
		FileID:  payload.FileID,
		Content: payload.Content,
	})
	if err != nil {
		return
	}

	response := dto.WSEvent{
		Type:    dto.EventFileUpdated,
		Payload: payloadData,
	}

	data, _ := json.Marshal(response)

	s.Hub.Broadcast(websocket.Message{
		UserID: userID,
		Data:   data,
	})
}

func (s *WsService) handleChatMessage(userID uint, event dto.WSEvent) {
	if userID == 0 {
		return
	}

	var payload dto.ChatMessagePayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		s.logger.Error("invalid chat message payload", "error", err)
		return
	}

	// result := s.db.Model(&models.)
}
