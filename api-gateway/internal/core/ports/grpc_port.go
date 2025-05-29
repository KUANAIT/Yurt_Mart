package ports

import (
	"context"
	"google.golang.org/grpc"
	"user-service/proto/user"
)

type UserServiceClient interface {
	RegisterUser(ctx context.Context, in *user.RegisterUserRequest, opts ...grpc.CallOption) (*user.RegisterUserResponse, error)
	GetUser(ctx context.Context, in *user.GetUserRequest, opts ...grpc.CallOption) (*user.GetUserResponse, error)
	UpdateUser(ctx context.Context, in *user.UpdateUserRequest, opts ...grpc.CallOption) (*user.UpdateUserResponse, error)
	DeleteUser(ctx context.Context, in *user.DeleteUserRequest, opts ...grpc.CallOption) (*user.DeleteUserResponse, error)
}

func NewUserServiceClient(conn grpc.ClientConnInterface) UserServiceClient {
	return user.NewUserServiceClient(conn)
}
