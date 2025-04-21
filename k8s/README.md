# Kubernetes Deployment Instructions

## Подготовка

1. Запустить Minikube:
   ```
   minikube start
   ```

2. Включить Ingress addon:
   ```
   minikube addons enable ingress
   ```

3. Собрать Docker-образы:
   ```
   make build
   ```

4. Загрузить образы в Minikube:
   ```
   bash k8s/load-images-to-minikube.sh
   ```

## Blue-Green Deployment

1. Выполнить деплой:
   ```
   bash k8s/green-blue/deploy-to-k8s.sh
   ```
   
   Это создаст namespace и все необходимые ресурсы. По умолчанию сервисы направлены на blue версию.

2. Добавить запись в /etc/hosts:
   ```
   echo "$(minikube ip) messenger.local" | sudo tee -a /etc/hosts
   ```

3. Для переключения на green версию:
   ```
   bash k8s/green-blue/apply-green.sh
   ```

4. Для переключения обратно на blue версию:
   ```
   bash k8s/green-blue/apply-blue.sh
   ```

## Canary Deployment

1. Выполнить деплой:
   ```
   bash k8s/canary/deploy-to-k8s.sh
   ```
   
   Это создаст namespace и развернет текущую (v1) и новую (v2) версии приложения. 
   75% трафика будет направлено на v1, 25% - на v2.

2. Для постепенного перехода на новую версию:
   ```
   bash k8s/canary/gradual-migration.sh
   ```
   
   Этот скрипт постепенно увеличивает процент трафика на новую версию.

3. Для быстрого переключения на новую версию:
   ```
   bash k8s/canary/apply-canary.sh
   ```

4. Для отката к предыдущей версии:
   ```
   bash k8s/canary/rollback-canary.sh
   ```

## Доступ к приложению

После деплоя, приложение доступно по адресу: http://messenger.local

Эндпоинты сервисов:
- Auth Service: http://messenger.local/auth
- Chat Service: http://messenger.local/chat
- Friends Service: http://messenger.local/friends
- Users Service: http://messenger.local/users