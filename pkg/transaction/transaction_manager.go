package transaction

type TransactionFunc func() error

type TransactionManager interface {
	RunInTransaction(fn TransactionFunc) error
}
