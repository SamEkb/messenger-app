package env

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	DefaultGRPCPort = 9003
	DefaultHTTPPort = 8003
)

type Config struct {
	AppName string
	Debug   string
	Server  *ServerConfig
	Clients *ClientsConfig
	DB      *DBConfig
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
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func (db *DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User, db.Password, db.Host, db.Port, db.Name)
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
		AppName: getEnv("APP_NAME", "FriendsService"),
		Debug:   getEnv("DEBUG", "dev"),
		Server:  &ServerConfig{},
		Clients: &ClientsConfig{
			Users:   &ServiceClientConfig{},
			Friends: &ServiceClientConfig{},
		},
		DB: &DBConfig{},
	}

	c.Server.GRPCHost = getEnv("GRPC_HOST", "0.0.0.0")
	c.Server.GRPCPort = getEnvAsInt("GRPC_PORT", DefaultGRPCPort)
	c.Server.HTTPHost = getEnv("HTTP_HOST", "0.0.0.0")
	c.Server.HTTPPort = getEnvAsInt("HTTP_PORT", DefaultHTTPPort)

	c.Clients.Users.Host = getEnv("USERS_SERVICE_HOST", "localhost")
	c.Clients.Users.Port = getEnvAsInt("USERS_SERVICE_PORT", 9004)

	c.DB = &DBConfig{
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnvAsInt("POSTGRES_PORT", 5432),
		User:     getEnv("POSTGRES_USER", "root"),
		Password: getEnv("POSTGRES_PASSWORD", "root"),
		Name:     getEnv("POSTGRES_DB", "friends_db"),
	}

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
