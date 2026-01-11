package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/dermot10/code-reviewer/backend_go/config"
	"github.com/dermot10/code-reviewer/backend_go/database"
	"github.com/dermot10/code-reviewer/backend_go/handlers"
	"github.com/dermot10/code-reviewer/backend_go/middleware"
	"github.com/dermot10/code-reviewer/backend_go/models"
	"github.com/dermot10/code-reviewer/backend_go/redis"
	"github.com/dermot10/code-reviewer/backend_go/services"
	"gorm.io/gorm"
)

type Dependencies struct {
	redis *redis.RedisClient
	db    *gorm.DB
	mux   *http.ServeMux
}

func setUpDependencies(ctx context.Context, cfg *config.Config) (*Dependencies, error) {
	db, err := database.Connect(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := setUpMigrations(db); err != nil {
		return nil, err
	}

	c, err := redis.NewRedisService(cfg)
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

func setUpMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.User{},
		&models.Review{},
		&models.Enhancement{},
	); err != nil {
		return fmt.Errorf("db migrate: %w", err)
	}
	return nil
}

func registerRoutes(logger *slog.Logger, mux *http.ServeMux, db *gorm.DB, redis *redis.RedisClient, jwtSecret string) {

	authService := services.NewAuthService(db, redis, logger, jwtSecret)
	reviewService := services.NewReviewService(db, redis, logger)

	codeReviewHandler := handlers.NewCodeReviewHandler(logger, db, redis, reviewService)
	authReviewHandler := handlers.NewAuthHandler(logger, authService)
	healthHandler := handlers.NewHealthHandler(logger, db, redis)
	// metricsHandler := handlers.NewMetricsHandler(logger, db, redis)

	mux.HandleFunc("/api/auth/register", authReviewHandler.CreateUser)
	mux.HandleFunc("/healthz", healthHandler.HealthCheck)
	// mux.HandleFunc("/metrics", metricsHandler.)
	mux.Handle(
		"/api/auth/login",
		middleware.RateLimitAuth(redis)(
			http.HandlerFunc(authReviewHandler.Login),
		),
	)

	// auth protected routes - may need to refactor, and include review handlers in too
	mux.Handle(
		"/api/users/me",
		middleware.AuthMiddleware(jwtSecret)(
			http.HandlerFunc(authReviewHandler.GetUser),
		),
	)

	mux.Handle(
		"/api/auth/logout",
		middleware.AuthMiddleware(jwtSecret)(
			http.HandlerFunc(authReviewHandler.Logout),
		),
	)

	mux.Handle(
		"/review-code",
		middleware.AuthMiddleware(jwtSecret)(
			middleware.RateLimiterReviews(redis)(
				http.HandlerFunc(codeReviewHandler.ReviewCode),
			)),
	)

	mux.Handle(
		"/enhance-code",
		middleware.AuthMiddleware(jwtSecret)(
			http.HandlerFunc(codeReviewHandler.EnhanceCode),
		),
	)

	mux.Handle(
		"/review-code/download",
		middleware.AuthMiddleware(jwtSecret)(
			http.HandlerFunc(codeReviewHandler.ExportReview),
		),
	)

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
