package main

import (
	"fmt"
	"strings"

	"github.com/gesellix/bose-soundtouch/pkg/models"
	"github.com/urfave/cli/v2"
)

func printSource(source models.SourceItem) {
	fmt.Printf("    • %s", source.GetDisplayName())

	if source.SourceAccount != "" && source.SourceAccount != source.Source {
		fmt.Printf(" (%s)", source.SourceAccount)
	}

	var attributes []string
	if source.IsLocalSource() {
		attributes = append(attributes, "Local")
		attributes = append(attributes, "Available")
	}

	if len(attributes) > 0 {
		fmt.Printf(" [%s]", strings.Join(attributes, ", "))
	}

	fmt.Println()
}

// listSources handles listing available audio sources
func listSources(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Getting available sources", clientConfig.Host, clientConfig.Port)

	sources, err := client.GetSources()
	if err != nil {
		return fmt.Errorf("failed to get sources: %w", err)
	}

	fmt.Printf("Available Audio Sources:\n")
	fmt.Printf("  Device ID: %s\n", sources.DeviceID)

	// Show ready sources first
	availableSources := sources.GetAvailableSources()
	if len(availableSources) > 0 {
		fmt.Printf("  Ready Sources:\n")

		for _, source := range availableSources {
			printSource(source)
		}
	}

	// Show all configured sources
	fmt.Printf("  All Sources:\n")

	for _, source := range sources.SourceItem {
		status := "Available"
		if !source.IsLocalSource() {
			status = "Remote"
		}

		fmt.Printf("    • %s (%s)\n", source.GetDisplayName(), status)

		if source.SourceAccount != "" && source.SourceAccount != source.Source {
			fmt.Printf("      Account: %s\n", source.SourceAccount)
		}
	}

	// Show streaming sources
	streamingSources := sources.GetStreamingSources()
	if len(streamingSources) > 0 {
		fmt.Printf("  Streaming Services:\n")

		for _, source := range streamingSources {
			fmt.Printf("    • %s", source.GetDisplayName())

			if source.SourceAccount != "" {
				fmt.Printf(" (%s)", source.SourceAccount)
			}

			fmt.Println()
		}
	}

	// Show service availability summary
	fmt.Println()
	checker := NewServiceAvailabilityChecker(client)
	checker.PrintServiceAvailabilitySummary()

	return nil
}

// selectSource handles selecting an audio source
func selectSource(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	sourceName := strings.ToUpper(c.String("source"))
	sourceAccount := c.String("account")

	// Check service availability
	checker := NewServiceAvailabilityChecker(client)

	actionDescription := fmt.Sprintf("select %s source", strings.ToLower(sourceName))
	if !checker.CheckSourceAvailable(sourceName, actionDescription) {
		return fmt.Errorf("source '%s' is not available", sourceName)
	}

	PrintDeviceHeader(fmt.Sprintf("Selecting source '%s'", sourceName), clientConfig.Host, clientConfig.Port)

	err = client.SelectSource(sourceName, sourceAccount)
	if err != nil {
		return fmt.Errorf("failed to select source: %w", err)
	}

	if sourceAccount != "" {
		PrintSuccess(fmt.Sprintf("Source '%s' with account '%s' selected", sourceName, sourceAccount))
	} else {
		PrintSuccess(fmt.Sprintf("Source '%s' selected", sourceName))
	}

	return nil
}

// selectSpotify handles selecting Spotify source
func selectSpotify(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	// Check Spotify availability
	checker := NewServiceAvailabilityChecker(client)
	if !checker.ValidateSpotifyAvailable("select Spotify source") {
		return fmt.Errorf("spotify is not available on this device")
	}

	PrintDeviceHeader("Selecting Spotify source", clientConfig.Host, clientConfig.Port)

	err = client.SelectSpotify("")
	if err != nil {
		return fmt.Errorf("failed to select Spotify: %w", err)
	}

	PrintSuccess("Spotify source selected")

	return nil
}

// selectBluetooth handles selecting Bluetooth source
func selectBluetooth(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	// Check Bluetooth availability
	checker := NewServiceAvailabilityChecker(client)
	if !checker.ValidateBluetoothAvailable("select Bluetooth source") {
		return fmt.Errorf("bluetooth is not available on this device")
	}

	PrintDeviceHeader("Selecting Bluetooth source", clientConfig.Host, clientConfig.Port)

	err = client.SelectBluetooth()
	if err != nil {
		return fmt.Errorf("failed to select Bluetooth: %w", err)
	}

	PrintSuccess("Bluetooth source selected")

	return nil
}

// selectAux handles selecting AUX input source
func selectAux(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Selecting AUX input source", clientConfig.Host, clientConfig.Port)

	err = client.SelectAux()
	if err != nil {
		return fmt.Errorf("failed to select AUX: %w", err)
	}

	PrintSuccess("AUX input source selected")

	return nil
}

