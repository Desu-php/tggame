package pkg

import (
	"example.com/v2/pkg/image"
	"example.com/v2/pkg/logging"
	"example.com/v2/pkg/telegram"
	"example.com/v2/pkg/transaction"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	transaction.NewTransactionManager,
	logging.NewLogger,
	image.NewImage,
	telegram.ProvideTelegramBot,
)
