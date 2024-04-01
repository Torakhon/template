package main

import (
	"google.golang.org/grpc"
	"net"
	"user_service/config"
	pb "user_service/genproto/users"
	"user_service/pkg/db"
	"user_service/pkg/logger"
	"user_service/services"
	grpcient "user_service/services/grpc_clients"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "users-service")
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

	//client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://mongodb:27017"))
	//if err != nil {
	//	log.Fatal("connection to mongosh error", logger.Error(err))
	//}

	grpcClient, err := grpcient.New(cfg)
	if err != nil {
		log.Fatal("grpc client dial error", logger.Error(err))
	}
	// agar mongo ishlatmoqchi bo'linsa connDB orniga client beriladi
	userService := services.NewUserService(connDB, log, grpcClient)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, userService)
	log.Info("main: server running",
		logger.String("port", cfg.RPCPort))
	///
	//con, err := kaf.NewKafkaConsumer([]string{"localhost:9092"}, "test-topic", "1")
	//if err != nil {
	//	log.Fatal("cannot create a kafka producer", logger.Error(err))
	//}
	//defer con.Close()

	//err = con.ConsumeMessages(kaf.ConsumeHandler)
	//if err != nil {
	//	log.Fatal("ConsumeMessages: ", logger.Error(err))
	//}

	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}
