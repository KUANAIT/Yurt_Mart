package domain

import "context"

type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) error
	GetByID(ctx context.Context, id string) (*Payment, error)
	GetByOrderID(ctx context.Context, orderID string) ([]*Payment, error)
	GetByUserID(ctx context.Context, userID string, page, limit int) ([]*Payment, int, error)
	Update(ctx context.Context, payment *Payment) error
	UpdateStatus(ctx context.Context, paymentID string, status PaymentStatus) error
}

type CreditCardProcessor interface {
	ProcessPayment(ctx context.Context, payment *Payment, cardInfo *CreditCardInfo) error
	RefundPayment(ctx context.Context, payment *Payment) error
	ValidateCard(ctx context.Context, cardInfo *CreditCardInfo) error
}

type MetaMaskProcessor interface {
	InitiateTransaction(ctx context.Context, payment *Payment, walletAddress string) (*MetaMaskInfo, error)
	VerifyTransaction(ctx context.Context, payment *Payment, transactionHash string) error
	GetTransactionStatus(ctx context.Context, transactionHash string) (string, error)
}

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ttl int) error
	Get(ctx context.Context, key string) (interface{}, error)
	Delete(ctx context.Context, key string) error
}

type EventPublisher interface {
	PublishPaymentCreated(ctx context.Context, payment *Payment) error
	PublishPaymentStatusUpdated(ctx context.Context, payment *Payment) error
	PublishPaymentCompleted(ctx context.Context, payment *Payment) error
	PublishPaymentFailed(ctx context.Context, payment *Payment) error
	PublishPaymentRefunded(ctx context.Context, payment *Payment) error
}

type EmailNotifier interface {
	SendPaymentConfirmation(ctx context.Context, payment *Payment) error
	SendPaymentFailure(ctx context.Context, payment *Payment) error
	SendRefundConfirmation(ctx context.Context, payment *Payment) error
} 