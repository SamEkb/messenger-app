.PHONY: build up down restart ps logs help k8s-blue-green k8s-apply-green k8s-apply-blue k8s-canary k8s-canary-migrate k8s-canary-apply k8s-canary-rollback k8s-load-images k8s-deploy-all

DOCKER_COMPOSE = docker-compose

# Docker commands
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

# Kubernetes commands - Blue/Green deployment
k8s-blue-green:
	@echo "Deploying to Kubernetes with Blue-Green strategy..."
	@bash k8s/green-blue/deploy-to-k8s.sh
	@echo "Blue version deployed and active"

k8s-apply-green:
	@echo "Switching traffic to Green version..."
	@bash k8s/green-blue/apply-green.sh
	@echo "Green version is now active"

k8s-apply-blue:
	@echo "Switching traffic to Blue version..."
	@bash k8s/green-blue/apply-blue.sh
	@echo "Blue version is now active"

# Kubernetes commands - Canary deployment
k8s-canary:
	@echo "Deploying to Kubernetes with Canary release strategy..."
	@bash k8s/canary/deploy-to-k8s.sh
	@echo "Canary deployment complete: 75% traffic to v1, 25% to v2"

k8s-canary-migrate:
	@echo "Starting gradual migration to new version..."
	@bash k8s/canary/gradual-migration.sh

k8s-canary-apply:
	@echo "Switching all traffic to new version..."
	@bash k8s/canary/apply-canary.sh
	@echo "All traffic now routed to v2"

k8s-canary-rollback:
	@echo "Rolling back to previous version..."
	@bash k8s/canary/rollback-canary.sh
	@echo "All traffic now routed back to v1"

k8s-load-images:
	@echo "Loading images to Minikube..."
	@bash k8s/load-images-to-minikube.sh
	@echo "Images loaded successfully"

# Complete deployment to Kubernetes
k8s-deploy-all: build
	@echo "\n====== STARTING COMPLETE KUBERNETES DEPLOYMENT ======"
	@echo "\n[1/6] Checking if Minikube is running..."
	@minikube status > /dev/null 2>&1 || (echo "Minikube is not running. Starting Minikube..." && minikube start)
	
	@echo "\n[2/6] Enabling Ingress addon..."
	@minikube addons enable ingress
	
	@echo "\n[3/6] Loading Docker images to Minikube..."
	@bash k8s/load-images-to-minikube.sh
	
	@echo "\n[4/6] Deploying to Kubernetes with Blue-Green strategy..."
	@bash k8s/green-blue/deploy-to-k8s.sh
	
	@echo "\n[5/6] Getting Minikube IP..."
	@MINIKUBE_IP=$$(minikube ip); \
	echo "\n[6/6] Adding host entry (may require password for sudo)..."; \
	echo "$$MINIKUBE_IP messenger.local" | sudo tee -a /etc/hosts || \
	(echo "Failed to update /etc/hosts. Please manually add:"; \
	echo "$$MINIKUBE_IP messenger.local")
	
	@echo "\n====== DEPLOYMENT COMPLETE ======"
	@echo "\nApplication is now available at: http://messenger.local"
	@echo "- Auth Service:    http://messenger.local/auth"
	@echo "- Chat Service:    http://messenger.local/chat"
	@echo "- Friends Service: http://messenger.local/friends"
	@echo "- Users Service:   http://messenger.local/users"
	@echo "\nBlue version is currently active. To switch to Green version, run: make k8s-apply-green"