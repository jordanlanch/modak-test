package domain

import (
	"context"
	"fmt"
)

type RateLimitError struct {
	MessageType string
	Recipient   string
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limit exceeded for %s messages to %s", e.MessageType, e.Recipient)
}

type RateLimiter interface {
	Allow(ctx context.Context, recipient, messageType string) bool
}
