#!/bin/bash

kubectl patch service auth-service -p '{"spec":{"selector":{"app":"auth","color":"green"}}}'
kubectl patch service chat-service -p '{"spec":{"selector":{"app":"chat","color":"green"}}}'
kubectl patch service friends-service -p '{"spec":{"selector":{"app":"friends","color":"green"}}}'
kubectl patch service users-service -p '{"spec":{"selector":{"app":"users","color":"green"}}}'