# syntax=docker/dockerfile:1

# Stage 1: build
FROM golang:1.18-alpine as builder

RUN mkdir /build
WORKDIR /build

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o trade-utils ./cmd/trade-utils

# Stage 2: run binary
FROM alpine:3.11

COPY --from=builder /build/trade-utils /usr/local/bin/trade-utils
WORKDIR /usr/local/bin

CMD ["/usr/local/bin/trade-utils"]
