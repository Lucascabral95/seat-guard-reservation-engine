package messaging

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestSQSClient_Integration_Send(t *testing.T) {
	region := os.Getenv("BOOKING_IT_AWS_REGION")
	queueURL := os.Getenv("BOOKING_IT_SQS_QUEUE_URL")
	if region == "" || queueURL == "" {
		t.Skip("SQS integration env vars missing; set BOOKING_IT_AWS_REGION and BOOKING_IT_SQS_QUEUE_URL")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := NewSQSClient(ctx, region, queueURL)
	if err != nil {
		t.Fatalf("failed to create sqs client: %v", err)
	}

	id, err := client.Send(ctx, fmt.Sprintf(`{"type":"it_test","ts":%d}`, time.Now().UnixNano()), "booking-it", fmt.Sprintf("it-%d", time.Now().UnixNano()))
	if err != nil {
		t.Fatalf("failed to send message: %v", err)
	}
	if id == "" {
		t.Fatalf("expected non-empty message id")
	}
}
