package resources

import (
	"example.com/v2/internal/http/responses"
	"example.com/v2/pkg/image"
	"example.com/v2/pkg/utils"
)

type AspectWithStatsResource struct {
	image *image.Image
}

func NewAspectWithStatsResource(image *image.Image) *AspectWithStatsResource {
	return &AspectWithStatsResource{image}
}

func (r *AspectWithStatsResource) Map(object *responses.AspectWithStatsResponse) *responses.AspectWithStatsResponse {
	object.Image = r.image.Url(object.Image)

	if object.UserLevel != 0 {
		object.Amount = uint(utils.GrowthIncrease(float64(object.Amount), object.AmountGrowthFactor))
	}

	return object
}
