package services

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/dermot10/code-reviewer/backend_go/dto"
	"github.com/dermot10/code-reviewer/backend_go/models"
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

func NewAssistantService(db *gorm.DB, redis *redis.RedisClient, logger *slog.Logger, wsHub *websocket.Hub) *AssistantService {
	return &AssistantService{db: db, redis: redis, logger: logger, wsHub: wsHub}
}

func (s *AssistantService) SendPrompt(userID uint, payload dto.PromptPayload) error {
	// persist user msg, push to queue for worker
	ctx := context.Background()

	msg := models.ChatMessage{
		ConversationID: payload.ConversationID,
		Role:           "user",
		Content:        payload.Prompt,
	}

	if err := s.db.Create(&msg).Error; err != nil {
		return err
	}

	task := dto.AssistantTask{
		UserID:         userID,
		ConversationID: payload.ConversationID,
		Prompt:         payload.Prompt,
	}

	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return s.redis.PushQueue(ctx, data)
}

func (s *AssistantService) StreamResponse(userID, conversationID uint, chunk string, done bool) {
	// stream outbound response back to client
	payload := dto.AssistantStreamPayload{
		ConversationID: conversationID,
		Chunk:          chunk,
		Done:           done,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		s.logger.Error("failed to marshal assistant stream payload", "error", err)
		return
	}

	event := dto.WSEvent{
		Type:    dto.EventAssistantStream,
		Payload: payloadJSON,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		s.logger.Error("failed to marshal ws event", "error", err)
		return
	}

	s.wsHub.Broadcast(websocket.Message{
		UserID: userID,
		Data:   eventJSON,
	})
}

func (s *AssistantService) ListenForAssistantEvents(ctx context.Context) {
	pubsub := s.redis.Rdb.Subscribe(ctx, "assistant.events")
	defer pubsub.Close()

	s.logger.Info("listening for AI assistant messages")

	ch := pubsub.Channel()

	for msg := range ch {
		var event dto.AssistantTaskEvent
		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			s.logger.Error("invalid assistant event", "error", err)
			continue
		}

		switch event.Type {
		case "assistant.chunk":
			s.StreamResponse(event.UserID, event.ConversationID, event.Chunk, false)

		case "assistant.completed":
			finalMsg := models.ChatMessage{
				ConversationID: event.ConversationID,
				Role:           "assistant",
				Content:        event.Content,
			}
			if err := s.db.Create(&finalMsg).Error; err != nil {
				s.logger.Error("failed to save assistant message", "error", err)
			}
			s.StreamResponse(event.UserID, event.ConversationID, event.Content, true)
		}
	}
}
