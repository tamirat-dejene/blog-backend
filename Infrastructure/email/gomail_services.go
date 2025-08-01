package email

import (
	"context"
	"gopkg.in/gomail.v2"
)

type GomailEmailService struct {
	dialer *gomail.Dialer
	from   string
}

func NewGomailEmailService(smtpHost string, smtpPort int, from, username, password string) *GomailEmailService {
	dialer := gomail.NewDialer(smtpHost, smtpPort, username, password)
	return &GomailEmailService{
		dialer: dialer,
		from:   from,
	}
}

func (s *GomailEmailService) SendEmail(ctx context.Context, to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	return s.dialer.DialAndSend(m)
}
