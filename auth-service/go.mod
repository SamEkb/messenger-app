module github.com/SamEkb/messenger-app/auth-service

go 1.24

require (
	buf.build/go/protovalidate v0.12.0
	github.com/SamEkb/messenger-app/pkg/api v0.0.0-00010101000000-000000000000
	github.com/SamEkb/messenger-app/pkg/platform/errors v0.0.0-00010101000000-000000000000
	github.com/SamEkb/messenger-app/pkg/platform/logger v0.0.0-00010101000000-000000000000
	github.com/SamEkb/messenger-app/pkg/platform/middleware v0.0.0-00010101000000-000000000000
	github.com/SamEkb/messenger-app/pkg/platform/postgres v0.0.0-00010101000000-000000000000
	github.com/Shopify/sarama v1.38.1
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	github.com/stretchr/testify v1.10.0
	golang.org/x/crypto v0.37.0
	google.golang.org/grpc v1.72.1
	google.golang.org/protobuf v1.36.6
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.6-20250425153114-8976f5be98c1.1 // indirect
	cel.dev/expr v0.23.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/eapache/go-resiliency v1.3.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230111030713-bf00bc1b83b6 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/go-redsync/redsync/v4 v4.13.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/cel-go v0.25.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.3 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jmoiron/sqlx v1.4.0 // indirect
	github.com/klauspost/compress v1.15.14 // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/redis/go-redis/v9 v9.7.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/sony/gobreaker/v2 v2.1.0 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/exp v0.0.0-20240325151524-a685a6edb6d8 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	golang.org/x/time v0.11.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250425173222-7b384671a197 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250425173222-7b384671a197 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/SamEkb/messenger-app/pkg/api => ../pkg/api

replace github.com/SamEkb/messenger-app/pkg/platform/logger => ../pkg/platform/logger

replace github.com/SamEkb/messenger-app/pkg/platform/errors => ../pkg/platform/errors

replace github.com/SamEkb/messenger-app/pkg/platform/postgres => ../pkg/platform/postgres

replace github.com/SamEkb/messenger-app/pkg/platform/middleware => ../pkg/platform/middleware
