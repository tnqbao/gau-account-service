FROM golang:1.23-alpine AS builder
WORKDIR /gau_account

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /gau_account

COPY --from=builder /gau_account/main .

EXPOSE 8080
CMD ["./main"]
