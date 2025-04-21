#!/bin/bash

kubectl patch service auth -n messenger -p '{"spec":{"selector":{"app":"auth","color":"green"}}}'
kubectl patch service chat -n messenger -p '{"spec":{"selector":{"app":"chat","color":"green"}}}'
kubectl patch service friends -n messenger -p '{"spec":{"selector":{"app":"friends","color":"green"}}}'
kubectl patch service users -n messenger -p '{"spec":{"selector":{"app":"users","color":"green"}}}'