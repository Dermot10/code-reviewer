package services

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9" // Redis client for service
	"github.com/stretchr/testify/require"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	rediscontainer "github.com/testcontainers/testcontainers-go/modules/redis"
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
	)
	require.NoError(t, err)

	redisC, err := rediscontainer.Run(ctx, "redis:7")
	require.NoError(t, err)

	// Cleanup
	t.Cleanup(func() {
		pg.Terminate(ctx)
		redisC.Terminate(ctx)
	})

	// Connect GORM to Postgres container
	pgConnStr, _ := pg.ConnectionString(ctx, "sslmode=disable")
	db, err := gorm.Open(gormpostgres.Open(pgConnStr), &gorm.Config{})
	require.NoError(t, err)

	// Connect Redis client to Redis container
	addr, err := redisC.Endpoint(ctx, "")
	require.NoError(t, err)

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return db, rdb
}
