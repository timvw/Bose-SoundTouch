// Package main provides an example of discovering SoundTouch devices using all three mechanisms.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gesellix/bose-soundtouch/pkg/config"
	"github.com/gesellix/bose-soundtouch/pkg/discovery"
)

func main() {
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	timeout := flag.Duration("timeout", 5*time.Second, "Discovery timeout")
	showConfig := flag.Bool("show-config", false, "Show configuration details")

	flag.Parse()

	// Configure logging
	if *verbose {
		log.SetOutput(os.Stdout)
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	} else {
		log.SetOutput(os.Stderr)
	}

	fmt.Println("SoundTouch Unified Discovery Example")
	fmt.Println("===================================")

	fmt.Printf("Timeout: %v, Verbose: %v\n", *timeout, *verbose)
	fmt.Println()

	// Load configuration
	cfg, err := config.LoadFromEnv()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Override timeout from command line
	cfg.DiscoveryTimeout = *timeout

	if *showConfig {
		printConfiguration(cfg)
	}

	fmt.Println("Testing individual discovery mechanisms:")
	fmt.Println("--------------------------------------")

	testSSDP(cfg, *timeout, *verbose)
	testMDNS(cfg, *timeout, *verbose)
	testConfig(cfg, *verbose)
	testUnified(cfg, *timeout, *verbose)
}

func printConfiguration(cfg *config.Config) {
	fmt.Println("Configuration:")
	fmt.Printf("  UPnP Enabled: %v\n", cfg.UPnPEnabled)
	fmt.Printf("  mDNS Enabled: %v\n", cfg.MDNSEnabled)
	fmt.Printf("  Cache Enabled: %v\n", cfg.CacheEnabled)
	fmt.Printf("  Discovery Timeout: %v\n", cfg.DiscoveryTimeout)
	fmt.Printf("  Preferred Devices: %d\n", len(cfg.PreferredDevices))

	for i, device := range cfg.PreferredDevices {
		fmt.Printf("    %d. %s at %s:%d\n", i+1, device.Name, device.Host, device.Port)
	}

	fmt.Println()
}

func testSSDP(cfg *config.Config, timeout time.Duration, verbose bool) {
	// Test SSDP discovery
	fmt.Println("1. SSDP/UPnP Discovery:")

	if cfg.UPnPEnabled {
		// Create fresh context for SSDP test
		ssdpCtx, ssdpCancel := context.WithTimeout(context.Background(), timeout+2*time.Second)
		defer ssdpCancel()

		ssdpService := discovery.NewServiceWithConfig(cfg)
		start := time.Now()
		ssdpDevices, ssdpErr := ssdpService.DiscoverDevices(ssdpCtx)
		duration := time.Since(start)

		if ssdpErr != nil {
			fmt.Printf("   Error: %v\n", ssdpErr)
		} else {
			fmt.Printf("   Found %d devices in %v\n", len(ssdpDevices), duration)

			for _, device := range ssdpDevices {
				fmt.Printf("   - %s (%s:%d)\n", device.Name, device.Host, device.Port)

				if verbose {
					fmt.Printf("     Info URL: %s\n", device.InfoURL)

					if device.UPnPLocation != "" {
						fmt.Printf("     UPnP Location: %s\n", device.UPnPLocation)
					}
				}
			}
		}
	} else {
		fmt.Println("   Disabled in configuration")
	}

	fmt.Println()
}

func testMDNS(cfg *config.Config, timeout time.Duration, verbose bool) {
	// Test mDNS discovery
	fmt.Println("2. mDNS/Bonjour Discovery:")

	if cfg.MDNSEnabled {
		// Create fresh context for mDNS test
		mdnsCtx, mdnsCancel := context.WithTimeout(context.Background(), timeout+2*time.Second)
		defer mdnsCancel()

		mdnsService := discovery.NewMDNSDiscoveryService(timeout)
		start := time.Now()
		mdnsDevices, mdnsErr := mdnsService.DiscoverDevices(mdnsCtx)
		duration := time.Since(start)

		if mdnsErr != nil {
			fmt.Printf("   Error: %v\n", mdnsErr)
		} else {
			fmt.Printf("   Found %d devices in %v\n", len(mdnsDevices), duration)

			for _, device := range mdnsDevices {
				fmt.Printf("   - %s (%s:%d)\n", device.Name, device.Host, device.Port)

				if verbose {
					fmt.Printf("     Info URL: %s\n", device.InfoURL)

					if device.MDNSHostname != "" {
						fmt.Printf("     mDNS Hostname: %s\n", device.MDNSHostname)
					}
				}
			}
		}
	} else {
		fmt.Println("   Disabled in configuration")
	}

	fmt.Println()
}

