package repository

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/jordanlanch/modak-test/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewNotificationRepository(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("an error '%s' occurred when starting miniredis", err)
	}
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	repo := NewNotificationRepository(rdb)
	assert.NotNil(t, repo)
}

func TestGetCurrentCount(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("an error '%s' occurred when starting miniredis", err)
	}
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	repo := NewNotificationRepository(rdb)
	ctx := context.Background()

	// Test when key does not exist
	count, err := repo.GetCurrentCount(ctx, "recipient1", "Status")
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

	// Test with existing key
	mr.Set("recipient1:Status", "3")
	count, err = repo.GetCurrentCount(ctx, "recipient1", "Status")
	assert.NoError(t, err)
	assert.Equal(t, 3, count)

	// Test Redis error
	mr.Close()
	_, err = repo.GetCurrentCount(ctx, "recipient1", "Status")
	assert.Error(t, err)
}

func TestRecordNotificationSent(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("an error '%s' occurred when starting miniredis", err)
	}
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	repo := NewNotificationRepository(rdb)
	ctx := context.Background()

	notification := domain.Notification{
		Recipient:   "recipient1",
		MessageType: "Status",
	}

	// Test incrementing the count
	err = repo.RecordNotificationSent(ctx, notification)
	assert.NoError(t, err)

	count, err := mr.Get("recipient1:Status")
	assert.NoError(t, err)
	assert.Equal(t, "1", count)

	// Test Redis error
	mr.Close()
	err = repo.RecordNotificationSent(ctx, notification)
	assert.Error(t, err)
}
