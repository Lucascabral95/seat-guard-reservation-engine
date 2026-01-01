package repositories

import (
	"booking-service/pkg/domain"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	netmail "net/mail"
	"net/smtp"
	"time"

	mail "gopkg.in/jordan-wright/email.v3"
)

type EmailRepository interface {
	SendEmail(ctx context.Context, email *domain.Email) error
}

type emailRepository struct {
	host      string
	port      string
	from      string
	fromAddr  string
	auth      smtp.Auth
	tlsConfig *tls.Config
	sem       chan struct{}
}

func NewEmailRepository(host string, port string, user, pass, from string) (EmailRepository, error) {
	// Validaciones
	if host == "" || port == "" {
		return nil, fmt.Errorf("SMTP host and port are required")
	}
	if user == "" || pass == "" {
		return nil, fmt.Errorf("SMTP credentials are required")
	}
	if from == "" {
		return nil, fmt.Errorf("valid SMTP_FROM/EMAIL_FROM is required")
	}
	parsedFrom, err := netmail.ParseAddress(from)
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP_FROM/EMAIL_FROM: %w", err)
	}

	auth := smtp.PlainAuth("", user, pass, host)
	tlsCfg := &tls.Config{ServerName: host, MinVersion: tls.VersionTLS12}

	log.Printf("✅ Email repository initialized - From: %s", from)

	return &emailRepository{
		host:      host,
		port:      port,
		from:      from,
		fromAddr:  parsedFrom.Address,
		auth:      auth,
		tlsConfig: tlsCfg,
		sem:       make(chan struct{}, 4),
	}, nil
}

func (r *emailRepository) SendEmail(ctx context.Context, email *domain.Email) error {
	if len(email.To) == 0 {
		return fmt.Errorf("email recipients are required")
	}

	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
	}

	select {
	case r.sem <- struct{}{}:
		defer func() { <-r.sem }()
	case <-ctx.Done():
		return ctx.Err()
	}

	log.Printf("📧 Sending email FROM: %s TO: %v", r.from, email.To)

	e := mail.NewEmail()
	e.From = r.from
	e.To = email.To
	e.Subject = email.Subject
	e.HTML = []byte(email.Body)

	msg, err := e.Bytes()
	if err != nil {
		return fmt.Errorf("failed to build email: %w", err)
	}

	addr := fmt.Sprintf("%s:%s", r.host, r.port)
	d := net.Dialer{}
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		log.Printf("❌ Email send failed: %v", err)
		return fmt.Errorf("SMTP dial error: %w", err)
	}
	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	}

	c, err := smtp.NewClient(conn, r.host)
	if err != nil {
		_ = conn.Close()
		log.Printf("❌ Email send failed: %v", err)
		return fmt.Errorf("SMTP client error: %w", err)
	}
	defer func() { _ = c.Close() }()

	if ok, _ := c.Extension("STARTTLS"); ok {
		if err := c.StartTLS(r.tlsConfig); err != nil {
			log.Printf("❌ Email send failed: %v", err)
			return fmt.Errorf("SMTP STARTTLS error: %w", err)
		}
	} else {
		return fmt.Errorf("SMTP server does not support STARTTLS")
	}

	if err := c.Auth(r.auth); err != nil {
		log.Printf("❌ Email send failed: %v", err)
		return fmt.Errorf("SMTP auth error: %w", err)
	}

	if err := c.Mail(r.fromAddr); err != nil {
		log.Printf("❌ Email send failed: %v", err)
		return fmt.Errorf("SMTP MAIL FROM error: %w", err)
	}
	for _, to := range email.To {
		if err := c.Rcpt(to); err != nil {
			log.Printf("❌ Email send failed: %v", err)
			return fmt.Errorf("SMTP RCPT TO %s error: %w", to, err)
		}
	}

	w, err := c.Data()
	if err != nil {
		log.Printf("❌ Email send failed: %v", err)
		return fmt.Errorf("SMTP DATA error: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		_ = w.Close()
		log.Printf("❌ Email send failed: %v", err)
		return fmt.Errorf("SMTP write error: %w", err)
	}
	if err := w.Close(); err != nil {
		log.Printf("❌ Email send failed: %v", err)
		return fmt.Errorf("SMTP DATA close error: %w", err)
	}

	if err := c.Quit(); err != nil {
		log.Printf("❌ Email send failed: %v", err)
		return fmt.Errorf("SMTP QUIT error: %w", err)
	}

	log.Printf("✅ Email sent successfully to %v", email.To)
	return nil
}
