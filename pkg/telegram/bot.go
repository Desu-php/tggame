package telegram

import (
	"bytes"
	"encoding/json"
	"example.com/v2/config"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Bot struct {
	token string
}

func NewBot(token string) *Bot {
	return &Bot{token: token}
}

func ProvideTelegramBot(cfg *config.Config) *Bot {
	return NewBot(cfg.Telegram.NotificationToken)
}

func (b *Bot) SendMessage(message string, chatId uint64) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.token)

	log.Println(url, message, chatId)

	payload := map[string]interface{}{
		"chat_id": chatId,
		"text":    message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code from Telegram API: %d, response: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
