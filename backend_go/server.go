package main

import (
	"github.com/dermot10/code-reviewer/backend_go/handlers"

	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func setUpMux() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}

func registerRoutes(mux *http.ServeMux) {
	CodeReviewHandler := handlers.NewCodeReviewHandler()
	mux.HandleFunc("/review-code", CodeReviewHandler.ReviewCode)
	mux.HandleFunc("/review-code/download-md", CodeReviewHandler.ExportReview)

}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func server(ctx context.Context, mux *http.ServeMux) error {
	server := &http.Server{
		Addr:              ":" + "8080",
		Handler:           corsMiddleware(mux),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() error {
		log.Println("Server runnning on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
		return nil
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}
	return nil
}
