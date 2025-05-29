package integration_test

import (
	"context"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc"
	pb "shopping-cart-service/shopping-cart-service/proto/cartpb"

	"github.com/stretchr/testify/assert"
)

const grpcAddr = "localhost:50052"

func getClientConn(t *testing.T) pb.CartServiceClient {
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial gRPC: %v", err)
	}
	return pb.NewCartServiceClient(conn)
}

func TestCartIntegration(t *testing.T) {
	client := getClientConn(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userID := "test-user"
	productID := "test-product"

	addReq := &pb.AddToCartRequest{
		UserId: userID,
		Items: []*pb.CartItem{
			{ProductId: productID, Quantity: 2},
		},
	}
	_, err := client.AddToCart(ctx, addReq)
	assert.NoError(t, err, "AddToCart failed")

	getResp, err := client.GetCart(ctx, &pb.GetCartRequest{UserId: userID})
	assert.NoError(t, err, "GetCart failed")
	assert.NotEmpty(t, getResp.Items, "GetCart should return items")

	_, err = client.RemoveFromCart(ctx, &pb.RemoveFromCartRequest{
		UserId:    userID,
		ProductId: productID,
	})
	assert.NoError(t, err, "RemoveFromCart failed")

	getRespAfter, _ := client.GetCart(ctx, &pb.GetCartRequest{UserId: userID})
	assert.Equal(t, 0, len(getRespAfter.Items), "Cart should be empty after removal")

	log.Println("✅ Integration test passed for CartService")
}
