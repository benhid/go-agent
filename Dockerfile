FROM golang:1.16.7 as builder

RUN mkdir /build
ADD ../.. /build/

WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .
