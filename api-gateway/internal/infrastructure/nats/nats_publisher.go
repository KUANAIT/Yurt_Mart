package nats

import (
	"encoding/json"

	"api-gateway/internal/core/ports"
	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn
}

func NewPublisher(conn *nats.Conn) *Publisher {
	return &Publisher{conn: conn}
}

func (p *Publisher) PublishUserEvent(event *ports.UserEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.conn.Publish(event.Type, data)
}
