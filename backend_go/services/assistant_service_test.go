package services

import (
	"log/slog"

	"github.com/dermot10/code-reviewer/backend_go/redis"
	"github.com/dermot10/code-reviewer/backend_go/websocket"
	"gorm.io/gorm"
)

type assistantTestSuite struct {
	service *AssistantService
	db      *gorm.DB
	redis   *redis.RedisClient
	logger  *slog.Logger
	wsHub   *websocket.Hub
}
