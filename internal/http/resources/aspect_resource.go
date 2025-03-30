package resources

import (
	"example.com/v2/internal/http/responses"
	"example.com/v2/internal/models"
	"example.com/v2/pkg/image"
)

type AspectResource struct {
	image *image.Image
}

func NewAspectResource(image *image.Image) *AspectResource {
	return &AspectResource{image: image}
}

func (r *AspectResource) Map(object *models.Aspect) *responses.AspectResponse {
	return &responses.AspectResponse{
		ID:          object.ID,
		Name:        object.Name,
		Image:       r.image.Url(object.Image),
		Description: object.Description,
	}
}
