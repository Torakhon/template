package grpcClient

import (
	"fmt"
	"google.golang.org/grpc"
	"post_service/config"
	commentPb "post_service/genproto/comment"
	userPb "post_service/genproto/users"
)

type IserviveManager interface {
	UserService() userPb.UserServiceClient
	CommentService() commentPb.CommentServiceClient
}

type serviceManager struct {
	cfg            config.Config
	userService    userPb.UserServiceClient
	commentService commentPb.CommentServiceClient
}

func New(cfg config.Config) (IserviveManager, error) {
	userConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.UserServiceHost, cfg.UserServicePost),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("comment service dial host : %s comment :%d",
			cfg.UserServiceHost, cfg.UserServicePost)
	}

	commentConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.CommentServiceHost, cfg.CommentServicePost),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("comment service dial host : %s comment :%d",
			cfg.CommentServiceHost, cfg.CommentServicePost)
	}

	return &serviceManager{
		cfg:            cfg,
		userService:    userPb.NewUserServiceClient(userConn),
		commentService: commentPb.NewCommentServiceClient(commentConn),
	}, nil
}

func (s *serviceManager) UserService() userPb.UserServiceClient {
	return s.userService
}

func (s *serviceManager) CommentService() commentPb.CommentServiceClient {
	return s.commentService
}
