package events

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/hsibAD/order-service/internal/domain"
)

const (
	OrderCreatedSubject       = "order.created"
	OrderStatusUpdatedSubject = "order.status.updated"
	OrderCancelledSubject     = "order.cancelled"
)

type NATSPublisher struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

type OrderEvent struct {
	ID              string                 `json:"id"`
	UserID          string                 `json:"user_id"`
	Status          string                 `json:"status"`
	TotalPrice      float64               `json:"total_price"`
	Currency        string                 `json:"currency"`
	DeliveryAddress *domain.DeliveryAddress `json:"delivery_address,omitempty"`
	Items           []domain.OrderItem      `json:"items"`
	EventType       string                 `json:"event_type"`
	Timestamp       int64                  `json:"timestamp"`
}

func NewNATSPublisher(url string) (*NATSPublisher, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	// Create the stream if it doesn't exist
	stream := &nats.StreamConfig{
		Name:     "ORDERS",
		Subjects: []string{"order.*", "order.status.*"},
	}

	if _, err := js.AddStream(stream); err != nil {
		if err != nats.ErrStreamNameAlreadyInUse {
			return nil, err
		}
	}

	return &NATSPublisher{
		nc: nc,
		js: js,
	}, nil
}

func (p *NATSPublisher) PublishOrderCreated(ctx context.Context, order *domain.Order) error {
	event := OrderEvent{
		ID:              order.ID,
		UserID:          order.UserID,
		Status:          string(order.Status),
		TotalPrice:      order.TotalPrice,
		Currency:        order.Currency,
		DeliveryAddress: order.DeliveryAddress,
		Items:           order.Items,
		EventType:       "OrderCreated",
		Timestamp:       order.CreatedAt.Unix(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = p.js.Publish(OrderCreatedSubject, data)
	return err
}

func (p *NATSPublisher) PublishOrderStatusUpdated(ctx context.Context, order *domain.Order) error {
	event := OrderEvent{
		ID:              order.ID,
		UserID:          order.UserID,
		Status:          string(order.Status),
		TotalPrice:      order.TotalPrice,
		Currency:        order.Currency,
		DeliveryAddress: order.DeliveryAddress,
		Items:           order.Items,
		EventType:       "OrderStatusUpdated",
		Timestamp:       order.UpdatedAt.Unix(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = p.js.Publish(OrderStatusUpdatedSubject, data)
	return err
}

func (p *NATSPublisher) PublishOrderCancelled(ctx context.Context, order *domain.Order) error {
	event := OrderEvent{
		ID:              order.ID,
		UserID:          order.UserID,
		Status:          string(order.Status),
		TotalPrice:      order.TotalPrice,
		Currency:        order.Currency,
		DeliveryAddress: order.DeliveryAddress,
		Items:           order.Items,
		EventType:       "OrderCancelled",
		Timestamp:       order.UpdatedAt.Unix(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = p.js.Publish(OrderCancelledSubject, data)
	return err
}

func (p *NATSPublisher) Close() error {
	p.nc.Close()
	return nil
} 