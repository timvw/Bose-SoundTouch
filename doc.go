// Package soundtouch provides a comprehensive Go library and CLI tool for controlling Bose SoundTouch devices.
//
// This library implements the complete Bose SoundTouch Web API, enabling programmatic control
// of SoundTouch speakers including playback control, volume management, source selection,
// multiroom zone management, and real-time event monitoring via WebSocket connections.
//
// # Quick Start
//
// Install the library:
//
//	go get github.com/gesellix/bose-soundtouch
//
// Basic usage example:
//
//	package main
//
//	import (
//		"fmt"
//		"log"
//
//		"github.com/gesellix/bose-soundtouch/pkg/client"
//	)
//
//	func main() {
//		// Create a client for your SoundTouch device
//		config := &client.Config{
//			Host: "192.168.1.100",
//			Port: 8090,
//		}
//		client := client.NewClient(config)
//
//		// Get device information
//		info, err := client.GetInfo()
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Printf("Device: %s\n", info.Name)
//
//		// Control playback
//		err = client.Play()
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		// Set volume
//		err = client.SetVolume(50)
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//
// # Device Discovery
//
// Automatically discover SoundTouch devices on your network:
//
//	import "github.com/gesellix/bose-soundtouch/pkg/discovery"
//
//	// Discover devices using UPnP/SSDP
//	service := discovery.NewService(5*time.Second)
//	devices, err := service.DiscoverDevices(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for _, device := range devices {
//		fmt.Printf("Found device: %s at %s\n", device.Name, device.Host)
//	}
//
// # Real-time Events
//
// Monitor device state changes in real-time using WebSocket connections:
//
//	// Subscribe to device events
//	events, err := client.SubscribeToEvents(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for event := range events {
//		switch e := event.(type) {
//		case *models.NowPlayingUpdated:
//			fmt.Printf("Now playing: %s by %s\n", e.Track, e.Artist)
//		case *models.VolumeUpdated:
//			fmt.Printf("Volume changed to: %d\n", e.ActualVolume)
//		}
//	}
//
// # Multiroom Zone Management
//
// Create and manage multiroom zones:
//
//	// Create a zone with multiple speakers
//	zone := &models.Zone{
//		Master: "192.168.1.100",
//		Members: []models.ZoneMember{
//			{IPAddress: "192.168.1.101"},
//			{IPAddress: "192.168.1.102"},
//		},
//	}
//	err = client.SetZone(zone)
//
// # CLI Tool
//
// The package includes a comprehensive CLI tool for device control:
//
//	# Install the CLI
//	go install github.com/gesellix/bose-soundtouch/cmd/soundtouch-cli@latest
//
//	# Discover devices
//	soundtouch-cli discover devices
//
//	# Control a device
//	soundtouch-cli --host 192.168.1.100 play start
//	soundtouch-cli --host 192.168.1.100 volume set --level 50
//	soundtouch-cli --host 192.168.1.100 source select --source SPOTIFY
//
// # Supported Features
//
//   - ✅ Device Information & Capabilities
//   - ✅ Playback Control (Play/Pause/Stop/Next/Previous)
//   - ✅ Volume, Bass, and Balance Control
//   - ✅ Source Selection (Spotify, Bluetooth, AUX, etc.)
//   - ✅ Preset Management
//   - ✅ Clock/Time Management
//   - ✅ Network Information
//   - ✅ Real-time WebSocket Events
//   - ✅ Multiroom Zone Management
//   - ✅ Device Discovery (UPnP/SSDP and mDNS)
//   - ✅ Cross-platform Support (Windows, macOS, Linux)
//
// # Package Structure
//
//   - client: HTTP client for SoundTouch Web API
//   - discovery: Device discovery using UPnP/SSDP and mDNS
//   - models: Data structures for API requests/responses
//   - config: Configuration management
//   - cmd/soundtouch-cli: Command-line interface tool
//
// # Hardware Compatibility
//
// This library has been tested with real Bose SoundTouch hardware and supports
// all SoundTouch-compatible devices including:
//   - SoundTouch 10, 20, 30 series
//   - SoundTouch Portable
//   - Wave SoundTouch music system
//   - And other SoundTouch-enabled Bose speakers
//
// # Implementation Notes
//
// This implementation is based on the official Bose SoundTouch Web API documentation
// and provides 90% coverage of all available endpoints. It is an independent project
// and is not affiliated with or endorsed by Bose Corporation.
//
// For detailed API documentation, examples, and advanced usage patterns, visit:
// https://pkg.go.dev/github.com/gesellix/bose-soundtouch
package soundtouch
