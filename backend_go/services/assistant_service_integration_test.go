package services

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/dermot10/code-reviewer/backend_go/dto"
	"github.com/dermot10/code-reviewer/backend_go/models"
	"github.com/dermot10/code-reviewer/backend_go/redis"
	"github.com/dermot10/code-reviewer/backend_go/websocket"
	"github.com/stretchr/testify/require"
)

func TestAssistantService_SendPrompt(t *testing.T) {
	setup := SetUpTest(t)
	db := setup.DB
	rdb := setup.RDB

	// Create user and conversation that the FK requires
	user := models.User{Username: "testuser", Email: "test@test.com", HashedPassword: "hash", IsActive: true}
	require.NoError(t, db.Create(&user).Error)
	conv := models.Conversation{UserID: user.ID}
	require.NoError(t, db.Create(&conv).Error)

	logger := newTestLogger()
	wsHub := websocket.NewHub()
	rc := redis.NewRedisClientFromClient(rdb)
	service := NewAssistantService(db, rc, logger, wsHub)

	payload := dto.PromptPayload{
		ConversationID: 1,
		Prompt:         "Hello AI",
	}

	err := service.SendPrompt(1, payload)
	require.NoError(t, err)

	// DB row created
	var msg models.ChatMessage
	require.NoError(t, db.Where("conversation_id = ? AND role = ?", 1, "user").First(&msg).Error)
	require.Equal(t, "Hello AI", msg.Content)

	// Redis queue contains task
	val, err := rdb.RPop(context.Background(), "queue:tasks").Result()
	require.NoError(t, err)
	require.Contains(t, val, `"prompt":"Hello AI"`)
}

func TestAssistantService_ListenForAssistantEvents(t *testing.T) {
	setup := SetUpTest(t)
	db := setup.DB
	rdb := setup.RDB

	logger := newTestLogger()
	wsHub := websocket.NewHub()
	rc := redis.NewRedisClientFromClient(rdb)
	service := NewAssistantService(db, rc, logger, wsHub)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	go service.ListenForAssistantEvents(ctx)
	time.Sleep(100 * time.Millisecond)

	// Simulate "assistant.completed" event
	event := dto.AssistantTaskEvent{
		Type:           "assistant.completed",
		UserID:         1,
		ConversationID: 1,
		Content:        "AI response",
	}
	eventJSON, _ := json.Marshal(event)

	require.NoError(t, rdb.Publish(ctx, "assistant.events", eventJSON).Err())

	time.Sleep(200 * time.Millisecond) // allow listener to process

	// Check DB saved assistant message
	var msg models.ChatMessage
	err := db.Where("conversation_id = ? AND role = ?", 1, "assistant").First(&msg).Error
	require.NoError(t, err)
	require.Equal(t, "AI response", msg.Content)
}
