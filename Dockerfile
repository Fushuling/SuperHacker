FROM golang:latest

ADD . ./

WORKDIR $GOPATH/Test

ENV go env -w GO111MODULE=on

ENV GOPROXY=https://goproxy.cn,direct

RUN go build ./

expose 9999

ENTRYPOINT ["./main"]

