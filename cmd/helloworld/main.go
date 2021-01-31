package main

import (
	"monster/adapter"
	service "monster/service/helloworld"
)

func main() {
	ser := adapter.NewService()
	if err := ser.Serve(&service.HelloWorldGRPCAdapter{}, &service.HelloWorldHTTPApadter{}); err != nil {
		panic(err)
	}
	// if err := ser.Serve(); err != nil {
	// 	panic(err)
	// }
}
