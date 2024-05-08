package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedisRateLimiter_Allow(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("an error '%s' occurred when starting miniredis", err)
	}
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	limiter := NewRedisRateLimiter(rdb)
	ctx := context.Background()

	t.Run("Allow Status message under limit", func(t *testing.T) {
		mr.FlushAll()
		allowed := limiter.Allow(ctx, "user@example.com", "Status")
		assert.True(t, allowed)

		allowed = limiter.Allow(ctx, "user@example.com", "Status")
		assert.True(t, allowed)

		allowed = limiter.Allow(ctx, "user@example.com", "Status")
		assert.False(t, allowed)
	})

	t.Run("Allow Status message after expiration", func(t *testing.T) {
		mr.FlushAll()
		allowed := limiter.Allow(ctx, "user@example.com", "Status")
		assert.True(t, allowed)

		allowed = limiter.Allow(ctx, "user@example.com", "Status")
		assert.True(t, allowed)

		mr.FastForward(time.Minute)
		allowed = limiter.Allow(ctx, "user@example.com", "Status")
		assert.True(t, allowed)
	})

	t.Run("Allow News message under limit", func(t *testing.T) {
		mr.FlushAll()
		allowed := limiter.Allow(ctx, "user@example.com", "News")
		assert.True(t, allowed)

		allowed = limiter.Allow(ctx, "user@example.com", "News")
		assert.False(t, allowed)
	})

	t.Run("Allow News message after expiration", func(t *testing.T) {
		mr.FlushAll()
		allowed := limiter.Allow(ctx, "user@example.com", "News")
		assert.True(t, allowed)

		mr.FastForward(24 * time.Hour)
		allowed = limiter.Allow(ctx, "user@example.com", "News")
		assert.True(t, allowed)
	})

	t.Run("Allow Marketing message under limit", func(t *testing.T) {
		mr.FlushAll()
		allowed := limiter.Allow(ctx, "user@example.com", "Marketing")
		assert.True(t, allowed)

		allowed = limiter.Allow(ctx, "user@example.com", "Marketing")
		assert.True(t, allowed)

		allowed = limiter.Allow(ctx, "user@example.com", "Marketing")
		assert.True(t, allowed)

		allowed = limiter.Allow(ctx, "user@example.com", "Marketing")
		assert.False(t, allowed)
	})

	t.Run("Allow Marketing message after expiration", func(t *testing.T) {
		mr.FlushAll()
		allowed := limiter.Allow(ctx, "user@example.com", "Marketing")
		assert.True(t, allowed)

		allowed = limiter.Allow(ctx, "user@example.com", "Marketing")
		assert.True(t, allowed)

		allowed = limiter.Allow(ctx, "user@example.com", "Marketing")
		assert.True(t, allowed)

		mr.FastForward(time.Hour)
		allowed = limiter.Allow(ctx, "user@example.com", "Marketing")
		assert.True(t, allowed)
	})

	t.Run("Handle Redis error gracefully", func(t *testing.T) {
		mr.Close()
		allowed := limiter.Allow(ctx, "user@example.com", "Status")
		assert.False(t, allowed)
	})
}
