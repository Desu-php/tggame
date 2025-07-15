package commands

import (
	"example.com/v2/internal/models"
	"example.com/v2/internal/services/notification"
	"example.com/v2/pkg/db"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"strconv"
)

type DailyRewardReminderCommand struct {
	telegramChannel *notification.TelegramChannel
	db              *db.DB
	logger          *logrus.Logger
}

func NewDailyRewardReminderCommand(
	telegramChannel *notification.TelegramChannel,
	db *db.DB,
	logger *logrus.Logger,
) *DailyRewardReminderCommand {
	return &DailyRewardReminderCommand{
		telegramChannel,
		db,
		logger,
	}
}

func (c *DailyRewardReminderCommand) Execute() error {
	const chunkSize = 1000
	channels := []notification.Channel{
		c.telegramChannel,
	}

	notifyService := notification.NewService(channels)

	offset := 0

	for {
		var users []models.User

		err := c.db.DB.Model(&models.User{}).
			Select("users.*").
			Joins("inner join user_tasks ut on ut.user_id = users.id and ut.completed_at is null and ut.is_notified = false").
			Joins("inner join tasks t on t.id = ut.task_id and t.target_value <= ut.progress").
			Group("users.id").
			Limit(chunkSize).
			Offset(offset).
			Find(&users).Error

		if err != nil {
			return fmt.Errorf("failed to fetch users: %w", err)
		}

		if len(users) == 0 {
			break
		}

		var userIDs []uint
		for _, user := range users {
			log.Printf("user: %+v", user)
			err = notifyService.NotifyAll(
				"You have completed tasks! Hurry up and claim your reward.",
				&notification.Receiver{ID: strconv.FormatUint(user.TelegramID, 10)},
			)
			if err != nil {
				c.logger.WithError(err).Error("DailyRewardReminderCommand::Execute")
			}

			userIDs = append(userIDs, user.ID)
		}

		if len(userIDs) > 0 {
			if err = c.db.DB.Model(&models.UserTask{}).
				Where("user_id IN ? AND completed_at is null AND is_notified = false", userIDs).
				Update("is_notified", true).Error; err != nil {
				return fmt.Errorf("failed to update is_notified: %w", err)
			}
		}

		if len(users) < chunkSize {
			break
		}
		offset += chunkSize
	}

	return nil
}
