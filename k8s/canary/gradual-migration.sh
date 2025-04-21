#!/bin/bash

# Скрипт для постепенного перехода на новую версию (canary)

# Начальное состояние - 75% старая версия, 25% новая версия
echo "Step 1: Routing 25% traffic to new version (v2)"
kubectl scale deployment auth-v1 -n messenger --replicas=3
kubectl scale deployment auth-v2 -n messenger --replicas=1

kubectl scale deployment chat-v1 -n messenger --replicas=3
kubectl scale deployment chat-v2 -n messenger --replicas=1

kubectl scale deployment friends-v1 -n messenger --replicas=3
kubectl scale deployment friends-v2 -n messenger --replicas=1

kubectl scale deployment users-v1 -n messenger --replicas=3
kubectl scale deployment users-v2 -n messenger --replicas=1

echo "Wait 5 minutes to check stability..."
sleep 300

# 50% старая версия, 50% новая версия
echo "Step 2: Routing 50% traffic to new version (v2)"
kubectl scale deployment auth-v1 -n messenger --replicas=2
kubectl scale deployment auth-v2 -n messenger --replicas=2

kubectl scale deployment chat-v1 -n messenger --replicas=2
kubectl scale deployment chat-v2 -n messenger --replicas=2

kubectl scale deployment friends-v1 -n messenger --replicas=2
kubectl scale deployment friends-v2 -n messenger --replicas=2

kubectl scale deployment users-v1 -n messenger --replicas=2
kubectl scale deployment users-v2 -n messenger --replicas=2

echo "Wait 5 minutes to check stability..."
sleep 300

# 25% старая версия, 75% новая версия
echo "Step 3: Routing 75% traffic to new version (v2)"
kubectl scale deployment auth-v1 -n messenger --replicas=1
kubectl scale deployment auth-v2 -n messenger --replicas=3

kubectl scale deployment chat-v1 -n messenger --replicas=1
kubectl scale deployment chat-v2 -n messenger --replicas=3

kubectl scale deployment friends-v1 -n messenger --replicas=1
kubectl scale deployment friends-v2 -n messenger --replicas=3

kubectl scale deployment users-v1 -n messenger --replicas=1
kubectl scale deployment users-v2 -n messenger --replicas=3

echo "Wait 5 minutes for final check..."
sleep 300

# Если нет проблем, завершаем миграцию
echo "Final step: Complete migration to v2"
kubectl scale deployment auth-v1 -n messenger --replicas=0
kubectl scale deployment auth-v2 -n messenger --replicas=3

kubectl scale deployment chat-v1 -n messenger --replicas=0
kubectl scale deployment chat-v2 -n messenger --replicas=3

kubectl scale deployment friends-v1 -n messenger --replicas=0
kubectl scale deployment friends-v2 -n messenger --replicas=3

kubectl scale deployment users-v1 -n messenger --replicas=0
kubectl scale deployment users-v2 -n messenger --replicas=3

echo "Migration to v2 completed successfully!"