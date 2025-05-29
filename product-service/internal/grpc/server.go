package grpc

import (
	"google.golang.org/grpc"
	pb "product-service/client-service/proto/productpb"
	"product-service/internal/handler"
)

func RegisterProductServiceServer(s *grpc.Server) {
	pb.RegisterProductServiceServer(s, handler.NewProductHandler())
}
