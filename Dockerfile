FROM golang:1.11-alpine3.7 as golang

ADD . /go/src/kubesphere.io/alert
WORKDIR /go/src/kubesphere.io/alert

RUN apk add --update --no-cache ca-certificates git

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN mkdir -p /alert_bin
RUN go build -v -a -installsuffix cgo -ldflags '-w' -o /alert_bin/alert cmd/alert/main.go


FROM alpine:3.7
RUN apk add --no-cache bash ca-certificates
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# modify pod (container) timezone
RUN apk add -U tzdata && ls /usr/share/zoneinfo && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && apk del tzdata

COPY --from=golang /alert_bin/alert /alerting/alert

EXPOSE 9200
EXPOSE 9201
EXPOSE 8080
