package middleware

import (
	"context"
	"log"
	"runtime/debug"

	"google.golang.org/grpc"
)

// StreamRecoveryInterceptor stream recover
func StreamRecoveryInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("PANIC: ", err)
			log.Println("recover from ", string(debug.Stack()))
		}
	}()

	return handler(srv, stream)
}

// UnaryRecoveryInterceptor unary recover
func UnaryRecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("PANIC: ", err)
			log.Println("recover from ", string(debug.Stack()))
		}
	}()

	return handler(ctx, req)
}
