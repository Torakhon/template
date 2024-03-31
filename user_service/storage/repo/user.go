package repo

import (
	"context"
	pb "user_service/genproto/users"
)

// UserStorageI ...
type UserStorageI interface {
	CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.User, error)
	GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.User, error)
	UpdateUser(ctx context.Context, req *pb.UpdateUserReq) (*pb.User, error)
	DeleteUser(ctx context.Context, req *pb.DeleteUserReq) (*pb.DeleteUserRes, error)
	GetAllUsers(ctx context.Context, req *pb.GetAllUsersReq) (*pb.GetAllUsersRes, error)
	CheckUniques(ctx context.Context, req *pb.CheckUniqReq) (*pb.CheckUniqRes, error)
	Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error)
	UpdateRole(ctx context.Context, req *pb.UpdateRoleReq) (*pb.UpdateRoleRes, error)
	UpdateEmail(ctx context.Context, req *pb.UpdateEmailReq) (*pb.UpdateEmailRes, error)
}
