include .env

LOCAL_BIN:=$(CURDIR)/bin

GOOSE=goose
MIGRATIONS_DIR=./migrations
DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

# Сборка приложения
build:
	CGO_ENABLED=0 go build -v -o $(LOCAL_BIN) ./cmd

migrate-up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

migrate-status:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" status

.PHONY: build migrate-up migrate-down migrate-status 