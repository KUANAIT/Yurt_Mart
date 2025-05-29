package clients

import (
	"context"

	"api-gateway/internal/core/ports"
	"google.golang.org/grpc"
	"user-service/proto/user"
)

type UserServiceClient struct {
	client ports.UserServiceClient
	conn   *grpc.ClientConn
}

func NewUserServiceClient(address string) (*UserServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &UserServiceClient{
		client: ports.NewUserServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *UserServiceClient) RegisterUser(ctx context.Context, email, password, name string) (string, error) {
	resp, err := c.client.RegisterUser(ctx, &user.RegisterUserRequest{
		Email:    email,
		Password: password,
		Name:     name,
	})
	if err != nil {
		return "", err
	}
	return resp.UserId, nil
}

func (c *UserServiceClient) GetUser(ctx context.Context, userID string) (*ports.UserResponse, error) {
	resp, err := c.client.GetUser(ctx, &user.GetUserRequest{UserId: userID})
	if err != nil {
		return nil, err
	}

	return &ports.UserResponse{
		ID:    resp.UserId,
		Email: resp.Email,
		Name:  resp.Name,
	}, nil
}

func (c *UserServiceClient) UpdateUser(ctx context.Context, userID, email, name string) (*ports.UserResponse, error) {
	resp, err := c.client.UpdateUser(ctx, &user.UpdateUserRequest{
		UserId: userID,
		Email:  email,
		Name:   name,
	})
	if err != nil {
		return nil, err
	}

	return &ports.UserResponse{
		ID:    resp.UserId,
		Email: resp.Email,
		Name:  resp.Name,
	}, nil
}

func (c *UserServiceClient) DeleteUser(ctx context.Context, userID string) error {
	_, err := c.client.DeleteUser(ctx, &user.DeleteUserRequest{
		UserId: userID,
	})
	return err
}

func (c *UserServiceClient) Close() error {
	return c.conn.Close()
}
