package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/poupardm-GhostWrath/GoAdventure/internal/config"
)

func main() {
	fmt.Println("Hello World")
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	status := run(ctx, cancel)
	cancel()
	defer func() {
		if err := config.Cfg.DB.Close(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close db connection: %v\n", err)
		}
	}()
	os.Exit(status)
}

func run(ctx context.Context, cancel context.CancelFunc) int {
	<-ctx.Done()
	_, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return 0
}
