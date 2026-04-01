package services

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/dermot10/code-reviewer/backend_go/models"
	"github.com/dermot10/code-reviewer/backend_go/redis"
	"github.com/stretchr/testify/require"
)

func TestReviewService_CreateReview_Listener(t *testing.T) {
	setup := SetUpTest(t)
	db := setup.DB
	rdb := setup.RDB

	// Create a new user
	user := models.User{Username: "testuser", Email: "test@test.com", HashedPassword: "hash", IsActive: true}
	require.NoError(t, db.Create(&user).Error)

	logger := newTestLogger()
	rc := redis.NewRedisClientFromClient(rdb)
	service := NewReviewService(db, rc, logger)

	// Create a review
	review, err := service.CreateReview(1, "print('hello')")
	require.NoError(t, err)
	require.Equal(t, "processing", review.Status)

	// Prepare listener
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	go service.ListenForCodeCompletions(ctx)
	time.Sleep(100 * time.Millisecond)

	// Simulate completion event in Redis
	event := map[string]interface{}{
		"type":      "review.completed",
		"review_id": review.ID,
	}
	eventJSON, _ := json.Marshal(event)

	// Store a fake result in Redis (the listener will pick it up)
	resultKey := fmt.Sprintf("review:%d:result", review.ID)
	require.NoError(t, rdb.Set(ctx, resultKey, "review result", 0).Err())

	// Publish completion event
	require.NoError(t, rdb.Publish(ctx, "review.completed", eventJSON).Err())

	// Wait for listener to process
	time.Sleep(200 * time.Millisecond)

	// Assert review updated in DB
	var updated models.Review
	require.NoError(t, db.First(&updated, review.ID).Error)
	require.Equal(t, "completed", updated.Status)
	require.Equal(t, "review result", updated.Result)

	// Assert Redis keys cleaned up
	_, err = rdb.Get(ctx, resultKey).Result()
	require.Error(t, err) // should be deleted
}
