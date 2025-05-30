package payment

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/charge"
	"github.com/stripe/stripe-go/v74/refund"
	"github.com/hsibAD/payment-service/internal/domain"
)

var (
	ErrInvalidCardNumber   = errors.New("invalid card number")
	ErrInvalidExpiryMonth = errors.New("invalid expiry month")
	ErrInvalidExpiryYear  = errors.New("invalid expiry year")
	ErrInvalidCVV        = errors.New("invalid CVV")
	ErrCardExpired       = errors.New("card has expired")
	ErrPaymentFailed     = errors.New("payment failed")
)

type CreditCardProcessor struct {
	stripeSecretKey string
}

func NewCreditCardProcessor(stripeSecretKey string) *CreditCardProcessor {
	stripe.Key = stripeSecretKey
	return &CreditCardProcessor{
		stripeSecretKey: stripeSecretKey,
	}
}

func (p *CreditCardProcessor) ProcessPayment(ctx context.Context, payment *domain.Payment, cardInfo *domain.CreditCardInfo) error {
	if err := p.ValidateCard(ctx, cardInfo); err != nil {
		return err
	}

	// Create Stripe token
	token, err := p.createStripeToken(cardInfo)
	if err != nil {
		return fmt.Errorf("failed to create stripe token: %w", err)
	}

	// Create charge parameters
	params := &stripe.ChargeParams{
		Amount:      stripe.Int64(int64(payment.Amount * 100)), // Convert to cents
		Currency:    stripe.String(string(payment.Currency)),
		Source:      &token.ID,
		Description: stripe.String(fmt.Sprintf("Payment for order %s", payment.OrderID)),
		Metadata: map[string]string{
			"order_id":    payment.OrderID,
			"payment_id":  payment.ID,
			"customer_id": payment.UserID,
		},
	}

	// Create charge
	ch, err := charge.New(params)
	if err != nil {
		return fmt.Errorf("failed to create charge: %w", err)
	}

	if !ch.Paid {
		return ErrPaymentFailed
	}

	payment.TransactionID = ch.ID
	return nil
}

func (p *CreditCardProcessor) RefundPayment(ctx context.Context, payment *domain.Payment) error {
	if payment.TransactionID == "" {
		return errors.New("no transaction ID found")
	}

	params := &stripe.RefundParams{
		Charge: stripe.String(payment.TransactionID),
		Metadata: map[string]string{
			"order_id":    payment.OrderID,
			"payment_id":  payment.ID,
			"customer_id": payment.UserID,
		},
	}

	_, err := refund.New(params)
	if err != nil {
		return fmt.Errorf("failed to create refund: %w", err)
	}

	return nil
}

func (p *CreditCardProcessor) ValidateCard(ctx context.Context, cardInfo *domain.CreditCardInfo) error {
	// Validate card number (Luhn algorithm)
	if !isValidCardNumber(cardInfo.CardNumber) {
		return ErrInvalidCardNumber
	}

	// Validate expiry month
	month, err := strconv.Atoi(cardInfo.ExpiryMonth)
	if err != nil || month < 1 || month > 12 {
		return ErrInvalidExpiryMonth
	}

	// Validate expiry year
	year, err := strconv.Atoi(cardInfo.ExpiryYear)
	if err != nil || year < time.Now().Year() {
		return ErrInvalidExpiryYear
	}

	// Check if card is expired
	currentTime := time.Now()
	if year < currentTime.Year() || (year == currentTime.Year() && month < int(currentTime.Month())) {
		return ErrCardExpired
	}

	// Validate CVV
	if !isValidCVV(cardInfo.CVV) {
		return ErrInvalidCVV
	}

	return nil
}

func (p *CreditCardProcessor) createStripeToken(cardInfo *domain.CreditCardInfo) (*stripe.Token, error) {
	params := &stripe.TokenParams{
		Card: &stripe.CardParams{
			Number:   stripe.String(cardInfo.CardNumber),
			ExpMonth: stripe.String(cardInfo.ExpiryMonth),
			ExpYear:  stripe.String(cardInfo.ExpiryYear),
			CVC:     stripe.String(cardInfo.CVV),
			Name:    stripe.String(cardInfo.CardholderName),
		},
	}

	return stripe.Tokens.New(params)
}

// Helper functions

func isValidCardNumber(number string) bool {
	// Remove spaces and dashes
	re := regexp.MustCompile(`[\s-]`)
	number = re.ReplaceAllString(number, "")

	// Check if number contains only digits
	if match, _ := regexp.MatchString(`^\d+$`, number); !match {
		return false
	}

	// Check length (13-19 digits)
	if len(number) < 13 || len(number) > 19 {
		return false
	}

	// Luhn algorithm
	sum := 0
	isEven := false

	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')

		if isEven {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		isEven = !isEven
	}

	return sum%10 == 0
}

func isValidCVV(cvv string) bool {
	// CVV should be 3 or 4 digits
	if match, _ := regexp.MatchString(`^\d{3,4}$`, cvv); !match {
		return false
	}
	return true
} 