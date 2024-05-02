package test

import (
	"net/http"
	"testing"

	"github.com/jordanlanch/modak-test/domain"
	"github.com/sirupsen/logrus"

	"github.com/jordanlanch/modak-test/test/e2e"
)

const (
	statusOK              = http.StatusOK
	StatusTooManyRequests = http.StatusTooManyRequests
)

func TestE2E(t *testing.T) {
	expect, teardown := e2e.Setup(t, "fixtures/notification_test")
	defer teardown()

	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	tests := []struct {
		name         string
		messageType  string
		limit        int
		expectStatus []int // Use slice to anticipate different responses for each call
	}{
		{
			name:         "Status, limit testing",
			messageType:  "Status",
			limit:        3,
			expectStatus: []int{statusOK, statusOK, StatusTooManyRequests},
		},
		{
			name:         "News, limit testing",
			messageType:  "News",
			limit:        2,
			expectStatus: []int{statusOK, StatusTooManyRequests},
		},
		{
			name:         "Marketing, limit testing",
			messageType:  "Marketing",
			limit:        4,
			expectStatus: []int{statusOK, statusOK, statusOK, StatusTooManyRequests},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i < tc.limit; i++ {
				notification := domain.Notification{
					Recipient:   "test@example.com",
					MessageType: tc.messageType,
					Content:     "Test message for " + tc.name,
				}
				response := expect.POST("/api/notification").
					WithJSON(notification).
					Expect()

				response.Status(tc.expectStatus[i]) // Check the response status as expected per call
			}
		})
	}
}
