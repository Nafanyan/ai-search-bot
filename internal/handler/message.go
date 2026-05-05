package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"ai-search-bot/internal/services"
	"ai-search-bot/internal/storage"
)

type MessageHandler struct {
	search  services.SearchService
	storage storage.ChatStorage
}

func NewMessageHandler(search services.SearchService, store storage.ChatStorage) *MessageHandler {
	return &MessageHandler{
		search:  search,
		storage: store,
	}
}

func (h *MessageHandler) HandleCommand(cmd string, msg *tgbotapi.Message) (string, error) {
	switch cmd {
	case CommandSearch:
		h.storage.SetState(msg.Chat.ID, storage.StateSearching)
		return "Я переключился в режим поиска ресурсов в интернете. Напиши, что ты хочешь найти.", nil
	}
	return "", nil
}

func (h *MessageHandler) HandleMessage(msg *tgbotapi.Message) (string, error) {
	if h.storage.GetState(msg.Chat.ID) == storage.StateSearching {
		h.storage.SetState(msg.Chat.ID, storage.StateIdle)
		return h.search.Search(msg.Text)
	}

	return msg.Text, nil
}
