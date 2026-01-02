package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	cache "github.com/dermot10/code-reviewer/backend_go/Cache"
	"github.com/dermot10/code-reviewer/backend_go/dto"
	"github.com/dermot10/code-reviewer/backend_go/models"
	"github.com/dermot10/code-reviewer/backend_go/utils"
	"gorm.io/gorm"
)

// auth handlers for sign up, sign in

// org id
// user id
// user password and email

// gonna inject cache into handler like db

type AuthHandler struct {
	logger *slog.Logger
	db     *gorm.DB
	cache  *cache.RedisClient
}

func NewAuthHandler(logger *slog.Logger, db *gorm.DB, cache *cache.RedisClient) *AuthHandler {
	return &AuthHandler{
		logger: logger,
		db:     db,
		cache:  cache,
	}
}

func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}

	hashedPwd, err := utils.HashedPassword(req.Password)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: hashedPwd,
	}

	if err := h.db.Create(&user).Error; err != nil {
		http.Error(w, "could not create user", http.StatusConflict)
		return
	}

	resp := dto.CreateUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// take request
	// check cache for user session
	// if not send query to db
	// db response
	// cache user session
}

func (h *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

}
