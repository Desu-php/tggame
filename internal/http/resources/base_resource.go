package resources

type Mapper[T any, R any] interface {
	Map(object *T) *R
}

type BaseResource[T any, R any] struct {
	resource Mapper[T, R]
}

func NewBaseResource[T any, R any](resource Mapper[T, R]) *BaseResource[T, R] {
	return &BaseResource[T, R]{resource: resource}
}

func (r *BaseResource[T, R]) One(object *T) *R {
	return r.resource.Map(object)
}

func (r *BaseResource[T, R]) All(objects []T) []*R {
	var responses []*R

	for _, value := range objects {
		mapped := r.resource.Map(&value)
		responses = append(responses, mapped)
	}

	return responses
}
