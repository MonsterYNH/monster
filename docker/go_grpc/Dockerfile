FROM --platform=arm64 golang:latest as builder
# 编译protoc
RUN mkdir /protoc
RUN wget -O /protoc/protobuf-all-3.14.0.tar.gz https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/protobuf-all-3.14.0.tar.gz
RUN tar -zxvf /protoc/protobuf-all-3.14.0.tar.gz -C /protoc
WORKDIR /protoc/protobuf-3.14.0
RUN ./configure && make && make install && ldconfig
RUN mv /usr/local/bin/protoc /usr/bin
# go环境变量
ENV GO111MODULE on
# 构建环境
RUN go get -u github.com/jteeuwen/go-bindata/...
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u github.com/protocolbuffers/protobuf
RUN go get -u google.golang.org/protobuf/cmd/protoc-gen-go
RUN go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc