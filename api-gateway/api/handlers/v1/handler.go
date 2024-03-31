package v1

import (
	t "api-gateway/api/tokens"
	"api-gateway/config"
	"api-gateway/pkg/logger"
	"api-gateway/services"
	"github.com/casbin/casbin/v2"
)

type HandlerV1 struct {
	log            logger.Logger
	enforcer       *casbin.Enforcer
	serviceManager services.IServiceManager
	cfg            config.Config
	jwthandler     t.JWTHandler
	//writer         producer.KafkaProducer
}

// HandlerV1Config ...
type HandlerV1Config struct {
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
	JWTHandler     t.JWTHandler
	Enforcer       *casbin.Enforcer
	//Writer         producer.KafkaProducer
}

// New ...
func New(c *HandlerV1Config) *HandlerV1 {
	return &HandlerV1{
		log:            c.Logger,
		serviceManager: c.ServiceManager,
		jwthandler:     c.JWTHandler,
		cfg:            c.Cfg,
		//writer:         c.Writer,
	}
}
