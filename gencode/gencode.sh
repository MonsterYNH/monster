# gengerate helloworld api file
protoc -I ../protocs --go-grpc_out=./helloworld_api --grpc-gateway_out=logtostderr=true:./helloworld_api --go_out=./helloworld_api --swagger_out=logtostderr=true:./swagger_json ../protocs/helloworld.proto