package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dermot10/code-reviewer/backend_go/redis"
	"gorm.io/gorm"
)

type HealthHandler struct {
	logger *slog.Logger
	redis  *redis.RedisClient
	db     *gorm.DB
}

type MetricsHandler struct {
	logger *slog.Logger
	redis  *redis.RedisClient
	db     *gorm.DB
}

func NewHealthHandler(logger *slog.Logger, db *gorm.DB, redis *redis.RedisClient) *HealthHandler {
	return &HealthHandler{
		logger: logger,
		redis:  redis,
		db:     db,
	}
}

func NewMetricsHandler(logger *slog.Logger, db *gorm.DB, redis *redis.RedisClient) *MetricsHandler {
	return &MetricsHandler{
		logger: logger,
		redis:  redis,
		db:     db,
	}
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := h.redis.Rdb.Ping(ctx).Err(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy", "error": "redis down"})
		return
	}

	sqlDB, _ := h.db.DB()
	if err := sqlDB.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy", "error": "database down"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
