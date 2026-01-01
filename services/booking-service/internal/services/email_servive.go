package services

import (
	"booking-service/internal/repositories"
	"booking-service/pkg/domain"
	"context"
	"fmt"
	"sync"
	"time"
)

type EmailService interface {
	SendAsync(email *domain.Email) error
	SendBulk(emails []*domain.Email)
	Shutdown()

	SendPurchaseEmail(ctx context.Context, to string, name string, orderId string, amount float64) error
}

type emailService struct {
	repo    repositories.EmailRepository
	jobs    chan *domain.Email
	wg      sync.WaitGroup
	workers int
}

func NewEmailService(repo repositories.EmailRepository, workers int) EmailService {
	s := &emailService{
		repo:    repo,
		jobs:    make(chan *domain.Email, 1000),
		workers: workers,
	}
	s.startWorkers()
	return s
}

// SendSync envía inmediatamente (bloqueante)
func (s *emailService) SendPurchaseEmail(ctx context.Context, to string, name string, orderId string, amount float64) error {
	subject := fmt.Sprintf("✅ Confirmación de Compra #%s", orderId[:8])

	body := fmt.Sprintf(`<!DOCTYPE html><html><head><meta charset="UTF-8"><style>body{font-family:Arial;background:#f4f4f4;margin:0;padding:20px}.card{max-width:600px;margin:0 auto;background:#fff;padding:40px;border-radius:8px}.header{background:#667eea;color:#fff;padding:30px;text-align:center;border-radius:8px 8px 0 0;margin:-40px -40px 30px}.amount{font-size:32px;font-weight:bold;color:#667eea;margin:20px 0}.divider{height:1px;background:#e5e7eb;margin:30px 0}</style></head><body><div class="card"><div class="header"><h1>¡Compra Confirmada!</h1></div><p>Hola %s 👋</p><p>Tu compra se procesó exitosamente.</p><div class="amount">$%.2f USD</div><p><strong>Orden:</strong> %s</p><div class="divider"></div><p style="color:#666;font-size:14px;text-align:center">Ingresa a tu cuenta de SeatGuards para ver y descargar tu comprobante de pago</p><p style="color:#999;margin-top:30px;font-size:13px">Gracias por tu compra. Si tienes preguntas, contáctanos en soporte@seatguards.com</p></div></body></html>`, name, amount/100, orderId)

	email := &domain.Email{
		To:      []string{to},
		Subject: subject,
		Body:    body,
	}

	return s.repo.SendEmail(ctx, email)
}

// Worker pool: procesa emails concurrentemente
func (s *emailService) startWorkers() {
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}
}

func (s *emailService) worker(id int) {
	defer s.wg.Done()
	for email := range s.jobs {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		if err := s.repo.SendEmail(ctx, email); err != nil {
			fmt.Printf("Worker %d error: %v\n", id, err)
		}
		cancel()
	}
}

// SendAsync agrega a la cola (no bloqueante)
func (s *emailService) SendAsync(email *domain.Email) error {
	select {
	case s.jobs <- email:
		return nil
	default:
		return fmt.Errorf("queue full")
	}
}

// SendBulk procesa múltiples emails en la cola
func (s *emailService) SendBulk(emails []*domain.Email) {
	go func() {
		for _, email := range emails {
			s.SendAsync(email)
		}
	}()
}

// Shutdown cierra workers gracefully
func (s *emailService) Shutdown() {
	close(s.jobs)
	s.wg.Wait()
}
