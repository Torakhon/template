package services

import (
	"context"
	"github.com/jmoiron/sqlx"
	pb "user_service/genproto/users"
	l "user_service/pkg/logger"
	grpcClient "user_service/services/grpc_clients"
	"user_service/storage"
)

// UserService ...
type UserService struct {
	storage    storage.IStorage
	logger     l.Logger
	grpcClient grpcClient.IServiceManager
}

// NewUserService ...
func NewUserService(db *sqlx.DB, log l.Logger, grpcClient grpcClient.IServiceManager) *UserService {
	return &UserService{
		storage:    storage.NewStoragePg(db),
		logger:     log,
		grpcClient: grpcClient,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.User, error) {
	return s.storage.User().CreateUser(ctx, req)
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.User, error) {
	return s.storage.User().GetUser(ctx, req)
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserReq) (*pb.User, error) {
	return s.storage.User().UpdateUser(ctx, req)
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserReq) (*pb.DeleteUserRes, error) {
	return s.storage.User().DeleteUser(ctx, req)
}

func (s *UserService) GetAllUsers(ctx context.Context, req *pb.GetAllUsersReq) (*pb.GetAllUsersRes, error) {
	return s.storage.User().GetAllUsers(ctx, req)
}

func (s *UserService) CheckUniques(ctx context.Context, req *pb.CheckUniqReq) (*pb.CheckUniqRes, error) {
	return s.storage.User().CheckUniques(ctx, req)
}

func (s *UserService) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	return s.storage.User().Login(ctx, req)
}

func (s *UserService) UpdateRole(ctx context.Context, req *pb.UpdateRoleReq) (*pb.UpdateRoleRes, error) {
	return s.storage.User().UpdateRole(ctx, req)
}

func (s *UserService) UpdateEmail(ctx context.Context, req *pb.UpdateEmailReq) (*pb.UpdateEmailRes, error) {
	return s.storage.User().UpdateEmail(ctx, req)
}
