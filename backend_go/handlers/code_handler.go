package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/dermot10/code-reviewer/backend_go/dto"
	"github.com/dermot10/code-reviewer/backend_go/middleware"
	"github.com/dermot10/code-reviewer/backend_go/models"
)

type CodeReviewHandler struct {
	logger      *slog.Logger
	codeService CodeService
}

type CodeService interface {
	CreateReview(userID uint, code string) (*models.Review, error)
	CreateEnhancement(userID uint) (*models.Enhancement, error)
	GetReview(userID uint, reviewID string) (*models.Review, error)
	ListenForCodeCompletions(ctx context.Context)
}

func NewCodeHandler(logger *slog.Logger, codeService CodeService) *CodeReviewHandler {
	return &CodeReviewHandler{
		logger:      logger,
		codeService: codeService,
	}
}

func (h *CodeReviewHandler) ReviewCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		http.Error(w, "unathorized", http.StatusUnauthorized)
		return
	}

	var requestCode dto.ReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&requestCode); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	h.logger.Info("received code submission for review")

	if requestCode.Code == "" {
		http.Error(w, "code cannot be empty", http.StatusBadRequest)
		return
	}

	review, err := h.codeService.CreateReview(userID, requestCode.Code)
	if err != nil {
		http.Error(w, "failed to create review", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"review_id": review.ID,
		"status":    review.Status,
		"message":   "Review queued for processing",
	})
}

func (h *CodeReviewHandler) EnhanceCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		http.Error(w, "unathorized", http.StatusUnauthorized)
		return
	}

	var requestCode dto.ReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&requestCode); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	h.logger.Info("received code submission for enhancement")

	if requestCode.Code == "" {
		http.Error(w, "code cannot be empty", http.StatusBadRequest)
		return
	}

	enhancement, err := h.codeService.CreateEnhancement(userID)
	if err != nil {
		http.Error(w, "failed to create enhancement", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"enhancement_id": enhancement.ID,
		"status":         enhancement.Status,
		"message":        "Code enhancements queued for processing",
	})

}

func (h *CodeReviewHandler) ExportReview(w http.ResponseWriter, r *http.Request) {
	exportType := r.URL.Query().Get("type")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	client := &http.Client{}
	url := fmt.Sprintf("http://127.0.0.1:8000/analyse/export-%s", exportType)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("Content-Disposition", resp.Header.Get("Content-Disposition"))
	io.Copy(w, resp.Body)
}

func (h *CodeReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	reviewID := r.PathValue("id")

	// presence check for the userID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// ensures any nil review returns early
	review, err := h.codeService.GetReview(userID, reviewID)
	if err != nil || review == nil {
		http.Error(w, "review not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}
