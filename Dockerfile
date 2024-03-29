FROM golang:1.20-alpine3.17
RUN apk add build-base

WORKDIR /github.com/TelefonicaTC2Tech/golium

COPY go.mod .
COPY go.sum .

COPY . .

RUN go build -v ./...
