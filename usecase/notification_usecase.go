package usecase

import (
	"context"
	"fmt"

	"github.com/jordanlanch/modak-test/domain"
)

type NotificationUseCase struct {
	Repo        domain.NotificationRepository
	RateLimiter domain.RateLimiter
}

func NewNotificationUseCase(repo domain.NotificationRepository, limiter domain.RateLimiter) domain.NotificationService {
	return &NotificationUseCase{
		Repo:        repo,
		RateLimiter: limiter,
	}
}

func (n *NotificationUseCase) SendNotification(ctx context.Context, notification domain.Notification) error {
	if !n.RateLimiter.Allow(ctx, notification.Recipient, notification.MessageType) {
		return &domain.RateLimitError{MessageType: notification.MessageType, Recipient: notification.Recipient}
	}

	if err := n.Repo.RecordNotificationSent(ctx, notification); err != nil {
		return fmt.Errorf("error recording notification: %w", err)
	}

	return nil
}
