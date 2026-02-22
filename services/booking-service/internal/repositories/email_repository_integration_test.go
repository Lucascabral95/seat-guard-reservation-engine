package repositories

import (
	"booking-service/pkg/domain"
	"context"
	"os"
	"testing"
	"time"
)

func TestEmailRepository_Integration_SendEmail(t *testing.T) {
	host := os.Getenv("BOOKING_IT_SMTP_HOST")
	port := os.Getenv("BOOKING_IT_SMTP_PORT")
	user := os.Getenv("BOOKING_IT_SMTP_USER")
	pass := os.Getenv("BOOKING_IT_SMTP_PASS")
	from := os.Getenv("BOOKING_IT_SMTP_FROM")
	to := os.Getenv("BOOKING_IT_SMTP_TO")

	if host == "" || port == "" || user == "" || pass == "" || from == "" || to == "" {
		t.Skip("SMTP integration env vars missing; set BOOKING_IT_SMTP_HOST/PORT/USER/PASS/FROM/TO")
	}

	repo, err := NewEmailRepository(host, port, user, pass, from)
	if err != nil {
		t.Fatalf("failed to init email repo: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = repo.SendEmail(ctx, &domain.Email{
		To:      []string{to},
		Subject: "SeatGuard IT SMTP Test",
		Body:    "<p>Integration test email</p>",
	})
	if err != nil {
		t.Fatalf("send email failed: %v", err)
	}
}