// getServiceAvailability handles displaying service availability information
func getServiceAvailability(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Getting service availability", clientConfig.Host, clientConfig.Port)

	serviceAvailability, err := client.GetServiceAvailability()
	if err != nil {
		return fmt.Errorf("failed to get service availability: %w", err)
	}

	fmt.Printf("Service Availability Report:\n")
	fmt.Printf("  Total Services: %d\n", serviceAvailability.GetServiceCount())
	fmt.Printf("  Available Services: %d\n", serviceAvailability.GetAvailableServiceCount())
	fmt.Printf("  Unavailable Services: %d\n", serviceAvailability.GetUnavailableServiceCount())

	// Show available services
	fmt.Printf("\n✅ Available Services:\n")
	availableServices := serviceAvailability.GetAvailableServices()
	if len(availableServices) == 0 {
		fmt.Printf("    None\n")
	} else {
		for _, service := range availableServices {
			fmt.Printf("    • %s\n", formatServiceTypeForDisplay(models.ServiceType(service.Type)))
		}
	}

	// Show unavailable services with reasons
	fmt.Printf("\n❌ Unavailable Services:\n")
	unavailableServices := serviceAvailability.GetUnavailableServices()
	if len(unavailableServices) == 0 {
		fmt.Printf("    None\n")
	} else {
		for _, service := range unavailableServices {
			reason := ""
			if service.Reason != "" {
				reason = fmt.Sprintf(" (%s)", service.Reason)
			}
			fmt.Printf("    • %s%s\n", formatServiceTypeForDisplay(models.ServiceType(service.Type)), reason)
		}
	}

	// Show service categories
	fmt.Printf("\n🎵 Streaming Services:\n")
	streamingServices := serviceAvailability.GetStreamingServices()
	availableCount := 0
	for _, service := range streamingServices {
		status := "❌"
		if service.IsAvailable {
			status = "✅"
			availableCount++
		}
		fmt.Printf("    %s %s\n", status, formatServiceTypeForDisplay(models.ServiceType(service.Type)))
	}
	fmt.Printf("    Summary: %d/%d streaming services available\n", availableCount, len(streamingServices))

	fmt.Printf("\n🔗 Local Input Services:\n")
	localServices := serviceAvailability.GetLocalServices()
	localAvailableCount := 0
	for _, service := range localServices {
		status := "❌"
		if service.IsAvailable {
			status = "✅"
			localAvailableCount++
		}
		fmt.Printf("    %s %s\n", status, formatServiceTypeForDisplay(models.ServiceType(service.Type)))
	}
	fmt.Printf("    Summary: %d/%d local services available\n", localAvailableCount, len(localServices))

	return nil
}

// compareSourcesAndAvailability compares configured sources with service availability
func compareSourcesAndAvailability(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Comparing sources and service availability", clientConfig.Host, clientConfig.Port)

	// Get both sources and service availability
	sources, err := client.GetSources()
	if err != nil {
		return fmt.Errorf("failed to get sources: %w", err)
	}

	serviceAvailability, err := client.GetServiceAvailability()
	if err != nil {
		return fmt.Errorf("failed to get service availability: %w", err)
	}

	fmt.Printf("Source vs Availability Comparison:\n\n")

	// Check key services
	comparisons := []struct {
		name                 string
		configuredCheck      func() bool
		availableCheck       func() bool
		getConfiguredSources func() []models.SourceItem
	}{
		{
			"Spotify",
			sources.HasSpotify,
			serviceAvailability.HasSpotify,
			sources.GetSpotifySources,
		},
		{
			"Bluetooth",
			sources.HasBluetooth,
			serviceAvailability.HasBluetooth,
			func() []models.SourceItem { return sources.GetSourcesByType("BLUETOOTH") },
		},
	}

	for _, comp := range comparisons {
		configured := comp.configuredCheck()
		available := comp.availableCheck()

		fmt.Printf("🔍 %s:\n", comp.name)
		fmt.Printf("    Configured: %s\n", boolToStatus(configured))
		fmt.Printf("    Available: %s\n", boolToStatus(available))

		switch {
		case available && !configured:
			fmt.Printf("    💡 %s is available but not configured - consider setting it up\n", comp.name)
		case configured && !available:
			fmt.Printf("    ⚠️  %s is configured but not available - check device status\n", comp.name)

			// Show specific reason if available
			switch comp.name {
			case "Spotify":
				service := serviceAvailability.GetServiceByType(models.ServiceTypeSpotify)
				if service != nil && service.Reason != "" {
					fmt.Printf("    📝 Reason: %s\n", service.Reason)
				}
			case "Bluetooth":
				service := serviceAvailability.GetServiceByType(models.ServiceTypeBluetooth)
				if service != nil && service.Reason != "" {
					fmt.Printf("    📝 Reason: %s\n", service.Reason)
				}
			}
		case configured && available:
			fmt.Printf("    ✅ %s is properly configured and available\n", comp.name)
		default:
			fmt.Printf("    ➖ %s is neither configured nor available\n", comp.name)
		}
		fmt.Println()
	}

	// Summary
	fmt.Printf("📊 Summary:\n")
	fmt.Printf("    Total configured sources: %d\n", sources.GetSourceCount())
	fmt.Printf("    Ready configured sources: %d\n", sources.GetReadySourceCount())
	fmt.Printf("    Total available services: %d\n", serviceAvailability.GetAvailableServiceCount())
	fmt.Printf("    Total possible services: %d\n", serviceAvailability.GetServiceCount())

	return nil
}

// boolToStatus converts boolean to user-friendly status
func boolToStatus(b bool) string {
	if b {
		return "✅ Yes"
	}
	return "❌ No"
}
