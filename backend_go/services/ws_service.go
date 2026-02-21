package services

import (
	"log"
	"log/slog"

	"github.com/dermot10/code-reviewer/backend_go/dto"
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

func (s *WsService) handleFileUpload(userID uint, event dto.WSEvent) {
	//validate
	//persist
	//broadcast
}

func (s *WsService) handleChatMessage(userID uint, event dto.WSEvent) {
	// validate
	// persist if needed
	// broadcast
}
