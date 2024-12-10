package services

import (
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

func (s *ItemService) GetRandomItem() (*models.Item, error) {
	rarity, err := s.rarityService.GetRandom()

	if err != nil {
		return nil, fmt.Errorf("ItemService::GetRandomItem %w", err)
	}

	items, err := s.itemRepository.GetAllByRarity(rarity)

	if err != nil {
		return nil, fmt.Errorf("ItemService::GetRandomItem %w", err)
	}

	selectedItem := items[rand.Intn(len(items))]

	return &selectedItem, nil
}
