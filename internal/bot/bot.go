package bot

import (
	"context"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"ai-search-bot/internal/config"
)

func Run(ctx context.Context, cfg *config.Config, logger *slog.Logger) error {
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		return fmt.Errorf("failed to create bot: %w", err)
	}

	bot.Debug = cfg.Server.Debug
	logger.Info("bot authorized", "username", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			logger.Info("shutting down", "reason", ctx.Err())
			bot.StopReceivingUpdates()
			return nil
		case update := <-updates:
			if update.Message == nil || update.Message.Text == "" {
				continue
			}
			logger.Debug("message received",
				"from", update.Message.From.UserName,
				"text", update.Message.Text,
			)
			reply := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			reply.ReplyToMessageID = update.Message.MessageID
			if _, err := bot.Send(reply); err != nil {
				logger.Error("send failed", "err", err)
			}
		}
	}
}
