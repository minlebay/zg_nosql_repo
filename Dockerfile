FROM golang:1.22-alpine AS builder
LABEL maintainer="ilshatminnibaev@gmail.com"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/main.go
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/config.yaml .
EXPOSE 8888
EXPOSE 29092
EXPOSE 9092
CMD ["./main"]
