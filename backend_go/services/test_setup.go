package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dermot10/code-reviewer/backend_go/models"
	"github.com/redis/go-redis/v9"
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

	// Start PostgreSQL container
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

	// Start Redis container
	redisC, err := rediscontainer.Run(ctx,
		"redis:7",
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("6379/tcp"),
		),
	)
	require.NoError(t, err)

	// Cleanup containers after test
	t.Cleanup(func() {
		pg.Terminate(ctx)
		redisC.Terminate(ctx)
	})

	// Get mapped host port for Postgres
	pgPort, err := pg.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Retry connecting to Postgres until ready
	var db *gorm.DB
	for i := 0; i < 5; i++ {
		connStr := fmt.Sprintf(
			"host=localhost port=%d user=test dbname=testdb password=test sslmode=disable",
			pgPort.Int(),
		)
		db, err = gorm.Open(gormpostgres.Open(connStr), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	require.NoError(t, err, "failed to connect to Postgres")

	// Auto-migrate tables
	err = db.AutoMigrate(
		&models.User{},
		&models.File{},
		&models.Review{},
		&models.Enhancement{},
		&models.ChatMessage{},
		&models.Conversation{},
	)
	require.NoError(t, err)

	// Insert default test fixtures (users) to prevent FK violations
	err = db.Create(&models.User{
		ID:   1,
		Name: "Test User",
	}).Error
	require.NoError(t, err)

	// Connect Redis client using mapped endpoint
	addr, err := redisC.Endpoint(ctx, "")
	require.NoError(t, err)

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return db, rdb
}
