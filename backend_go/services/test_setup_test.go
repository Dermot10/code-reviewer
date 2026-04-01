package services

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/dermot10/code-reviewer/backend_go/models"
	"github.com/docker/go-connections/nat"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	rediscontainer "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestSetup struct {
	DB  *gorm.DB
	RDB *redis.Client
}

var (
	pgContainer    *postgrescontainer.PostgresContainer
	redisContainer *rediscontainer.RedisContainer
	pgConnStr      string
	redisAddr      string
	initOnce       sync.Once
	initErr        error
)

func init() {
	_ = godotenv.Load(".env")
}

// TestMain owns the full container lifecycle.
// All containers start before any test runs.
// All containers terminate after every test in the package finishes.
func TestMain(m *testing.M) {
	ctx := context.Background()

	initOnce.Do(func() {
		initErr = startContainers(ctx)
	})

	if initErr != nil {
		fmt.Fprintf(os.Stderr, "failed to start test containers: %v\n", initErr)
		os.Exit(1)
	}

	code := m.Run()

	if pgContainer != nil {
		_ = pgContainer.Terminate(ctx)
	}
	if redisContainer != nil {
		_ = redisContainer.Terminate(ctx)
	}

	os.Exit(code)
}

func startContainers(ctx context.Context) error {
	var err error

	pgContainer, err = postgrescontainer.Run(ctx,
		"postgres:15",
		postgrescontainer.WithDatabase("testdb"),
		postgrescontainer.WithUsername("test"),
		postgrescontainer.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			// Wait until postgres actually accepts queries, not just TCP
			wait.ForSQL("5432/tcp", "pgx", func(host string, port nat.Port) string {
				return fmt.Sprintf("host=%s port=%s user=test dbname=testdb password=test sslmode=disable", host, port.Port())
			}).WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return fmt.Errorf("start postgres: %w", err)
	}

	pgPort, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		return fmt.Errorf("get postgres port: %w", err)
	}

	pgConnStr = fmt.Sprintf(
		"host=localhost port=%d user=test dbname=testdb password=test sslmode=disable",
		pgPort.Int(),
	)

	var db *gorm.DB
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(gormpostgres.Open(pgConnStr), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}

	if err = db.AutoMigrate(
		&models.User{},
		&models.File{},
		&models.Review{},
		&models.Enhancement{},
		&models.ChatMessage{},
		&models.Conversation{},
	); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	// Verify tables actually exist before proceeding
	var tables []string
	db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables)
	fmt.Printf("Tables created: %v\n", tables)

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(10)
	sqlDB.Close()

	redisContainer, err = rediscontainer.Run(ctx,
		"redis:7",
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("6379/tcp").WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return fmt.Errorf("start redis: %w", err)
	}

	redisAddr, err = redisContainer.Endpoint(ctx, "")
	if err != nil {
		return fmt.Errorf("get redis endpoint: %w", err)
	}

	return nil
}

// SetUpTest is called at the start of each individual test.
// Wipes all tables and resets sequences so every test starts with id=1.
func SetUpTest(t *testing.T) *TestSetup {
	t.Helper()

	db, err := gorm.Open(gormpostgres.Open(pgConnStr), &gorm.Config{})
	require.NoError(t, err)

	// Child tables first to respect FK constraints
	tables := []string{
		"chat_messages", "conversations",
		"enhancements", "reviews",
		"files", "users",
	}
	for _, table := range tables {
		require.NoError(t, db.Exec("TRUNCATE "+table+" CASCADE").Error)
	}
	for _, table := range tables {
		db.Exec(fmt.Sprintf("ALTER SEQUENCE %s_id_seq RESTART WITH 1", table))
	}

	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		rdb.Close()
	})

	return &TestSetup{DB: db, RDB: rdb}
}
