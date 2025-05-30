package proxy

import (
	"context"
	"time"

	pb "github.com/hsibAD/order-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderServiceClient struct {
	client pb.OrderServiceClient
	conn   *grpc.ClientConn
}

func NewOrderServiceClient(serviceURL string) (*OrderServiceClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serviceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return &OrderServiceClient{
		client: pb.NewOrderServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *OrderServiceClient) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	return c.client.CreateOrder(ctx, req)
}

func (c *OrderServiceClient) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	return c.client.GetOrder(ctx, req)
}

func (c *OrderServiceClient) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.Order, error) {
	return c.client.UpdateOrderStatus(ctx, req)
}

func (c *OrderServiceClient) AddDeliveryAddress(ctx context.Context, req *pb.DeliveryAddress) (*pb.DeliveryAddress, error) {
	return c.client.AddDeliveryAddress(ctx, req)
}

func (c *OrderServiceClient) UpdateDeliveryAddress(ctx context.Context, req *pb.DeliveryAddress) (*pb.DeliveryAddress, error) {
	return c.client.UpdateDeliveryAddress(ctx, req)
}

func (c *OrderServiceClient) DeleteDeliveryAddress(ctx context.Context, req *pb.DeleteAddressRequest) error {
	_, err := c.client.DeleteDeliveryAddress(ctx, req)
	return err
}

func (c *OrderServiceClient) ListDeliveryAddresses(ctx context.Context, req *pb.ListAddressesRequest) (*pb.ListAddressesResponse, error) {
	return c.client.ListDeliveryAddresses(ctx, req)
}

func (c *OrderServiceClient) SetDeliveryTime(ctx context.Context, req *pb.SetDeliveryTimeRequest) (*pb.Order, error) {
	return c.client.SetDeliveryTime(ctx, req)
}

func (c *OrderServiceClient) GetAvailableDeliverySlots(ctx context.Context, req *pb.DeliverySlotsRequest) (*pb.DeliverySlotsResponse, error) {
	return c.client.GetAvailableDeliverySlots(ctx, req)
}

func (c *OrderServiceClient) Close() error {
	return c.conn.Close()
}
