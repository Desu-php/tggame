package cron

import (
	"context"
	"example.com/v2/internal/commands"
	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
	"log"
)

func NewCron() *cron.Cron {
	return cron.New(cron.WithSeconds()) // поддержка секундами, если нужно
}

func StartTasks(
	lc fx.Lifecycle,
	c *cron.Cron,
	dailyRewardReminderCommand *commands.DailyRewardReminderCommand,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_, err := c.AddFunc("0 */1 * * * *", func() {
				log.Println("⏰ CRON: Проверка непринятых наград")

				err := dailyRewardReminderCommand.Execute()

				if err != nil {
					log.Println(err)
				}
			})

			if err != nil {
				return err
			}

			log.Println("✅ CRON: Запуск задач")
			c.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("🛑 CRON: Остановка")
			ctxStop := c.Stop()
			select {
			case <-ctxStop.Done():
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	})
}
