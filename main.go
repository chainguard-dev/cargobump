package main

import (
	"context"
	"github.com/hectorj2f/cargobump/cmd/cargobump"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt)
	defer done()

	if err := cargobump.New().ExecuteContext(ctx); err != nil {
		log.Fatalf("error during command execution: %v", err)
	}
}
