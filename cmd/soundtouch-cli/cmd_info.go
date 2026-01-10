package main

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

// getDeviceInfo handles the device info command
func getDeviceInfo(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Getting device information", clientConfig.Host, clientConfig.Port)

	deviceInfo, err := client.GetDeviceInfo()
	if err != nil {
		return fmt.Errorf("failed to get device info: %w", err)
	}

	// Display basic device information
	fmt.Printf("Device Information:\n")
	fmt.Printf("  Name: %s\n", deviceInfo.Name)
	fmt.Printf("  Type: %s\n", deviceInfo.Type)
	fmt.Printf("  Device ID: %s\n", deviceInfo.DeviceID)

	if deviceInfo.MargeAccountUUID != "" {
		fmt.Printf("  Account UUID: %s\n", deviceInfo.MargeAccountUUID)
	}

	if len(deviceInfo.NetworkInfo) > 0 {
		fmt.Printf("  Network Info:\n")

		for _, net := range deviceInfo.NetworkInfo {
			fmt.Printf("    - Type: %s\n", net.Type)
			fmt.Printf("      MAC Address: %s\n", net.MacAddress)
			fmt.Printf("      IP Address: %s\n", net.IPAddress)
		}
	}

	if len(deviceInfo.Components) > 0 {
		fmt.Printf("  Components:\n")

		for _, component := range deviceInfo.Components {
			fmt.Printf("    - Category: %s\n", component.ComponentCategory)

			if component.SoftwareVersion != "" {
				fmt.Printf("      Software Version: %s\n", component.SoftwareVersion)
			}

			if component.SerialNumber != "" {
				fmt.Printf("      Serial Number: %s\n", component.SerialNumber)
			}
		}
	}

	return nil
}

// getDeviceName handles getting the device name
func getDeviceName(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Getting device name", clientConfig.Host, clientConfig.Port)

	name, err := client.GetName()
	if err != nil {
		return fmt.Errorf("failed to get device name: %w", err)
	}

	fmt.Printf("Device Name: %s\n", name)

	return nil
}

// setDeviceName handles setting the device name
func setDeviceName(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	newName := c.String("value")
	if newName == "" {
		return fmt.Errorf("device name cannot be empty")
	}

	PrintDeviceHeader(fmt.Sprintf("Setting device name to '%s'", newName), clientConfig.Host, clientConfig.Port)

	err = client.SetName(newName)
	if err != nil {
		return fmt.Errorf("failed to set device name: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Device name set to '%s'", newName))

	return nil
}

// getCapabilities handles getting device capabilities
func getCapabilities(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Getting device capabilities", clientConfig.Host, clientConfig.Port)

	capabilities, err := client.GetCapabilities()
	if err != nil {
		return fmt.Errorf("failed to get capabilities: %w", err)
	}

	fmt.Printf("Device Capabilities:\n")
	fmt.Printf("  Device ID: %s\n", capabilities.DeviceID)

	// Network capabilities
	networkCaps := capabilities.GetNetworkCapabilities()
	if len(networkCaps) > 0 {
		fmt.Printf("  Network Capabilities:\n")

		for _, cap := range networkCaps {
			fmt.Printf("    - %s\n", cap)
		}
	}

	// Extended capabilities
	capNames := capabilities.GetCapabilityNames()
	if len(capNames) > 0 {
		fmt.Printf("  Extended Capabilities:\n")

		for _, capName := range capNames {
			capability := capabilities.GetCapabilityByName(capName)
			fmt.Printf("    - %s", capName)

			if capability.URL != "" {
				fmt.Printf(" (%s)", capability.URL)
			}

			fmt.Println()
		}
	}

	return nil
}

// getPresets handles getting device presets
func getPresets(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Getting device presets", clientConfig.Host, clientConfig.Port)

	presets, err := client.GetPresets()
	if err != nil {
		return fmt.Errorf("failed to get presets: %w", err)
	}

	fmt.Printf("Device Presets:\n")

	if len(presets.Preset) == 0 {
		fmt.Printf("  No presets configured\n")
		return nil
	}

	fmt.Printf("  Configured Presets:\n")

	for _, preset := range presets.Preset {
		fmt.Printf("    %d. %s\n", preset.ID, preset.GetDisplayName())
		fmt.Printf("       Source: %s\n", preset.ContentItem.Source)

		if preset.ContentItem.SourceAccount != "" && preset.ContentItem.SourceAccount != preset.ContentItem.Source {
			fmt.Printf("       Account: %s\n", preset.ContentItem.SourceAccount)
		}

		if preset.ContentItem.Location != "" {
			fmt.Printf("       Location: %s\n", preset.ContentItem.Location)
		}

		// Show preset creation time if available
		if preset.CreatedOn != nil && *preset.CreatedOn != 0 {
			createdTime := time.Unix(*preset.CreatedOn, 0)
			fmt.Printf("       Created: %s\n", createdTime.Format("2006-01-02 15:04:05"))
		}

		fmt.Println()
	}

	return nil
}

// selectPreset selects a preset by number (1-6)
func selectPreset(c *cli.Context) error {
	presetNum := c.Int("preset")
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader(fmt.Sprintf("Selecting preset %d", presetNum), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.SelectPreset(presetNum)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to select preset: %v", err))
		return err
	}

	PrintSuccess(fmt.Sprintf("Preset %d selected", presetNum))

	return nil
}

// getTrackInfo gets the track information
func getTrackInfo(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Getting track information", clientConfig.Host, clientConfig.Port)

	fmt.Println("⚠️  WARNING: /trackInfo endpoint times out on real devices.")
	fmt.Println("   Use 'soundtouch-cli now' (playback status) command instead for track information.")

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	trackInfo, err := client.GetTrackInfo()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get track info: %v", err))
		return err
	}

	fmt.Println("Track Information:")
	fmt.Printf("  Source: %s\n", trackInfo.Source)

	if trackInfo.Track != "" {
		fmt.Printf("  Track: %s\n", trackInfo.Track)
	}

	if trackInfo.Artist != "" {
		fmt.Printf("  Artist: %s\n", trackInfo.Artist)
	}

	if trackInfo.Album != "" {
		fmt.Printf("  Album: %s\n", trackInfo.Album)
	}

	if trackInfo.StationName != "" {
		fmt.Printf("  Station: %s\n", trackInfo.StationName)
	}

	fmt.Printf("  Play Status: %s\n", trackInfo.PlayStatus)

	return nil
}
