package handlers

import (
	"bytes"
	"io"
	"net/http"
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

func (h *CodeReviewHandler) ExportReview(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8000/analyse/export-md", bytes.NewReader(body))
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
