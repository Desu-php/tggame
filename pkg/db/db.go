package db

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type txKey string

var ctxWithTx = txKey("tx")

type Options struct {
	withoutPrometheus bool
}

type GormOptions func(o *Options)

type DB struct {
	*gorm.DB
}

func NewDB(dsn string, options ...GormOptions) (*DB, error) {
	setupGorm, err := setupGorm(dsn, options...)

	if err != nil {
		return nil, err
	}

	return &DB{setupGorm}, nil
}

func (d *DB) WithContext(ctx context.Context) *gorm.DB {
	tx, ok := ExtractTxFromContext(ctx)

	if ok {
		return tx.WithContext(ctx)
	}

	return d.DB.WithContext(ctx)
}

func setupGorm(dsn string, options ...GormOptions) (*gorm.DB, error) {
	o := &Options{}
	for _, f := range options {
		f(o)
	}
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

func ExtractTxFromContext(ctx context.Context) (*gorm.DB, bool) {
	tx := ctx.Value(ctxWithTx)

	if t, ok := tx.(*gorm.DB); ok {
		return t, true
	}

	return nil, false
}

func (d *DB) PutTxToContext(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxWithTx, tx)
}
