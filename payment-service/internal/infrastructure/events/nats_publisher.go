package events

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/hsibAD/payment-service/internal/domain"
)

const (
	PaymentCreatedSubject       = "payment.created"
	PaymentStatusUpdatedSubject = "payment.status.updated"
	PaymentCompletedSubject     = "payment.completed"
	PaymentFailedSubject        = "payment.failed"
	PaymentRefundedSubject      = "payment.refunded"
)

type NATSPublisher struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

type PaymentEvent struct {
	ID            string  `json:"id"`
	OrderID       string  `json:"order_id"`
	UserID        string  `json:"user_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
	PaymentMethod string  `json:"payment_method"`
	TransactionID string  `json:"transaction_id,omitempty"`
	ErrorMessage  string  `json:"error_message,omitempty"`
	EventType     string  `json:"event_type"`
	Timestamp     int64   `json:"timestamp"`
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
		Name:     "PAYMENTS",
		Subjects: []string{"payment.*", "payment.status.*"},
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

func (p *NATSPublisher) PublishPaymentCreated(ctx context.Context, payment *domain.Payment) error {
	event := PaymentEvent{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		UserID:        payment.UserID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		Status:        string(payment.Status),
		PaymentMethod: string(payment.PaymentMethod),
		EventType:     "PaymentCreated",
		Timestamp:     payment.CreatedAt.Unix(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = p.js.Publish(PaymentCreatedSubject, data)
	return err
}

func (p *NATSPublisher) PublishPaymentStatusUpdated(ctx context.Context, payment *domain.Payment) error {
	event := PaymentEvent{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		UserID:        payment.UserID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		Status:        string(payment.Status),
		PaymentMethod: string(payment.PaymentMethod),
		TransactionID: payment.TransactionID,
		ErrorMessage:  payment.ErrorMessage,
		EventType:     "PaymentStatusUpdated",
		Timestamp:     payment.UpdatedAt.Unix(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = p.js.Publish(PaymentStatusUpdatedSubject, data)
	return err
}

func (p *NATSPublisher) PublishPaymentCompleted(ctx context.Context, payment *domain.Payment) error {
	event := PaymentEvent{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		UserID:        payment.UserID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		Status:        string(payment.Status),
		PaymentMethod: string(payment.PaymentMethod),
		TransactionID: payment.TransactionID,
		EventType:     "PaymentCompleted",
		Timestamp:     payment.UpdatedAt.Unix(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = p.js.Publish(PaymentCompletedSubject, data)
	return err
}

func (p *NATSPublisher) PublishPaymentFailed(ctx context.Context, payment *domain.Payment) error {
	event := PaymentEvent{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		UserID:        payment.UserID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		Status:        string(payment.Status),
		PaymentMethod: string(payment.PaymentMethod),
		ErrorMessage:  payment.ErrorMessage,
		EventType:     "PaymentFailed",
		Timestamp:     payment.UpdatedAt.Unix(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = p.js.Publish(PaymentFailedSubject, data)
	return err
}

func (p *NATSPublisher) PublishPaymentRefunded(ctx context.Context, payment *domain.Payment) error {
	event := PaymentEvent{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		UserID:        payment.UserID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		Status:        string(payment.Status),
		PaymentMethod: string(payment.PaymentMethod),
		TransactionID: payment.TransactionID,
		EventType:     "PaymentRefunded",
		Timestamp:     payment.UpdatedAt.Unix(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = p.js.Publish(PaymentRefundedSubject, data)
	return err
}

func (p *NATSPublisher) Close() error {
	p.nc.Close()
	return nil
} 