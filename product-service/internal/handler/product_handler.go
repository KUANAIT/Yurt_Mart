package handler

import (
	"context"
	pb "product-service/client-service/proto/productpb"
	"product-service/internal/domain"
	"product-service/internal/repository"
	"product-service/internal/usecase"
)

type ProductHandler struct {
	pb.UnimplementedProductServiceServer
	usecase usecase.ProductUsecase
}

func NewProductHandler() *ProductHandler {
	repo := repository.NewMongoProductRepository()
	use := usecase.NewProductUsecase(repo)
	return &ProductHandler{usecase: use}
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	product := &domain.Product{
		Name:        req.Product.Name,
		Description: req.Product.Description,
		Category:    req.Product.Category,
		Price:       req.Product.Price,
		Quantity:    req.Product.Quantity,
		UserID:      req.Product.UserId,
	}
	id, err := h.usecase.Create(product)
	if err != nil {
		return nil, err
	}
	return &pb.CreateProductResponse{Id: id}, nil
}

func (h *ProductHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	product, err := h.usecase.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetProductResponse{
		Product: &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Category:    product.Category,
			Price:       product.Price,
			Quantity:    product.Quantity,
			UserId:      product.UserID,
		},
	}, nil
}

func (h *ProductHandler) GetProductsByCategory(ctx context.Context, req *pb.GetProductsByCategoryRequest) (*pb.GetProductsByCategoryResponse, error) {
	products, err := h.usecase.GetProductsByCategory(req.Category)
	if err != nil {
		return nil, err
	}

	var pbProducts []*pb.Product
	for _, p := range products {
		pbProducts = append(pbProducts, &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Category:    p.Category,
			Price:       p.Price,
			Quantity:    p.Quantity,
			UserId:      p.UserID,
		})
	}

	return &pb.GetProductsByCategoryResponse{Products: pbProducts}, nil
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	err := h.usecase.Delete(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteProductResponse{Message: "Product deleted"}, nil
}

func (h *ProductHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, err := h.usecase.List()
	if err != nil {
		return nil, err
	}

	var pbProducts []*pb.Product
	for _, p := range products {
		pbProducts = append(pbProducts, &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Category:    p.Category,
			Price:       p.Price,
			Quantity:    p.Quantity,
			UserId:      p.UserID,
		})
	}

	return &pb.ListProductsResponse{Products: pbProducts}, nil
}
