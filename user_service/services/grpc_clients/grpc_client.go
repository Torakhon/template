package grpcClient

import (
	"fmt"
	"google.golang.org/grpc"
	"user_service/config"
	commentPb "user_service/genproto/comment"
	postPb "user_service/genproto/post"
)

type IServiceManager interface {
	PostService() postPb.PostServiceClient
	CommentService() commentPb.CommentServiceClient
}

type serviceManager struct {
	cfg            config.Config
	productService postPb.PostServiceClient
	commentService commentPb.CommentServiceClient
}

func New(cfg config.Config) (IServiceManager, error) {
	postConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.PostServiceHost, cfg.PostServicePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("post service dial host : %s comment :%d",
			cfg.PostServiceHost, cfg.PostServicePort)
	}
	commentConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.CommentServiceHost, cfg.CommentServicePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("comment service dial host : %s comment :%d",
			cfg.CommentServiceHost, cfg.CommentServicePort)
	}

	return &serviceManager{
		cfg:            cfg,
		productService: postPb.NewPostServiceClient(postConn),
		commentService: commentPb.NewCommentServiceClient(commentConn),
	}, nil
}

func (s *serviceManager) PostService() postPb.PostServiceClient {
	return s.productService
}

func (s *serviceManager) CommentService() commentPb.CommentServiceClient {
	return s.commentService
}
