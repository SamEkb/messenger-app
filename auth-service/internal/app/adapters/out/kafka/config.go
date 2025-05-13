package kafka

import (
	"log/slog"

	"github.com/SamEkb/messenger-app/auth-service/config/env"
	"github.com/Shopify/sarama"
)

const serviceName = "auth-service"

func NewSaramaConfig(kafkaConfig *env.KafkaConfig, logger *slog.Logger) *sarama.Config {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	config.Producer.Retry.Max = kafkaConfig.MaxRetry
	config.Producer.Retry.Backoff = kafkaConfig.RetryInterval

	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true

	config.ClientID = serviceName
	config.Version = sarama.V2_8_0_0

	if logger != nil {
		sarama.Logger = SaramaLoggerAdapter{logger: logger}
	}

	return config
}

type SaramaLoggerAdapter struct {
	logger *slog.Logger
}

func (s SaramaLoggerAdapter) Print(v ...interface{}) {
	s.logger.Debug("sarama internal", "message", v)
}

func (s SaramaLoggerAdapter) Printf(format string, v ...interface{}) {
	s.logger.Debug("sarama internal", "format", format, "args", v)
}

func (s SaramaLoggerAdapter) Println(v ...interface{}) {
	s.logger.Debug("sarama internal", "message", v)
}
