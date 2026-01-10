package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// addZoneSlave adds a device to an existing zone using the official /addZoneSlave endpoint
func addZoneSlave(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	masterID := c.String("master")
	slaveID := c.String("slave")
	slaveIP := c.String("slave-ip")

	if masterID == "" {
		return fmt.Errorf("master device ID is required (use --master)")
	}

	if slaveID == "" {
		return fmt.Errorf("slave device ID is required (use --slave)")
	}

	PrintDeviceHeader(fmt.Sprintf("Adding slave '%s' to zone master '%s'", slaveID, masterID), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	if slaveIP != "" {
		err = client.AddZoneSlave(masterID, slaveID, slaveIP)
	} else {
		err = client.AddZoneSlaveByDeviceID(masterID, slaveID)
	}

	if err != nil {
		PrintError(fmt.Sprintf("Failed to add zone slave: %v", err))
		return err
	}

	fmt.Printf("✅ Successfully added device '%s' to zone master '%s'\n", slaveID, masterID)

	if slaveIP != "" {
		fmt.Printf("   Slave IP: %s\n", slaveIP)
	}

	return nil
}

// removeZoneSlave removes a device from an existing zone using the official /removeZoneSlave endpoint
func removeZoneSlave(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	masterID := c.String("master")
	slaveID := c.String("slave")
	slaveIP := c.String("slave-ip")

	if masterID == "" {
		return fmt.Errorf("master device ID is required (use --master)")
	}

	if slaveID == "" {
		return fmt.Errorf("slave device ID is required (use --slave)")
	}

	PrintDeviceHeader(fmt.Sprintf("Removing slave '%s' from zone master '%s'", slaveID, masterID), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	if slaveIP != "" {
		err = client.RemoveZoneSlave(masterID, slaveID, slaveIP)
	} else {
		err = client.RemoveZoneSlaveByDeviceID(masterID, slaveID)
	}

	if err != nil {
		PrintError(fmt.Sprintf("Failed to remove zone slave: %v", err))
		return err
	}

	fmt.Printf("✅ Successfully removed device '%s' from zone master '%s'\n", slaveID, masterID)

	if slaveIP != "" {
		fmt.Printf("   Slave IP: %s\n", slaveIP)
	}

	return nil
}
