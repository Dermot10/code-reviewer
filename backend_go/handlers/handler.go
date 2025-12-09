package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CodeReviewHandler struct{}

func NewCodeReviewHandler() *CodeReviewHandler {
	return &CodeReviewHandler{}
}

func submitCode(w http.ResponseWriter, r *http.Request, url string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read body", http.StatusBadRequest)
		return
	}

	code := string(body)
	if code == "" {
		http.Error(w, "code cannot be empty", http.StatusBadRequest)
		return
	}

	payload := map[string]string{"submitted_code": code}
	b, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
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

func (h *CodeReviewHandler) ReviewCode(w http.ResponseWriter, r *http.Request) {
	url := "http://127.0.0.1:8000/analyse/code"
	submitCode(w, r, url)
}

func (h *CodeReviewHandler) EnhanceCode(w http.ResponseWriter, r *http.Request) {
	url := "http://127.0.0.1:8000/analyse/code-quality"
	submitCode(w, r, url)
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
