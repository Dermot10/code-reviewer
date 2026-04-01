package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/dermot10/code-reviewer/backend_go/models"
	"github.com/redis/go-redis/v9" // Redis client for service
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	rediscontainer "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setUp(t *testing.T) (*gorm.DB, *redis.Client) {
	ctx := context.Background()

	// Start PostgreSQL and Redis containers
	pg, err := postgrescontainer.Run(ctx,
		"postgres:15",
		postgrescontainer.WithDatabase("testdb"),
		postgrescontainer.WithUsername("test"),
		postgrescontainer.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp"),
		),
	)
	require.NoError(t, err)

	redisC, err := rediscontainer.Run(ctx,
		"redis:7",
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("6379/tcp"),
		),
	)
	require.NoError(t, err)

	// Cleanup
	t.Cleanup(func() {
		pg.Terminate(ctx)
		redisC.Terminate(ctx)
	})

	// Connect GORM to Postgres container
	pgPort, err := pg.MappedPort(ctx, "5432")
	require.NoError(t, err)

	connStr := fmt.Sprintf(
		"host=localhost port=%d user=test dbname=testdb password=test sslmode=disable",
		pgPort.Int(),
	)

	db, err := gorm.Open(gormpostgres.Open(connStr), &gorm.Config{})
	require.NoError(t, err)

	// Migrate DB table to testcontainer for test usage
	db.AutoMigrate(
		&models.User{},
		&models.File{},
		&models.Review{},
		&models.Enhancement{},
		&models.ChatMessage{},
		&models.Conversation{},
	)

	// Connect Redis client to Redis container
	addr, err := redisC.Endpoint(ctx, "")
	require.NoError(t, err)

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return db, rdb
}
