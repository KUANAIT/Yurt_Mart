package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"user-service/internal/core/ports"
)

type SMTPSender struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewSMTPSender(host string, port int, username, password, from string) ports.EmailService {
	return &SMTPSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *SMTPSender) SendWelcomeEmail(to, name string) error {
	subject := "Welcome to Our Store"
	body := "Hello " + name + ",\n\nWelcome to our store!"
	return s.sendEmail(to, subject, body)
}

func (s *SMTPSender) SendPasswordResetEmail(to, token string) error {
	subject := "Password Reset"
	body := "Your password reset token is: " + token
	return s.sendEmail(to, subject, body)
}

func (s *SMTPSender) sendEmail(to, subject, body string) error {
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         s.host,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", s.host, s.port), tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Close()

	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	if err := client.Mail(s.from); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to create email writer: %v", err)
	}
	defer w.Close()

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write email: %v", err)
	}

	return nil
}
