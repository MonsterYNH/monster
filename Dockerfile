# FROM --platform=arm64 golang:latest as builder
# # 编译protoc
# RUN mkdir /protoc
# RUN wget -O /protoc/protobuf-all-3.14.0.tar.gz https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/protobuf-all-3.14.0.tar.gz
# RUN tar -zxvf /protoc/protobuf-all-3.14.0.tar.gz -C /protoc
# WORKDIR /protoc/protobuf-3.14.0
# RUN ./configure && make && make install && ldconfig
# RUN mv /usr/local/bin/protoc /usr/bin
# # go环境变量
# ENV GO111MODULE on
# # 构建环境
# RUN go get -u github.com/jteeuwen/go-bindata/...
# RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
# RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
# RUN go get -u github.com/golang/protobuf/protoc-gen-go
# RUN go get -u github.com/protocolbuffers/protobuf
# RUN go get -u google.golang.org/protobuf/cmd/protoc-gen-go
# RUN go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
FROM runnermonster/go_grpc as builder

# 添加包依赖
WORKDIR /go/cache

ADD go.mod .
ADD go.sum .

# 下载包
RUN go mod download

WORKDIR /go/release

ADD . .

# 生成文件
WORKDIR /go/release/gencode
RUN ./gencode.sh

# 编译生成可执行文件
WORKDIR /go/release
# 编译helloworld http api
RUN go build -o cmd/helloworld/payload cmd/helloworld/main.go

# helloworld-----------------------------------------------------------------------------------
FROM --platform=arm64 ubuntu:latest as helloworld
# 创建工作目录
RUN mkdir -p /app
# 复制可执行文件和swagger-ui静态资源文件
COPY --from=builder /go/release/cmd/helloworld/payload /app
# 指定工作目录
WORKDIR /app


# etcd-----------------------------------------------------------------------------------
FROM --platform=arm64 golang:latest as etcd

RUN cd /go/src && git clone https://github.com/etcd-io/etcd.git
RUN cd /go/src/etcd && ./build
RUN mv /go/src/etcd/bin/etcd /go/bin/etcd && mv /go/src/etcd/bin/etcdctl /go/bin


# gateway-----------------------------------------------------------------------------------
FROM --platform=arm64 nginx:latest as gateway

RUN mkdir /usr/share/nginx/html/swagger-ui/ && mkdir /usr/share/nginx/html/json/
COPY --from=builder /go/release/docker/gateway/dist/ /usr/share/nginx/html/swagger-ui
COPY --from=builder /go/release/gencode/swagger_json/ /usr/share/nginx/html/json
COPY ./docker/gateway/docker-entrypoint.d/prepare.sh /docker-entrypoint.d/
RUN chmod +x /docker-entrypoint.d/prepare.sh