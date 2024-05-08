package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jordanlanch/modak-test/domain"
)

type RateLimitRule struct {
	Limit    int64
	Duration time.Duration
}

type RedisRateLimiter struct {
	rdb   *redis.Client
	rules map[string]RateLimitRule
}

func NewRedisRateLimiter(rdb *redis.Client) domain.RateLimiter {
	return &RedisRateLimiter{
		rdb: rdb,
		rules: map[string]RateLimitRule{
			"Status":    {Limit: 2, Duration: time.Minute},
			"News":      {Limit: 1, Duration: 24 * time.Hour},
			"Marketing": {Limit: 3, Duration: time.Hour},
		},
	}
}

func (rl *RedisRateLimiter) Allow(ctx context.Context, recipient, messageType string) bool {
	rule, exists := rl.rules[messageType]
	if !exists {
		return true
	}

	key := fmt.Sprintf("rate_limit:%s:%s", recipient, messageType)

	newCount, err := rl.rdb.Incr(ctx, key).Result()
	if err != nil {
		log.Printf("Redis error: %v", err)
		return false
	}

	if newCount == 1 {
		if _, err := rl.rdb.Expire(ctx, key, rule.Duration).Result(); err != nil {
			log.Printf("Failed to set expiration for %s: %v", key, err)
			return false
		}
	}

	if newCount > rule.Limit {
		return false
	}

	return true
}
