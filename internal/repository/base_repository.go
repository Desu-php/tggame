package repository

import (
	"gorm.io/gorm"
)

type BaseRepository struct {
	DB *gorm.DB
}


func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{DB: db}
}
