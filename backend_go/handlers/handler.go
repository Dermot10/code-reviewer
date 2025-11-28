package handlers

import "net/http"

type CodeReviewHandler struct{}

func NewCodeReviewHandler() *CodeReviewHandler {
	return &CodeReviewHandler{}
}

func (h *CodeReviewHandler) ReviewCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// another api call
}

func (h *CodeReviewHandler) ReviewFile(w http.ResponseWriter, r *http.Request) {

}
