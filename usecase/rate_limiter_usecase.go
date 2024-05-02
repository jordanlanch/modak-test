package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/jordanlanch/modak-test/domain"
)

type RedisRateLimiter struct {
	rdb *redis.Client
}

func NewRedisRateLimiter(rdb *redis.Client) domain.RateLimiter {
	return &RedisRateLimiter{rdb}
}

func (rl *RedisRateLimiter) Allow(ctx context.Context, recipient, messageType string) bool {
	key := fmt.Sprintf("%s:%s", recipient, messageType)
	result, err := rl.rdb.Get(ctx, key).Int()
	if err == redis.Nil {
		// If no key exists, no messages have been sent yet.
		return true
	} else if err != nil {
		// Handle Redis errors appropriately.
		log.Printf("Redis error: %v", err)
		return false
	}

	// Check if the result exceeds the limit.
	if result >= rl.getLimit(messageType) {
		// Here, rather than just returning false, you might consider returning an error or setting an error state.
		return false
	}

	return true
}

func (rl *RedisRateLimiter) getLimit(messageType string) int {
	switch messageType {
	case "Status":
		return 2
	case "News":
		return 1
	case "Marketing":
		return 3
	default:
		return 5 // Default limit
	}
}
