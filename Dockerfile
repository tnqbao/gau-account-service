FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build both services
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o http-service ./main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o consumer-service ./consumer/main.go

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install dependencies
RUN apk add --no-cache \
    bash \
    ca-certificates \
    curl \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz \
    | tar xvz -C /tmp \
    && mv /tmp/migrate /usr/local/bin/migrate \
    && chmod +x /usr/local/bin/migrate \
    && rm -rf /tmp/*

COPY --from=builder /app/http-service .
COPY --from=builder /app/consumer-service .
COPY deploy/migrations ./migrations
COPY shared/config ./config
COPY entrypoint.sh .

RUN chmod +x entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
