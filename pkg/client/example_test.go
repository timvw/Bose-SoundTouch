package client_test

import (
	"fmt"
	"log"
	"time"

	"github.com/gesellix/bose-soundtouch/pkg/client"
	"github.com/gesellix/bose-soundtouch/pkg/models"
)

// Example demonstrates basic device control operations.
func Example() {
	// Create a client for your SoundTouch device
	config := &client.Config{
		Host:    "192.168.1.100",
		Port:    8090,
		Timeout: 10 * time.Second,
	}
	c := client.NewClient(config)

	// Get device information
	info, err := c.GetDeviceInfo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Device: %s\n", info.Name)

	// Control playback
	err = c.Play()
	if err != nil {
		log.Fatal(err)
	}

	// Set volume to 50%
	err = c.SetVolume(50)
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// Device: Living Room Speaker
}

// ExampleClient_GetNowPlaying demonstrates how to get current playback information.
func ExampleClient_GetNowPlaying() {
	config := &client.Config{Host: "192.168.1.100"}
	c := client.NewClient(config)

	nowPlaying, err := c.GetNowPlaying()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Track: %s\n", nowPlaying.Track)
	fmt.Printf("Artist: %s\n", nowPlaying.Artist)
	fmt.Printf("Album: %s\n", nowPlaying.Album)
	fmt.Printf("Source: %s\n", nowPlaying.Source)

	// Output:
	// Track: Bohemian Rhapsody
	// Artist: Queen
	// Album: A Night at the Opera
	// Source: SPOTIFY
}

// ExampleClient_SetVolume demonstrates volume control with validation.
func ExampleClient_SetVolume() {
	config := &client.Config{Host: "192.168.1.100"}
	c := client.NewClient(config)

	// Set volume to 75%
	err := c.SetVolume(75)
	if err != nil {
		log.Fatal(err)
	}

	// Get current volume
	volume, err := c.GetVolume()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Volume: %d\n", volume.ActualVolume)
	fmt.Printf("Muted: %t\n", volume.MuteEnabled)

	// Output:
	// Volume: 75
	// Muted: false
}

// ExampleClient_SelectSource demonstrates how to change audio sources.
func ExampleClient_SelectSource() {
	config := &client.Config{Host: "192.168.1.100"}
	c := client.NewClient(config)

	// Switch to Spotify
	err := c.SelectSource("SPOTIFY", "")
	if err != nil {
		log.Fatal(err)
	}

	// Switch to Bluetooth
	err = c.SelectSource("BLUETOOTH", "")
	if err != nil {
		log.Fatal(err)
	}

	// Switch to AUX input
	err = c.SelectSource("AUX", "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Source changed successfully")

	// Output:
	// Source changed successfully
}

// ExampleClient_SetBass demonstrates bass control.
func ExampleClient_SetBass() {
	config := &client.Config{Host: "192.168.1.100"}
	c := client.NewClient(config)

	// Set bass to +3 (range: -9 to +9)
	err := c.SetBass(3)
	if err != nil {
		log.Fatal(err)
	}

	// Get current bass level
	bass, err := c.GetBass()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Bass level: %d\n", bass.ActualBass)

	// Output:
	// Bass level: 3
}

// ExampleClient_SetBalance demonstrates balance control.
func ExampleClient_SetBalance() {
	config := &client.Config{Host: "192.168.1.100"}
	c := client.NewClient(config)

	// Set balance slightly to the right (range: -50 to +50)
	err := c.SetBalance(10)
	if err != nil {
		log.Fatal(err)
	}

	// Get current balance
	balance, err := c.GetBalance()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Balance: %d\n", balance.ActualBalance)

	// Output:
	// Balance: 10
}

// ExampleClient_SetZone demonstrates multiroom zone management.
func ExampleClient_SetZone() {
	config := &client.Config{Host: "192.168.1.100"}
	c := client.NewClient(config)

	// Create a zone with multiple speakers
	zone := &models.ZoneRequest{
		Master: "192.168.1.100",
		Members: []models.MemberEntry{
			{IP: "192.168.1.101"},
			{IP: "192.168.1.102"},
		},
	}

	err := c.SetZone(zone)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Zone created successfully")

	// Output:
	// Zone created successfully
}

// ExampleClient_GetPresets demonstrates how to retrieve configured presets.
func ExampleClient_GetPresets() {
	config := &client.Config{Host: "192.168.1.100"}
	c := client.NewClient(config)

	presets, err := c.GetPresets()
	if err != nil {
		log.Fatal(err)
	}

	for _, preset := range presets.Preset {
		fmt.Printf("Preset %d: %s (%s)\n", preset.ID, preset.GetDisplayName(), preset.GetSource())
	}

	// Output:
	// Preset 1: Morning Jazz (SPOTIFY)
	// Preset 2: Classic Rock (SPOTIFY)
	// Preset 3: NPR News (INTERNET_RADIO)
}

// ExampleClient_NewWebSocketClient demonstrates WebSocket client creation.
func ExampleClient_NewWebSocketClient() {
	config := &client.Config{Host: "192.168.1.100"}
	c := client.NewClient(config)

	// Create WebSocket client for real-time events
	wsClient := c.NewWebSocketClient(nil)

	// Connect to device WebSocket
	err := wsClient.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer wsClient.Disconnect()

	fmt.Printf("WebSocket connected: %t\n", wsClient.IsConnected())

	// Output:
	// WebSocket connected: true
}

// ExampleClient_SendKey demonstrates sending key commands.
func ExampleClient_SendKey() {
	config := &client.Config{Host: "192.168.1.100"}
	c := client.NewClient(config)

	// Send various key commands
	commands := []string{"PLAY", "PAUSE", "NEXT_TRACK", "PREV_TRACK", "MUTE"}

	for _, cmd := range commands {
		err := c.SendKey(cmd)
		if err != nil {
			log.Printf("Failed to send %s: %v", cmd, err)
			continue
		}
		fmt.Printf("Sent command: %s\n", cmd)
	}

	// Output:
	// Sent command: PLAY
	// Sent command: PAUSE
	// Sent command: NEXT_TRACK
	// Sent command: PREV_TRACK
	// Sent command: MUTE
}

// ExampleClient_GetCapabilities demonstrates how to check device capabilities.
func ExampleClient_GetCapabilities() {
	config := &client.Config{Host: "192.168.1.100"}
	c := client.NewClient(config)

	capabilities, err := c.GetCapabilities()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Device supports %d capabilities\n", len(capabilities.Capability))
	for _, capability := range capabilities.Capability {
		fmt.Printf("- %s (URL: %s)\n", capability.Name, capability.URL)
	}

	// Output:
	// Device supports 5 capabilities
	// - VOLUME (/volume)
	// - BASS (/bass)
	// - SOURCES (/sources)
	// - PRESETS (/presets)
	// - ZONE (/getZone)
}
