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
		ID:        object.ID,
		Item:      NewItemResource(r.image).Map(&object.Item),
		CreatedAt: time.Time{},
	}
}
