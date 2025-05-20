package repository

import (
	"context"
	"example.com/v2/internal/models"
	"example.com/v2/pkg/db"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	Progress(ctx context.Context, user *models.User, taskType models.TaskType, progress uint) error
	FindUserTask(ctx context.Context, user *models.User, userTaskID uint) (*models.UserTask, error)
	Complete(ctx context.Context, userTaskID uint) error
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

func (repo *taskRepository) Progress(ctx context.Context, user *models.User, taskType models.TaskType, progress uint) error {
	var tasks []models.Task

	if err := repo.db.WithContext(ctx).
		Where("type = ?", taskType).
		Find(&tasks).Error; err != nil {
		return fmt.Errorf("taskRepository::Progress %w", err)
	}

	if len(tasks) == 0 {
		return nil
	}

	today := time.Now().Truncate(24 * time.Hour)
	userTasks := make([]models.UserTask, 0, len(tasks))

	for _, task := range tasks {
		userTasks = append(userTasks, models.UserTask{
			UserID:   user.ID,
			TaskID:   task.ID,
			Date:     today,
			Progress: progress,
		})
	}

	err := repo.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "task_id"}, {Name: "user_id"}, {Name: "date"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"progress": gorm.Expr("user_tasks.progress + EXCLUDED.progress"),
		}),
	}).Create(&userTasks).Error

	if err != nil {
		return fmt.Errorf("taskRepository::Progress %w", err)
	}

	return nil
}

func (repo *taskRepository) FindUserTask(ctx context.Context, user *models.User, userTaskID uint) (*models.UserTask, error) {
	var userTask models.UserTask

	err := repo.db.WithContext(ctx).Model(&models.UserTask{}).
		Preload("Task").
		Where("user_id = ? AND id = ?", user.ID, userTaskID).
		First(&userTask).Error

	if err != nil {
		return nil, fmt.Errorf("taskRepository::FindUserTask %w", err)
	}

	return &userTask, nil
}

func (repo *taskRepository) Complete(ctx context.Context, userTaskID uint) error {
	err := repo.db.WithContext(ctx).Model(&models.UserTask{}).
		Where("id = ?", userTaskID).
		Update("completed_at", time.Now()).Error

	if err != nil {
		return fmt.Errorf("taskRepository::Complete %w", err)
	}

	return nil
}
