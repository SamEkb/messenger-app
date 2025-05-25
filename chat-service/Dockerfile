FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/chat-service

RUN go mod download
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o chat-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/chat-service/chat-service .

EXPOSE 8002 9002

CMD ["./chat-service"]