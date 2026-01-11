package middleware

import (
	"fmt"
	"net/http"

	"github.com/dermot10/code-reviewer/backend_go/redis"
	"github.com/go-redis/redis_rate/v10"
)

const userID contextKey = "user_id"

func RateLimiterReviews(redis *redis.RedisClient) func(http.Handler) http.Handler {
	limiter := redis_rate.NewLimiter(redis.Rdb)
	return rateLimitMiddleware(limiter, "reviews", redis_rate.PerHour(10), byUser)
}

func RateLimitAuth(redis *redis.RedisClient) func(http.Handler) http.Handler {
	limiter := redis_rate.NewLimiter(redis.Rdb)
	return rateLimitMiddleware(limiter, "auth", redis_rate.PerMinute(5), byIP)
}

// layer 1 captures rl input creates limiter - returns handler
// layer 2 - configures behaviour - limit amount, resource name
// layer 3 - executes per request (check and enforce)

func rateLimitMiddleware(
	limiter *redis_rate.Limiter,
	resource string,
	limit redis_rate.Limit,
	keyFunc func(*http.Request) string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFunc(r)
			res, _ := limiter.Allow(r.Context(), key, limit)

			if res.Allowed == 0 {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
func byUser(r *http.Request) string {
	userID := r.Context().Value(UserIDKey).(uint)
	return fmt.Sprintf("user:%d", userID)
}

func byIP(r *http.Request) string {
	ip := r.RemoteAddr
	return fmt.Sprintf("ip:%s", ip)
}
