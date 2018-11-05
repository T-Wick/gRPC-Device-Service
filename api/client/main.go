package main

import (
	"context"
	"fmt"
	"projects/gRPC-rest-api/api/proto"

	"google.golang.org/grpc"
)

var (
	deviceOne = &pb.Device{
		Id:       1234,
		Hardware: "PC",
		Name:     "Trece-PC",
		Location: "Washington",
		Type:     pb.Device_onOff,
		Unit:     "Uncategorized",
		State:    1,
	}

	deviceTwo = &pb.Device{
		Id:       8801,
		Hardware: "Laptop",
		Name:     "Elizabeth-Laptop",
		Location: "Washington",
		Type:     pb.Device_onOff,
		Unit:     "Uncategorized",
		State:    1,
	}

	deviceThree = &pb.Device{
		Id:       9085,
		Hardware: "XBOX",
		Name:     "Next-Gen",
		Location: "Washington",
		Type:     pb.Device_onOff,
		Unit:     "Uncategorized",
		State:    1,
	}

	deviceFour = &pb.Device{
		Id:       3341,
		Hardware: "Lights",
		Name:     "Front-Lights",
		Location: "Washington",
		Type:     pb.Device_onOff,
		Unit:     "Uncategorized",
		State:    1,
	}
)

func main() {
	serverAddr := "localhost:8082"
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Error Connecting: ", err)
	}
	defer conn.Close()

	client := pb.NewDeviceServiceClient(conn)

	client.RegisterDevice(context.Background(), deviceOne)
	client.RegisterDevice(context.Background(), deviceThree)
	client.RegisterDevice(context.Background(), deviceTwo)
	client.RegisterDevice(context.Background(), deviceFour)

	printDevices(client)

	/*device, err := client.RegisterDevice(context.Background(), deviceTwo)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(device)*/

	foundDevice, err := client.GetDeviceByID(context.Background(), &pb.ID{Id: 1234})
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("\n---------------------------------------------Printing Device By ID----------------------------------------------\n")
	fmt.Printf("Device With ID 1234: %v\n", foundDevice)
	fmt.Printf("----------------------------------------------------------------------------------------------------------------\n\n")

	_, err = client.SwitchDevice(context.Background(), &pb.UpdateDevice{Id: 1234, Value: 13})
	if err != nil {
		fmt.Println(err.Error())
	}
	//fmt.Printf("After Update: %v", updatedDevice)

	printDevices(client)
}

func printDevices(client pb.DeviceServiceClient) error {
	devices, err := client.GetAllDevices(context.Background(), &pb.Empty{})
	if err != nil {
		return fmt.Errorf("An Error Occured Getting Devices")
	}

	fmt.Println("---------------------------------------------Printing Devices---------------------------------------------")
	for _, device := range devices.Device {
		fmt.Printf("Name:     %s\nID:       %v\nHardware: %s\nLocation: %s\nType:     %v\nUnit:     %s\nState:    %v\n\n",
			device.Name, device.Id, device.Hardware, device.Location, device.Type, device.Unit, device.State)
	}
	fmt.Println("-----------------------------------------------------------------------------------------------------------")

	return nil
}
