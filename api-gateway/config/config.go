package config

import (
	"os"

	"github.com/spf13/cast"
)

// Config ...
type Config struct {
	Environment string // develop, staging, production

	UserServiceHost string
	UserServicePort int

	ProductServiceHost string
	ProductServicePort int

	MongoHost     string
	MongoPort     string
	MongoDatabase string

	SigningKey string

	AuthConfigPath string
	CSVFilePath    string

	// context timeout in seconds
	CtxTimeout int

	LogLevel string
	HTTPPort string
}

// Load loads environment vars and inflates Config
func Load() Config {
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":9091"))

	c.UserServiceHost = cast.ToString(getOrReturnDefault("USER_SERVICE_HOST", "user-service_ali"))
	c.UserServicePort = cast.ToInt(getOrReturnDefault("USER_SERVICE_PORT", 9008))

	c.ProductServiceHost = cast.ToString(getOrReturnDefault("PRODUCT_SERVICE_HOST", "product-service_ali"))
	c.ProductServicePort = cast.ToInt(getOrReturnDefault("PRODUCT_SERVICE_PORT", 8008))

	c.AuthConfigPath = cast.ToString(getOrReturnDefault("AUTH_CONFIG_PATH", "./config/auth.conf"))
	c.CSVFilePath = cast.ToString(getOrReturnDefault("CSV_FILE_PATH", "./config/auth.csv"))

	c.CtxTimeout = cast.ToInt(getOrReturnDefault("CTX_TIMEOUT", 7))
	c.SigningKey = cast.ToString(getOrReturnDefault("SIGNINGKEY", "mrx"))

	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
