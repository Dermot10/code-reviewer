package handlers

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dermot10/code-reviewer/backend_go/middleware"
	"github.com/dermot10/code-reviewer/backend_go/models"
)

type mockCodeService struct {
	// change behaviour per test, strategy pattern
	createReviewFn      func(userID uint, code string) (*models.Review, error)
	createEnhancementFn func(userID uint) (*models.Enhancement, error)
	getReviewFn         func(userID uint, reviewID string) (*models.Review, error)
}

func NewtestCodeHandler() *CodeReviewHandler {
	mockService := &mockCodeService{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	return NewCodeHandler(logger, mockService)
}

func (m *mockCodeService) CreateReview(userID uint, code string) (*models.Review, error) {
	if m.createReviewFn != nil {
		return m.createReviewFn(userID, code)
	}
	// default success if function not set
	return &models.Review{ID: 1, Status: "pending"}, nil
}

func (m *mockCodeService) CreateEnhancement(userID uint) (*models.Enhancement, error) {
	if m.createEnhancementFn != nil {
		return m.createEnhancementFn(userID)
	}
	return &models.Enhancement{ID: 1, Status: "pending"}, nil
}

func (m *mockCodeService) GetReview(userID uint, reviewID string) (*models.Review, error) {
	if m.getReviewFn != nil {
		return m.getReviewFn(userID, reviewID)
	}
	return &models.Review{ID: 1, Status: "pending", Result: "mocked"}, nil
}

func (m *mockCodeService) ListenForCodeCompletions(ctx context.Context) {}

// success
// invalid input
// service failure
// missing auth

func TestCodeHandler_CreateReview_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mock := &mockCodeService{
		createReviewFn: func(user uint, code string) (*models.Review, error) {
			return &models.Review{ID: 1, Status: "pending"}, nil
		},
	}

	handler := NewCodeHandler(logger, mock)

	body := `{ "code": "print('the early bird catches the worm')" }`
	req := httptest.NewRequest("POST", "/review-code", strings.NewReader(body))

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uint(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ReviewCode(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusAccepted {
		t.Errorf("expected 202, got %d", res.StatusCode)
	}
}

func TestCodeHandler_CreateReview_InvalidJSON(t *testing.T) {
	handler := NewCodeHandler(slog.New(slog.NewTextHandler(io.Discard, nil)), &mockCodeService{})

	req := httptest.NewRequest("POST", "/review-code", strings.NewReader("{bad json"))

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uint(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ReviewCode(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Result().StatusCode)
	}
}

func TestCodeHandler_CreateReview_EmptyCode(t *testing.T) {
	handler := NewCodeHandler(slog.New(slog.NewTextHandler(io.Discard, nil)), &mockCodeService{})

	body := `{"code": ""}`
	req := httptest.NewRequest("POST", "/review-code", strings.NewReader(body))

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uint(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ReviewCode(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Result().StatusCode)
	}
}

func TestCodeHandler_CreateReview_ServiceError(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mock := &mockCodeService{
		createReviewFn: func(userID uint, code string) (*models.Review, error) {
			return nil, fmt.Errorf("service failure")
		},
	}

	handler := NewCodeHandler(logger, mock)

	body := `{"code": "print('hello')"}`
	req := httptest.NewRequest("POST", "/review-code", strings.NewReader(body))

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uint(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ReviewCode(w, req)

	if w.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Result().StatusCode)
	}
}

func TestCodeHandler_CreateReview_NoContext(t *testing.T) {
	handler := NewCodeHandler(slog.New(slog.NewTextHandler(io.Discard, nil)), &mockCodeService{})

	body := `{"code": "print('hello')"}`
	req := httptest.NewRequest("POST", "/review-code", strings.NewReader(body))

	w := httptest.NewRecorder()
	handler.ReviewCode(w, req)

	if w.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Result().StatusCode)
	}
}

func TestCodeHandler_EnhanceCode_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mock := &mockCodeService{
		createEnhancementFn: func(userID uint) (*models.Enhancement, error) {
			return &models.Enhancement{ID: 42, Status: "pending"}, nil
		},
	}

	handler := NewCodeHandler(logger, mock)

	body := `{"code": "print('enhance me')"}`
	req := httptest.NewRequest("POST", "/enhance-code", strings.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uint(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.EnhanceCode(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusAccepted {
		t.Errorf("expected 202, got %d", res.StatusCode)
	}
}

func TestCodeHandler_EnhanceCode_ServiceError(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mock := &mockCodeService{
		createEnhancementFn: func(userID uint) (*models.Enhancement, error) {
			return nil, fmt.Errorf("fail")
		},
	}

	handler := NewCodeHandler(logger, mock)

	body := `{"code": "print('enhance me')"}`
	req := httptest.NewRequest("POST", "/enhance-code", strings.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uint(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.EnhanceCode(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", res.StatusCode)
	}
}

func TestCodeHandler_GetReview_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mock := &mockCodeService{
		getReviewFn: func(userID uint, reviewID string) (*models.Review, error) {
			return &models.Review{ID: 10, Status: "completed", Result: "All good"}, nil
		},
	}

	handler := NewCodeHandler(logger, mock)

	req := httptest.NewRequest("GET", "/review?id=10", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uint(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetReview(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", res.StatusCode)
	}
}

func TestCodeHandler_GetReview_NotFound(t *testing.T) {
	handler := NewCodeHandler(slog.New(slog.NewTextHandler(io.Discard, nil)), &mockCodeService{
		getReviewFn: func(userID uint, reviewID string) (*models.Review, error) {
			return nil, fmt.Errorf("not found")
		},
	})

	req := httptest.NewRequest("GET", "/review?id=99", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uint(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetReview(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Result().StatusCode)
	}
}

func TestCodeHandler_GetReview_NoContext(t *testing.T) {
	handler := NewCodeHandler(slog.New(slog.NewTextHandler(io.Discard, nil)), &mockCodeService{})

	req := httptest.NewRequest("GET", "/review?id=10", nil)
	w := httptest.NewRecorder()
	handler.GetReview(w, req)

	if w.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Result().StatusCode)
	}
}
