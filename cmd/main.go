package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"ai-search-bot/internal/bot"
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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := bot.Run(ctx, cfg, log); err != nil {
		log.Error("bot stopped with error", "err", err)
		os.Exit(1)
	}
}
