FROM golang:1.17-alpine AS builder

ENV GOPROXY=https://goproxy.io

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk add --no-app build-base libpcap-dev

COPY . /root/streak
WORKDIR /root/streak

RUN go build -ldflags="-s -w"

FROM alpine
ENV TZ=Asia/Shanghai
ENV LD_LIBRARY_PATH=/usr/lib

COPY --from=builder /root/streak/streak /root/
WORKDIR /root/

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk add --no-app tzdata libpcap-dev

CMD "/root/streak"