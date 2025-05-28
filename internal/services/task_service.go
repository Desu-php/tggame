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

type TaskService struct {
	repository     repository.TaskRepository
	balanceService *BalanceService
	trx            transaction.TransactionManager
}

type ProgressOneTimeDto struct {
	User     *models.User
	TaskID   uint
	TaskType models.TaskType
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

func (t *TaskService) Progress(ctx context.Context, dto *repository.TaskProgressDto) error {
	err := t.repository.Progress(ctx, dto)

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

func (t *TaskService) ProgressOnTimeTask(ctx context.Context, dto *ProgressOneTimeDto) (bool, error) {
	var ok bool

	err := t.trx.RunInTransaction(ctx, func(ctx context.Context) error {
		userTask, err := t.repository.FindTaskById(ctx, dto.User, dto.TaskID)

		if err != nil {
			ok = false
			return err
		}

		if userTask.CompletedAt != nil {
			ok = false
			return nil
		}

		if userTask.Progress >= userTask.Task.TargetValue {
			ok = true
			return nil
		}

		err = t.repository.ProgressTask(ctx, userTask, 1)

		if err != nil {
			ok = false
			return err
		}

		ok = true
		return nil
	})

	if err != nil {
		return ok, fmt.Errorf("TaskService::ProgressOnTimeTask %w", err)
	}

	return ok, nil
}
