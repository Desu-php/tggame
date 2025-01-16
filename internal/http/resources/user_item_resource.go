package resources

import (
	"example.com/v2/internal/http/responses"
	"example.com/v2/internal/models"
	"example.com/v2/pkg/image"
	"time"
)

type UserItemResource struct {
	image *image.Image
}

func NewUserItemResource(image *image.Image) *UserItemResource {
	return &UserItemResource{image}
}

func (r *UserItemResource) Map(object *models.UserItem) *responses.UserItemResponse {
	return &responses.UserItemResponse{
		ID: object.ID,
		Item: &responses.ItemResponse{
			ID:     object.Item.ID,
			Name:   object.Item.Name,
			Rarity: NewRarityResource().Map(&object.Item.Rarity),
			Image:  r.image.Url(object.Item.Image),
			Type: &responses.ItemTypeResponse{
				ID:          object.Item.Type.ID,
				Name:        object.Item.Type.Name,
				Description: object.Item.Type.Description,
			},
			Description: object.Item.Description,
		},
		CreatedAt: time.Time{},
	}
}
