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
		allUsers: map[string]*pb.User{},
		listOfDevices: &pb.Devices{
			Device: []*pb.Device{},
		},
		mapOfDevices: map[int32]*pb.Device{},
	}

	grpcPort = ":8082"
	httpPort = ":8080"
	grpcAddr = fmt.Sprintf("localhost%s", grpcPort)
	httpAddr = fmt.Sprintf("localhost%s", httpPort)
)

func runHTTP() {

	opts := []grpc.DialOption{grpc.WithInsecure()}
	gwmux := runtime.NewServeMux()
	if err := pb.RegisterDeviceServiceHandlerFromEndpoint(context.Background(), gwmux, grpcAddr, opts); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
	log.Printf("HTTP Listening on %s\n", httpAddr)
	log.Fatal(http.ListenAndServe(httpAddr, gwmux))
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

// DeviceService is a struct that holds a list of
// devices, specifically a Devices type.
type DeviceService struct {
	allUsers      map[string]*pb.User
	listOfDevices *pb.Devices
	mapOfDevices  map[int32]*pb.Device
}

// CreateUser cadds a user to the server
func (s *DeviceService) CreateUser(ctx context.Context, userReq *pb.CreateUserRequest) (*pb.CreateUserRequest, error) {
	return nil, nil
}

// GetAllDevices returns all the devices that are
// registered in the server.
func (s *DeviceService) GetAllDevices(ctx context.Context, req *pb.Empty) (*pb.Devices, error) {
	return s.listOfDevices, nil
}

// GetDeviceByID returns the device specified by the
// ID passed in.
func (s *DeviceService) GetDeviceByID(ctx context.Context, id *pb.ID) (*pb.Device, error) {
	return s.mapOfDevices[id.Id], nil
}

// SwitchDevice updates the status of an already existing
// device.
func (s *DeviceService) SwitchDevice(ctx context.Context, updatedDevice *pb.UpdateDevice) (*pb.Device, error) {
	toUpdate, found := s.mapOfDevices[updatedDevice.Id]
	if !found {
		return s.mapOfDevices[updatedDevice.Id], fmt.Errorf("Device (ID: %v) was NOT Found as a Registered Device", updatedDevice.Id)
	}
	toUpdate.State = updatedDevice.Value

	return s.mapOfDevices[updatedDevice.Id], nil
}

// RegisterDevice adds the desired device to the server. Each
// device must have a unique ID.
func (s *DeviceService) RegisterDevice(ctx context.Context, device *pb.Device) (*pb.Device, error) {
	if device == nil {
		return nil, fmt.Errorf("Can NOT Send a nil Request")
	}

	_, found := s.mapOfDevices[device.Id]
	if found {
		return device, fmt.Errorf("Device (ID: %v) is Already Registered", device.Id)
	}

	s.mapOfDevices[device.Id] = device
	s.listOfDevices.Device = append(s.listOfDevices.Device, device)
	return device, nil
}
