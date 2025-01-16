package resources

import (
	"example.com/v2/internal/http/responses"
	"example.com/v2/internal/models"
	"example.com/v2/pkg/image"
)

type UserResource struct {
	image *image.Image
}

func NewUserResource(image *image.Image) *UserResource {
	return &UserResource{image}
}

func (r *UserResource) Map(object *models.User) *responses.UserResponse {
	return &responses.UserResponse{
		ID:         object.ID,
		Username:   object.Username,
		Session:    object.Session,
		TelegramId: object.TelegramID,
		UserChest:  NewUserChestResource(r.image).Map(&object.UserChest),
	}
}
