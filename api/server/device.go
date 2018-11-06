package main

import (
	"context"
	"fmt"
	pb "projects/grpc-device-service/api/proto"
)

// DeviceService is a struct that holds a list of
// devices, specifically a Devices type.
type DeviceService struct {
	listOfDevices *pb.Devices
	mapOfDevices  map[int32]*pb.Device
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
