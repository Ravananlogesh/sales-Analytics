# Dockerfile

# Build stage
FROM golang:1.23 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o salesmanage main.go

FROM alpine:latest

WORKDIR /root/


COPY --from=builder /app/salesmanage .

COPY toml/config.toml /root/config.toml

EXPOSE 28090

CMD ["./salesmanage"]
