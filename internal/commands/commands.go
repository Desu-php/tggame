package commands

import "go.uber.org/fx"

var Module = fx.Provide(
	NewDailyRewardReminderCommand,
)
