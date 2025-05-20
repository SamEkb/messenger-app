module github.com/SamEkb/messenger-app/chat-service

go 1.24

require (
	github.com/SamEkb/messenger-app/pkg/api v0.0.0-00010101000000-000000000000
	github.com/SamEkb/messenger-app/pkg/platform/errors v0.0.0-00010101000000-000000000000
	github.com/SamEkb/messenger-app/pkg/platform/logger v0.0.0-00010101000000-000000000000
	github.com/SamEkb/messenger-app/pkg/platform/middleware v0.0.0-00010101000000-000000000000
	github.com/SamEkb/messenger-app/pkg/platform/mongodb v0.0.0-00010101000000-000000000000
	github.com/bufbuild/protovalidate-go v0.10.0
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3
	github.com/joho/godotenv v1.5.1
	go.mongodb.org/mongo-driver v1.17.3
	google.golang.org/grpc v1.72.1
	google.golang.org/protobuf v1.36.6
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.6-20250423154025-7712fb530c57.1 // indirect
	cel.dev/expr v0.23.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/cel-go v0.25.0 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/exp v0.0.0-20240325151524-a685a6edb6d8 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250425173222-7b384671a197 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250425173222-7b384671a197 // indirect
)

replace github.com/SamEkb/messenger-app/pkg/api => ../pkg/api

replace github.com/SamEkb/messenger-app/pkg/platform/logger => ../pkg/platform/logger

replace github.com/SamEkb/messenger-app/pkg/platform/errors => ../pkg/platform/errors

replace github.com/SamEkb/messenger-app/pkg/platform/mongodb => ../pkg/platform/mongodb

replace github.com/SamEkb/messenger-app/pkg/platform/middleware => ../pkg/platform/middleware
