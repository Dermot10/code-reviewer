package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dermot10/code-reviewer/backend_go/config"
	"github.com/dermot10/code-reviewer/backend_go/websocket"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := godotenv.Load(); err != nil {
		logger.Warn("no .env file found", "error", err)
	}

	cfg, err := config.LoadConfig()

	if err != nil {
		logger.Error("error loading config", "error", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	deps, err := setUpDependencies(ctx, cfg)
	if err != nil {
		logger.Error("failed to setup dependencies", "error", err)
		os.Exit(1)
	}

	wsHub := websocket.NewHub()

	defer func() {
		logger.Info("cleaning up resources")
		if err := deps.redis.Close(); err != nil {
			logger.Error("error closing redis", "error", err)
		}

		if sqlDB, err := deps.db.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	registerRoutes(logger, deps, cfg.JWTSecret)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return setUpServer(ctx, deps.mux)
	})

	g.Go(func() error {
		deps.reviewService.ListenForCompletions(ctx)
		return nil
	})

	g.Go(func() error {
		wsHub.Run(ctx)
		return nil
	})

	if err := g.Wait(); err != nil {
		return
	}
}
