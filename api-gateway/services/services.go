package services

import (
	"fmt"

	"api-gateway/config"
	commentPb "api-gateway/genproto/comment"
	postPb "api-gateway/genproto/post"
	userPb "api-gateway/genproto/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

type IServiceManager interface {
	UserService() userPb.UserServiceClient
	PostService() postPb.PostServiceClient
	CommentService() commentPb.CommentServiceClient
}

type serviceManager struct {
	userService    userPb.UserServiceClient
	postService    postPb.PostServiceClient
	commentService commentPb.CommentServiceClient
}

func (s *serviceManager) UserService() userPb.UserServiceClient {
	return s.userService
}
func (s *serviceManager) PostService() postPb.PostServiceClient {
	return s.postService
}
func (s *serviceManager) CommentService() commentPb.CommentServiceClient {
	return s.commentService
}

func NewServiceManager(conf *config.Config) (IServiceManager, error) {
	resolver.SetDefaultScheme("dns")

	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%d", conf.UserServiceHost, conf.UserServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	connPost, err := grpc.Dial(
		fmt.Sprintf("%s:%d", conf.PostServiceHost, conf.PostServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	connComment, err := grpc.Dial(
		fmt.Sprintf("%s:%d", conf.CommentServiceHost, conf.CommentServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	serviceManager := &serviceManager{
		userService:    userPb.NewUserServiceClient(connUser),
		postService:    postPb.NewPostServiceClient(connPost),
		commentService: commentPb.NewCommentServiceClient(connComment),
	}

	return serviceManager, nil
}
