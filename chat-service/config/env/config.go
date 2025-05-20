package env

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	DefaultGRPCPort = 9002
	DefaultHTTPPort = 8002
)

type Config struct {
	AppName string
	Debug   string
	Server  *ServerConfig
	Clients *ClientsConfig
	MongoDB *MongoDBConfig
}

type ServerConfig struct {
	GRPCHost string
	GRPCPort int
	HTTPHost string
	HTTPPort int
}

type ClientsConfig struct {
	Users      *ServiceClientConfig
	Friends    *ServiceClientConfig
	MaxRetries int
	RetryDelay time.Duration

	CircuitBreaker *CircuitBreakerConfig
}

type CircuitBreakerConfig struct {
	Enabled          bool
	Name             string
	MaxRequests      uint32
	Interval         time.Duration
	Timeout          time.Duration
	MinRequests      uint32
	FailureRatio     float64
	ServerErrorCodes []string
}

type MongoDBConfig struct {
	URI      string
	Database string
}

func (m *MongoDBConfig) ConnectionString() string {
	return fmt.Sprintf("%s/%s", m.URI, m.Database)
}

type ServiceClientConfig struct {
	Host string
	Port int
}

func (s *ServerConfig) GrpcAddr() string {
	return s.GRPCHost + ":" + strconv.Itoa(s.GRPCPort)
}

func (s *ServerConfig) HttpAddr() string {
	return s.HTTPHost + ":" + strconv.Itoa(s.HTTPPort)
}

func (c *ServiceClientConfig) Addr() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Info: .env file not found or couldn't be loaded; using environment variables")
	}

	c := &Config{
		AppName: getEnv("APP_NAME", "ChatService"),
		Debug:   getEnv("DEBUG", "dev"),
		Server:  &ServerConfig{},
		Clients: &ClientsConfig{
			Users:      &ServiceClientConfig{},
			Friends:    &ServiceClientConfig{},
			MaxRetries: 3,
			RetryDelay: 100 * time.Millisecond,
			CircuitBreaker: &CircuitBreakerConfig{
				Enabled:      getEnvAsBool("CB_ENABLED", true),
				Name:         getEnv("CB_NAME", "grpc_circuit_breaker"),
				MaxRequests:  uint32(getEnvAsInt("CB_MAX_REQUESTS", 10)),
				Interval:     time.Duration(getEnvAsInt("CB_INTERVAL_SEC", 60)) * time.Second,
				Timeout:      time.Duration(getEnvAsInt("CB_TIMEOUT_SEC", 300)) * time.Second,
				MinRequests:  uint32(getEnvAsInt("CB_MIN_REQUESTS", 40)),
				FailureRatio: getEnvAsFloat("CB_FAILURE_RATIO", 0.6),
				ServerErrorCodes: getEnvAsStringSlice("CB_SERVER_ERROR_CODES", []string{
					"INTERNAL", "UNAVAILABLE", "DATA_LOSS", "DEADLINE_EXCEEDED",
					"RESOURCE_EXHAUSTED", "UNKNOWN", "ABORTED",
				}),
			},
		},
		MongoDB: &MongoDBConfig{},
	}

	c.Server.GRPCHost = getEnv("GRPC_HOST", "0.0.0.0")
	c.Server.GRPCPort = getEnvAsInt("GRPC_PORT", DefaultGRPCPort)
	c.Server.HTTPHost = getEnv("HTTP_HOST", "0.0.0.0")
	c.Server.HTTPPort = getEnvAsInt("HTTP_PORT", DefaultHTTPPort)

	c.Clients.Users.Host = getEnv("USERS_SERVICE_HOST", "localhost")
	c.Clients.Users.Port = getEnvAsInt("USERS_SERVICE_PORT", 9004)

	c.Clients.Friends.Host = getEnv("FRIENDS_SERVICE_HOST", "localhost")
	c.Clients.Friends.Port = getEnvAsInt("FRIENDS_SERVICE_PORT", 9003)

	c.Clients.MaxRetries = getEnvAsInt("MAX_RETRIES", 3)
	c.Clients.RetryDelay = time.Duration(getEnvAsInt("RETRY_DELAY", 100)) * time.Millisecond

	c.MongoDB.URI = getEnv("MONGODB_URI", "mongodb://localhost:27017")
	c.MongoDB.Database = getEnv("MONGODB_DATABASE", "chat_db")

	return c, nil
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if v := os.Getenv(key); v != "" {
		val, err := strconv.Atoi(v)
		if err == nil {
			return val
		}
	}
	return defaultValue
}

func getEnvAsStringSlice(key string, defaultValue []string) []string {
	if v := os.Getenv(key); v != "" {
		return strings.Split(v, ",")
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if v := os.Getenv(key); v != "" {
		val, err := strconv.ParseBool(v)
		if err == nil {
			return val
		}
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if v := os.Getenv(key); v != "" {
		val, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return val
		}
	}
	return defaultValue
}
