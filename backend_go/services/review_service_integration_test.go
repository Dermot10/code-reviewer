package services

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/dermot10/code-reviewer/backend_go/models"
	"github.com/dermot10/code-reviewer/backend_go/redis"
	"github.com/stretchr/testify/require"
)

func TestReviewService_CreateReview_Listener(t *testing.T) {
	db, rdb := setUp(t) // your Testcontainers setup

	logger := newTestLogger()
	rc := redis.NewRedisClientFromClient(rdb)
	service := NewReviewService(db, rc, logger)

	// 1️ Create a review
	review, err := service.CreateReview(1, "print('hello')")
	require.NoError(t, err)
	require.Equal(t, "processing", review.Status)

	// 2️Prepare listener
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	go service.ListenForReviewCompletions(ctx)

	// 3️ Simulate completion event in Redis
	event := map[string]interface{}{
		"type":      "review.completed",
		"review_id": review.ID,
	}
	eventJSON, _ := json.Marshal(event)

	// Store a fake result in Redis (the listener will pick it up)
	resultKey := "review:" + string(rune(review.ID)) + ":result"
	require.NoError(t, rdb.Set(ctx, resultKey, "review result", 0).Err())

	// Publish completion event
	require.NoError(t, rdb.Publish(ctx, "review.completed", eventJSON).Err())

	// 4️ Wait for listener to process
	time.Sleep(200 * time.Millisecond)

	// 5️ Assert review updated in DB
	var updated models.Review
	require.NoError(t, db.First(&updated, review.ID).Error)
	require.Equal(t, "completed", updated.Status)
	require.Equal(t, "review result", updated.Result)

	// 6️ Assert Redis keys cleaned up
	_, err = rdb.Get(ctx, resultKey).Result()
	require.Error(t, err) // should be deleted
}
