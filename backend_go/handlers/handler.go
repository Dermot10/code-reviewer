package handlers

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

type CodeReviewHandler struct{}

func NewCodeReviewHandler() *CodeReviewHandler {
	return &CodeReviewHandler{}
}

func (h *CodeReviewHandler) ReviewCode(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	code := string(body)
	if code == "" {
		http.Error(w, "code cannot be empty", http.StatusBadRequest)
		return
	}

	// payload := map[string]string{"code": code}
	// payloadBytes, _ := json.Marshal(payload)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8000/analyse/code", bytes.NewReader(body))
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

	respBody, _ := io.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

func (h *CodeReviewHandler) ReviewFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // set max memory to ~10MB
	if err != nil {
		http.Error(w, "failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file not provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := fileHeader.Filename
	if !strings.HasSuffix(fileName, "py") {
		http.Error(w, "only .py files are allowed for review", http.StatusBadRequest)
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil || len(fileBytes) == 0 {
		http.Error(w, "file is empty or unreadable", http.StatusBadRequest)
		return
	}
}
