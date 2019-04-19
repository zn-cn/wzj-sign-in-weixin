FROM golang:1.10

LABEL maintainer="yun_tofar@qq.com"
LABEL version="1.0"
LABEL description="wechat-service-center weixin mini program"

WORKDIR /app
RUN mkdir /app/log && mkdir /app/src
ENV GOPATH "/app"
ENTRYPOINT [ "go", "run", "src/main.go" ]