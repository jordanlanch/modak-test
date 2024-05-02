package domain

import "context"

type Notification struct {
	Recipient   string `json:"recipient" jsonschema:"required"`
	MessageType string `json:"message_type" jsonschema:"required"`
	Content     string `json:"content" jsonschema:"required"`
}

type NotificationRepository interface {
	GetCurrentCount(ctx context.Context, recipient string, messageType string) (int, error)
	RecordNotificationSent(ctx context.Context, notification Notification) error
}

type NotificationService interface {
	SendNotification(ctx context.Context, notification Notification) error
}
