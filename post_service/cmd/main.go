package main

import (
	"google.golang.org/grpc"
	"net"
	"post_service/config"
	pb "post_service/genproto/post"
	"post_service/pkg/db"
	"post_service/pkg/logger"
	"post_service/service"
	grpcient "post_service/service/grpc_client"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "post_service")
	defer func(l logger.Logger) {
		err := logger.Cleanup(l)
		if err != nil {
			return
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

	//client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://mongodb:27017"))
	//if err != nil {
	//	log.Fatal("connection to mongosh error", logger.Error(err))
	//}

	grpcClient, err := grpcient.New(cfg)
	if err != nil {
		log.Fatal("grpc client dial error", logger.Error(err))
	}

	// agar mongo ishlatmoqchi bo'linsa connDB orniga client beriladi

	postService := service.NewPostService(connDB, log, grpcClient)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterPostServiceServer(s, postService)
	log.Info("main: server running",
		logger.String("port", cfg.RPCPort))

	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}
