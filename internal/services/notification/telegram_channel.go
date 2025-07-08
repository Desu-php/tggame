package notification

import (
	"example.com/v2/pkg/telegram"
	"fmt"
	"strconv"
)

type TelegramChannel struct {
	bot *telegram.Bot
}

func NewTelegramChannel(bot *telegram.Bot) *TelegramChannel {
	return &TelegramChannel{bot}
}

func (t *TelegramChannel) Send(message string, receiver *Receiver) error {

	id, err := strconv.ParseUint(receiver.ID, 10, 64)

	if err != nil {
		return fmt.Errorf("TelegramChannel invalid receiver ID: %v", err)
	}

	err = t.bot.SendMessage(message, id)

	if err != nil {
		return fmt.Errorf("TelegramChannel::Send: %w", err)
	}

	return nil
}
