package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"

	monitoring "go-grpc-stream/proto"

	"google.golang.org/grpc"
)

// Server is our struct that will handle the Hardware monitoring Logic
// It will fulfill the gRPC interface generated
type HardwareServer struct {
	monitoring.UnimplementedHardwareMonitorServer
}

type InventoryServer struct {
	monitoring.UnimplementedInventoryMonitorServer
}

// Monitor is used to start a stream of HardwareStats
func (hs *HardwareServer) Monitor(req *monitoring.MonitorRequest, stream monitoring.HardwareMonitor_MonitorServer) error {
	// Start a ticker that executes each 2 seconds
	timer := time.NewTicker(2 * time.Second)

	for {
		select {
		// Exit on stream context done
		case <-stream.Context().Done():
			return nil
		case <-timer.C:
			// Grab stats and output
			hwStats, err := hs.GetStats()
			if err != nil {
				log.Println(err.Error())
			} else {

			}
			// Send the Hardware stats on the stream
			err = stream.Send(hwStats)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}

func (is *InventoryServer) Status(req *monitoring.InventoryRequest, stream monitoring.InventoryMonitor_StatusServer) error {
	// Start a ticker that executes each 2 seconds
	timer := time.NewTicker(2 * time.Second)

	for {
		select {
		// Exit on stream context done
		case <-stream.Context().Done():
			return nil
		case <-timer.C:

			// random inventory item
			inventory_type_list := []string{"hardware", "software"}
			inventory_index := rand.Intn(len(inventory_type_list))
			inventory_item := inventory_type_list[inventory_index]
			switch inventory_item {
			case "hardware":
				inventory, err := is.GetHardwareInventory()
				if err != nil {
					log.Println(err.Error())
				} else {
					err = stream.Send(inventory)
					if err != nil {
						log.Println(err.Error())
					}
				}
			case "software":
				inventory, err := is.GetSoftwareInventory()
				if err != nil {
					log.Println(err.Error())
				} else {
					err = stream.Send(inventory)
					if err != nil {
						log.Println(err.Error())
					}
				}
			}

		}
	}
}

// GetStats will extract system stats and output a Hardware Object, or an error
// if extraction fails
func (hs *HardwareServer) GetStats() (*monitoring.HardwareStats, error) {
	// Extarcyt Memory statas
	mem, err := memory.Get()
	if err != nil {
		return nil, err
	}
	// Extract CPU stats
	cpu, err := cpu.Get()
	if err != nil {
		return nil, err
	}
	// Create our response object
	hwStats := &monitoring.HardwareStats{
		Cpu:        int32(cpu.Total),
		MemoryFree: int32(mem.Free),
		MemoryUsed: int32(mem.Used),
	}

	return hwStats, nil
}

func GetHardwareItem() (*monitoring.Inventory_HardwareItem, error) {

	year := int32(rand.Intn(2020-2000) + 2000)

	inventory_item := &monitoring.Inventory_HardwareItem{
		&monitoring.HardwareItem{
			Name:        "CPU",
			Description: "i7",
			Year:        year,
		},
	}

	return inventory_item, nil

}

func GetSoftwareItem() (*monitoring.Inventory_SoftwareItem, error) {

	year := int32(rand.Intn(2020-2000) + 2000)

	inventory_item := &monitoring.Inventory_SoftwareItem{
		&monitoring.SoftwareItem{
			Name:        "MSWord",
			Description: "office",
			Year:        year,
		},
	}

	return inventory_item, nil

}

func (is *InventoryServer) GetHardwareInventory() (*monitoring.Inventory, error) {

	// random action
	action := []string{"sell", "buy", "destroy"}
	action_index := rand.Intn(len(action))
	action_item := action[action_index]

	// random count
	random_count := rand.Intn(10) + 1

	hardware_item, _ := GetHardwareItem()

	inventory := &monitoring.Inventory{
		Count:  int32(random_count),
		Action: action_item,
		Item:   hardware_item,
	}
	return inventory, nil

}

func (is *InventoryServer) GetSoftwareInventory() (*monitoring.Inventory, error) {

	// random action
	action := []string{"sell", "buy", "destroy"}
	action_index := rand.Intn(len(action))
	action_item := action[action_index]

	// random count
	random_count := rand.Intn(10) + 1

	software_item, _ := GetSoftwareItem()

	inventory := &monitoring.Inventory{
		Count:  int32(random_count),
		Action: action_item,
		Item:   software_item,
	}

	return inventory, nil

}

func main() {

	fmt.Println("Welcome to streaming HW monitoring")
	// Setup a tcp connection to port 7777
	lis, err := net.Listen("tcp", ":7777")
	if err != nil {
		panic(err)
	}
	// Create a gRPC server
	gRPCserver := grpc.NewServer()

	// Create a server object of the type we created in server.go
	hs := &HardwareServer{}
	is := &InventoryServer{}

	// Regiser our server as a gRPC server
	monitoring.RegisterHardwareMonitorServer(gRPCserver, hs)
	monitoring.RegisterInventoryMonitorServer(gRPCserver, is)

	log.Println(gRPCserver.Serve(lis))

}
