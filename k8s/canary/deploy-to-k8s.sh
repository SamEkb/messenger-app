#!/bin/bash

deploy_app() {
    echo "Deploying to Kubernetes..."
    kubectl apply -f k8s/canary/namespace.yaml
    kubectl apply -f k8s/canary/
    echo "Application deployed. Waiting for readiness..."
    kubectl wait --for=condition=available --timeout=120s deployments --all -n messenger 2>/dev/null || true
    echo "Deployment completed"
    echo "Access via: http://$(minikube ip) (Add to /etc/hosts: $(minikube ip) messenger.local)"
}


deploy_app