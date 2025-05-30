package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/hsibAD/payment-service/internal/domain"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type SMTPNotifier struct {
	config SMTPConfig
	auth   smtp.Auth
}

func NewSMTPNotifier(config SMTPConfig) *SMTPNotifier {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	return &SMTPNotifier{
		config: config,
		auth:   auth,
	}
}

func (n *SMTPNotifier) SendPaymentConfirmation(ctx context.Context, payment *domain.Payment, email string) error {
	subject := "Payment Confirmation"
	body := n.generatePaymentConfirmationEmail(payment)

	return n.sendEmail(email, subject, body)
}

func (n *SMTPNotifier) SendPaymentFailure(ctx context.Context, payment *domain.Payment, email string) error {
	subject := "Payment Failed"
	body := n.generatePaymentFailureEmail(payment)

	return n.sendEmail(email, subject, body)
}

func (n *SMTPNotifier) SendRefundConfirmation(ctx context.Context, payment *domain.Payment, email string) error {
	subject := "Refund Confirmation"
	body := n.generateRefundConfirmationEmail(payment)

	return n.sendEmail(email, subject, body)
}

func (n *SMTPNotifier) sendEmail(to, subject, body string) error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := fmt.Sprintf("Subject: %s\n%s\n%s", subject, mime, body)

	addr := fmt.Sprintf("%s:%d", n.config.Host, n.config.Port)
	return smtp.SendMail(addr, n.auth, n.config.From, []string{to}, []byte(msg))
}

func (n *SMTPNotifier) generatePaymentConfirmationEmail(payment *domain.Payment) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f8f9fa; padding: 20px; text-align: center; }
        .payment-details { margin: 20px 0; }
        .total { font-weight: bold; margin-top: 20px; }
        .footer { text-align: center; margin-top: 30px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Payment Confirmation</h1>
            <p>Payment ID: {{.ID}}</p>
        </div>
        <div class="payment-details">
            <h2>Payment Details</h2>
            <p>Order ID: {{.OrderID}}</p>
            <p>Status: {{.Status}}</p>
            <p>Method: {{.PaymentMethod}}</p>
            <p>Amount: {{.Currency}} {{.Amount}}</p>
            {{if .TransactionID}}
            <p>Transaction ID: {{.TransactionID}}</p>
            {{end}}
        </div>
        <div class="footer">
            <p>Thank you for your payment!</p>
        </div>
    </div>
</body>
</html>`

	t, err := template.New("payment_confirmation").Parse(tmpl)
	if err != nil {
		return "Error generating email template"
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, payment); err != nil {
		return "Error executing email template"
	}

	return buf.String()
}

func (n *SMTPNotifier) generatePaymentFailureEmail(payment *domain.Payment) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f8f9fa; padding: 20px; text-align: center; }
        .payment-details { margin: 20px 0; }
        .error { color: #dc3545; margin: 20px 0; }
        .footer { text-align: center; margin-top: 30px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Payment Failed</h1>
            <p>Payment ID: {{.ID}}</p>
        </div>
        <div class="payment-details">
            <h2>Payment Details</h2>
            <p>Order ID: {{.OrderID}}</p>
            <p>Method: {{.PaymentMethod}}</p>
            <p>Amount: {{.Currency}} {{.Amount}}</p>
        </div>
        <div class="error">
            <h3>Error Details</h3>
            <p>{{.ErrorMessage}}</p>
        </div>
        <div class="footer">
            <p>Please try again or contact support if the problem persists.</p>
        </div>
    </div>
</body>
</html>`

	t, err := template.New("payment_failure").Parse(tmpl)
	if err != nil {
		return "Error generating email template"
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, payment); err != nil {
		return "Error executing email template"
	}

	return buf.String()
}

func (n *SMTPNotifier) generateRefundConfirmationEmail(payment *domain.Payment) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f8f9fa; padding: 20px; text-align: center; }
        .refund-details { margin: 20px 0; }
        .total { font-weight: bold; margin-top: 20px; }
        .footer { text-align: center; margin-top: 30px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Refund Confirmation</h1>
            <p>Payment ID: {{.ID}}</p>
        </div>
        <div class="refund-details">
            <h2>Refund Details</h2>
            <p>Order ID: {{.OrderID}}</p>
            <p>Original Payment Method: {{.PaymentMethod}}</p>
            <p>Refund Amount: {{.Currency}} {{.Amount}}</p>
            {{if .TransactionID}}
            <p>Transaction ID: {{.TransactionID}}</p>
            {{end}}
        </div>
        <div class="footer">
            <p>The refund has been processed successfully.</p>
            <p>Please allow 5-10 business days for the refund to appear in your account.</p>
        </div>
    </div>
</body>
</html>`

	t, err := template.New("refund_confirmation").Parse(tmpl)
	if err != nil {
		return "Error generating email template"
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, payment); err != nil {
		return "Error executing email template"
	}

	return buf.String()
} 