package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"ai-search-bot/internal/app"
	"ai-search-bot/internal/config"
	"ai-search-bot/internal/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	log, cleanup, err := logger.New(env, cfg.Log)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger error: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	a, err := app.New(cfg, log)
	if err != nil {
		log.Error("failed to create app", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := a.Run(ctx); err != nil {
		log.Error("app stopped with error", "err", err)
		os.Exit(1)
	}
}
