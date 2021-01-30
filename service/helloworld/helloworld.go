package service

import (
	"context"
	"fmt"
	"monster/gencode/helloworld_api"
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
