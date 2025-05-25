FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/users-service

RUN go mod download
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o users-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/users-service/users-service .

EXPOSE 8004 9004

CMD ["./users-service"]