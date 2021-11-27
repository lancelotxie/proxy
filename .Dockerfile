FROM golang:latest

ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE="on"

WORKDIR /go/src

COPY . .

RUN go mod download -x

RUN go build -o proxy_client /go/src/proxy.client/main.go

ENTRYPOINT ["./proxy.client/proxy_client]
