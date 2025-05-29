package handlers

import (
	//"context"
	"user-service/internal/core/domain"
	"user-service/internal/core/ports"
)

type EventHandler struct {
	emailService ports.EmailService
}

func NewEventHandler(emailService ports.EmailService) domain.EventHandler {
	return &EventHandler{
		emailService: emailService,
	}
}

func (h *EventHandler) HandleUserCreated(user *domain.User) error {
	return h.emailService.SendWelcomeEmail(user.Email, user.Name)
}

func (h *EventHandler) HandleUserUpdated(user *domain.User) error {
	return nil
}

func (h *EventHandler) HandleUserDeleted(userID string) error {
	// Implement if needed
	return nil
}
