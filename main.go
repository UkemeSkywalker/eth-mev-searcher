package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	sub := mempool.NewSubscriber(wsURL)
	
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
