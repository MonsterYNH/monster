package main

import (
	"context"
	"fmt"
	"monster/helloworld/api"
	"monster/helloworld/api/swagger"
	"monster/helloworld/config"
	"monster/helloworld/service"
	"net"
	"net/http"
	"path"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {
	go func() {
		if err := startHTTP(config.HTTPEndpoint, config.GRPCEndpoint); err != nil {
			panic(fmt.Sprintf("ERROR: failed to start http serve, error: %s", err.Error()))
		}
	}()
	// go func() {
	// 	if err := startSwagger(config.SwaggerEndpoint); err != nil {
	// 		panic(fmt.Sprintf("ERROR: failed to start swagger serve, error: %s", err.Error()))
	// 	}
	// }()
	if err := startGRPC(config.GRPCEndpoint); err != nil {
		panic(fmt.Sprintf("ERROR: failed, to start grpc server, error: %s", err.Error()))
	}
}

func startHTTP(httpURI, grpcURI string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := api.RegisterHelloWorldServiceHandlerFromEndpoint(ctx, grpcMux, grpcURI, opts); err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	mux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "swagger.json") {
			fmt.Printf("Not Found: %s\r\n", r.URL.Path)
			http.NotFound(w, r)
			return
		}

		p := strings.TrimPrefix(r.URL.Path, "/swagger/")
		p = path.Join("api/protos", p)
		fmt.Println(p)

		fmt.Printf("Serving swagger-file: %s\r\n", p)

		http.ServeFile(w, r, p)
	})
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "third_party/swagger-ui/",
	})
	prefix := "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
	return http.ListenAndServe(httpURI, mux)
}

func startGRPC(uri string) error {
	lis, err := net.Listen("tcp", uri)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.ChainStreamInterceptor(
			grpc_recovery.StreamServerInterceptor(),
			// middleware.StreamRecoveryInterceptor,
		),
		grpc.ChainUnaryInterceptor(
			grpc_recovery.UnaryServerInterceptor(),
			// middleware.UnaryRecoveryInterceptor,
		),
	)
	api.RegisterHelloWorldServiceServer(grpcServer, &service.HelloWorldService{})
	return grpcServer.Serve(lis)
}
