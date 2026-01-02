package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dermot10/code-reviewer/backend_go/config"
	"golang.org/x/sync/errgroup"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("error loading config")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	deps, err := setUpDependencies(ctx, cfg)
	if err != nil {
		logger.Error("failed to setup dependencies", "error", err)
		os.Exit(1)
	}

	registerRoutes(logger, deps.mux, deps.db, deps.redis)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return setUpServer(ctx, deps.mux)
	})

	if err := g.Wait(); err != nil {
		return
	}
}
