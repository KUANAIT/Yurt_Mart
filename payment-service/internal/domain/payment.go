package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidPaymentID     = errors.New("invalid payment ID")
	ErrInvalidOrderID       = errors.New("invalid order ID")
	ErrInvalidUserID        = errors.New("invalid user ID")
	ErrInvalidAmount        = errors.New("invalid amount")
	ErrInvalidCurrency      = errors.New("invalid currency")
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
)

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "PENDING"
	PaymentStatusProcessing PaymentStatus = "PROCESSING"
	PaymentStatusCompleted  PaymentStatus = "COMPLETED"
	PaymentStatusFailed     PaymentStatus = "FAILED"
	PaymentStatusCancelled  PaymentStatus = "CANCELLED"
	PaymentStatusRefunded   PaymentStatus = "REFUNDED"
)

type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "CREDIT_CARD"
	PaymentMethodMetaMask   PaymentMethod = "METAMASK"
)

type Payment struct {
	ID            string
	OrderID       string
	UserID        string
	Amount        float64
	Currency      string
	Status        string
	PaymentMethod string
	TransactionID string
	ErrorMessage  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreditCardInfo struct {
	CardNumber     string
	ExpiryMonth    string
	ExpiryYear     string
	CVV            string
	CardholderName string
}

type MetaMaskInfo struct {
	WalletAddress   string
	TransactionHash string
	ContractAddress string
	PaymentData     string
}

func NewPayment(
	orderID string,
	userID string,
	amount float64,
	currency string,
	method PaymentMethod,
) (*Payment, error) {
	if orderID == "" {
		return nil, ErrInvalidOrderID
	}

	if userID == "" {
		return nil, ErrInvalidUserID
	}

	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	if currency == "" {
		return nil, ErrInvalidCurrency
	}

	if method != PaymentMethodCreditCard && method != PaymentMethodMetaMask {
		return nil, ErrInvalidPaymentMethod
	}

	return &Payment{
		OrderID:       orderID,
		UserID:        userID,
		Amount:        amount,
		Currency:      currency,
		Status:        string(PaymentStatusPending),
		PaymentMethod: string(method),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

func (p *Payment) UpdateStatus(status PaymentStatus) {
	p.Status = string(status)
	p.UpdatedAt = time.Now()
}

func (p *Payment) SetTransactionID(txID string) {
	p.TransactionID = txID
	p.UpdatedAt = time.Now()
}

func (p *Payment) SetError(err string) {
	p.ErrorMessage = err
	p.Status = string(PaymentStatusFailed)
	p.UpdatedAt = time.Now()
}

func (p *Payment) IsCompleted() bool {
	return p.Status == string(PaymentStatusCompleted)
}

func (p *Payment) IsPending() bool {
	return p.Status == string(PaymentStatusPending) || p.Status == string(PaymentStatusProcessing)
}

func (p *Payment) CanBeRetried() bool {
	return p.Status == string(PaymentStatusFailed) || p.Status == string(PaymentStatusCancelled)
}

func (p *Payment) MarkAsProcessing() {
	if p.Status == string(PaymentStatusPending) {
		p.Status = string(PaymentStatusProcessing)
		p.UpdatedAt = time.Now()
	}
}

func (p *Payment) MarkAsCompleted(transactionID string) {
	if p.IsPending() {
		p.Status = string(PaymentStatusCompleted)
		p.TransactionID = transactionID
		p.UpdatedAt = time.Now()
	}
}

func (p *Payment) Cancel() {
	if p.IsPending() {
		p.Status = string(PaymentStatusCancelled)
		p.UpdatedAt = time.Now()
	}
}

func (p *Payment) Refund() error {
	if p.Status != string(PaymentStatusCompleted) {
		return errors.New("only completed payments can be refunded")
	}

	p.Status = string(PaymentStatusRefunded)
	p.UpdatedAt = time.Now()
	return nil
} 