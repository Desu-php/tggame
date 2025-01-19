package resources

import (
	"example.com/v2/internal/http/responses"
	"example.com/v2/internal/models"
)

type ReferralUserResource struct {
}

func NewReferralUserResource() *ReferralUserResource {
	return &ReferralUserResource{}
}

func (r *ReferralUserResource) Map(object *models.ReferralUser) *responses.ReferralUserResponse {
	return &responses.ReferralUserResponse{
		Username:  object.ReferredUser.Username,
		CreatedAt: object.CreatedAt,
	}
}
