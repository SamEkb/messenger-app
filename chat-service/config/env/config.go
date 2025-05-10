package env

import (
	"log"
	"os"
	"strconv"
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
}

type ServerConfig struct {
	GRPCHost string
	GRPCPort int
	HTTPHost string
	HTTPPort int
}

type ClientsConfig struct {
	Users   *ServiceClientConfig
	Friends *ServiceClientConfig
	GRPC    *GRPCClientConfig
}

type ServiceClientConfig struct {
	Host string
	Port int
}

type GRPCClientConfig struct {
	ConnectionTimeout time.Duration
	RetryAttempts     int
	RetryDelay        time.Duration
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
			Users:   &ServiceClientConfig{},
			Friends: &ServiceClientConfig{},
			GRPC:    &GRPCClientConfig{},
		},
	}

	c.Server.GRPCHost = getEnv("GRPC_HOST", "0.0.0.0")
	c.Server.GRPCPort = getEnvAsInt("GRPC_PORT", DefaultGRPCPort)
	c.Server.HTTPHost = getEnv("HTTP_HOST", "0.0.0.0")
	c.Server.HTTPPort = getEnvAsInt("HTTP_PORT", DefaultHTTPPort)

	c.Clients.Users.Host = getEnv("USERS_SERVICE_HOST", "localhost")
	c.Clients.Users.Port = getEnvAsInt("USERS_SERVICE_PORT", 9004)

	c.Clients.Friends.Host = getEnv("FRIENDS_SERVICE_HOST", "localhost")
	c.Clients.Friends.Port = getEnvAsInt("FRIENDS_SERVICE_PORT", 9003)

	c.Clients.GRPC.ConnectionTimeout = time.Duration(getEnvAsInt("GRPC_CONNECTION_TIMEOUT_SEC", 5)) * time.Second
	c.Clients.GRPC.RetryAttempts = getEnvAsInt("GRPC_RETRY_ATTEMPTS", 3)
	c.Clients.GRPC.RetryDelay = time.Duration(getEnvAsInt("GRPC_RETRY_DELAY_SEC", 1)) * time.Second

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
