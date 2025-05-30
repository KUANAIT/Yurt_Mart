package proxy

import (
	"context"
	"time"

	pb "github.com/hsibAD/payment-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentServiceClient struct {
	client pb.PaymentServiceClient
	conn   *grpc.ClientConn
}

func NewPaymentServiceClient(serviceURL string) (*PaymentServiceClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serviceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return &PaymentServiceClient{
		client: pb.NewPaymentServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *PaymentServiceClient) InitiatePayment(ctx context.Context, req *pb.InitiatePaymentRequest) (*pb.Payment, error) {
	return c.client.InitiatePayment(ctx, req)
}

func (c *PaymentServiceClient) ProcessCreditCardPayment(ctx context.Context, req *pb.CreditCardPaymentRequest) (*pb.Payment, error) {
	return c.client.ProcessCreditCardPayment(ctx, req)
}

func (c *PaymentServiceClient) InitiateMetaMaskPayment(ctx context.Context, req *pb.MetaMaskPaymentRequest) (*pb.MetaMaskPaymentResponse, error) {
	return c.client.InitiateMetaMaskPayment(ctx, req)
}

func (c *PaymentServiceClient) ConfirmMetaMaskPayment(ctx context.Context, req *pb.ConfirmMetaMaskPaymentRequest) (*pb.Payment, error) {
	return c.client.ConfirmMetaMaskPayment(ctx, req)
}

func (c *PaymentServiceClient) GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.Payment, error) {
	return c.client.GetPayment(ctx, req)
}

func (c *PaymentServiceClient) GetPaymentsByOrder(ctx context.Context, req *pb.GetPaymentsByOrderRequest) (*pb.GetPaymentsByOrderResponse, error) {
	return c.client.GetPaymentsByOrder(ctx, req)
}

func (c *PaymentServiceClient) UpdatePaymentStatus(ctx context.Context, req *pb.UpdatePaymentStatusRequest) (*pb.Payment, error) {
	return c.client.UpdatePaymentStatus(ctx, req)
}

func (c *PaymentServiceClient) GetPendingPayments(ctx context.Context, req *pb.GetPendingPaymentsRequest) (*pb.GetPendingPaymentsResponse, error) {
	return c.client.GetPendingPayments(ctx, req)
}

func (c *PaymentServiceClient) RetryPayment(ctx context.Context, req *pb.RetryPaymentRequest) (*pb.Payment, error) {
	return c.client.RetryPayment(ctx, req)
}

func (c *PaymentServiceClient) Close() error {
	return c.conn.Close()
}
