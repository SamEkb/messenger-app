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
	DefaultGRPCPort           = 9001
	DefaultHTTPPort           = 8001
	DefaultKafkaBroker        = "localhost:9092"
	DefaultKafkaTopic         = "user-events"
	DefaultKafkaRetryInterval = 5 * time.Second
	DefaultKafkaMaxRetry      = 3
	DefaultTokenTTL           = 24 * time.Hour
)

type Config struct {
	AppName string
	Debug   string
	Server  *ServerConfig
	Kafka   *KafkaConfig
	Auth    *AuthConfig
	DB      *DBConfig
}

type ServerConfig struct {
	GRPCHost string
	GRPCPort int
	HTTPHost string
	HTTPPort int
}

type KafkaConfig struct {
	Brokers       []string
	Topic         string
	ConsumerGroup string
	MaxRetry      int
	RetryInterval time.Duration
}

type AuthConfig struct {
	TokenTTL time.Duration
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Name, c.Password)
}

func (s *ServerConfig) GrpcAddr() string {
	return s.GRPCHost + ":" + strconv.Itoa(s.GRPCPort)
}

func (s *ServerConfig) HttpAddr() string {
	return s.HTTPHost + ":" + strconv.Itoa(s.HTTPPort)
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Info: .env file not found or couldn't be loaded; using environment variables")
	}

	c := &Config{
		AppName: getEnv("APP_NAME", "AuthService"),
		Debug:   getEnv("DEBUG", "dev"),
		Server:  &ServerConfig{},
		Kafka:   &KafkaConfig{},
		Auth:    &AuthConfig{},
	}

	c.Server.GRPCHost = getEnv("GRPC_HOST", "0.0.0.0")
	c.Server.GRPCPort = getEnvAsInt("GRPC_PORT", DefaultGRPCPort)
	c.Server.HTTPHost = getEnv("HTTP_HOST", "0.0.0.0")
	c.Server.HTTPPort = getEnvAsInt("HTTP_PORT", DefaultHTTPPort)

	c.Kafka.Brokers = getEnvAsSlice("KAFKA_BROKERS", []string{DefaultKafkaBroker})
	c.Kafka.Topic = getEnv("KAFKA_PRODUCER_TOPIC", DefaultKafkaTopic)
	c.Kafka.ConsumerGroup = getEnv("KAFKA_CONSUMER_GROUP", "auth-service")
	c.Kafka.MaxRetry = getEnvAsInt("KAFKA_MAX_RETRY", DefaultKafkaMaxRetry)
	c.Kafka.RetryInterval = getEnvAsDuration("KAFKA_RETRY_INTERVAL", DefaultKafkaRetryInterval)

	c.Auth.TokenTTL = getEnvAsDuration("AUTH_TOKEN_TTL", DefaultTokenTTL)

	c.DB = &DBConfig{
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnvAsInt("POSTGRES_PORT", 5432),
		User:     getEnv("POSTGRES_USER", "root"),
		Password: getEnv("POSTGRES_PASSWORD", "root"),
		Name:     getEnv("POSTGRES_DB", "auth_db"),
	}

	return c, nil
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
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

func getEnvAsInt(key string, defaultValue int) int {
	if v := os.Getenv(key); v != "" {
		val, err := strconv.Atoi(v)
		if err == nil {
			return val
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		val, err := time.ParseDuration(v)
		if err == nil {
			return val
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if v := os.Getenv(key); v != "" {
		return strings.Split(v, ",")
	}
	return defaultValue
}
