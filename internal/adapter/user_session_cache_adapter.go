package adapter

import (
	"fmt"

	"example.com/v2/internal/repository"
)

type UserSessionAdapter interface {
	Get(key uint64) (string ,error)
}

type userSessionAdapter struct {
	repo repository.UserRepository
}

func NewUserSessionCacheAdapter(repo repository.UserRepository) UserSessionAdapter{
	return &userSessionAdapter{repo: repo}
}

func(u *userSessionAdapter) Get(key uint64) (string,error) {
   user, err :=	u.repo.FindByTgId(key)

   if err != nil {
		return "", fmt.Errorf("Get: %w", err)
   }

   if user == nil {
		return "", nil
   }

   return user.Session, nil
}