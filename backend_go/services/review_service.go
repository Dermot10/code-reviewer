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

	// Ai service  - add queue consumption

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

	return enhancement, nil
}
