package nats

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
	"user-service/internal/core/domain"
)

type Subscriber struct {
	conn    *nats.Conn
	handler domain.EventHandler
}

func NewSubscriber(conn *nats.Conn, handler domain.EventHandler) *Subscriber {
	return &Subscriber{
		conn:    conn,
		handler: handler,
	}
}

func (s *Subscriber) Subscribe() error {
	if _, err := s.conn.Subscribe(domain.UserCreatedEvent, func(msg *nats.Msg) {
		var user domain.User
		if err := json.Unmarshal(msg.Data, &user); err != nil {
			return
		}
		s.handler.HandleUserCreated(&user)
	}); err != nil {
		return err
	}

	if _, err := s.conn.Subscribe(domain.UserUpdatedEvent, func(msg *nats.Msg) {
		var user domain.User
		if err := json.Unmarshal(msg.Data, &user); err != nil {
			return
		}
		s.handler.HandleUserUpdated(&user)
	}); err != nil {
		return err
	}

	if _, err := s.conn.Subscribe(domain.UserDeletedEvent, func(msg *nats.Msg) {
		var userID string
		if err := json.Unmarshal(msg.Data, &userID); err != nil {
			return
		}
		s.handler.HandleUserDeleted(userID)
	}); err != nil {
		return err
	}

	return nil
}
