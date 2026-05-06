package bot

import (
	"context"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"ai-search-bot/internal/config"
	"ai-search-bot/internal/handler"
)

type Handler interface {
	HandleCommand(cmd string, msg *tgbotapi.Message) (string, error)
	HandleMessage(msg *tgbotapi.Message) (string, error)
}

type Bot struct {
	api     *tgbotapi.BotAPI
	logger  *slog.Logger
	handler Handler
}

func New(cfg *config.Config, logger *slog.Logger, h Handler) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}
	api.Debug = cfg.Server.Debug
	logger.Info("bot authorized", "username", api.Self.UserName)

	b := &Bot{
		api:     api,
		logger:  logger,
		handler: h,
	}

	if err := b.registerCommands(); err != nil {
		return nil, fmt.Errorf("failed to register commands: %w", err)
	}

	return b, nil
}

func (b *Bot) registerCommands() error {
	commands := []tgbotapi.BotCommand{
		{Command: handler.CommandSearch, Description: "Поиск ресурсов в интернете"},
	}
	_, err := b.api.Request(tgbotapi.NewSetMyCommands(commands...))
	return err
}

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			b.logger.Info("shutting down", "reason", ctx.Err())
			b.api.StopReceivingUpdates()
			return nil
		case update := <-updates:
			if update.Message == nil || update.Message.Text == "" {
				continue
			}
			b.dispatch(update.Message)
		}
	}
}

func (b *Bot) dispatch(msg *tgbotapi.Message) {
	b.logger.Debug("message received", "from", msg.From.UserName, "text", msg.Text)

	var (
		text string
		err  error
	)

	if msg.IsCommand() {
		text, err = b.handler.HandleCommand(msg.Command(), msg)
	} else {
		text, err = b.handler.HandleMessage(msg)
	}

	if err != nil {
		b.logger.Error("handler error", "err", err)
		text = "Произошла ошибка. Попробуй ещё раз."
	}

	if text == "" {
		return
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ReplyToMessageID = msg.MessageID

	if _, err := b.api.Send(reply); err != nil {
		b.logger.Error("send failed", "err", err)
	}
}
