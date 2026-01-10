package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gesellix/bose-soundtouch/pkg/config"
	"github.com/gesellix/bose-soundtouch/pkg/discovery"
	"github.com/urfave/cli/v2"
)

// discoverDevices handles device discovery command
func discoverDevices(c *cli.Context) error {
	httpTimeout := c.Duration("timeout")
	showAll := c.Bool("all")

	fmt.Printf("Discovering SoundTouch devices...\n")

	// Load configuration
	cfg, err := config.LoadFromEnv()
	if err != nil {
		cfg = config.DefaultConfig()
	}

	// Only override discovery timeout if user explicitly provided --timeout flag
	// This respects DISCOVERY_TIMEOUT from .env file when no flag is provided
	if c.IsSet("timeout") {
		cfg.HTTPTimeout = httpTimeout
		// Set discovery timeout to be 2x HTTP timeout (min 5s, max 30s)
		discoveryTimeout := httpTimeout * 2
		if discoveryTimeout < 5*time.Second {
			discoveryTimeout = 5 * time.Second
		}

		if discoveryTimeout > 30*time.Second {
			discoveryTimeout = 30 * time.Second
		}

		cfg.DiscoveryTimeout = discoveryTimeout
	}

	if showAll {
		fmt.Printf("HTTP Timeout: %v\n", cfg.HTTPTimeout)
		fmt.Printf("Discovery Timeout: %v\n", cfg.DiscoveryTimeout)
		fmt.Printf("Mode: Detailed information\n")
	}

	fmt.Println()

	// Create discovery service
	discoveryService := discovery.NewUnifiedDiscoveryService(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.DiscoveryTimeout+5*time.Second)
	defer cancel()

	// Perform discovery
	devices, err := discoveryService.DiscoverDevices(ctx)
	if err != nil {
		return fmt.Errorf("discovery failed: %w", err)
	}

	if len(devices) == 0 {
		fmt.Println("No SoundTouch devices found on the network.")
		fmt.Println()
		fmt.Println("This could mean:")
		fmt.Println("- No SoundTouch devices are powered on")
		fmt.Println("- Devices are on a different network segment")
		fmt.Println("- Network blocks multicast traffic")
		fmt.Println("- Firewall is blocking discovery ports")

		return nil
	}

	// Display results
	fmt.Printf("Found %d SoundTouch device(s):\n\n", len(devices))

	for i, device := range devices {
		fmt.Printf("%d. %s\n", i+1, device.Name)
		fmt.Printf("   Host: %s:%d\n", device.Host, device.Port)
		fmt.Printf("   Model: %s\n", device.ModelID)

		if device.SerialNo != "" {
			fmt.Printf("   Serial: %s\n", device.SerialNo)
		}

		if device.Location != "" {
			fmt.Printf("   Location: %s\n", device.Location)
		}

		if showAll {
			fmt.Printf("   Last Seen: %s\n", device.LastSeen.Format("2006-01-02 15:04:05"))
		}

		// Add spacing between devices
		if i < len(devices)-1 {
			fmt.Println()
		}
	}

	fmt.Println()
	fmt.Printf("Use any of these hosts with other commands:\n")
	fmt.Printf("Example: soundtouch-cli info --host %s\n", devices[0].Host)

	return nil
}
