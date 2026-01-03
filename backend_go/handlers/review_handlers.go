package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/dermot10/code-reviewer/backend_go/cache"
	"gorm.io/gorm"
)

type CodeReviewHandler struct {
	logger *slog.Logger
	db     *gorm.DB
	cache  *cache.RedisClient
}

type InputCode struct {
	SubmittedCode string `json:"submitted_code"`
}

func NewCodeReviewHandler(logger *slog.Logger, db *gorm.DB, cache *cache.RedisClient) *CodeReviewHandler {
	return &CodeReviewHandler{
		logger: logger,
		db:     db,
		cache:  cache,
	}
}

func submitCode(w http.ResponseWriter, r *http.Request, url string, logger *slog.Logger) {
	var input InputCode

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	logger.Info("received code submission")

	if input.SubmittedCode == "" {
		http.Error(w, "code cannot be empty", http.StatusBadRequest)
		return
	}

	payload := map[string]string{
		"submitted_code": input.SubmittedCode,
	}
	b, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}
	logger.Info("request forwarded to ai service")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	logger.Info("response received from ai service")

	respBody, _ := io.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(respBody)

	logger.Info("http response returned")
}

func (h *CodeReviewHandler) ReviewCode(w http.ResponseWriter, r *http.Request) {
	url := "http://127.0.0.1:8000/analyse/code"
	submitCode(w, r, url, h.logger)
}

func (h *CodeReviewHandler) EnhanceCode(w http.ResponseWriter, r *http.Request) {
	url := "http://127.0.0.1:8000/analyse/code-quality"
	submitCode(w, r, url, h.logger)
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
