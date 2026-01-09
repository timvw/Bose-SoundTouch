package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/user_account/bose-soundtouch/pkg/models"
	"github.com/urfave/cli/v2"
)

// getZone retrieves the current zone configuration
func getZone(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Getting zone information", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	zone, err := client.GetZone()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get zone: %v", err))
		return err
	}

	if zone.Master == "" {
		fmt.Println("Device is not in a zone")
		return nil
	}

	fmt.Println("Zone Configuration:")
	fmt.Printf("  Master: %s\n", zone.Master)

	if len(zone.Members) > 0 {
		fmt.Printf("  Members (%d):\n", len(zone.Members))
		for _, member := range zone.Members {
			fmt.Printf("    - %s", member.DeviceID)
			if member.IP != "" {
				fmt.Printf(" (IP: %s)", member.IP)
			}
			fmt.Println()
		}
	} else {
		fmt.Println("  Members: none (standalone device)")
	}

	return nil
}

// getZoneStatus retrieves the zone status
func getZoneStatus(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Getting zone status", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	status, err := client.GetZoneStatus()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get zone status: %v", err))
		return err
	}

	fmt.Printf("Zone Status: %s\n", status)

	inZone, err := client.IsInZone()
	if err != nil {
		PrintWarning(fmt.Sprintf("Could not determine zone membership: %v", err))
	} else {
		fmt.Printf("In Zone: %t\n", inZone)
	}

	return nil
}

// getZoneMembers lists all zone members
func getZoneMembers(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Getting zone members", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	members, err := client.GetZoneMembers()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get zone members: %v", err))
		return err
	}

	if len(members) == 0 {
		fmt.Println("No zone members found")
		return nil
	}

	fmt.Printf("Zone Members (%d):\n", len(members))
	for i, member := range members {
		fmt.Printf("  %d. %s", i+1, member)
		if member == clientConfig.Host {
			fmt.Print(" (this device)")
		}
		fmt.Println()
	}

	return nil
}

// createZone creates a new zone with specified members
func createZone(c *cli.Context) error {
	members := c.StringSlice("members")
	clientConfig := GetClientConfig(c)

	if len(members) == 0 {
		PrintError("At least one member must be specified")
		return fmt.Errorf("no members specified")
	}

	PrintDeviceHeader(fmt.Sprintf("Creating zone with %d members", len(members)), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	// Parse and validate member IPs
	var memberIPs []net.IP
	for _, member := range members {
		ip := net.ParseIP(member)
		if ip == nil {
			PrintError(fmt.Sprintf("Invalid IP address: %s", member))
			return fmt.Errorf("invalid IP address: %s", member)
		}
		memberIPs = append(memberIPs, ip)
	}

	// For simplicity, use the first member as master and rest as members
	// In a real scenario, you might want to specify the master separately
	masterDeviceID := "master" // This would need to be a real device ID
	memberMap := make(map[string]string)
	for i, ip := range memberIPs {
		deviceID := fmt.Sprintf("device_%d", i+1)
		memberMap[deviceID] = ip.String()
	}

	err = client.CreateZoneWithIPs(masterDeviceID, memberMap)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create zone: %v", err))
		return err
	}

	PrintSuccess(fmt.Sprintf("Zone created with members: %s", strings.Join(members, ", ")))
	return nil
}

// addToZone adds a device to the current zone
func addToZone(c *cli.Context) error {
	memberIP := c.String("member")
	clientConfig := GetClientConfig(c)

	if memberIP == "" {
		PrintError("Member IP address is required")
		return fmt.Errorf("member IP is required")
	}

	PrintDeviceHeader(fmt.Sprintf("Adding %s to zone", memberIP), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	// Validate IP address
	ip := net.ParseIP(memberIP)
	if ip == nil {
		PrintError(fmt.Sprintf("Invalid IP address: %s", memberIP))
		return fmt.Errorf("invalid IP address: %s", memberIP)
	}

	// For this example, we'll use the IP as the device ID
	// In practice, you'd need the actual device ID
	deviceID := memberIP
	err = client.AddToZone(deviceID, memberIP)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to add to zone: %v", err))
		return err
	}

	PrintSuccess(fmt.Sprintf("Added %s to zone", memberIP))
	return nil
}

// removeFromZone removes a device from the current zone
func removeFromZone(c *cli.Context) error {
	memberIP := c.String("member")
	clientConfig := GetClientConfig(c)

	if memberIP == "" {
		PrintError("Member IP address is required")
		return fmt.Errorf("member IP is required")
	}

	PrintDeviceHeader(fmt.Sprintf("Removing %s from zone", memberIP), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	// Validate IP address
	ip := net.ParseIP(memberIP)
	if ip == nil {
		PrintError(fmt.Sprintf("Invalid IP address: %s", memberIP))
		return fmt.Errorf("invalid IP address: %s", memberIP)
	}

	// For this example, we'll use the IP as the device ID
	// In practice, you'd need the actual device ID
	deviceID := memberIP
	err = client.RemoveFromZone(deviceID)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to remove from zone: %v", err))
		return err
	}

	PrintSuccess(fmt.Sprintf("Removed %s from zone", memberIP))
	return nil
}

// dissolveZone dissolves the current zone
func dissolveZone(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Dissolving zone", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.DissolveZone()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to dissolve zone: %v", err))
		return err
	}

	PrintSuccess("Zone dissolved")
	return nil
}

// setZoneConfig sets zone configuration
func setZoneConfig(c *cli.Context) error {
	master := c.String("master")
	members := c.StringSlice("members")
	clientConfig := GetClientConfig(c)

	if master == "" {
		PrintError("Master device IP is required")
		return fmt.Errorf("master IP is required")
	}

	PrintDeviceHeader(fmt.Sprintf("Setting zone with master %s", master), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	// Create zone configuration using ZoneRequest
	zoneRequest := models.NewZoneRequest(master)

	// Add members if specified
	if len(members) > 0 {
		for _, memberIP := range members {
			// Validate IP address
			if net.ParseIP(memberIP) == nil {
				PrintError(fmt.Sprintf("Invalid member IP address: %s", memberIP))
				return fmt.Errorf("invalid IP address: %s", memberIP)
			}

			// Use IP as device ID for simplicity - in practice you'd need real device IDs
			deviceID := memberIP
			zoneRequest.AddMember(deviceID, memberIP)
		}
	}

	err = client.SetZone(zoneRequest)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to set zone: %v", err))
		return err
	}

	if len(members) > 0 {
		PrintSuccess(fmt.Sprintf("Zone configured with master %s and members: %s", master, strings.Join(members, ", ")))
	} else {
		PrintSuccess(fmt.Sprintf("Zone configured with master %s", master))
	}

	return nil
}
