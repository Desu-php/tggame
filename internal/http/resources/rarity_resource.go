package resources

import (
	"example.com/v2/internal/http/responses"
	"example.com/v2/internal/models"
)

type RarityResource struct {
}

func NewRarityResource() *RarityResource {
	return &RarityResource{}
}

func (r *RarityResource) Map(object *models.Rarity) *responses.RarityResponse {
	return &responses.RarityResponse{
		ID:          object.ID,
		Name:        object.Name,
		Color:       object.Color,
		Description: object.Description,
	}
}
