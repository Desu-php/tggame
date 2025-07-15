package services

import (
	"context"
	"fmt"
	"math/rand"

	"example.com/v2/internal/models"
	"example.com/v2/internal/repository/item"
)

type ItemService struct {
	rarityService  *RarityService
	itemRepository repository.ItemRepository
}

func NewItemService(
	rarityService *RarityService,
	itemRepository repository.ItemRepository,
) *ItemService {
	return &ItemService{
		rarityService:  rarityService,
		itemRepository: itemRepository,
	}
}

func (s *ItemService) GetRandomItem(ctx context.Context, minRarity *models.Rarity, maxRarity *models.Rarity) (*models.Item, error) {
	rarity, err := s.rarityService.GetRandom(ctx, minRarity, maxRarity)
	if err != nil {
		return nil, fmt.Errorf("ItemService::GetRandomItem %w", err)
	}

	items, err := s.itemRepository.GetAllByRarity(ctx, rarity)
	if err != nil {
		return nil, fmt.Errorf("ItemService::GetRandomItem %w", err)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("ItemService::GetRandomItem: no items found with rarity %v", rarity)
	}

	// Вычисляем суммарный шанс
	var totalChance float32
	for _, item := range items {
		totalChance += item.DropChance
	}

	if totalChance == 0 {
		return nil, fmt.Errorf("ItemService::GetRandomItem: total drop chance is zero")
	}

	// Генерируем случайное число в пределах от 0 до totalChance
	randVal := rand.Float32() * totalChance

	// Выбираем предмет по весу
	var cumulative float32
	for _, item := range items {
		cumulative += item.DropChance
		if randVal <= cumulative {
			return &item, nil
		}
	}

	// Фолбэк (на случай ошибок округления)
	return &items[len(items)-1], nil
}
