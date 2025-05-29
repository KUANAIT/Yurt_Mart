package handler

import (
	"context"
	"shopping-cart-service/events"
	"shopping-cart-service/internal/model"
	"shopping-cart-service/internal/service"
	pb "shopping-cart-service/shopping-cart-service/proto/cartpb"
)

type CartHandler struct {
	pb.UnimplementedCartServiceServer
	service *service.CartService
}

func NewCartHandler(s *service.CartService) *CartHandler {
	return &CartHandler{service: s}
}

func (h *CartHandler) AddToCart(ctx context.Context, req *pb.AddToCartRequest) (*pb.CartResponse, error) {
	for _, item := range req.Items {
		err := h.service.AddToCart(ctx, model.CartItem{
			UserID:    req.UserId,
			ProductID: item.ProductId,
			Quantity:  int32(item.Quantity),
		})
		if err != nil {
			return nil, err
		}

		events.PublishCartItemAdded(events.CartItemAddedEvent{
			UserID:    req.UserId, // в этом случае здесь будет email
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		})

	}

	items, _ := h.service.GetCart(ctx, req.UserId)
	return toCartResponse(items), nil
}

func (h *CartHandler) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.CartResponse, error) {
	items, err := h.service.GetCart(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return toCartResponse(items), nil
}

func (h *CartHandler) RemoveFromCart(ctx context.Context, req *pb.RemoveFromCartRequest) (*pb.CartResponse, error) {
	err := h.service.RemoveFromCart(ctx, req.UserId, req.ProductId)
	if err != nil {
		return nil, err
	}
	items, _ := h.service.GetCart(ctx, req.UserId)
	return toCartResponse(items), nil
}

func toCartResponse(items []model.CartItem) *pb.CartResponse {
	var respItems []*pb.CartItem
	for _, item := range items {
		respItems = append(respItems, &pb.CartItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		})
	}
	return &pb.CartResponse{Items: respItems}
}
