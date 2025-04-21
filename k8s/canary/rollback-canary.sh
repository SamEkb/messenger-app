#!/bin/bash

kubectl scale deployment auth-v2 --replicas=0
kubectl scale deployment auth-v1 --replicas=3

kubectl scale deployment chat-v2 --replicas=0
kubectl scale deployment chat-v1 --replicas=3

kubectl scale deployment friends-v2 --replicas=0
kubectl scale deployment friends-v1 --replicas=3

kubectl scale deployment users-v2 --replicas=0
kubectl scale deployment users-v1 --replicas=3