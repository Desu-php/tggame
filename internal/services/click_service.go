package services

import (
	"fmt"
	"example.com/v2/internal/adapter"
)

type ClickService struct {
	cache adapter.UserClickCacheAdapter
}

func NewClickService(cache adapter.UserClickCacheAdapter) *ClickService {
	return &ClickService{
		cache: cache,
	}
}

func (s *ClickService) Store(userId uint, count uint) error {
	err := s.cache.Increment(userId, count)

	if err != nil {
		return fmt.Errorf("ClickService::store %w", err)
	}

	return nil
}
