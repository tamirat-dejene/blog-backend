package security

import(
	"fmt"
	"net/smtp"
)

type Mailer struct {
	Host	 string
	Port	 string
	Username string	
	Password string
	From     string
}

func NewMailer(host, port, username, password, from string) *Mailer {
	return &Mailer{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}
}

func (m *Mailer) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)
	err := smtp.SendMail(addr, auth, m.From, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}