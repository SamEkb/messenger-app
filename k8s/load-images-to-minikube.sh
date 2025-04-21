#!/bin/bash

echo "Loading Docker images to Minikube..."
minikube image load auth-service:latest
minikube image load chat-service:latest
minikube image load friends-service:latest
minikube image load users-service:latest
echo "All images loaded successfully"