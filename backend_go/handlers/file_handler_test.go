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
	"time"

	"github.com/dermot10/code-reviewer/backend_go/dto"
	"github.com/dermot10/code-reviewer/backend_go/middleware"
	"github.com/dermot10/code-reviewer/backend_go/models"
)

type mockFileService struct{}

type failingFileService struct {
	mockFileService
}

func newTestFileHandler() *FileHandler {
	mockService := &mockFileService{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	return NewFileHandler(logger, mockService)
}

func (m *mockFileService) CreateFile(userID uint, name, content string) (*models.File, error) {
	return &models.File{
		ID:      1,
		UserID:  userID,
		Name:    name,
		Content: content,
	}, nil
}

func (m *mockFileService) GetFile(userID, fileID uint) (*models.File, error) {
	return &models.File{
		ID:        fileID,
		UserID:    userID,
		Name:      "file1",
		Content:   "hello",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *mockFileService) ListFiles(userID uint) ([]models.File, error) {
	return []models.File{
		{
			ID:        1,
			UserID:    userID,
			Name:      "file1",
			Content:   "hello",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			UserID:    userID,
			Name:      "file2",
			Content:   "world",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil
}

func (m *mockFileService) UpdateFile(userID, fileID uint, content string) (*models.File, error) {
	return &models.File{
		ID:        fileID,
		UserID:    userID,
		Name:      "file1",
		Content:   "hello world",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *mockFileService) DeleteFile(userID, fileID uint) error {
	return nil
}

func (m *failingFileService) CreateFile(userID uint, name, content string) (*models.File, error) {
	return nil, fmt.Errorf("service failure")
}

func (m *failingFileService) GetFile(userID, fileID uint) (*models.File, error) {
	return nil, fmt.Errorf("service failure")
}

func (m *failingFileService) ListFiles(userID uint) ([]models.File, error) {
	return nil, fmt.Errorf("service failure")
}

func (m *failingFileService) UpdateFile(userID, fileID uint, content string) (*models.File, error) {
	return nil, fmt.Errorf("service failure")
}

func (m *failingFileService) DeleteFile(userID, fileID uint) error {
	return fmt.Errorf("service failure")
}

func TestFileHandler_CreateFile(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := NewFileHandler(logger, &mockFileService{})

	body := `{"name": "test.txt", "content": "hello"}`
	req := httptest.NewRequest("POST", "/files", strings.NewReader(body))

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uint(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.CreateFile(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", res.StatusCode)
	}
}

func TestFileHandler_GetFile(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := NewFileHandler(logger, &mockFileService{})

	req := httptest.NewRequest("GET", "/files", nil)

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uint(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetFile(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	var resp dto.FileResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.ID != 1 {
		t.Errorf("expected file ID 1, got %d", resp.ID)
	}

	if resp.Name != "file1" {
		t.Errorf("expected name 'file1', got %s", resp.Name)
	}

	if resp.Content != "hello" {
		t.Errorf("expected content 'hello', got %s", resp.Content)
	}
}

func TestFileHandler_ListFiles(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := NewFileHandler(logger, &mockFileService{})

	req := httptest.NewRequest("GET", "/files", nil)

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uint(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ListFiles(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	var resp []dto.FileResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatal("failed to decode response: %v", err)
	}

	if len(resp) != 2 {
		t.Errorf("expected 2 files, got %d", len(resp))
	}

	if resp[0].Name != "file1" || resp[1].Name != "file2" {
		t.Errorf("expected files 'file1' and 'file2', got '%s' and '%s'", resp[0].Name, resp[1].Name)
	}
}

func TestFileHandler_UpdateFile(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := NewFileHandler(logger, &mockFileService{})

	body := `{"name": "test.txt", "content": "hello world"}`
	req := httptest.NewRequest("PUT", "/files", strings.NewReader(body))

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uint(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.UpdateFile(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	var resp dto.FileResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.ID != 1 {
		t.Errorf("expected file ID 1, got %d", resp.ID)
	}

	if resp.Name != "file1" {
		t.Errorf("expected name 'file1', got %s", resp.Name)
	}

	if resp.Content != "hello world" {
		t.Errorf("expected content 'hello world', got %s", resp.Content)
	}
}

func TestFileHandler_DeleteFile(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := NewFileHandler(logger, &mockFileService{})

	req := httptest.NewRequest("DELETE", "/files", nil)

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uint(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.DeleteFile(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

}
