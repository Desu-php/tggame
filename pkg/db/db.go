package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Options struct {
	withoutPrometheus bool
}

type GormOptions func(o *Options)

func WithoutPrometheus() GormOptions {
	return func(o *Options) {
		o.withoutPrometheus = true
	}
}

func SetupGorm(dsn string, options ...GormOptions) (*gorm.DB, error) {
	o := &Options{}
	for _, f := range options {
		f(o)
	}
	gormDB, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}