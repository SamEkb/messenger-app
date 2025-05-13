package env

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	DefaultGRPCPort           = 9004
	DefaultHTTPPort           = 8004
	DefaultKafkaBroker        = "localhost:9092"
	DefaultKafkaTopic         = "user-events"
	DefaultKafkaRetryInterval = 5 * time.Second
	DefaultKafkaMaxRetry      = 3
)

type Config struct {
	AppName string
	Debug   string
	Server  *ServerConfig
	Kafka   *KafkaConfig
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

func (s *ServerConfig) GrpcAddr() string {
	return s.GRPCHost + ":" + strconv.Itoa(s.GRPCPort)
}

func (s *ServerConfig) HttpAddr() string {
	return s.HTTPHost + ":" + strconv.Itoa(s.HTTPPort)
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Info: .env file not found or couldn't be loaded; using environment variables")
	}

	c := &Config{
		AppName: getEnv("APP_NAME", "AuthService"),
		Debug:   getEnv("DEBUG", "dev"),
		Server:  &ServerConfig{},
		Kafka:   &KafkaConfig{},
	}

	c.Server.GRPCHost = getEnv("GRPC_HOST", "0.0.0.0")
	c.Server.GRPCPort = getEnvAsInt("GRPC_PORT", DefaultGRPCPort)
	c.Server.HTTPHost = getEnv("HTTP_HOST", "0.0.0.0")
	c.Server.HTTPPort = getEnvAsInt("HTTP_PORT", DefaultHTTPPort)

	c.Kafka.Brokers = getEnvAsSlice("KAFKA_BROKERS", []string{DefaultKafkaBroker})
	c.Kafka.Topic = getEnv("KAFKA_PRODUCER_TOPIC", DefaultKafkaTopic)
	c.Kafka.ConsumerGroup = getEnv("KAFKA_CONSUMER_GROUP", "users-service-group")
	c.Kafka.MaxRetry = getEnvAsInt("KAFKA_MAX_RETRY", DefaultKafkaMaxRetry)
	c.Kafka.RetryInterval = getEnvAsDuration("KAFKA_RETRY_INTERVAL", DefaultKafkaRetryInterval)

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
