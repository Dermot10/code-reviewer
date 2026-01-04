package services

import (
	"errors"
	"time"

	"github.com/dermot10/code-reviewer/backend_go/cache"
	"github.com/dermot10/code-reviewer/backend_go/dto"
	"github.com/dermot10/code-reviewer/backend_go/models"
	"github.com/dermot10/code-reviewer/backend_go/utils"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type AuthService struct {
	db    *gorm.DB
	cache *cache.RedisClient
}

func NewAuthService(db *gorm.DB, cache *cache.RedisClient) *AuthService {
	return &AuthService{db: db, cache: cache}
}

func (s *AuthService) CreateUser(username, email, password string) (*dto.CreateUserResponse, error) {
	if username == "" || email == "" || password == "" {
		return nil, errors.New("invalid credentials")
	}

	hashedPwd, err := utils.HashedPassword(password)
	if err != nil {
		return nil, nil
	}

	user := models.User{
		Username:       username,
		Email:          email,
		HashedPassword: hashedPwd,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, errors.New("cannot create user in db")
	}

	resp := &dto.CreateUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	return resp, nil
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

	tokenString, err := token.SignedString([]byte(""))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
