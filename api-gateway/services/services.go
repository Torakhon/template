package services

import (
	"fmt"

	"api-gateway/config"
	prPb "api-gateway/genproto/product"
	userPb "api-gateway/genproto/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

type IServiceManager interface {
	UserService() userPb.UserServiceClient
	ProductService() prPb.ProductServiceClient
}

type serviceManager struct {
	userService    userPb.UserServiceClient
	productService prPb.ProductServiceClient
}

func (s *serviceManager) UserService() userPb.UserServiceClient {
	return s.userService
}
func (s *serviceManager) ProductService() prPb.ProductServiceClient {
	return s.productService
}

func NewServiceManager(conf *config.Config) (IServiceManager, error) {
	resolver.SetDefaultScheme("dns")

	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%d", conf.UserServiceHost, conf.UserServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	connProduct, err := grpc.Dial(
		fmt.Sprintf("%s:%d", conf.ProductServiceHost, conf.ProductServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	serviceManager := &serviceManager{
		userService:    userPb.NewUserServiceClient(connUser),
		productService: prPb.NewProductServiceClient(connProduct),
	}

	return serviceManager, nil
}
