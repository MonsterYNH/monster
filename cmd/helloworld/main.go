package main

import (
	"context"
	"log"
	"monster/config"
	"monster/gencode/helloworld_api"
	service "monster/service/helloworld"
	"monster/utils"
	"net"
	"net/http"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func main() {
	if err := utils.MultiRun(serveHTTP(config.GRPCEndpoint, config.HTTPEndpoint), serveGRPC(config.GRPCEndpoint)); err != nil {
		log.Println(err)
	}
}

func serveHTTP(grpcEndpoint, httpEndpoint string) utils.MultiServe {
	return func() error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		grpcMux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithInsecure()}
		if err := helloworld_api.RegisterHelloWorldServiceHandlerFromEndpoint(ctx, grpcMux, grpcEndpoint, opts); err != nil {
			return err
		}

		mux := http.NewServeMux()
		mux.Handle("/", grpcMux)
		log.Println("Http service start to listen: ", httpEndpoint)
		return http.ListenAndServe(httpEndpoint, mux)
	}
}

func serveGRPC(endpoint string) utils.MultiServe {
	return func() error {
		lis, err := net.Listen("tcp", endpoint)
		if err != nil {
			return err
		}

		grpcServer := grpc.NewServer(
			grpc.ChainStreamInterceptor(
				grpc_recovery.StreamServerInterceptor(),
			),
			grpc.ChainUnaryInterceptor(
				grpc_recovery.UnaryServerInterceptor(),
			),
		)
		log.Println("Grpc service start to listen: ", endpoint)
		helloworld_api.RegisterHelloWorldServiceServer(grpcServer, &service.HelloWorldService{})
		return grpcServer.Serve(lis)
	}
}
