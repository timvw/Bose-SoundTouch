package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// getNetworkInfo retrieves network information from the device
func getNetworkInfo(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Getting network information", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	networkInfo, err := client.GetNetworkInfo()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get network info: %v", err))
		return err
	}

	fmt.Println("Network Information:")

	if networkInfo.GetWifiProfileCount() > 0 {
		fmt.Printf("  WiFi Profiles: %d\n", networkInfo.GetWifiProfileCount())
	}

	interfaces := networkInfo.GetInterfaces()
	if len(interfaces) == 0 {
		fmt.Println("  No network interfaces found")
		return nil
	}

	fmt.Printf("  Interfaces (%d):\n", len(interfaces))

	for i := range interfaces {
		iface := &interfaces[i]
		fmt.Printf("\n  Interface %d:\n", i+1)
		fmt.Printf("    Type: %s\n", iface.GetType())

		if iface.GetName() != "" {
			fmt.Printf("    Name: %s\n", iface.GetName())
		}

		if iface.GetIPAddress() != "" {
			fmt.Printf("    IP Address: %s\n", iface.GetIPAddress())
		}

		if iface.GetMacAddress() != "" {
			fmt.Printf("    MAC Address: %s\n", iface.GetMacAddress())
		}

		fmt.Printf("    State: %s\n", iface.GetStateDescription())

		if iface.IsWiFi() {
			if iface.GetSSID() != "" {
				fmt.Printf("    SSID: %s\n", iface.GetSSID())
			}

			if iface.GetSignal() != "" {
				fmt.Printf("    Signal: %s (%d%%)\n", iface.GetSignalDescription(), iface.GetSignalQuality())
			}

			if iface.GetFrequencyKHz() > 0 {
				fmt.Printf("    Frequency: %s (%s)\n", iface.FormatFrequency(), iface.GetFrequencyBand())
			}

			if iface.GetMode() != "" {
				fmt.Printf("    Mode: %s\n", iface.GetModeDescription())
			}
		}
	}

	// Show active connections summary
	activeInterfaces := networkInfo.GetActiveInterfaces()
	if len(activeInterfaces) > 0 {
		fmt.Println("\n  Active Connections:")

		for i := range activeInterfaces {
			iface := &activeInterfaces[i]
			fmt.Printf("    - %s: %s\n", iface.GetType(), iface.GetNetworkSummary())
		}
	}

	return nil
}

// pingDevice pings the device to test connectivity
func pingDevice(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Pinging device", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.Ping()
	if err != nil {
		PrintError(fmt.Sprintf("Ping failed: %v", err))
		return err
	}

	PrintSuccess("Device is reachable")

	return nil
}

// getDeviceURL displays the device's base URL
func getDeviceURL(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	baseURL := client.BaseURL()
	fmt.Printf("Device URL: %s\n", baseURL)

	return nil
}
