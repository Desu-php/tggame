package resources

import (
	"example.com/v2/internal/http/responses"
	"example.com/v2/internal/models"
	"example.com/v2/pkg/image"
)

type UserChestResource struct {
	image *image.Image
}

func NewUserChestResource(image *image.Image) *UserChestResource {
	return &UserChestResource{image}
}

func (r *UserChestResource) Map(object *models.UserChest) *responses.UserChestResponse {
	return &responses.UserChestResponse{
		ID:            object.ID,
		Health:        object.Health,
		CurrentHealth: object.CurrentHealth,
		Level:         object.Level,
		Amount:        object.Amount,
		Chest: &responses.ChestResponse{
			ID:     object.Chest.ID,
			Name:   object.Chest.Name,
			Image:  r.image.Url(object.Chest.Image),
			Rarity: NewRarityResource().Map(&object.Chest.Rarity),
		},
	}
}
