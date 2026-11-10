package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sinon273/kvstore/internal/server"
	"github.com/sinon273/kvstore/internal/storage"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	store := storage.New()
	if err := server.Run(ctx, ":6379", store); err != nil {
		log.Fatal(err)
	}
}
