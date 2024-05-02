package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jordanlanch/modak-test/domain"
	"github.com/jordanlanch/modak-test/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetCurrentCount(ctx context.Context, recipient string, messageType string) (int, error) {
	args := m.Called(ctx, recipient, messageType)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) RecordNotificationSent(ctx context.Context, notification domain.Notification) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}

type MockRateLimiter struct {
	mock.Mock
}

func (m *MockRateLimiter) Allow(ctx context.Context, recipient string, messageType string) bool {
	args := m.Called(ctx, recipient, messageType)
	return args.Bool(0)
}

func TestNotificationUseCase_SendNotification(t *testing.T) {
	ctx := context.Background()
	notification := domain.Notification{Recipient: "user@example.com", MessageType: "Status"}

	t.Run("Send notification successfully", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockLimiter := new(MockRateLimiter)

		uc := usecase.NewNotificationUseCase(mockRepo, mockLimiter)

		mockLimiter.On("Allow", ctx, notification.Recipient, notification.MessageType).Return(true)
		mockRepo.On("RecordNotificationSent", ctx, notification).Return(nil)

		err := uc.SendNotification(ctx, notification)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockLimiter.AssertExpectations(t)
	})

	t.Run("Send notification with rate limit error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockLimiter := new(MockRateLimiter)

		uc := usecase.NewNotificationUseCase(mockRepo, mockLimiter)

		mockLimiter.On("Allow", ctx, notification.Recipient, notification.MessageType).Return(false)

		err := uc.SendNotification(ctx, notification)
		assert.Error(t, err)
		var rle *domain.RateLimitError
		assert.True(t, errors.As(err, &rle))


		mockRepo.AssertNotCalled(t, "RecordNotificationSent")
		mockLimiter.AssertExpectations(t)
	})

	t.Run("Send notification with repository error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockLimiter := new(MockRateLimiter)

		uc := usecase.NewNotificationUseCase(mockRepo, mockLimiter)

		mockLimiter.On("Allow", ctx, notification.Recipient, notification.MessageType).Return(true)
		mockRepo.On("RecordNotificationSent", ctx, notification).Return(errors.New("repository error"))

		err := uc.SendNotification(ctx, notification)
		assert.Error(t, err)
		assert.EqualError(t, err, "error recording notification: repository error")

		mockRepo.AssertExpectations(t)
		mockLimiter.AssertExpectations(t)
	})
}

func TestNewNotificationUseCase(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLimiter := new(MockRateLimiter)

	uc := usecase.NewNotificationUseCase(mockRepo, mockLimiter)

	assert.NotNil(t, uc)
	assert.Equal(t, mockRepo, uc.(*usecase.NotificationUseCase).Repo)
	assert.Equal(t, mockLimiter, uc.(*usecase.NotificationUseCase).RateLimiter)
}
