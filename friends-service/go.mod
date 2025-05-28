module github.com/SamEkb/messenger-app/friends-service

go 1.24

require (
	github.com/SamEkb/messenger-app/pkg/api v0.0.0-00010101000000-000000000000
	github.com/SamEkb/messenger-app/pkg/platform/errors v0.0.0-00010101000000-000000000000
	github.com/SamEkb/messenger-app/pkg/platform/logger v0.0.0-00010101000000-000000000000
	github.com/SamEkb/messenger-app/pkg/platform/middleware v0.0.0-00010101000000-000000000000
	github.com/SamEkb/messenger-app/pkg/platform/postgres v0.0.0-00010101000000-000000000000
	github.com/bufbuild/protovalidate-go v0.10.0
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	google.golang.org/grpc v1.72.1
	google.golang.org/protobuf v1.36.6
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.6-20250423154025-7712fb530c57.1 // indirect
	cel.dev/expr v0.23.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redsync/redsync/v4 v4.13.0 // indirect
	github.com/google/cel-go v0.25.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jmoiron/sqlx v1.4.0 // indirect
	github.com/redis/go-redis/v9 v9.7.0 // indirect
	github.com/sony/gobreaker/v2 v2.1.0 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	golang.org/x/exp v0.0.0-20240325151524-a685a6edb6d8 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	golang.org/x/time v0.11.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250425173222-7b384671a197 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250519155744-55703ea1f237 // indirect
)

replace github.com/SamEkb/messenger-app/pkg/api => ../pkg/api

replace github.com/SamEkb/messenger-app/pkg/platform/logger => ../pkg/platform/logger

replace github.com/SamEkb/messenger-app/pkg/platform/errors => ../pkg/platform/errors

replace github.com/SamEkb/messenger-app/pkg/platform/postgres => ../pkg/platform/postgres

replace github.com/SamEkb/messenger-app/pkg/platform/middleware => ../pkg/platform/middleware
