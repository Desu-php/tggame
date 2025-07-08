package commands

import (
	"example.com/v2/internal/models"
	"example.com/v2/internal/services/notification"
	"example.com/v2/pkg/db"
	"fmt"
)

type DailyRewardReminderCommand struct {
	telegramChannel *notification.TelegramChannel
	db              *db.DB
}

func NewDailyRewardReminderCommand(
	telegramChannel *notification.TelegramChannel,
	db *db.DB,
) *DailyRewardReminderCommand {
	return &DailyRewardReminderCommand{
		telegramChannel,
		db,
	}
}

func (c *DailyRewardReminderCommand) Execute() error {
	channels := []notification.Channel{
		c.telegramChannel,
	}

	var users []models.User

	c.db.DB.Model(&models.User{}).Where("daily_reward_sent_at IS NULL").Find(&users)

	err := notification.NewService(channels).NotifyAll("Test", &notification.Receiver{ID: "5326644263"})

	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	return nil
}
