package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// getBass retrieves the current bass level from the device
func getBass(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Getting bass level", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	bass, err := client.GetBass()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get bass: %v", err))
		return err
	}

	fmt.Printf("Current bass level: %d\n", bass.ActualBass)

	if bass.TargetBass != bass.ActualBass {
		fmt.Printf("Target bass level: %d\n", bass.TargetBass)
	}

	return nil
}

// setBass sets the bass level on the device
func setBass(c *cli.Context) error {
	level := c.Int("level")
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader(fmt.Sprintf("Setting bass level to %d", level), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.SetBassSafe(level)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to set bass: %v", err))
		return err
	}

	PrintSuccess(fmt.Sprintf("Bass level set to %d", level))
	return nil
}

// bassUp increases the bass level
func bassUp(c *cli.Context) error {
	amount := c.Int("amount")
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader(fmt.Sprintf("Increasing bass by %d", amount), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	// Get current bass level first
	currentBass, err := client.GetBass()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get current bass: %v", err))
		return err
	}

	newLevel := currentBass.ActualBass + amount
	err = client.SetBassSafe(newLevel)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to increase bass: %v", err))
		return err
	}

	PrintSuccess(fmt.Sprintf("Bass increased from %d to %d", currentBass.ActualBass, newLevel))
	return nil
}

// bassDown decreases the bass level
func bassDown(c *cli.Context) error {
	amount := c.Int("amount")
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader(fmt.Sprintf("Decreasing bass by %d", amount), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	// Get current bass level first
	currentBass, err := client.GetBass()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get current bass: %v", err))
		return err
	}

	newLevel := currentBass.ActualBass - amount
	err = client.SetBassSafe(newLevel)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to decrease bass: %v", err))
		return err
	}

	PrintSuccess(fmt.Sprintf("Bass decreased from %d to %d", currentBass.ActualBass, newLevel))
	return nil
}

// getBassCapabilities retrieves the bass capabilities of the device
func getBassCapabilities(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Getting bass capabilities", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	capabilities, err := client.GetBassCapabilities()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get bass capabilities: %v", err))
		return err
	}

	fmt.Println("Bass Capabilities:")
	fmt.Printf("  Available: %t\n", capabilities.BassAvailable)

	if capabilities.BassAvailable {
		fmt.Printf("  Range: %d to %d\n", capabilities.BassMin, capabilities.BassMax)
		fmt.Printf("  Default: %d\n", capabilities.BassDefault)
	}

	return nil
}
