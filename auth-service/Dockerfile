FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/auth-service

RUN go mod download
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/auth-service/auth-service .

EXPOSE 8001 9001

CMD ["./auth-service"]