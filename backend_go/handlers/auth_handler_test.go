package handlers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dermot10/code-reviewer/backend_go/dto"
)

type mockAuthService struct{}

func newTestAuthHandler() *AuthHandler {
	mockService := &mockAuthService{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	return NewAuthHandler(logger, mockService)
}

func (m *mockAuthService) CreateUser(username, email, password string) (*dto.CreateUserResponse, error) {
	return &dto.CreateUserResponse{
		ID:       1,
		Username: username,
		Email:    email,
	}, nil
}

func (m *mockAuthService) GetUser(userID uint) (*dto.UserResponse, error) {
	return &dto.UserResponse{
		ID:       userID,
		Username: "Marth",
		Email:    "marth@test.com",
	}, nil
}

func (m *mockAuthService) Login(email, password string) (string, error) {
	return "test-token", nil
}

func (m *mockAuthService) Logout(userID int) error {
	return nil
}

func TestAuthHandler_CreateUser(t *testing.T) {
	handler := newTestAuthHandler()

	body := `{"username":"marth", "email":"marth@test.com", "password":"pass123"}`
	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", res.StatusCode)
	}

	var resp dto.CreateUserResponse
	json.NewDecoder(res.Body).Decode(&resp)
	if resp.Username != "marth" {
		t.Errorf("expected username 'marth', got %s", resp.Username)
	}
}
