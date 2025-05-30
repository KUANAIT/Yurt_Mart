package auth

import (
	"context"
	"google.golang.org/grpc"
	pb "user-service/proto/auth"
)

type AuthServiceClient interface {
	IsTokenBlacklisted(token string) (bool, error)
}

type authServiceClient struct {
	client pb.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthServiceClient(address string) (AuthServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &authServiceClient{
		client: pb.NewAuthServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *authServiceClient) IsTokenBlacklisted(token string) (bool, error) {
	// For now, we'll just check if the token is valid
	// In a real implementation, you would make a gRPC call to the auth service
	// to check if the token is blacklisted
	_, err := c.client.Logout(context.Background(), &pb.LogoutRequest{
		Token: token,
	})

	if err != nil {
		return false, err
	}

	return false, nil
}

func (c *authServiceClient) Close() error {
	return c.conn.Close()
}
