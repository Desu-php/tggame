package services

import (
	"fmt"
	"math/rand"
	"time"

	"example.com/v2/internal/models"
	"example.com/v2/internal/repository/item"
)

type RarityService struct {
	rarityRepository repository.RarityRepository
}

func NewRarityService(rarityRepository repository.RarityRepository) *RarityService {
	return &RarityService{rarityRepository: rarityRepository}
}

func (s *RarityService) GetRandom() (*models.Rarity, error) {
	rarities, err := s.rarityRepository.GetAll()

	if err != nil {
		return nil, fmt.Errorf("RarityService::GetRandom %w", err)
	}

	totalWeight := 0
	for _, rarity := range rarities {
		totalWeight += rarity.DropWeight
	}

	randomWeight := getRandomWeight(totalWeight)

	var selectedRarity models.Rarity

	currentWeight := 0
	for _, rarity := range rarities {
		currentWeight += rarity.DropWeight
		if randomWeight < currentWeight {
			selectedRarity = rarity
			break
		}
	}

	return &selectedRarity, nil
}

func getRandomWeight(max int) int {
	source := rand.NewSource(time.Now().UnixNano())
	localRand := rand.New(source)

	return localRand.Intn(max)
}
