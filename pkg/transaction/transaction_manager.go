package transaction

import (
	"context"
)

type TransactionFunc func(ctx context.Context) error

type TransactionManager interface {
	RunInTransaction(ctx context.Context, fn TransactionFunc) error
}
