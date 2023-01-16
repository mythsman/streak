FROM golang:1.17

ENV GOPROXY=https://goproxy.io
ENV LIBPCAP_VERSION=1.10.3

RUN sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list && apt update && apt install flex bison -y && rm -rf /var/lib/apt/lists

RUN cd /root && wget https://www.tcpdump.org/release/libpcap-${LIBPCAP_VERSION}.tar.gz && tar -zxvf  libpcap-${LIBPCAP_VERSION}.tar.gz && cd /root/libpcap-${LIBPCAP_VERSION} && ./configure && make

COPY . /root/streak

#WORKDIR /root/streak

#RUN go build -ldflags="-s -w"
