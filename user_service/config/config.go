package config

import (
	"github.com/spf13/cast"
	"os"
)

// Config ...
type Config struct {
	Environment      string // develop, staging, production
	PostgresHost     string
	PostgresPort     int
	PostgresDatabase string
	PostgresUser     string
	PostgresPassword string

	MongoHost     string
	MongoPort     string
	MongoDatabase string

	LogLevel string
	RPCPort  string

	PostServiceHost string
	PostServicePort int

	CommentServiceHost string
	CommentServicePort int
}

// Load loads environment vars and inflates Config
func Load() Config {
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "db"))
	c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "v1"))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "20030505"))

	c.MongoHost = cast.ToString(getOrReturnDefault("MONGO_HOST", "mongoDB"))
	c.MongoPort = cast.ToString(getOrReturnDefault("MONGO_PORT", "27017"))
	c.MongoDatabase = cast.ToString(getOrReturnDefault("MONGO_DATABASE", "user"))

	c.PostServiceHost = cast.ToString(getOrReturnDefault("POST_SERVICE_HOST", "post_service"))
	c.PostServicePort = cast.ToInt(getOrReturnDefault("POST_SERVICE_PORT", "9009"))

	c.CommentServiceHost = cast.ToString(getOrReturnDefault("COMMENT_SERVICE_HOST", "comment_service"))
	c.CommentServicePort = cast.ToInt(getOrReturnDefault("COMMENT_SERVICE_PORT", "9008"))

	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))

	c.RPCPort = cast.ToString(getOrReturnDefault("RPC_PORT", ":9000"))

	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
