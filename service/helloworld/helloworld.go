package service

import (
	"context"
	"fmt"
	"log"
	"monster/config"
	"monster/gencode/helloworld_api"
	"net"
	"net/http"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// HelloWorldService api struct
type HelloWorldService struct {
	helloworld_api.UnimplementedHelloWorldServiceServer
}

// Greating api
func (service *HelloWorldService) Greating(ctx context.Context, request *helloworld_api.HelloWorldRequest) (*helloworld_api.HelloWorldResponse, error) {
	return &helloworld_api.HelloWorldResponse{
		Code: 0,
		Data: fmt.Sprintf("Hello %s ~!", request.Name),
	}, nil
}

// HelloWorldHTTPApadter helloworld service adapter
type HelloWorldHTTPApadter struct {
	ser *http.Server
}

// Run listen and serve
func (service *HelloWorldHTTPApadter) Run(config *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	helloworldGRPCURI, exist := config.SubServeInfos["helloworld_grpc"]
	if !exist {
		return fmt.Errorf("con not find sub uri: %s", "helloworld_grpc")
	}
	if err := helloworld_api.RegisterHelloWorldServiceHandlerFromEndpoint(ctx, grpcMux, helloworldGRPCURI, opts); err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	service.ser = &http.Server{
		Addr:    config.Endpoint,
		Handler: mux,
	}
	log.Println("Http service start to listen: ", config.Endpoint)
	return service.ser.ListenAndServe()
}

// Close stop the service
func (service *HelloWorldHTTPApadter) Close() error {
	return service.ser.Close()
}

// GetName get service name
func (service *HelloWorldHTTPApadter) GetName() string {
	return "helloworld_http"
}

// HelloWorldGRPCAdapter grpc service adapter
type HelloWorldGRPCAdapter struct {
	ser *grpc.Server
}

// Run start grpc server and listen
func (service *HelloWorldGRPCAdapter) Run(config *config.Config) error {
	lis, err := net.Listen("tcp", config.Endpoint)
	if err != nil {
		return err
	}

	service.ser = grpc.NewServer(
		grpc.ChainStreamInterceptor(
			grpc_recovery.StreamServerInterceptor(),
		),
		grpc.ChainUnaryInterceptor(
			grpc_recovery.UnaryServerInterceptor(),
		),
	)
	log.Println("Grpc service start to listen: ", config.Endpoint)
	helloworld_api.RegisterHelloWorldServiceServer(service.ser, &HelloWorldService{})
	return service.ser.Serve(lis)
}

// Close stop grpc serve and listen
func (service *HelloWorldGRPCAdapter) Close() error {
	service.ser.Stop()
	return nil
}

// GetName get service name
func (service *HelloWorldGRPCAdapter) GetName() string {
	return "helloworld_grpc"
}
