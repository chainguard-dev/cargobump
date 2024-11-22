/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/chainguard-dev/cargobump/cmd/cargobump"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt)
	defer done()

	if err := cargobump.New().ExecuteContext(ctx); err != nil {
		log.Fatalf("error during command execution: %v", err) //nolint: gocritic
	}
}
