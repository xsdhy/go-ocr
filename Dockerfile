# 编译
FROM golang:1.20.9-alpine  as builder
#ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=off
ENV GOPATH="/go/release:/go/release/src/gopathlib/"
#安装编译需要的环境gcc等
RUN apk add build-base

WORKDIR /go/release
#将上层整个文件夹拷贝到/go/release
ADD . /go/release/src
WORKDIR /go/release/src
#交叉编译，需要制定CGO_ENABLED=1，默认是关闭的
RUN  GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o ./bin/ocr main.go

#正式环境
FROM alpine

COPY --from=builder  /go/release/src/bin/ocr /app/ocr
COPY --from=builder  /go/release/src/models /app/models
COPY --from=builder  /go/release/src/lib /app/lib

WORKDIR /app

# ENV LD_LIBRARY_PATH= /app/lib/lib

CMD ["/app/ocr"]
EXPOSE 8080