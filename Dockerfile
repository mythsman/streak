FROM golang:1.17-alpine AS builder

ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=on

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk add --no-cache build-base libpcap-dev

WORKDIR /root/streak

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -v -ldflags="-s -w"

FROM alpine
ENV TZ=Asia/Shanghai
ENV LD_LIBRARY_PATH=/usr/lib

WORKDIR /srv/

COPY --from=builder /root/streak/streak ./
COPY application.yml.sample ./application.yml

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk add --no-cache tzdata libpcap-dev

CMD "./streak"