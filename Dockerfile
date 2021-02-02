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
WORKDIR /go/release/grpc_protos/gencode
RUN ./gencode.sh

# 编译生成可执行文件
WORKDIR /go/release
# 编译helloworld http api
RUN go build -o cmd/helloworld/payload cmd/helloworld/main.go

# helloworld-----------------------------------------------------------------------------------
FROM --platform=arm64 ubuntu:latest as helloworld
# 添加普通用户
RUN groupadd -g 666666 bsafe && \
useradd -u 666666 -g bsafe -d /app bsafe -m
# 复制可执行文件
COPY --from=builder /go/release/cmd/helloworld/payload /app
# 指定工作目录
WORKDIR /app


# etcd-----------------------------------------------------------------------------------
FROM --platform=arm64 golang:latest as etcd
# 添加普通用户
RUN groupadd -g 666666 bsafe && \
useradd -u 666666 -g bsafe -d /app bsafe -m

RUN cd /go/src && git clone https://github.com/etcd-io/etcd.git
RUN cd /go/src/etcd && ./build
RUN mv /go/src/etcd/bin/etcd /go/bin/etcd && mv /go/src/etcd/bin/etcdctl /go/bin


# gateway-----------------------------------------------------------------------------------
FROM --platform=arm64 nginx:latest as gateway
RUN mkdir /usr/share/nginx/html/swagger-ui /usr/share/nginx/html/json
COPY --from=builder /go/release/docker/gateway/dist/ /usr/share/nginx/html/swagger-ui
COPY --from=builder /go/release/grpc_protos/gencode/swagger_json/ /usr/share/nginx/html/json
COPY ./docker/gateway/docker-entrypoint.d/prepare.sh /docker-entrypoint.d/
RUN chmod +x /docker-entrypoint.d/prepare.sh
