package resources

import (
	"example.com/v2/internal/http/responses"
	"example.com/v2/pkg/image"
)

type AspectWithStatsResource struct {
	image *image.Image
}

func NewAspectWithStatsResource(image *image.Image) *AspectWithStatsResource {
	return &AspectWithStatsResource{image}
}

func (r *AspectWithStatsResource) Map(object *responses.AspectWithStatsResponse) *responses.AspectWithStatsResponse {
	object.Image = r.image.Url(object.Image)

	return object
}
