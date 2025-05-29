package services

import (
	"fmt"
	"net/smtp"

	"user-service/internal/core/ports"
)

type emailService struct {
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	fromEmail    string
}

func NewEmailService(host string, port int, username, password, from string) ports.EmailService {
	return &emailService{
		smtpHost:     host,
		smtpPort:     port,
		smtpUsername: username,
		smtpPassword: password,
		fromEmail:    from,
	}
}

func (s *emailService) SendWelcomeEmail(email, name string) error {
	to := []string{email}
	subject := "Welcome to Our Store!"
	body := fmt.Sprintf("Hello %s,\n\nThank you for registering with our store!", name)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", email, subject, body))

	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)
	addr := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)

	return smtp.SendMail(addr, auth, s.fromEmail, to, msg)
}

func (s *emailService) SendPasswordResetEmail(email, token string) error {
	to := []string{email}
	subject := "Password Reset Request"
	body := fmt.Sprintf("Please use the following token to reset your password: %s", token)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", email, subject, body))

	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)
	addr := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)

	return smtp.SendMail(addr, auth, s.fromEmail, to, msg)
}
