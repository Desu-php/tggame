package services

import (
	"errors"
	"example.com/v2/internal/models"
	"example.com/v2/internal/repository"
	"example.com/v2/pkg/transaction"
	"fmt"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"time"
)

type TaskProgressDto struct {
	Progress uint
	Type     models.TaskType
	User     *models.User
}

type TaskService struct {
	repository     repository.TaskRepository
	balanceService *BalanceService
	trx            transaction.TransactionManager
}

func NewTaskService(
	taskRepository repository.TaskRepository,
	balanceService *BalanceService,
	trx transaction.TransactionManager,
) *TaskService {
	return &TaskService{
		taskRepository,
		balanceService,
		trx,
	}
}

func (t *TaskService) Progress(ctx context.Context, dto *TaskProgressDto) error {
	err := t.repository.Progress(ctx, dto.User, dto.Type, dto.Progress)

	if err != nil {
		return fmt.Errorf("TaskService::Progress %w", err)
	}

	return nil
}

func (t *TaskService) ReceiveReward(ctx context.Context, userTaskID uint, user *models.User) (*models.UserTask, error) {
	userTask, err := t.repository.FindUserTask(ctx, user, userTaskID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, fmt.Errorf("TaskService::ReceiveReward %w", err)
	}

	if userTask.Progress >= userTask.Task.TargetValue && userTask.CompletedAt == nil {
		err = t.trx.RunInTransaction(ctx, func(ctx context.Context) error {
			err = t.balanceService.Replenish(ctx, &TransactionDto{
				Amount: userTask.Task.Amount,
				User:   user,
				Model:  userTask,
				Type:   models.TransactionTypeTaskCompleted,
			})

			if err != nil {
				return err
			}

			err = t.repository.Complete(ctx, userTask.ID)

			if err != nil {
				return err
			}

			now := time.Now()
			userTask.CompletedAt = &now

			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("TaskService::ReceiveReward %w", err)
		}
	}

	return userTask, nil
}
