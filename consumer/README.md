# Consumer Module

## Mô tả | Description

**Tiếng Việt:** Module Consumer xử lý các tác vụ bất đồng bộ qua message queues.

**English:** Consumer module handling asynchronous tasks via message queues.

## Cấu trúc | Structure

```
consumer/
└── main.go        # Consumer service entry point
```

## Features

### Current
- Service placeholder
- Configuration loading
- Service runner

### Planned
- Email queue consumer
- SMS queue consumer  
- Audit log consumer
- Analytics consumer

## Message Queues

### Supported Brokers
- RabbitMQ (planned)
- Apache Kafka (planned)
- Redis Pub/Sub (planned)

### Queue Types
```
account.email.send      # Email notifications
account.sms.send        # SMS notifications
account.audit.log       # Audit logging
account.analytics.track # Analytics data
```

## Deployment

### Docker
```bash
# Build
docker build -t gau-account-service .

# Run consumer
docker run gau-account-service consumer
```

### Kubernetes
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gau-account-consumer
spec:
  containers:
  - name: consumer
    image: gau-account-service:latest
    command: ["./entrypoint.sh", "consumer"]
```

## Usage

```go
cfg := config.NewConfig()
// TODO: Initialize message queue consumers
log.Printf("Consumer service starting")
select {} // Keep running
```
