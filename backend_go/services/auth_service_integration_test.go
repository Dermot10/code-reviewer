package services

import (
	"context"
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
	setup := SetUpTest(t)
	db := setup.DB
	rdb := setup.RDB

	logger := newTestLogger()
	rc := redis.NewRedisClientFromClient(rdb)

	service := NewAuthService(db, rc, logger, "testsecret")

	resp, err := service.CreateUser("Roxas", "roxas@test.com", "password")

	require.NoError(t, err)
	require.Equal(t, "Roxas", resp.Username)

	var user models.User
	err = db.First(&user, resp.ID).Error

	require.NoError(t, err)
	require.Equal(t, "roxas@test.com", user.Email)
}

func TestAuthService_GetUser_CacheMissThenHit(t *testing.T) {
	setup := SetUpTest(t)
	db := setup.DB
	rdb := setup.RDB

	logger := newTestLogger()
	rc := redis.NewRedisClientFromClient(rdb)

	service := NewAuthService(db, rc, logger, "testsecret")

	user := models.User{
		Username:       "sora",
		Email:          "sora@test.com",
		HashedPassword: "hash",
	}

	require.NoError(t, db.Create(&user).Error)

	// first call (cache miss)
	resp, err := service.GetUser(user.ID)
	require.NoError(t, err)
	require.Equal(t, "sora", resp.Username)

	// second call should hit cache (no DB query required)
	resp2, err := service.GetUser(user.ID)
	require.NoError(t, err)
	require.Equal(t, "sora", resp2.Username)
}

func TestAuthService_Login(t *testing.T) {
	setup := SetUpTest(t)
	db := setup.DB
	rdb := setup.RDB

	logger := newTestLogger()
	rc := redis.NewRedisClientFromClient(rdb)

	service := NewAuthService(db, rc, logger, "testsecret")

	resp, err := service.CreateUser("axel", "axel@test.com", "chakrams")
	require.NoError(t, err)

	token, err := service.Login("axel@test.com", "chakrams")

	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.Equal(t, resp.Email, "axel@test.com")
}

func TestAuthService_Logout(t *testing.T) {
	setup := SetUpTest(t)
	db := setup.DB
	rdb := setup.RDB

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
