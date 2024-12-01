package adapter

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

const prefix = "clicks:user"

type ChunkFunction func(clicks map[uint]uint) error

type UserClickCacheAdapter interface {
	Increment(userId uint, count uint) error
	Get(userId uint) (uint, error)
	ChunkAll(count int64, fn ChunkFunction) error
}

type userClickCacheAdapter struct {
	redis *redis.Client
}

func NewUserClickCacheAdapter(redis *redis.Client) UserClickCacheAdapter {
	log.Println("UserClickCacheAdapter")
	return &userClickCacheAdapter{redis: redis}
}

func (a *userClickCacheAdapter) Increment(userId uint, count uint) error {
	countBy := int64(count)

	log.Println("countBy", countBy)

	result := a.redis.IncrBy(context.Background(), generateCacheKey(userId), countBy)

	log.Println("result", result.Err())

	if err := result.Err(); err != nil {
		return fmt.Errorf("UserClickCacheAdapter::Increment failed to increment clicks for user %d: %w", userId, err)
	}

	return nil
}

func (a *userClickCacheAdapter) Get(userId uint) (uint, error) {
	count, err := a.redis.Get(context.Background(), generateCacheKey(userId)).Result()

	if err != nil {
		return 0, fmt.Errorf("UserClickCacheAdapter::Get failed to get clicks for user %d: %w", userId, err)
	}

	clicksCount, err := strconv.ParseUint(count, 10, 32)

	if err != nil {
		return 0, fmt.Errorf("UserClickCacheAdapter::Get failed to parse clicks for user %d: %w", userId, err)
	}

	return uint(clicksCount), nil
}

func (a *userClickCacheAdapter) ChunkAll(count int64, fn ChunkFunction) error {
	var cursor uint64
	clicks := make(map[uint]uint)

	for {
		keys, nextCursor, err := a.redis.Scan(context.Background(), cursor, fmt.Sprintf("%s:*", prefix), count).Result()

		if err != nil {
			return fmt.Errorf("UserClickCacheAdapter::GetAll %w", err)
		}

		for _, key := range keys {
			parts := strings.Split(key, ":")

			if len(parts) < 3 {
				log.Println("continue", parts)
				continue
			}

			userId, err := strconv.ParseUint(parts[2], 10, 64)

			if err != nil {
				log.Println("err parse", parts)
				return fmt.Errorf("failed to parse userId from key %s: %w", key, err)
			}

			value, err := a.Get(uint(userId))

			if err != nil {
				if err == redis.Nil {
					log.Println("nil", userId)
					continue
				}
				log.Println("err getting", parts)
				return fmt.Errorf("failed to get value for key %s: %w", key, err)
			}

			clicks[uint(userId)] = uint(value)

			if len(clicks) >= int(count) {
				fn(clicks)
				clicks = make(map[uint]uint)
			}
		}

		if nextCursor == 0 {
			break
		}

		cursor = nextCursor
	}

	if len(clicks) > 0 {
		fn(clicks)
	}

	return nil
}

func generateCacheKey(userId uint) string {
	return fmt.Sprintf("%s:%d", prefix, userId)
}
