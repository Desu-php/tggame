package resources

import (
	"example.com/v2/internal/http/responses"
	"example.com/v2/internal/repository"
	"example.com/v2/pkg/image"
)

type GroupedUserItemResource struct {
	image *image.Image
}

func NewGroupedUserItemResource(image *image.Image) *GroupedUserItemResource {
	return &GroupedUserItemResource{image}
}

func (r *GroupedUserItemResource) Map(object *repository.GroupedUserItem) *responses.GroupedUserItemResponse {
	return &responses.GroupedUserItemResponse{
		ID:     object.ID,
		Name:   object.Name,
		Count:  object.Count,
		Type:   object.Type,
		Rarity: object.Rarity,
		Image:  r.image.Url(object.Image),
	}
}
