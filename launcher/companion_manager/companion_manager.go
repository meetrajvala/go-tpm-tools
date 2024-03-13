package companion_manager

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/compute/v1"
)

func CreateInstance(projectID, instanceName, zone, machineType string) error {
	ctx := context.Background()
	fmt.Printf("Creating instance %s in project %s\n", instanceName, projectID)

	// Create a new Compute Engine service client
	service, err := compute.NewService(ctx)
	if err != nil {
		return fmt.Errorf("failed to create service: %v", err)
	}

	// Create an instance resource object with the instance details
	instance := &compute.Instance{
		Name:        instanceName,
		MachineType: fmt.Sprintf("zones/%s/machineTypes/%s", zone, machineType),
		Disks: []*compute.AttachedDisk{
			{
				AutoDelete: true,
				Boot:       true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/debian-cloud/global/images/family/debian-10",
				},
			},
		},
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				AccessConfigs: []*compute.AccessConfig{
					{
						Name: "External NAT",
						Type: "ONE_TO_ONE_NAT",
					},
				},
				Network: "global/networks/default",
			},
		},
	}

	// Call the Instances.Insert method to create the instance
	op, err := service.Instances.Insert(projectID, zone, instance).Do()
	if err != nil {
		return fmt.Errorf("failed to create instance: %v", err)
	}

	for {
		// Check operation status
		op, err = service.ZoneOperations.Get(projectID, zone, op.Name).Do()
		if err != nil {
			log.Fatalf("Failed to get operation: %v", err)
		}
		if op.Status == "DONE" {
			break
		}
		fmt.Printf("Waiting more 10 secs for creating instance %s\n", instanceName)
		time.Sleep(10 * time.Second)
	}

	fmt.Printf("Instance %s created successfully!\n", instanceName)
	return nil
}

func main() {
	projectID := "compute-engine-examples"
	instanceNameTemplate := "test-instance-%d"
	zone := "us-central1-a"
	machineType := "n1-standard-1"

	// We can use go routines to create instances in parallel
	for i := 0; i < 5; i++ {
		instanceName := fmt.Sprintf(instanceNameTemplate, i)
		err := createInstance(projectID, instanceName, zone, machineType)
		if err != nil {
			log.Fatalf("Failed to create instance: %v", err)
		}
	}
}
