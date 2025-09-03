FROM golang:1.22-alpine3.20
RUN apk add build-base

WORKDIR /github.com/TelefonicaTC2Tech/golium

COPY go.mod .
COPY go.sum .

COPY . .

RUN go build -v ./...
