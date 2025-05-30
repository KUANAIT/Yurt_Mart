package handler

import (
	"context"

	"github.com/hsibAD/payment-service/internal/domain"
	pb "github.com/hsibAD/payment-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentHandler struct {
	pb.UnimplementedPaymentServiceServer
	paymentRepo domain.PaymentRepository
}

func RegisterServices(s *grpc.Server, cfg interface{}) {
	pb.RegisterPaymentServiceServer(s, &PaymentHandler{})
}

func (h *PaymentHandler) InitiatePayment(ctx context.Context, req *pb.InitiatePaymentRequest) (*pb.Payment, error) {
	if req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "order ID is required")
	}

	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	if req.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "amount must be greater than 0")
	}

	// Create domain payment
	payment, err := domain.NewPayment(
		req.OrderId,
		req.UserId,
		req.Amount,
		req.Currency,
		domain.PaymentMethod(req.PaymentMethod.String()),
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Save payment using repository
	if err := h.paymentRepo.Create(ctx, payment); err != nil {
		return nil, status.Error(codes.Internal, "failed to create payment")
	}

	// Convert domain payment to proto payment
	return &pb.Payment{
		Id:            payment.ID,
		OrderId:       payment.OrderID,
		UserId:        payment.UserID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		Status:        pb.PaymentStatus(pb.PaymentStatus_value[string(payment.Status)]),
		PaymentMethod: req.PaymentMethod,
		CreatedAt:     nil, // TODO: Add timestamp conversion
		UpdatedAt:     nil, // TODO: Add timestamp conversion
	}, nil
}

func (h *PaymentHandler) ProcessCreditCardPayment(ctx context.Context, req *pb.CreditCardPaymentRequest) (*pb.Payment, error) {
	// TODO: Implement credit card payment processing logic
	return nil, status.Error(codes.Unimplemented, "method ProcessCreditCardPayment not implemented")
}

func (h *PaymentHandler) InitiateMetaMaskPayment(ctx context.Context, req *pb.MetaMaskPaymentRequest) (*pb.MetaMaskPaymentResponse, error) {
	// TODO: Implement MetaMask payment initiation logic
	return nil, status.Error(codes.Unimplemented, "method InitiateMetaMaskPayment not implemented")
}

func (h *PaymentHandler) ConfirmMetaMaskPayment(ctx context.Context, req *pb.ConfirmMetaMaskPaymentRequest) (*pb.Payment, error) {
	// TODO: Implement MetaMask payment confirmation logic
	return nil, status.Error(codes.Unimplemented, "method ConfirmMetaMaskPayment not implemented")
}

func (h *PaymentHandler) GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.Payment, error) {
	// TODO: Implement get payment logic
	return nil, status.Error(codes.Unimplemented, "method GetPayment not implemented")
}

func (h *PaymentHandler) GetPaymentsByOrder(ctx context.Context, req *pb.GetPaymentsByOrderRequest) (*pb.GetPaymentsByOrderResponse, error) {
	// TODO: Implement get payments by order logic
	return nil, status.Error(codes.Unimplemented, "method GetPaymentsByOrder not implemented")
}

func (h *PaymentHandler) GetPendingPayments(ctx context.Context, req *pb.GetPendingPaymentsRequest) (*pb.GetPendingPaymentsResponse, error) {
	// TODO: Implement get pending payments logic
	return nil, status.Error(codes.Unimplemented, "method GetPendingPayments not implemented")
}

func (h *PaymentHandler) UpdatePaymentStatus(ctx context.Context, req *pb.UpdatePaymentStatusRequest) (*pb.Payment, error) {
	// TODO: Implement update payment status logic
	return nil, status.Error(codes.Unimplemented, "method UpdatePaymentStatus not implemented")
}

func (h *PaymentHandler) RetryPayment(ctx context.Context, req *pb.RetryPaymentRequest) (*pb.Payment, error) {
	// TODO: Implement retry payment logic
	return nil, status.Error(codes.Unimplemented, "method RetryPayment not implemented")
} 