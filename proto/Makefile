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

# Установить buf глобально или локально
install-buf:
	@echo "Installing buf..."
	go install github.com/bufbuild/buf/cmd/buf@v1.36.0

vendor:
	$(info Downloading plugins...)
	make -f vendor.proto.mk vendor


.PHONY: buf-lint buf-breaking buf-generate install-buf


buf-lint: install-buf
	@echo "⇨ Launching buf lint"
	@chmod +x $(LOCAL_BIN)/buf || true
	@$(LOCAL_BIN)/buf lint --path auth-service --path chat-service --path events --path friends-service --path users-service

buf-breaking: install-buf
	@echo "⇨ Checking breaking changes"
	@$(LOCAL_BIN)/buf breaking --against '.git#branch=main,subdir=proto'

buf-generate: install-buf
	@echo "⇨ Generation buf"
	@export PATH="$(LOCAL_BIN):$$PATH" && \
	$(LOCAL_BIN)/buf generate --template buf.gen.yaml --path auth_service/v1 --path events/v1

.PHONY: all
all: .bin-deps vendor buf-lint buf-generate