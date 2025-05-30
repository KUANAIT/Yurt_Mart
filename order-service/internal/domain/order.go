package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidOrderID      = errors.New("invalid order ID")
	ErrInvalidUserID       = errors.New("invalid user ID")
	ErrInvalidCartID       = errors.New("invalid cart ID")
	ErrInvalidTotalPrice   = errors.New("invalid total price")
	ErrInvalidDeliveryTime = errors.New("invalid delivery time")
	ErrOrderNotFound       = errors.New("order not found")
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusPreparing  OrderStatus = "preparing"
	OrderStatusDelivering OrderStatus = "delivering"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	ID              string
	UserID          string
	CartID          string
	TotalPrice      float64
	Currency        string
	Status          OrderStatus
	DeliveryAddress *DeliveryAddress
	DeliveryTime    time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OrderItem struct {
	ProductID   string
	ProductName string
	Quantity    int32
	UnitPrice   float64
	TotalPrice  float64
}

func NewOrder(userID string, cartID string, address *DeliveryAddress, deliveryTime time.Time) (*Order, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	if cartID == "" {
		return nil, ErrInvalidCartID
	}

	if deliveryTime.Before(time.Now()) {
		return nil, ErrInvalidDeliveryTime
	}

	return &Order{
		UserID:          userID,
		CartID:          cartID,
		Status:          OrderStatusPending,
		DeliveryAddress: address,
		DeliveryTime:    deliveryTime,
		Currency:        "USD", // TODO: Make configurable
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil
}

func (o *Order) UpdateStatus(status OrderStatus) {
	o.Status = status
	o.UpdatedAt = time.Now()
}

func (o *Order) UpdateDeliveryTime(deliveryTime time.Time) error {
	if deliveryTime.Before(time.Now()) {
		return ErrInvalidDeliveryTime
	}

	o.DeliveryTime = deliveryTime
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) UpdateDeliveryAddress(address *DeliveryAddress) {
	o.DeliveryAddress = address
	o.UpdatedAt = time.Now()
}

func (o *Order) CanBePaid() bool {
	return o.Status == OrderStatusPending || o.Status == OrderStatusConfirmed
}

func (o *Order) CanBeCancelled() bool {
	return o.Status != OrderStatusDelivered && o.Status != OrderStatusCancelled
}

func (o *Order) MarkAsPaid() {
	if o.CanBePaid() {
		o.Status = OrderStatusConfirmed
		o.UpdatedAt = time.Now()
	}
}

func (o *Order) MarkAsAwaitingPayment() {
	if o.Status == OrderStatusPending {
		o.Status = OrderStatusConfirmed
		o.UpdatedAt = time.Now()
	}
}

func (o *Order) Cancel() error {
	if !o.CanBeCancelled() {
		return errors.New("order cannot be cancelled in current status")
	}

	o.Status = OrderStatusCancelled
	o.UpdatedAt = time.Now()
	return nil
} 