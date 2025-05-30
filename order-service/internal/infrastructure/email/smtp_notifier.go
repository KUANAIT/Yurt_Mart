package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/hsibAD/order-service/internal/domain"
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

func (n *SMTPNotifier) SendOrderConfirmation(ctx context.Context, order *domain.Order, email string) error {
	subject := "Order Confirmation"
	body := n.generateOrderConfirmationEmail(order)

	return n.sendEmail(email, subject, body)
}

func (n *SMTPNotifier) SendOrderStatusUpdate(ctx context.Context, order *domain.Order, email string) error {
	subject := fmt.Sprintf("Order Status Update: %s", order.Status)
	body := n.generateOrderStatusUpdateEmail(order)

	return n.sendEmail(email, subject, body)
}

func (n *SMTPNotifier) SendOrderCancellation(ctx context.Context, order *domain.Order, email string) error {
	subject := "Order Cancellation"
	body := n.generateOrderCancellationEmail(order)

	return n.sendEmail(email, subject, body)
}

func (n *SMTPNotifier) sendEmail(to, subject, body string) error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := fmt.Sprintf("Subject: %s\n%s\n%s", subject, mime, body)

	addr := fmt.Sprintf("%s:%d", n.config.Host, n.config.Port)
	return smtp.SendMail(addr, n.auth, n.config.From, []string{to}, []byte(msg))
}

func (n *SMTPNotifier) generateOrderConfirmationEmail(order *domain.Order) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f8f9fa; padding: 20px; text-align: center; }
        .order-details { margin: 20px 0; }
        .items { margin: 20px 0; }
        .item { margin: 10px 0; }
        .total { font-weight: bold; margin-top: 20px; }
        .footer { text-align: center; margin-top: 30px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Order Confirmation</h1>
            <p>Order ID: {{.ID}}</p>
        </div>
        <div class="order-details">
            <h2>Order Details</h2>
            <p>Status: {{.Status}}</p>
            <p>Total: {{.Currency}} {{.TotalPrice}}</p>
            
            <h3>Delivery Address:</h3>
            <p>
                {{.DeliveryAddress.FullName}}<br>
                {{.DeliveryAddress.StreetAddress}}<br>
                {{if .DeliveryAddress.Apartment}}{{.DeliveryAddress.Apartment}}<br>{{end}}
                {{.DeliveryAddress.City}}, {{.DeliveryAddress.State}} {{.DeliveryAddress.PostalCode}}<br>
                {{.DeliveryAddress.Country}}
            </p>
        </div>
        <div class="items">
            <h3>Items:</h3>
            {{range .Items}}
            <div class="item">
                <p>{{.ProductName}} x {{.Quantity}} - {{$.Currency}} {{.TotalPrice}}</p>
            </div>
            {{end}}
        </div>
        <div class="total">
            <p>Total Amount: {{.Currency}} {{.TotalPrice}}</p>
        </div>
        <div class="footer">
            <p>Thank you for your order!</p>
        </div>
    </div>
</body>
</html>`

	t, err := template.New("order_confirmation").Parse(tmpl)
	if err != nil {
		return "Error generating email template"
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, order); err != nil {
		return "Error executing email template"
	}

	return buf.String()
}

func (n *SMTPNotifier) generateOrderStatusUpdateEmail(order *domain.Order) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f8f9fa; padding: 20px; text-align: center; }
        .status-update { margin: 20px 0; }
        .footer { text-align: center; margin-top: 30px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Order Status Update</h1>
            <p>Order ID: {{.ID}}</p>
        </div>
        <div class="status-update">
            <h2>Status Update</h2>
            <p>Your order status has been updated to: <strong>{{.Status}}</strong></p>
            <p>Updated at: {{.UpdatedAt}}</p>
        </div>
        <div class="footer">
            <p>Thank you for your patience!</p>
        </div>
    </div>
</body>
</html>`

	t, err := template.New("order_status_update").Parse(tmpl)
	if err != nil {
		return "Error generating email template"
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, order); err != nil {
		return "Error executing email template"
	}

	return buf.String()
}

func (n *SMTPNotifier) generateOrderCancellationEmail(order *domain.Order) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f8f9fa; padding: 20px; text-align: center; }
        .cancellation { margin: 20px 0; }
        .footer { text-align: center; margin-top: 30px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Order Cancellation</h1>
            <p>Order ID: {{.ID}}</p>
        </div>
        <div class="cancellation">
            <h2>Cancellation Details</h2>
            <p>Your order has been cancelled.</p>
            <p>Cancelled at: {{.UpdatedAt}}</p>
        </div>
        <div class="footer">
            <p>We hope to serve you again soon!</p>
        </div>
    </div>
</body>
</html>`

	t, err := template.New("order_cancellation").Parse(tmpl)
	if err != nil {
		return "Error generating email template"
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, order); err != nil {
		return "Error executing email template"
	}

	return buf.String()
} 