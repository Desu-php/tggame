package repository

import (
	"context"
	"errors"
	"example.com/v2/internal/http/responses"
	"example.com/v2/internal/models"
	"example.com/v2/pkg/db"
	"fmt"
	"gorm.io/gorm"
)

type AspectRepository interface {
	GetUserAspectByType(ctx context.Context, user *models.User, aspectType models.AspectType) ([]responses.AspectWithStatsResponse, error)
	FindByID(ctx context.Context, id uint, aspectType models.AspectType) (*models.Aspect, error)
}

type aspectRepository struct {
	db *db.DB
}

func NewAspectRepository(db *db.DB) AspectRepository {
	return &aspectRepository{db: db}
}

func (a *aspectRepository) GetUserAspectByType(ctx context.Context, user *models.User, aspectType models.AspectType) ([]responses.AspectWithStatsResponse, error) {
	var aspects []responses.AspectWithStatsResponse

	err := a.db.WithContext(ctx).
		Raw(`
SELECT
  a.id,
  a.name,
  a.image,
  a.description,
  ua.level AS user_level,
  ast.id as aspect_stat_id,
  ast.damage,
  ast.critical_damage,
  ast.critical_chance,
  ast.gold_multiplier,
  ast.passive_damage,
  CASE
    WHEN ua.amount IS NULL THEN ast.amount
    WHEN ast.amount IS NULL THEN ua.amount
    ELSE GREATEST(ua.amount, ast.amount)
  END AS amount,
  ast.amount_growth_factor
FROM aspects a
LEFT JOIN user_aspects ua ON ua.aspect_id = a.id AND ua.user_id = ?
LEFT JOIN LATERAL (
  SELECT *
  FROM aspect_stats ast
  WHERE ast.aspect_id = a.id
    AND (
      (ua.level + 1 >= ast.start_level AND ua.level + 1 <= ast.end_level)
      OR ua.level IS NULL
    )
  ORDER BY
    CASE
      WHEN ua.level + 1 BETWEEN ast.start_level AND ast.end_level THEN 0
      ELSE 1
    END,
    start_level
  LIMIT 1
) ast ON true
WHERE a.type = ?
`, user.ID, aspectType).Scan(&aspects).Error

	if err != nil {
		return nil, fmt.Errorf("AspectRepository::GetUserAspectByType %w", err)

	}

	return aspects, nil
}

func (a *aspectRepository) FindByID(ctx context.Context, id uint, aspectType models.AspectType) (*models.Aspect, error) {
	var aspect models.Aspect

	err := a.db.WithContext(ctx).
		Where("type = ?", aspectType).
		First(&aspect, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("AspectRepository::FindByID %w", err)
	}

	return &aspect, nil
}
