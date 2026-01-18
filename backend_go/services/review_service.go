package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/dermot10/code-reviewer/backend_go/dto"
	"github.com/dermot10/code-reviewer/backend_go/models"
	"github.com/dermot10/code-reviewer/backend_go/redis"
	"gorm.io/gorm"
)

type ReviewService struct {
	db     *gorm.DB
	redis  *redis.RedisClient
	logger *slog.Logger
}

func NewReviewService(db *gorm.DB, redis *redis.RedisClient, logger *slog.Logger) *ReviewService {
	return &ReviewService{db: db, redis: redis, logger: logger}
}

func (s *ReviewService) CreateReview(userID uint, code string) (*models.Review, error) {
	ctx := context.Background()

	review := &models.Review{
		UserID: userID,
		Code:   code,
		Status: "pending",
	}

	if err := s.db.Create(review).Error; err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	task := &dto.ReviewTask{
		Type:     "review",
		UserID:   userID,
		ReviewID: review.ID,
		Code:     review.Code,
		Action:   "generate_summary",
	}

	data, err := json.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal review task: %w", err)
	}
	if err := s.redis.PushQueue(ctx, data); err != nil {
		return nil, fmt.Errorf("failed to push review task to queue: %w", err)
	}

	if err := s.db.Model(review).Update("status", "processing").Error; err != nil {
		return nil, err
	}

	return review, nil
}

func (s *ReviewService) CreateEnhancement(userID uint) (*models.Enhancement, error) {
	ctx := context.Background()

	enhancement := &models.Enhancement{
		UserID: userID,
		Status: "pending",
	}

	if err := s.db.Create(enhancement).Error; err != nil {
		return nil, fmt.Errorf("failed to create enhancement: %w", err)
	}

	task := &dto.EnhanceTask{
		Type:          "enhance",
		UserID:        userID,
		EnhancementID: enhancement.ID,
		Action:        "generate_enhancement",
	}

	data, err := json.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal enhancement task: %w", err)
	}
	if err := s.redis.PushQueue(ctx, data); err != nil {
		return nil, fmt.Errorf("failed to push enhancement task to queue: %w", err)
	}

	if err := s.db.Model(enhancement).Update("status", "processing").Error; err != nil {
		return nil, err
	}

	return enhancement, nil
}

func (s *ReviewService) ListenForCompletions(ctx context.Context) {
	pubsub := s.redis.Rdb.Subscribe(ctx, "review.completed")
	defer pubsub.Close()

	s.logger.Info("listening for review completions")

	ch := pubsub.Channel()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("completion listener shutting down")
			return

		case msg, ok := <-ch:
			if !ok {
				s.logger.Info("pubsub channel closed")
				return
			}

			var event map[string]interface{}
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				s.logger.Error("failed to parse completion event", "error", err)
				continue
			}

			reviewID := uint(event["review_id"].(float64))

			// Get result from Redis
			resultKey := fmt.Sprintf("review:%d:result", reviewID)
			resultData, err := s.redis.Rdb.Get(ctx, resultKey).Result()
			if err != nil {
				s.logger.Error("failed to get result", "review_id", reviewID, "error", err)
				continue
			}

			// Update database
			if err := s.db.Model(&models.Review{}).Where("id = ?", reviewID).Updates(map[string]interface{}{
				"status": "completed",
				"result": resultData,
			}).Error; err != nil {
				s.logger.Error("failed to update review", "review_id", reviewID, "error", err)
				continue
			}

			s.redis.Rdb.Del(
				ctx,
				resultKey,
				fmt.Sprintf("review:%d:lock", reviewID),
			)

			s.logger.Info("review completed and redis deps cleaned up", "review_id", reviewID)
		}
	}
}
