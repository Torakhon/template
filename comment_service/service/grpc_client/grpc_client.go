package grpcClient

import (
	"comment_service/config"
	postPb "comment_service/genproto/post"
	userPb "comment_service/genproto/users"
	"fmt"
	"google.golang.org/grpc"
)

type IserviveManager interface {
	UserService() userPb.UserServiceClient
	PostService() postPb.PostServiceClient
}

type serviceManager struct {
	cfg         config.Config
	userService userPb.UserServiceClient
	postService postPb.PostServiceClient
}

func New(cfg config.Config) (IserviveManager, error) {
	userConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.UserServiceHost, cfg.UserServicePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("comment service dial host : %s comment :%d",
			cfg.UserServiceHost, cfg.UserServicePort)
	}

	postConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.PostServiceHost, cfg.PostServicePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("comment service dial host : %s comment :%d",
			cfg.PostServiceHost, cfg.PostServicePort)
	}

	return &serviceManager{
		cfg:         cfg,
		userService: userPb.NewUserServiceClient(userConn),
		postService: postPb.NewPostServiceClient(postConn),
	}, nil
}

func (s *serviceManager) UserService() userPb.UserServiceClient {
	return s.userService
}

func (s *serviceManager) PostService() postPb.PostServiceClient {
	return s.postService
}
