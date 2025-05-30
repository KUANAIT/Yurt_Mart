package handler

import (
	"context"
	"time"

	"github.com/hsibAD/order-service/internal/domain"
	"github.com/hsibAD/order-service/internal/events"
	pb "github.com/hsibAD/order-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	orderRepo domain.OrderRepository
	cartSub  *events.CartSubscriber
}

func NewOrderHandler(orderRepo domain.OrderRepository, cartSub *events.CartSubscriber) *OrderHandler {
	return &OrderHandler{
		orderRepo: orderRepo,
		cartSub:  cartSub,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	// Extract user ID from context or authentication token
	userID := "test-user" // TODO: Get from context

	if req.CartId == "" {
		return nil, status.Error(codes.InvalidArgument, "cart ID is required")
	}

	// Get cart information from NATS
	cartInfo, err := h.cartSub.GetCartInfo(ctx, req.CartId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get cart information")
	}

	// Convert proto delivery address to domain delivery address
	var deliveryAddr *domain.DeliveryAddress
	if req.DeliveryAddress != nil {
		var err error
		deliveryAddr, err = domain.NewDeliveryAddress(
			userID,
			"", // FullName - not provided in proto
			req.DeliveryAddress.Street,
			"", // Apartment - not provided in proto
			req.DeliveryAddress.City,
			req.DeliveryAddress.State,
			req.DeliveryAddress.PostalCode,
			req.DeliveryAddress.Country,
			"", // Phone - not provided in proto
			false, // IsDefault - not provided in proto
		)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		deliveryAddr.ID = req.DeliveryAddress.Id
	}

	// Parse delivery time
	deliveryTime := time.Now().Add(24 * time.Hour) // Default to 24 hours from now
	if req.DeliveryTime != nil {
		deliveryTime = req.DeliveryTime.AsTime()
	}

	// Create domain order
	order, err := domain.NewOrder(userID, req.CartId, deliveryAddr, deliveryTime)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Set price and currency from cart
	order.TotalPrice = cartInfo.TotalPrice
	order.Currency = cartInfo.Currency

	// Save order using repository
	if err := h.orderRepo.Create(ctx, order); err != nil {
		return nil, status.Error(codes.Internal, "failed to create order")
	}

	// Convert domain order back to proto order
	protoAddr := req.DeliveryAddress
	if order.DeliveryAddress != nil {
		protoAddr = &pb.DeliveryAddress{
			Id:         order.DeliveryAddress.ID,
			UserId:     order.DeliveryAddress.UserID,
			Street:     order.DeliveryAddress.StreetAddress,
			City:       order.DeliveryAddress.City,
			State:      order.DeliveryAddress.State,
			Country:    order.DeliveryAddress.Country,
			PostalCode: order.DeliveryAddress.PostalCode,
		}
	}

	return &pb.Order{
		Id:              order.ID,
		UserId:          order.UserID,
		CartId:          order.CartID,
		Status:          string(order.Status),
		TotalPrice:      order.TotalPrice,
		Currency:        order.Currency,
		DeliveryAddress: protoAddr,
		DeliveryTime:    timestamppb.New(order.DeliveryTime),
		CreatedAt:       timestamppb.New(order.CreatedAt),
		UpdatedAt:       timestamppb.New(order.UpdatedAt),
	}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	// TODO: Implement get order logic
	return nil, status.Error(codes.Unimplemented, "method GetOrder not implemented")
}

func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.Order, error) {
	// TODO: Implement update order status logic
	return nil, status.Error(codes.Unimplemented, "method UpdateOrderStatus not implemented")
}

func (h *OrderHandler) AddDeliveryAddress(ctx context.Context, req *pb.DeliveryAddress) (*pb.DeliveryAddress, error) {
	// TODO: Implement add delivery address logic
	return nil, status.Error(codes.Unimplemented, "method AddDeliveryAddress not implemented")
}

func (h *OrderHandler) ListDeliveryAddresses(ctx context.Context, req *pb.ListAddressesRequest) (*pb.ListAddressesResponse, error) {
	// TODO: Implement list delivery addresses logic
	return nil, status.Error(codes.Unimplemented, "method ListDeliveryAddresses not implemented")
}

func (h *OrderHandler) GetAvailableDeliverySlots(ctx context.Context, req *pb.DeliverySlotsRequest) (*pb.DeliverySlotsResponse, error) {
	if req.PostalCode == "" {
		return nil, status.Error(codes.InvalidArgument, "postal code is required")
	}

	if req.Date == nil {
		return nil, status.Error(codes.InvalidArgument, "date is required")
	}

	date := req.Date.AsTime()
	if date.Before(time.Now()) {
		return nil, status.Error(codes.InvalidArgument, "date must be in the future")
	}

	// В реальном приложении здесь будет логика проверки доступности слотов
	// на основе существующих заказов, графика работы и т.д.
	// Для демонстрации вернем фиксированные слоты
	slots := []*pb.DeliverySlot{
		{
			StartTime: timestamppb.New(req.Date.AsTime().Add(9 * time.Hour)),
			EndTime:   timestamppb.New(req.Date.AsTime().Add(11 * time.Hour)),
			Available: true,
		},
		{
			StartTime: timestamppb.New(req.Date.AsTime().Add(11 * time.Hour)),
			EndTime:   timestamppb.New(req.Date.AsTime().Add(13 * time.Hour)),
			Available: true,
		},
		{
			StartTime: timestamppb.New(req.Date.AsTime().Add(13 * time.Hour)),
			EndTime:   timestamppb.New(req.Date.AsTime().Add(15 * time.Hour)),
			Available: true,
		},
		{
			StartTime: timestamppb.New(req.Date.AsTime().Add(15 * time.Hour)),
			EndTime:   timestamppb.New(req.Date.AsTime().Add(17 * time.Hour)),
			Available: true,
		},
		{
			StartTime: timestamppb.New(req.Date.AsTime().Add(17 * time.Hour)),
			EndTime:   timestamppb.New(req.Date.AsTime().Add(19 * time.Hour)),
			Available: true,
		},
	}

	return &pb.DeliverySlotsResponse{
		PostalCode: req.PostalCode,
		Date:       req.Date,
		Slots:      slots,
	}, nil
} 