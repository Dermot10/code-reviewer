package services

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"testing"

	"github.com/dermot10/code-reviewer/backend_go/models"
	"github.com/dermot10/code-reviewer/backend_go/redis"
	"github.com/stretchr/testify/require"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestAuthService_CreateUser(t *testing.T) {

	db, rdb := setUp(t)

	logger := newTestLogger()
	rc := redis.NewRedisClientFromClient(rdb)

	service := NewAuthService(db, rc, logger, "testsecret")

	resp, err := service.CreateUser("alice", "alice@test.com", "password")

	require.NoError(t, err)
	require.Equal(t, "alice", resp.Username)

	var user models.User
	err = db.First(&user, resp.ID).Error

	require.NoError(t, err)
	require.Equal(t, "alice@test.com", user.Email)
}

func TestAuthService_GetUser_CacheMissThenHit(t *testing.T) {

	db, rdb := setUp(t)

	logger := newTestLogger()
	rc := redis.NewRedisClientFromClient(rdb)

	service := NewAuthService(db, rc, logger, "testsecret")

	user := models.User{
		Username:       "bob",
		Email:          "bob@test.com",
		HashedPassword: "hash",
	}

	require.NoError(t, db.Create(&user).Error)

	// first call (cache miss)
	resp, err := service.GetUser(int(user.ID))
	require.NoError(t, err)
	require.Equal(t, "bob", resp.Username)

	// ensure cached
	key := "user:" + string(rune(user.ID)) + ":profile"

	val, err := rdb.Get(context.Background(), key).Result()
	require.NoError(t, err)

	var cached map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(val), &cached))

	require.Equal(t, "bob", cached["Username"])
}

func TestAuthService_Login(t *testing.T) {

	db, rdb := setUp(t)

	logger := newTestLogger()
	rc := redis.NewRedisClientFromClient(rdb)

	service := NewAuthService(db, rc, logger, "testsecret")

	resp, err := service.CreateUser("carol", "carol@test.com", "password")
	require.NoError(t, err)

	token, err := service.Login("carol@test.com", "password")

	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.Equal(t, resp.Email, "carol@test.com")
}

func TestAuthService_Logout(t *testing.T) {

	db, rdb := setUp(t)

	logger := newTestLogger()
	rc := redis.NewRedisClientFromClient(rdb)

	service := NewAuthService(db, rc, logger, "testsecret")

	ctx := context.Background()

	rdb.Set(ctx, "user:1:profile", "test", 0)
	rdb.Set(ctx, "session:1", "token", 0)

	err := service.Logout(1)

	require.NoError(t, err)

	_, err = rdb.Get(ctx, "user:1:profile").Result()
	require.Error(t, err)

	_, err = rdb.Get(ctx, "session:1").Result()
	require.Error(t, err)
}
