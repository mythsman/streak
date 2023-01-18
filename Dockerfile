FROM golang:1.17

ENV GOPROXY=https://goproxy.io

RUN sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list && apt update && apt install -y libpcap-dev && rm -rf /var/lib/apt/lists

COPY . /root/streak

WORKDIR /root/streak

RUN go build -ldflags="-s -w"