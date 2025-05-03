FROM golang:1.23-alpine AS builder
WORKDIR /gau_user 
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /gau_user
COPY --from=builder /gau_user/main .
EXPOSE 8080
CMD ["./main"]
