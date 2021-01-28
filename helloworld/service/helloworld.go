package service

import (
	"context"
	"fmt"
	"monster/helloworld/api"
)

// HelloWorldService api struct
type HelloWorldService struct {
	api.UnimplementedHelloWorldServiceServer
}

// Greating api
func (service *HelloWorldService) Greating(ctx context.Context, request *api.HelloWorldRequest) (*api.HelloWorldResponse, error) {
	return &api.HelloWorldResponse{
		Code: 0,
		Data: fmt.Sprintf("Hello %s ~!", request.Name),
	}, nil
}