func testConfig(cfg *config.Config, verbose bool) {
	// Test configuration-based devices
	fmt.Println("3. Configuration-based Devices:")

	configDevices := cfg.GetPreferredDevicesAsDiscovered()
	if len(configDevices) > 0 {
		fmt.Printf("   Found %d configured devices\n", len(configDevices))

		for _, device := range configDevices {
			fmt.Printf("   - %s (%s:%d)\n", device.Name, device.Host, device.Port)

			if verbose {
				fmt.Printf("     Info URL: %s\n", device.InfoURL)
			}
		}
	} else {
		fmt.Println("   No devices configured in .env file")
	}

	fmt.Println()
}

func testUnified(cfg *config.Config, timeout time.Duration, verbose bool) {
	// Test unified discovery
	fmt.Println("4. Unified Discovery (combines all methods):")
	// Create fresh context for unified test
	unifiedCtx, unifiedCancel := context.WithTimeout(context.Background(), timeout+2*time.Second)
	defer unifiedCancel()

	unifiedService := discovery.NewUnifiedDiscoveryService(cfg)
	start := time.Now()
	allDevices, err := unifiedService.DiscoverDevices(unifiedCtx)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("   Error: %v\n", err)
		return
	}

	fmt.Printf("   Found %d total devices in %v\n", len(allDevices), duration)
	fmt.Println()

	if len(allDevices) == 0 {
		fmt.Println("No SoundTouch devices found via any discovery method")
		fmt.Println()
		fmt.Println("This could mean:")
		fmt.Println("- No SoundTouch devices on network")
		fmt.Println("- All discovery methods are disabled")
		fmt.Println("- Network blocks multicast traffic")
		fmt.Println("- Devices are not advertising services")

		return
	}

	fmt.Println("Unified Device List:")
	fmt.Println("-------------------")

	for i, device := range allDevices {
		fmt.Printf("%d. %s\n", i+1, device.Name)
		fmt.Printf("   Host: %s\n", device.Host)
		fmt.Printf("   Port: %d\n", device.Port)
		fmt.Printf("   API Base URL: %s\n", device.APIBaseURL)
		fmt.Printf("   Info URL: %s\n", device.InfoURL)
		fmt.Printf("   Discovery Method: %s\n", device.DiscoveryMethod)
		fmt.Printf("   Last seen: %s\n", device.LastSeen.Format("2006-01-02 15:04:05"))

		if verbose {
			if device.ModelID != "" {
				fmt.Printf("   Model ID: %s\n", device.ModelID)
			}

			if device.SerialNo != "" {
				fmt.Printf("   Serial No: %s\n", device.SerialNo)
			}

			// Show protocol-specific details
			if device.UPnPLocation != "" {
				fmt.Printf("   UPnP Location: %s\n", device.UPnPLocation)

				if device.UPnPUSN != "" {
					fmt.Printf("   UPnP USN: %s\n", device.UPnPUSN)
				}
			}

			if device.MDNSHostname != "" {
				fmt.Printf("   mDNS Hostname: %s\n", device.MDNSHostname)

				if device.MDNSService != "" {
					fmt.Printf("   mDNS Service: %s\n", device.MDNSService)
				}
			}

			if device.ConfigName != "" {
				fmt.Printf("   Config Name: %s\n", device.ConfigName)
			}
		}

		fmt.Println()
	}

	fmt.Printf("✓ Unified discovery completed successfully!\n")
	fmt.Printf("✓ Found %d unique device(s) in %v\n", len(allDevices), duration)

	if verbose {
		fmt.Println()
		fmt.Println("Technical Details:")
		fmt.Printf("- SSDP multicast address: 239.255.255.250:1900\n")
		fmt.Printf("- mDNS service type: _soundtouch._tcp.local\n")
		fmt.Printf("- Discovery timeout: %v\n", timeout)
		fmt.Printf("- Configuration file: .env (if present)\n")
	}
}
