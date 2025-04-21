.PHONY: build up down restart ps logs help

DOCKER_COMPOSE = docker-compose

build:
	@echo "Building all services..."
	@$(DOCKER_COMPOSE) build

up:
	@echo "Start all services..."
	@$(DOCKER_COMPOSE) up -d
	@echo "All services successfully started!"
	@echo "Auth Service:   http://localhost:8001"
	@echo "Chat Service:   http://localhost:8002"
	@echo "Friends Service: http://localhost:8003"
	@echo "Users Service:  http://localhost:8004"

down:
	@echo "Stopping all services..."
	@$(DOCKER_COMPOSE) down

restart: down up

ps:
	@$(DOCKER_COMPOSE) ps

logs:
	@$(DOCKER_COMPOSE) logs -f