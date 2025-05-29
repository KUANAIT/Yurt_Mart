package tests

import (
	"testing"
	"user-service/internal/infrastructure/email"
)

func TestSMTPSender(t *testing.T) {
	sender := email.NewSMTPSender(
		"smtp.mail.ru",
		465,
		"aytzhanovk@internet.ru",
		"Q1mfUaA5BcuryV4Fo3Gq",
		"aytzhanovk@internet.ru",
	)

	err := sender.SendWelcomeEmail("aytzhanovk@internet.ru", "Test User")
	if err != nil {
		t.Errorf("Failed to send welcome email: %v", err)
	}

	err = sender.SendPasswordResetEmail("aytzhanovk@internet.ru", "test-reset-token")
	if err != nil {
		t.Errorf("Failed to send password reset email: %v", err)
	}
}
