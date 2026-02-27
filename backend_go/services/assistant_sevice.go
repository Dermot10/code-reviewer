package services

import (
	"log/slog"

	"github.com/dermot10/code-reviewer/backend_go/dto"
	"github.com/dermot10/code-reviewer/backend_go/redis"
	"github.com/dermot10/code-reviewer/backend_go/websocket"
	"gorm.io/gorm"
)

type AssistantService struct {
	db     *gorm.DB
	redis  *redis.RedisClient
	logger *slog.Logger
	wsHub  *websocket.Hub
}

func NewAssistantService(redis *redis.RedisClient, logger *slog.Logger) *AssistantService {
	return &AssistantService{redis: redis, logger: logger}
}

func (s *AssistantService) SendPrompt(userID, conversationID uint, prompt string) error {
	// Push to AI worker via Redis queue

}

func (s *AssistantService) StreamResponse(userID, conversationID uint, chunk string) {
	event := dto.WSEvent{
		Type: dto.EventAssistantStream,
		Payload: marshal(dto.AssistantStreamPayload{
			ConversationID: conversationID,
			Chunk:          chunk,
		}),
	}
	data := marshal(event)
	hub.Broadcast(websocket.Message{
		UserID: userID,
		Data:   data,
	})
}
