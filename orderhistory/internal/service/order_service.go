package service

import (
	"context"
	"log"
	"time"

	"meaningfullname/Yurt_Mart/orderhistory/internal/model"
	"meaningfullname/Yurt_Mart/orderhistory/internal/repository"
	"meaningfullname/Yurt_Mart/orderhistory/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderHistoryService struct {
	proto.UnimplementedOrderHistoryServiceServer
	repo repository.OrderRepository
}

func NewOrderHistoryService(repo repository.OrderRepository) *OrderHistoryService {
	return &OrderHistoryService{
		repo: repo,
	}
}

func (s *OrderHistoryService) AddOrderHistory(ctx context.Context, req *proto.AddOrderHistoryRequest) (*proto.AddOrderHistoryResponse, error) {
	order := &model.Order{
		ID:         uuid.New().String(),
		UserID:     req.Order.UserId,
		Username:   req.Username,
		ProductIDs: req.Order.ProductIds,
		Total:      req.Order.Total,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if s.repo == nil {
		log.Println("s.repo is nil in AddOrderHistory!")
		return nil, status.Error(codes.Internal, "repository not initialized")
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, status.Error(codes.Internal, "failed to add order history")
	}

	return &proto.AddOrderHistoryResponse{
		Order: &proto.Order{
			OrderId:    order.ID,
			UserId:     order.UserID,
			Username:   order.Username,
			ProductIds: order.ProductIDs,
			Total:      order.Total,
			Timestamp:  order.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *OrderHistoryService) GetOrderHistory(ctx context.Context, req *proto.GetOrderHistoryRequest) (*proto.GetOrderHistoryResponse, error) {
	if s.repo == nil {
		log.Println("s.repo is nil in GetOrderHistory!")
		return nil, status.Error(codes.Internal, "repository not initialized")
	}
	orders, err := s.repo.GetByUserID(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get order history")
	}

	protoOrders := make([]*proto.Order, len(orders))
	for i, order := range orders {
		protoOrders[i] = &proto.Order{
			OrderId:    order.ID,
			UserId:     order.UserID,
			Username:   order.Username,
			ProductIds: order.ProductIDs,
			Total:      order.Total,
			Timestamp:  order.CreatedAt.Format(time.RFC3339),
		}
	}

	return &proto.GetOrderHistoryResponse{
		Orders: protoOrders,
	}, nil
}

func (s *OrderHistoryService) GetOrderById(ctx context.Context, req *proto.GetOrderByIdRequest) (*proto.GetOrderByIdResponse, error) {
	if s.repo == nil {
		log.Println("s.repo is nil in GetOrderById!")
		return nil, status.Error(codes.Internal, "repository not initialized")
	}
	order, err := s.repo.GetByID(ctx, req.OrderId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "order not found")
	}

	return &proto.GetOrderByIdResponse{
		Order: &proto.Order{
			OrderId:    order.ID,
			UserId:     order.UserID,
			Username:   order.Username,
			ProductIds: order.ProductIDs,
			Total:      order.Total,
			Timestamp:  order.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *OrderHistoryService) Reorder(ctx context.Context, req *proto.ReorderRequest) (*proto.ReorderResponse, error) {
	if s.repo == nil {
		log.Println("s.repo is nil in Reorder!")
		return nil, status.Error(codes.Internal, "repository not initialized")
	}
	originalOrder, err := s.repo.GetByID(ctx, req.OrderId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "order not found")
	}

	newOrder := &model.Order{
		ID:         uuid.New().String(),
		UserID:     originalOrder.UserID,
		Username:   originalOrder.Username,
		ProductIDs: originalOrder.ProductIDs,
		Total:      originalOrder.Total,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.repo.Create(ctx, newOrder); err != nil {
		return nil, status.Error(codes.Internal, "failed to create reorder")
	}

	return &proto.ReorderResponse{
		Order: &proto.Order{
			OrderId:    newOrder.ID,
			UserId:     newOrder.UserID,
			Username:   newOrder.Username,
			ProductIds: newOrder.ProductIDs,
			Total:      newOrder.Total,
			Timestamp:  newOrder.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *OrderHistoryService) DeleteOrderHistory(ctx context.Context, req *proto.DeleteOrderHistoryRequest) (*proto.DeleteOrderHistoryResponse, error) {
	if s.repo == nil {
		log.Println("s.repo is nil in DeleteOrderHistory!")
		return nil, status.Error(codes.Internal, "repository not initialized")
	}
	if err := s.repo.Delete(ctx, req.OrderId); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete order history")
	}

	return &proto.DeleteOrderHistoryResponse{
		Success: true,
	}, nil
}

func (s *OrderHistoryService) GetRecentOrders(ctx context.Context, req *proto.GetRecentOrdersRequest) (*proto.GetRecentOrdersResponse, error) {
	if s.repo == nil {
		log.Println("s.repo is nil in GetRecentOrders!")
		return nil, status.Error(codes.Internal, "repository not initialized")
	}
	orders, err := s.repo.GetRecentByUserID(ctx, req.UserId, int(req.Limit))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get recent orders")
	}

	protoOrders := make([]*proto.Order, len(orders))
	for i, order := range orders {
		protoOrders[i] = &proto.Order{
			OrderId:    order.ID,
			UserId:     order.UserID,
			Username:   order.Username,
			ProductIds: order.ProductIDs,
			Total:      order.Total,
			Timestamp:  order.CreatedAt.Format(time.RFC3339),
		}
	}

	return &proto.GetRecentOrdersResponse{
		Orders: protoOrders,
	}, nil
}
