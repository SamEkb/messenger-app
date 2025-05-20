module github.com/SamEkb/messenger-app/pkg/platform/middleware

go 1.24

require (
	github.com/SamEkb/messenger-app/pkg/platform/logger v0.0.0-00010101000000-000000000000
	github.com/sony/gobreaker/v2 v2.1.0
	google.golang.org/grpc v1.72.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redsync/redsync/v4 v4.13.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/redis/go-redis/v9 v9.7.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)

replace github.com/SamEkb/messenger-app/pkg/platform/logger => ../logger
