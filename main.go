package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/flashbots-lab/searcher/bundle"
	"github.com/flashbots-lab/searcher/mempool"
)

func main() {
	wsURL := os.Getenv("GETH_WS_URL")
	if wsURL == "" {
		wsURL = "ws://127.0.0.1:32787"
	}

	log.Printf("Starting MEV searcher, connecting to %s", wsURL)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bundler := bundle.NewBundler()
	handler := mempool.NewHandler(bundler)
	sub := mempool.NewSubscriber(wsURL, handler)

	go func() {
		for b := range bundler.Start() {
			bundle.LogBundle(b)
		}
	}()
	
	errCh := make(chan error, 1)
	go func() {
		errCh <- sub.Start(ctx)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		log.Fatalf("Subscriber error: %v", err)
	case sig := <-sigCh:
		log.Printf("Received signal %v, shutting down", sig)
		cancel()
	}
}
