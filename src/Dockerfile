FROM golang:1.22.11-alpine AS builder

RUN apk add --no-cache git build-base

WORKDIR /build

COPY . .

RUN go build -buildmode=plugin -o certificate_validation.so certificate_validation.go
