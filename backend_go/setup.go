package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	cache "github.com/dermot10/code-reviewer/backend_go/Cache"
	"github.com/dermot10/code-reviewer/backend_go/config"
	"github.com/dermot10/code-reviewer/backend_go/database"
	"github.com/dermot10/code-reviewer/backend_go/handlers"
	"gorm.io/gorm"
)

type Dependencies struct {
	redis *cache.RedisClient
	db    *gorm.DB
	mux   *http.ServeMux
}

func setUpDependencies(ctx context.Context, cfg *config.Config) (*Dependencies, error) {
	// could add to set up layer and just evoke the wrapper
	db, err := database.Connect(ctx)
	if err != nil {
		return nil, err
	}

	c, err := cache.NewCacheService(cfg)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	return &Dependencies{
		db:    db,
		redis: c,
		mux:   mux,
	}, nil
}

func registerRoutes(logger *slog.Logger, mux *http.ServeMux, db *gorm.DB, cache *cache.RedisClient) {
	CodeReviewHandler := handlers.NewCodeReviewHandler(logger, db, cache)
	AuthReviewHandler := handlers.NewAuthHandler(logger, db, cache)

	mux.HandleFunc("/review-code", CodeReviewHandler.ReviewCode)
	mux.HandleFunc("/enhance-code", CodeReviewHandler.EnhanceCode)
	mux.HandleFunc("/review-code/download", CodeReviewHandler.ExportReview)
	mux.HandleFunc("/sign-up", AuthReviewHandler.CreateUser)

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

func setUpServer(ctx context.Context, mux *http.ServeMux) error {
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
