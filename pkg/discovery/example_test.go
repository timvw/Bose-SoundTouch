package discovery_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gesellix/bose-soundtouch/pkg/discovery"
)

// Example demonstrates basic device discovery.
func Example() {
	ctx := context.Background()
	timeout := 5 * time.Second

	// Discover all SoundTouch devices on the network
	devices, err := discovery.DiscoverDevices(ctx, timeout)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d devices:\n", len(devices))
	for _, device := range devices {
		fmt.Printf("- %s at %s:%d\n", device.Name, device.Host, device.Port)
	}

	// Output:
	// Found 2 devices:
	// - Living Room at 192.168.1.100:8090
	// - Kitchen at 192.168.1.101:8090
}

// ExampleDiscoverDevices demonstrates discovering devices with timeout.
func ExampleDiscoverDevices() {
	ctx := context.Background()

	// Quick discovery with 3 second timeout
	devices, err := discovery.DiscoverDevices(ctx, 3*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	if len(devices) == 0 {
		fmt.Println("No SoundTouch devices found")
		return
	}

	// Print detailed device information
	for _, device := range devices {
		fmt.Printf("Device: %s\n", device.Name)
		fmt.Printf("  Address: %s:%d\n", device.Host, device.Port)
		fmt.Printf("  MAC: %s\n", device.MACAddress)
		fmt.Printf("  Method: %s\n", device.DiscoveryMethod)
		fmt.Printf("  URL: %s\n", device.BaseURL)
		fmt.Println()
	}

	// Output:
	// Device: Living Room
	//   Address: 192.168.1.100:8090
	//   MAC: AA:BB:CC:DD:EE:FF
	//   Method: UPnP
	//   URL: http://192.168.1.100:8090
	//
	// Device: Kitchen
	//   Address: 192.168.1.101:8090
	//   MAC: BB:CC:DD:EE:FF:AA
	//   Method: mDNS
	//   URL: http://192.168.1.101:8090
}

// ExampleUnifiedDiscoveryService_DiscoverWithCache demonstrates caching functionality.
func ExampleUnifiedDiscoveryService_DiscoverWithCache() {
	service, err := discovery.NewUnifiedDiscoveryService()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	timeout := 5 * time.Second

	// First discovery scan
	fmt.Println("First scan:")
	devices, err := service.DiscoverWithCache(ctx, timeout)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d devices\n", len(devices))

	// Second scan (should use cache)
	fmt.Println("Second scan (cached):")
	devices, err = service.DiscoverWithCache(ctx, timeout)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d devices (from cache)\n", len(devices))

	// Output:
	// First scan:
	// Found 2 devices
	// Second scan (cached):
	// Found 2 devices (from cache)
}

// ExampleUnifiedDiscoveryService_DiscoverUPnP demonstrates UPnP-only discovery.
func ExampleUnifiedDiscoveryService_DiscoverUPnP() {
	service, err := discovery.NewUnifiedDiscoveryService()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	timeout := 3 * time.Second

	// Use only UPnP/SSDP discovery
	devices, err := service.DiscoverUPnP(ctx, timeout)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("UPnP discovered %d devices:\n", len(devices))
	for _, device := range devices {
		fmt.Printf("- %s (Method: %s)\n", device.Name, device.DiscoveryMethod)
	}

	// Output:
	// UPnP discovered 1 devices:
	// - Living Room (Method: UPnP)
}

// ExampleUnifiedDiscoveryService_DiscoverMDNS demonstrates mDNS-only discovery.
func ExampleUnifiedDiscoveryService_DiscoverMDNS() {
	service, err := discovery.NewUnifiedDiscoveryService()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	timeout := 3 * time.Second

	// Use only mDNS discovery
	devices, err := service.DiscoverMDNS(ctx, timeout)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("mDNS discovered %d devices:\n", len(devices))
	for _, device := range devices {
		fmt.Printf("- %s (Method: %s)\n", device.Name, device.DiscoveryMethod)
	}

	// Output:
	// mDNS discovered 1 devices:
	// - Kitchen (Method: mDNS)
}

// ExampleService_Discover demonstrates basic UPnP discovery service.
func ExampleService_Discover() {
	service := discovery.NewService()
	ctx := context.Background()
	timeout := 5 * time.Second

	devices, err := service.Discover(ctx, timeout)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("UPnP/SSDP found %d devices:\n", len(devices))
	for _, device := range devices {
		fmt.Printf("- %s at %s\n", device.Name, device.BaseURL)
	}

	// Output:
	// UPnP/SSDP found 1 devices:
	// - Living Room at http://192.168.1.100:8090
}

// ExampleMDNSDiscoveryService_Discover demonstrates mDNS discovery service.
func ExampleMDNSDiscoveryService_Discover() {
	service := discovery.NewMDNSDiscoveryService()
	ctx := context.Background()
	timeout := 5 * time.Second

	devices, err := service.Discover(ctx, timeout)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("mDNS found %d devices:\n", len(devices))
	for _, device := range devices {
		fmt.Printf("- %s at %s\n", device.Name, device.BaseURL)
	}

	// Output:
	// mDNS found 1 devices:
	// - Kitchen at http://192.168.1.101:8090
}

// Example_errorHandling demonstrates proper error handling in discovery.
func Example_errorHandling() {
	ctx := context.Background()

	// Very short timeout to demonstrate timeout handling
	shortTimeout := 100 * time.Millisecond

	devices, err := discovery.DiscoverDevices(ctx, shortTimeout)
	if err != nil {
		fmt.Printf("Discovery error: %v\n", err)
		return
	}

	if len(devices) == 0 {
		fmt.Println("No devices found - check network connectivity")
		return
	}

	fmt.Printf("Found %d devices despite short timeout\n", len(devices))

	// Output:
	// No devices found - check network connectivity
}

// Example_contextCancellation demonstrates context cancellation.
func Example_contextCancellation() {
	// Create a context that cancels after 2 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	devices, err := discovery.DiscoverDevices(ctx, 10*time.Second)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Discovery cancelled due to context timeout")
		} else {
			fmt.Printf("Discovery error: %v\n", err)
		}
		return
	}

	fmt.Printf("Found %d devices before context cancellation\n", len(devices))

	// Output:
	// Discovery cancelled due to context timeout
}
