package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mux := setUpMux()
	registerRoutes(mux, logger)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return server(ctx, mux)
	})

	if err := g.Wait(); err != nil {
		return
	}
}
