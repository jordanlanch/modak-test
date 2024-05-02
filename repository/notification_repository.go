package repository

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/jordanlanch/modak-test/domain"
)

type notificationRepository struct {
	rdb *redis.Client
}

func NewNotificationRepository(rdb *redis.Client) domain.NotificationRepository {
	return &notificationRepository{rdb}
}

// GetCurrentCount retrieves the current count for notifications to a recipient based on message type
func (r *notificationRepository) GetCurrentCount(ctx context.Context, recipient string, messageType string) (int, error) {
	key := fmt.Sprintf("%s:%s", recipient, messageType)
	result, err := r.rdb.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return result, nil
}

// RecordNotificationSent increases the count of sent notifications
func (r *notificationRepository) RecordNotificationSent(ctx context.Context, notification domain.Notification) error {
	key := fmt.Sprintf("%s:%s", notification.Recipient, notification.MessageType)
	if _, err := r.rdb.Incr(ctx, key).Result(); err != nil {
		return err
	}
	return nil
}
