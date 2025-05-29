package ports

import "user-service/proto/user"

type GRPCPort interface {
	user.UserServiceServer
}
