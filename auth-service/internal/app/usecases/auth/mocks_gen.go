package auth

//go:generate mockery --dir=../../ports --disable-version-string --with-expecter --name AuthRepository --output ./mocks --filename auth_repository_mock.go
//go:generate mockery --dir=../../ports --disable-version-string --with-expecter --name TokenRepository --output ./mocks --filename token_repository_mock.go
//go:generate mockery --dir=../../ports --disable-version-string --with-expecter --name UserEventsKafkaProducer --output ./mocks --filename user_events_kafka_producer_mock.go
