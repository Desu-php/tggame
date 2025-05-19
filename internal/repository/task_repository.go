package repository

import (
	"context"
	"example.com/v2/internal/models"
	"example.com/v2/pkg/db"
	"fmt"
	"time"
)

type UserTask struct {
	ID          uint64          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	TargetValue int             `json:"target_value"`
	Progress    int             `json:"progress"`
	CompletedAt *time.Time      `json:"completed_at"`
	Type        models.TaskType `json:"type"`
}

type TaskRepository interface {
	GetAll(ctx context.Context, user *models.User) ([]UserTask, error)
}

type taskRepository struct {
	db *db.DB
}

func NewTaskRepository(db *db.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (repo *taskRepository) GetAll(ctx context.Context, user *models.User) ([]UserTask, error) {
	var tasks []UserTask
	err := repo.db.WithContext(ctx).
		Table("tasks").
		Select(`
			tasks.id,
			tasks.name,
			tasks.description,
			tasks.target_value,
			COALESCE(ut.progress, 0) AS progress,
			ut.completed_at,
			tasks.type
		`).
		Joins(`
			LEFT JOIN user_tasks ut 
				ON ut.task_id = tasks.id 
				AND ut.user_id = ? 
				AND ut.date = CURRENT_DATE
		`, user.ID).
		Scan(&tasks).Error

	if err != nil {
		return nil, fmt.Errorf("taskRepository::GetAll %w", err)
	}

	return tasks, nil
}
