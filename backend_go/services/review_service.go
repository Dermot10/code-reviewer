package services

import (
	"fmt"
	"log/slog"

	"github.com/dermot10/code-reviewer/backend_go/cache"
	"github.com/dermot10/code-reviewer/backend_go/models"
	"gorm.io/gorm"
)

type ReviewService struct {
	db     *gorm.DB
	cache  *cache.RedisClient
	logger *slog.Logger

	// queue - *rabbitmq client
}

func NewReviewService(db *gorm.DB, cache *cache.RedisClient, logger *slog.Logger) *ReviewService {
	return &ReviewService{db: db, cache: cache, logger: logger}
}

func (s *ReviewService) CreateReview(userID uint, code string) (*models.Review, error) {
	review := &models.Review{
		UserID: userID,
		Code:   code,
		Status: "pending",
	}

	if err := s.db.Create(review).Error; err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	// TODO:
	// publish to rabbitMQ
	// Ai service  - add queue consumption

	return review, nil
}

func (s *ReviewService) CreateEnhancement(userID uint) (*models.Enhancement, error) {
	enhancement := &models.Enhancement{
		UserID: userID,
		Status: "pending",
	}

	if err := s.db.Create(enhancement).Error; err != nil {
		return nil, fmt.Errorf("failed to create enhancement: %w", err)
	}

	return enhancement, nil
}
