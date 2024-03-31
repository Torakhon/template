package main

import (
	"comment_service/config"
	pb "comment_service/genproto/comment"
	"comment_service/pkg/db"
	"comment_service/pkg/logger"
	"comment_service/service"
	grpclient "comment_service/service/grpc_client"
	"google.golang.org/grpc"
	"net"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "user-service")
	defer func(l logger.Logger) {
		err := logger.Cleanup(l)
		if err != nil {

		}
	}(log)

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	connDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}

	grpcClient, err := grpclient.New(cfg)
	if err != nil {
		log.Fatal("grpc client dial error", logger.Error(err))
	}

	commentService := service.NewCommentService(connDB, log, grpcClient)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterCommentServiceServer(s, commentService)
	log.Info("main: server running",
		logger.String("port", cfg.RPCPort))

	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}
