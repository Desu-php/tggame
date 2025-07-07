package services

import (
	"context"
	"example.com/v2/internal/models"
	"example.com/v2/internal/repository"
	"example.com/v2/pkg/utils"
	"fmt"
)

type ReferralService struct {
	balanceService     *BalanceService
	repository         repository.ReferralUserRepository
	userStatRepository repository.UserStatRepository
}

func NewReferralService(
	balance *BalanceService,
	repository repository.ReferralUserRepository,
	userStatRepository repository.UserStatRepository,
) *ReferralService {
	return &ReferralService{
		balanceService:     balance,
		repository:         repository,
		userStatRepository: userStatRepository,
	}
}

func (s *ReferralService) RewardForReferral(ctx context.Context, inviter *models.User, referredUser *models.User) error {
	countReferrals, err := s.repository.Count(ctx, inviter.ID)

	if err != nil {
		return fmt.Errorf("ReferralService::RewardForReferral %w", err)
	}

	reward := s.getAmount(countReferrals)

	userStat, err := s.userStatRepository.GetStat(ctx, inviter)

	if err != nil {
		return fmt.Errorf("ReferralService::RewardForReferral %w", err)
	}

	reward = int64(utils.GrowthIncrease(float64(reward), userStat.GoldMultiplier))

	err = s.balanceService.Replenish(ctx, &TransactionDto{
		Amount: reward,
		User:   inviter,
		Model:  referredUser,
		Type:   models.TransactionReferralReward,
	})

	if err != nil {
		return fmt.Errorf("ReferralService::RewardForReferral %w", err)
	}

	return nil
}

func (s *ReferralService) getAmount(count uint) int64 {
	var reward int64

	switch {
	case count < 10:
		reward = 10
	case count < 50:
		reward = 25
	case count < 100:
		reward = 50
	default:
		reward = 100
	}

	return reward
}
