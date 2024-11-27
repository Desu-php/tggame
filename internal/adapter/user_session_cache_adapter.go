package adapter

import (
	"context"
	"fmt"
	"strconv"
	"github.com/redis/go-redis/v9"
)

type UserSessionAdapter interface {
	Get(key uint64) (string ,error)
	Set(key uint64, value string) error 
}

type userSessionAdapter struct {
	repo *redis.Client
}

func NewUserSessionCacheAdapter(repo *redis.Client) UserSessionAdapter{
	return &userSessionAdapter{repo: repo}
}

func(u *userSessionAdapter) Get(key uint64) (string,error) {
	session, err := u.repo.Get(context.Background(), strconv.FormatUint(key, 10)).Result()

   if err != nil {
		return "", fmt.Errorf("Get: %w", err)
   }

   return session, nil
}

func (u *userSessionAdapter) Set(key uint64, value string) error{
	err :=	u.repo.Set(context.Background(), strconv.FormatUint(key, 10), value, 0).Err()

	if err != nil {
		return fmt.Errorf("Set: err %w", err)
	}

	return nil
}