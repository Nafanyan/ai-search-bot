package app

import (
	"context"
	"log/slog"

	"ai-search-bot/internal/bot"
	"ai-search-bot/internal/config"
	"ai-search-bot/internal/handler"
	"ai-search-bot/internal/services"
	"ai-search-bot/internal/storage"
)

type App struct {
	bot *bot.Bot
}

func New(cfg *config.Config, log *slog.Logger) (*App, error) {
	searchService := &services.MockSearchService{}
	chatStorage := storage.NewInMemoryChatStorage()
	h := handler.NewMessageHandler(searchService, chatStorage)

	b, err := bot.New(cfg, log, h)
	if err != nil {
		return nil, err
	}

	return &App{bot: b}, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.bot.Run(ctx)
}
