FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN CGO_ENABLED=0 go build -o /app/server ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .
COPY config.env config.env

CMD ["./server"]
