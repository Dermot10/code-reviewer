package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dermot10/code-reviewer/backend_go/dto"
	"github.com/dermot10/code-reviewer/backend_go/middleware"
	"github.com/dermot10/code-reviewer/backend_go/models"
)

type mockAuthService struct{}

type failingAuthService struct {
	mockAuthService
}

func newTestAuthHandler() *AuthHandler {
	mockService := &mockAuthService{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	return NewAuthHandler(logger, mockService)
}

func (m *mockAuthService) CreateUser(username, email, password string) (*models.User, error) {
	return &models.User{
		ID:       1,
		Username: username,
		Email:    email,
	}, nil
}

func (m *mockAuthService) GetUser(userID uint) (*models.User, error) {
	return &models.User{
		ID:       userID,
		Username: "Marth",
		Email:    "marth@test.com",
	}, nil
}

func (m *mockAuthService) Login(email, password string) (string, error) {
	if email != "marth@test.com" || password != "falchion" {
		return "", fmt.Errorf("invalid credentials")
	}
	return "test-token", nil
}

func (m *mockAuthService) Logout(userID int) error {
	return nil
}

func (f *failingAuthService) CreateUser(username, email, password string) (*models.User, error) {
	return nil, fmt.Errorf("service failure")
}

func (f *failingAuthService) GetUser(userID uint) (*models.User, error) {
	return nil, fmt.Errorf("service failure")
}

func (f *failingAuthService) Login(email, password string) (string, error) {
	return "", fmt.Errorf("service failure")
}

func (f *failingAuthService) Logout(userID int) error {
	return fmt.Errorf("service failure")
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

func TestAuthHandler_CreateUser_InvalidJSON(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := NewAuthHandler(logger, &failingAuthService{})

	req := httptest.NewRequest("POST", "/users", strings.NewReader("{bad json"))
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Result().StatusCode)
	}
}

func TestAuthHandler_CreateUser_ServiceError(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := NewAuthHandler(logger, &failingAuthService{})

	body := `{"username":"marth", "email":"marth@test.com", "password":"pass"}`
	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	if w.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Result().StatusCode)
	}
}

func TestAuthHandler_GetUser(t *testing.T) {
	handler := newTestAuthHandler()

	req := httptest.NewRequest("GET", "/users", nil)

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uint(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.GetUser(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	var resp dto.UserResponse
	json.NewDecoder(res.Body).Decode(&resp)

	if resp.Username != "Marth" {
		t.Errorf("expected username 'Marth', got %s", resp.Username)
	}
}

func TestAuthHandler_GetUser_NoContext(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := NewAuthHandler(logger, &failingAuthService{})

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	handler.GetUser(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 401, got %d", w.Result().StatusCode)
	}
}

func TestAuthHandler_Login(t *testing.T) {
	handler := newTestAuthHandler()

	body := `{"email": "marth@test.com", "password": "falchion"}`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.Login(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	var resp map[string]string
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["token"] != "test-token" {
		t.Errorf("expected token 'test-token', got %s", resp["token"])
	}
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := NewAuthHandler(logger, &failingAuthService{})

	body := `{"email": "marth@test.com, "password": "falchion"`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.Login(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Result().StatusCode)
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	handler := newTestAuthHandler()

	req := httptest.NewRequest("POST", "/logout", nil)

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uint(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.Logout(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}
}

func TestAuthHandler_Logout_NoContext(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := NewAuthHandler(logger, &failingAuthService{})

	req := httptest.NewRequest("POST", "/logout", nil)
	w := httptest.NewRecorder()

	handler.Logout(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Result().StatusCode)
	}
}

func TestAuthHandler_Logout_ServiceError(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := NewAuthHandler(logger, &failingAuthService{})

	req := httptest.NewRequest("POST", "/logout", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uint(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.Logout(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", res.StatusCode)
	}
}
