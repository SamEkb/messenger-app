# Kubernetes Deployment Instructions

## Полный автоматический деплой

Для автоматического выполнения всех шагов деплоя с помощью одной команды:

```
make k8s-deploy-all
```

Эта команда:
1. Проверяет и запускает Minikube (если не запущен)
2. Включает Ingress addon
3. Собирает и загружает Docker-образы
4. Выполняет Blue-Green деплой 
5. Добавляет запись в /etc/hosts
6. Выводит URL'ы для доступа к сервисам

## Подготовка (ручной процесс)

Шаги деплоя вручную:

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
   make k8s-load-images
   ```
   или
   ```
   bash k8s/load-images-to-minikube.sh
   ```

## Blue-Green Deployment

1. Выполнить деплой:
   ```
   make k8s-blue-green
   ```
   или
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
   make k8s-apply-green
   ```
   или
   ```
   bash k8s/green-blue/apply-green.sh
   ```

4. Для переключения обратно на blue версию:
   ```
   make k8s-apply-blue
   ```
   или
   ```
   bash k8s/green-blue/apply-blue.sh
   ```

## Canary Deployment

1. Выполнить деплой:
   ```
   make k8s-canary
   ```
   или
   ```
   bash k8s/canary/deploy-to-k8s.sh
   ```
   
   Это создаст namespace и развернет текущую (v1) и новую (v2) версии приложения. 
   75% трафика будет направлено на v1, 25% - на v2.

2. Для постепенного перехода на новую версию:
   ```
   make k8s-canary-migrate
   ```
   или
   ```
   bash k8s/canary/gradual-migration.sh
   ```
   
   Этот скрипт постепенно увеличивает процент трафика на новую версию.

3. Для быстрого переключения на новую версию:
   ```
   make k8s-canary-apply
   ```
   или
   ```
   bash k8s/canary/apply-canary.sh
   ```

4. Для отката к предыдущей версии:
   ```
   make k8s-canary-rollback
   ```
   или
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