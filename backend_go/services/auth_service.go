package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/dermot10/code-reviewer/backend_go/models"
	"github.com/dermot10/code-reviewer/backend_go/redis"
	"github.com/dermot10/code-reviewer/backend_go/utils"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type AuthService struct {
	db        *gorm.DB
	redis     *redis.RedisClient
	logger    *slog.Logger
	jwtSecret string
}

func NewAuthService(db *gorm.DB, redis *redis.RedisClient, logger *slog.Logger, jwtSecret string) *AuthService {
	return &AuthService{db: db, redis: redis, logger: logger, jwtSecret: jwtSecret}
}

func (s *AuthService) CreateUser(username, email, password string) (*models.User, error) {
	if username == "" || email == "" || password == "" {
		return nil, errors.New("invalid credentials")
	}

	hashedPwd, err := utils.HashedPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:       username,
		Email:          email,
		HashedPassword: hashedPwd,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, errors.New("cannot create user in db")
	}

	return user, nil
}

func (s *AuthService) GetUser(userID uint) (*models.User, error) {
	ctx := context.Background()

	cacheKey := fmt.Sprintf("user:%d:profile", userID)
	cached, err := s.redis.GetCache(ctx, cacheKey)

	// check cache, if miss, query db
	if err == nil {
		var user models.User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			s.logger.Info("user cached", "key", cacheKey)
			return &user, nil
		}
	}

	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// cache for next hr
	if data, err := json.Marshal(&user); err == nil {
		if err = s.redis.SetCache(ctx, cacheKey, data); err != nil {
			s.logger.Error("failed to set cache", "error", err)
		}
	}

	return &user, nil

}

func (s *AuthService) Login(email, password string) (string, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPassword(user.HashedPassword, password) {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) Logout(userID int) error {
	ctx := context.Background()

	// del from redis
	profileKey := fmt.Sprintf("user:%d:profile", userID)
	s.redis.DelKey(ctx, profileKey)

	sessionKey := fmt.Sprintf("session:%d", userID)
	s.redis.DelKey(ctx, sessionKey)

	return nil
}
