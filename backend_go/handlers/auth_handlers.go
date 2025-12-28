package handlers

import (
	"log/slog"
	"net/http"
)

// auth handlers for sign up, sign in

// org id
// user id
// user password and email

type AuthHandler struct {
	logger *slog.Logger
}

func NewAuthHandler(logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		logger: logger,
	}
}

func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

}
