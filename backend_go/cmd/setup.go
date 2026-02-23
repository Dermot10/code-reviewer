package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/dermot10/code-reviewer/backend_go/config"
	"github.com/dermot10/code-reviewer/backend_go/database"
	"github.com/dermot10/code-reviewer/backend_go/handlers"
	"github.com/dermot10/code-reviewer/backend_go/middleware"
	"github.com/dermot10/code-reviewer/backend_go/models"
	"github.com/dermot10/code-reviewer/backend_go/redis"
	"github.com/dermot10/code-reviewer/backend_go/services"
	"github.com/dermot10/code-reviewer/backend_go/websocket"
	"gorm.io/gorm"
)

type Dependencies struct {
	redis         *redis.RedisClient
	db            *gorm.DB
	wsHub         *websocket.Hub
	mux           *http.ServeMux
	authService   *services.AuthService
	reviewService *services.ReviewService
	fileService   *services.FileService
}

func setUpDependencies(ctx context.Context, cfg *config.Config) (*Dependencies, error) {
	db, err := database.Connect(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := setUpMigrations(db); err != nil {
		return nil, err
	}

	r, err := redis.NewRedisService(cfg)
	if err != nil {
		return nil, err
	}

	wsHub := websocket.NewHub()

	mux := http.NewServeMux()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	authService := services.NewAuthService(db, r, logger, cfg.JWTSecret)
	reviewService := services.NewReviewService(db, r, logger, wsHub)
	fileService := services.NewFileService(db, logger)

	return &Dependencies{
		db:            db,
		redis:         r,
		wsHub:         wsHub,
		mux:           mux,
		authService:   authService,
		reviewService: reviewService,
		fileService:   fileService,
	}, nil
}

func setUpMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.User{},
		&models.Review{},
		&models.Enhancement{},
		&models.File{},
	); err != nil {
		return fmt.Errorf("db migrate: %w", err)
	}
	return nil
}

func registerRoutes(logger *slog.Logger, deps *Dependencies, jwtSecret string) {

	codeReviewHandler := handlers.NewCodeReviewHandler(logger, deps.db, deps.redis, deps.reviewService)
	authReviewHandler := handlers.NewAuthHandler(logger, deps.authService)
	healthHandler := handlers.NewHealthHandler(logger, deps.db, deps.redis)
	fileHandler := handlers.NewFileHandler(logger, deps.db, deps.fileService)
	wsHandler := handlers.NewWSHandler(logger, deps.wsHub)
	// metricsHandler := handlers.NewMetricsHandler(logger, db, redis)

	deps.mux.HandleFunc("/api/auth/register", authReviewHandler.CreateUser)
	deps.mux.HandleFunc("/healthz", healthHandler.HealthCheck)
	// mux.HandleFunc("/metrics", metricsHandler.)
	deps.mux.Handle(
		"/api/auth/login",
		middleware.RateLimitAuth(deps.redis)(
			http.HandlerFunc(authReviewHandler.Login),
		),
	)

	// auth protected routes - may need to refactor,

	deps.mux.Handle(
		"/api/users/me",
		middleware.AuthMiddleware(jwtSecret)(
			http.HandlerFunc(authReviewHandler.GetUser),
		),
	)

	deps.mux.Handle(
		"/api/auth/logout",
		middleware.AuthMiddleware(jwtSecret)(
			http.HandlerFunc(authReviewHandler.Logout),
		),
	)

	deps.mux.Handle(
		"/review-code",
		middleware.AuthMiddleware(jwtSecret)(
			middleware.RateLimiterReviews(deps.redis)(
				http.HandlerFunc(codeReviewHandler.ReviewCode),
			)),
	)

	deps.mux.Handle(
		"/enhance-code",
		middleware.AuthMiddleware(jwtSecret)(
			http.HandlerFunc(codeReviewHandler.EnhanceCode),
		),
	)

	deps.mux.Handle(
		"/review-code/download",
		middleware.AuthMiddleware(jwtSecret)(
			http.HandlerFunc(codeReviewHandler.ExportReview),
		),
	)

	deps.mux.Handle(
		"/api/reviews/{id}",
		middleware.AuthMiddleware(jwtSecret)(
			http.HandlerFunc(codeReviewHandler.GetReview),
		),
	)

	deps.mux.Handle(
		"POST /api/files",
		middleware.AuthMiddleware(jwtSecret)(
			http.HandlerFunc(fileHandler.CreateFile),
		),
	)

	deps.mux.Handle(
		"GET /api/files",
		middleware.AuthMiddleware(jwtSecret)(
			http.HandlerFunc(fileHandler.ListFiles),
		),
	)

	deps.mux.Handle(
		"GET /api/files/{id}",
		middleware.AuthMiddleware(jwtSecret)(
			http.HandlerFunc(fileHandler.GetFile),
		),
	)

	deps.mux.Handle(
		"PUT /api/files/{id}",
		middleware.AuthMiddleware(jwtSecret)(
			http.HandlerFunc(fileHandler.UpdateFile),
		),
	)

	deps.mux.Handle(
		"DELETE /api/files/{id}",
		middleware.AuthMiddleware(jwtSecret)(
			http.HandlerFunc(fileHandler.DeleteFile),
		),
	)

	deps.mux.Handle(
		"/ws",
		middleware.AuthMiddleware(jwtSecret)(
			http.HandlerFunc(wsHandler.HandleWebSocket),
		),
	)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

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
