package usecase

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedisRateLimiter_Allow(t *testing.T) {
	// Start a miniredis server
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("an error '%s' occurred when starting miniredis", err)
	}
	defer mr.Close()

	// Connect to miniredis using go-redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	limiter := NewRedisRateLimiter(rdb)
	ctx := context.Background()

	t.Run("Allow with no existing key", func(t *testing.T) {
		allowed := limiter.Allow(ctx, "user@example.com", "Status")
		assert.True(t, allowed)
	})

	t.Run("Allow with existing key under limit", func(t *testing.T) {
		mr.Set("user@example.com:Status", "1")
		allowed := limiter.Allow(ctx, "user@example.com", "Status")
		assert.True(t, allowed)
	})

	t.Run("Disallow with key over limit", func(t *testing.T) {
		mr.Set("user@example.com:Status", "2")
		allowed := limiter.Allow(ctx, "user@example.com", "Status")
		assert.False(t, allowed)
	})

	t.Run("Redis error handling", func(t *testing.T) {
		mr.Close() // Simulate a Redis failure
		allowed := limiter.Allow(ctx, "user@example.com", "Status")
		assert.False(t, allowed)
	})
}
