FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY . .
RUN go build -o main cmd/main.go

CMD ["./main"]