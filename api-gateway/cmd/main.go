package main

import (
	"api-gateway/api"
	"api-gateway/config"
	"api-gateway/pkg/logger"
	"api-gateway/services"
	"github.com/casbin/casbin/v2"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "api-gateway_ali")

	serviceManager, err := services.NewServiceManager(&cfg)
	if err != nil {
		log.Error("gRPC dial error", logger.Error(err))
	}

	//writer, err := producer.NewKafkaProducer([]string{"kafka:9092"})
	//if err != nil {
	//	log.Error("NewKafkaProducer: ", logger.Error(err))
	//}
	//defer writer.Close()

	enforcer, err := casbin.NewEnforcer("config/auth.conf", "config/auth.csv")
	if err != nil {
		log.Error("NewEnforcer error", logger.Error(err))
	}

	server := api.New(api.Option{
		Conf:           cfg,
		Logger:         log,
		Enforcer:       enforcer,
		ServiceManager: serviceManager,
		//Writer:         writer,
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("failed to run http server", logger.Error(err))
		panic(err.Error())
	}
}
