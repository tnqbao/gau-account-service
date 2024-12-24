FROM golang:1.23-alpine AS builder
WORKDIR /gau_user 
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY .env .  
RUN go build -o main .

FROM alpine:latest
WORKDIR /gau_user
COPY --from=builder /gau_user/main .
COPY .env .  
EXPOSE 8080
CMD ["./main"]
