#!/bin/bash

kubectl scale deployment auth-v2 -n messenger --replicas=3
kubectl scale deployment auth-v1 -n messenger --replicas=0

kubectl scale deployment chat-v2 -n messenger --replicas=3
kubectl scale deployment chat-v1 -n messenger --replicas=0

kubectl scale deployment friends-v2 -n messenger --replicas=3
kubectl scale deployment friends-v1 -n messenger --replicas=0

kubectl scale deployment users-v2 -n messenger --replicas=3
kubectl scale deployment users-v1 -n messenger --replicas=0

echo "Canary deployments fully promoted to production"