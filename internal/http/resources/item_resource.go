package resources

import (
	"example.com/v2/internal/http/responses"
	"example.com/v2/internal/models"
	"example.com/v2/pkg/image"
)

type ItemResource struct {
	image *image.Image
}

func NewItemResource(image *image.Image) *ItemResource {
	return &ItemResource{image}
}

func (r *ItemResource) Map(object *models.Item) *responses.ItemResponse {
	return &responses.ItemResponse{
		ID:     object.ID,
		Name:   object.Name,
		Rarity: NewRarityResource().Map(&object.Rarity),
		Image:  r.image.Url(object.Image),
		Type: &responses.ItemTypeResponse{
			ID:          object.Type.ID,
			Name:        object.Type.Name,
			Description: object.Type.Description,
		},
		Description:    object.Description,
		Damage:         object.Damage,
		CriticalDamage: object.CriticalDamage,
		CriticalChance: object.CriticalChance,
		GoldMultiplier: object.GoldMultiplier,
		PassiveDamage:  object.PassiveDamage,
	}
}
