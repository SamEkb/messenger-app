include vendor.proto.mk

# Используем bin в текущей директории для установки плагинов protoc
LOCAL_BIN := $(CURDIR)/bin

# Корневая директория проекта (на один уровень выше)
PROJECT_ROOT := $(shell dirname $(CURDIR))

# Путь до protobuf файлов
PROTO_PATH := $(CURDIR)

# Путь до сгенеренных .pb.go файлов
PKG_PROTO_PATH := $(PROJECT_ROOT)/pkg/api

# Путь до swagger документации
SWAGGER_PATH := $(PROJECT_ROOT)/swagger

# Путь до завендореных protobuf файлов
VENDOR_PROTO_PATH := $(CURDIR)/vendor.protobuf

# устанавливаем необходимые плагины
.bin-deps: export GOBIN := $(LOCAL_BIN)
.bin-deps:
	$(info Installing binary dependencies...)
	mkdir -p $(LOCAL_BIN)
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/bufbuild/buf/cmd/buf@v1.32.2
	go install github.com/yoheimuta/protolint/cmd/protolint@latest
	chmod +x $(LOCAL_BIN)/* || true

install-buf:
	@echo "Installing buf..."
	go install github.com/bufbuild/buf/cmd/buf@v1.36.0

vendor:
	$(info Downloading plugins...)
	make -f vendor.proto.mk vendor


.PHONY: buf-format buf-lint buf-generate install-buf

buf-format:
	@echo "⇨ Launching buf fromat"
	$(LOCAL_BIN)/buf format -w --path auth_service --path chat_service --path events --path friends_service --path users_service

buf-lint: install-buf
	@echo "⇨ Launching buf lint"
	@chmod +x $(LOCAL_BIN)/buf || true
	@$(LOCAL_BIN)/buf lint --path auth_service --path chat_service --path events --path friends_service --path users_service

buf-generate: install-buf
	@echo "⇨ Generation buf"
	@export PATH="$(LOCAL_BIN):$$PATH" && \
	$(LOCAL_BIN)/buf generate --template buf.gen.yaml --path auth_service/v1 --path events/v1 --path users_service/v1 --path chat_service/v1 --path friends_service/v1

.PHONY: all
all: .bin-deps vendor buf-format buf-lint buf-generate