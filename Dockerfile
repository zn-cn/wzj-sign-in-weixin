FROM golang:1.12

WORKDIR /app
RUN mkdir /app/log && mkdir /app/src && mkdir /app/pkg
ENV GOPATH=/app GOPROXY=https://goproxy.io GO111MODULE=on