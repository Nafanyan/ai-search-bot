package handler

import (
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"ai-search-bot/internal/services"
)

type MessageHandler struct {
	search         services.SearchService
	searchingChats map[int64]bool
	mu             sync.Mutex
}

func NewMessageHandler(search services.SearchService) *MessageHandler {
	return &MessageHandler{
		search:         search,
		searchingChats: make(map[int64]bool),
	}
}

func (h *MessageHandler) HandleCommand(cmd string, msg *tgbotapi.Message) (string, error) {
	switch cmd {
	case CommandSearch:
		h.mu.Lock()
		h.searchingChats[msg.Chat.ID] = true
		h.mu.Unlock()
		return "Я переключился в режим поиска ресурсов в интернете. Напиши, что ты хочешь найти.", nil
	}
	return "", nil
}

func (h *MessageHandler) HandleMessage(msg *tgbotapi.Message) (string, error) {
	h.mu.Lock()
	searching := h.searchingChats[msg.Chat.ID]
	if searching {
		delete(h.searchingChats, msg.Chat.ID)
	}
	h.mu.Unlock()

	if searching {
		return h.search.Search(msg.Text)
	}

	return msg.Text, nil
}
