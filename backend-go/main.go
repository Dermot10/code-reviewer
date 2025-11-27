package main

import (
	"context"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	mux := setUpMux()
	registerRoutes(mux)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return server(ctx, mux)
	})

	if err := g.Wait(); err != nil {
		return
	}
}
