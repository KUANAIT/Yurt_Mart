package domain

import "context"

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	GetByUserID(ctx context.Context, userID string, page, limit int) ([]*Order, int, error)
	Update(ctx context.Context, order *Order) error
	UpdateStatus(ctx context.Context, orderID string, status OrderStatus) error
	Delete(ctx context.Context, id string) error
}

type DeliveryAddressRepository interface {
	Create(ctx context.Context, address *DeliveryAddress) error
	GetByID(ctx context.Context, id string) (*DeliveryAddress, error)
	GetByUserID(ctx context.Context, userID string) ([]*DeliveryAddress, error)
	Update(ctx context.Context, address *DeliveryAddress) error
	Delete(ctx context.Context, id string) error
	SetDefault(ctx context.Context, userID string, addressID string) error
}

type DeliverySlotRepository interface {
	GetAvailableSlots(ctx context.Context, date string) ([]*DeliverySlot, error)
	ReserveSlot(ctx context.Context, orderID string, slotID string) error
	ReleaseSlot(ctx context.Context, orderID string, slotID string) error
}

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ttl int) error
	Get(ctx context.Context, key string) (interface{}, error)
	Delete(ctx context.Context, key string) error
}

type EventPublisher interface {
	PublishOrderCreated(ctx context.Context, order *Order) error
	PublishOrderStatusUpdated(ctx context.Context, order *Order) error
	PublishOrderCancelled(ctx context.Context, order *Order) error
} 