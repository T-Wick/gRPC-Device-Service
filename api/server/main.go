package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	pb "projects/grpc-device-service/api/proto"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

var (
	service = &DeviceService{
		listOfDevices: &pb.Devices{
			Device: []*pb.Device{},
		},
		mapOfDevices: map[int32]*pb.Device{},
	}
	grpcPort = ":8082"
	httpPort = ":8080"
)

func runHTTP() {
	clientAddr := fmt.Sprintf("localhost%s", httpPort)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	gwmux := runtime.NewServeMux()
	if err := pb.RegisterDeviceServiceHandlerFromEndpoint(context.Background(), gwmux, clientAddr, opts); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
	log.Printf("HTTP Listening on localhost%s\n", httpPort)
	log.Fatal(http.ListenAndServe(httpPort, gwmux))
}

func runGRPC(listen net.Listener) {
	// Create new grpc server
	server := grpc.NewServer()

	// Register service
	pb.RegisterDeviceServiceServer(server, service)

	log.Printf("GRPC Listening on localhost%s\n", grpcPort)

	// Start serving requests
	server.Serve(listen)
}

func main() {
	// Start listening for grpc
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatal(err)
	}

	// Do a goroutine so that a grpc server is up
	// allowing for http server to run as well
	go runGRPC(listen)
	runHTTP()
}
