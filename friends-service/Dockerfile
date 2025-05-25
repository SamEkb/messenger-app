FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/friends-service

RUN go mod download
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o friends-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/friends-service/friends-service .

EXPOSE 8003 9003

CMD ["./friends-service"]