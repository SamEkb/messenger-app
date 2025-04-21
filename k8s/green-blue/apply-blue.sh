#!/bin/bash

kubectl patch service auth -n messenger -p '{"spec":{"selector":{"app":"auth","color":"blue"}}}'
kubectl patch service chat -n messenger -p '{"spec":{"selector":{"app":"chat","color":"blue"}}}'
kubectl patch service friends -n messenger -p '{"spec":{"selector":{"app":"friends","color":"blue"}}}'
kubectl patch service users -n messenger -p '{"spec":{"selector":{"app":"users","color":"blue"}}}'