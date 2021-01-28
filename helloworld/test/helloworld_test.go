package service_test

import (
	"context"
	"encoding/json"
	"fmt"
	"monster/helloworld/api"
	"monster/helloworld/config"
	"testing"

	"google.golang.org/grpc"
)

func TestHelloWorldGreating(t *testing.T) {
	conn, err := grpc.Dial(config.GRPCEndpoint, grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	client := api.NewHelloWorldServiceClient(conn)
	resp, err := client.Greating(context.Background(), &api.HelloWorldRequest{
		Name: "Bob",
	})
	if err != nil {
		t.Error(err)
	}
	bytes, _ := json.Marshal(resp)
	fmt.Println(string(bytes))
}
