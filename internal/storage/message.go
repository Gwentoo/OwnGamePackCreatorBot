package storage

import (
	tele "gopkg.in/telebot.v3"
	"sync"
)

type SavedMessage struct {
	Text     string
	Keyboard *tele.ReplyMarkup
}

var (
	userMessages = make(map[int64]SavedMessage) // chatID:message
	mu           sync.RWMutex
)

func SaveMessage(chatID int64, text string, keyboard *tele.ReplyMarkup) {
	mu.Lock()
	defer mu.Unlock()
	userMessages[chatID] = SavedMessage{text, keyboard}
}

func GetMessage(chatID int64) SavedMessage {
	mu.RLock()
	defer mu.RUnlock()
	msg, _ := userMessages[chatID]
	return msg
}
