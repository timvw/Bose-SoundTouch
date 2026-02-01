package main

import (
	"fmt"
	"time"

	"github.com/gesellix/bose-soundtouch/pkg/models"
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

// getSupportedURLs handles getting supported URLs/endpoints
func getSupportedURLs(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Getting supported URLs", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	supportedURLs, err := client.GetSupportedURLs()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get supported URLs: %v", err))
		return err
	}

	printSupportedURLs(supportedURLs, c)

	return nil
}

// printSupportedURLs formats and displays supported URLs information
func printSupportedURLs(supportedURLs *models.SupportedURLsResponse, c *cli.Context) {
	verbose := c.Bool("verbose")
	showFeatures := c.Bool("features")

	fmt.Printf("Device Supported URLs:\n")
	fmt.Printf("  Device ID: %s\n", supportedURLs.DeviceID)
	fmt.Printf("  Total Endpoints: %d\n", supportedURLs.GetURLCount())

	// Show feature completeness score
	completeness, supported, total := supportedURLs.GetFeatureCompleteness()
	fmt.Printf("  Feature Coverage: %d%% (%d/%d features)\n\n", completeness, supported, total)

	if showFeatures || (!verbose && !showFeatures) {
		// Show feature mapping (default view)
		printFeatureMapping(supportedURLs, verbose)
	}

	if verbose {
		fmt.Println()
		printDetailedEndpoints(supportedURLs)
	}

	if !showFeatures && !verbose {
		fmt.Printf("\n💡 Options:\n")
		fmt.Printf("   --features  Show detailed feature mapping and CLI commands\n")
		fmt.Printf("   --verbose   Show complete endpoint list\n")
	}
}

// printFeatureMapping displays the feature-to-endpoint mapping
func printFeatureMapping(supportedURLs *models.SupportedURLsResponse, verbose bool) {
	fmt.Printf("🎯 Device Feature Support:\n\n")

	// Get features organized by category
	featuresByCategory := supportedURLs.GetFeaturesByCategory()
	printFeatureCategories(featuresByCategory, supportedURLs, verbose)
	printMissingEssentialFeatures(supportedURLs)
	printPartiallyImplementedFeatures(supportedURLs, verbose)
}

func printFeatureCategories(featuresByCategory map[string][]models.EndpointFeature, supportedURLs *models.SupportedURLsResponse, verbose bool) {
	categoryInfo := map[string]string{
		"Core":      "⚡",
		"Audio":     "🔊",
		"Playback":  "▶️",
		"Sources":   "📱",
		"Content":   "📻",
		"Presets":   "⭐",
		"Multiroom": "🏠",
		"Network":   "🌐",
		"System":    "⚙️",
	}

	categoryOrder := []string{"Core", "Audio", "Playback", "Sources", "Content", "Presets", "Multiroom", "Network", "System"}

	for _, category := range categoryOrder {
		features := featuresByCategory[category]
		if len(features) == 0 {
			continue
		}

		emoji := categoryInfo[category]
		fmt.Printf("%s %s (%d features):\n", emoji, category, len(features))

		for _, feature := range features {
			printFeatureStatus(feature, supportedURLs, verbose)
		}

		fmt.Println()
	}
}

func printFeatureStatus(feature models.EndpointFeature, supportedURLs *models.SupportedURLsResponse, verbose bool) {
	supportedEndpoints := countSupportedEndpoints(feature, supportedURLs)

	status := "✅"
	if supportedEndpoints < len(feature.Endpoints) && len(feature.Endpoints) > 1 {
		status = "⚠️" // Partial support
	}

	fmt.Printf("    %s %s", status, feature.Name)

	if feature.Essential {
		fmt.Printf(" ⭐")
	}

	fmt.Printf("\n")

	if verbose {
		printVerboseFeatureDetails(feature, supportedEndpoints)
	}
}

func countSupportedEndpoints(feature models.EndpointFeature, supportedURLs *models.SupportedURLsResponse) int {
	supportedEndpoints := 0

	for _, endpoint := range feature.Endpoints {
		if supportedURLs.HasURL(endpoint) {
			supportedEndpoints++
		}
	}

	return supportedEndpoints
}

func printVerboseFeatureDetails(feature models.EndpointFeature, supportedEndpoints int) {
	fmt.Printf("        %s\n", feature.Description)
	fmt.Printf("        CLI: %s\n", feature.CLICommand)
	fmt.Printf("        Endpoints: %d/%d supported", supportedEndpoints, len(feature.Endpoints))

	if supportedEndpoints < len(feature.Endpoints) {
		fmt.Printf(" (partial)")
	}

	fmt.Printf("\n")
}

func printMissingEssentialFeatures(supportedURLs *models.SupportedURLsResponse) {
	missingEssential := supportedURLs.GetMissingEssentialFeatures()
	if len(missingEssential) > 0 {
		fmt.Printf("⚠️  Missing Essential Features:\n")

		for _, feature := range missingEssential {
			fmt.Printf("    ❌ %s - %s\n", feature.Name, feature.Description)
		}

		fmt.Println()
	}
}

func printPartiallyImplementedFeatures(supportedURLs *models.SupportedURLsResponse, verbose bool) {
	partial := supportedURLs.GetPartiallyImplementedFeatures()
	if len(partial) > 0 && verbose {
		fmt.Printf("⚠️  Partially Supported Features:\n")

		for _, feature := range partial {
			fmt.Printf("    🟡 %s\n", feature.Name)

			for _, endpoint := range feature.Endpoints {
				status := "❌"
				if supportedURLs.HasURL(endpoint) {
					status = "✅"
				}

				fmt.Printf("        %s %s\n", status, endpoint)
			}
		}

		fmt.Println()
	}
}

// printDetailedEndpoints shows the traditional endpoint listing
func printDetailedEndpoints(supportedURLs *models.SupportedURLsResponse) {
	fmt.Printf("📋 Detailed Endpoint Analysis:\n\n")

	// Show core functionality
	coreURLs := supportedURLs.GetCoreURLs()
	if len(coreURLs) > 0 {
		fmt.Printf("🎮 Core Functionality (%d endpoints):\n", len(coreURLs))

		for _, url := range coreURLs {
			fmt.Printf("    • %s\n", url)
		}

		fmt.Println()
	}

	// Show streaming functionality
	streamingURLs := supportedURLs.GetStreamingURLs()
	if len(streamingURLs) > 0 {
		fmt.Printf("📻 Streaming Services (%d endpoints):\n", len(streamingURLs))

		for _, url := range streamingURLs {
			fmt.Printf("    • %s\n", url)
		}

		fmt.Println()
	}

	// Show advanced audio functionality
	advancedURLs := supportedURLs.GetAdvancedURLs()
	if len(advancedURLs) > 0 {
		fmt.Printf("🔧 Advanced Audio (%d endpoints):\n", len(advancedURLs))

		for _, url := range advancedURLs {
			fmt.Printf("    • %s\n", url)
		}

		fmt.Println()
	}

	// Show network functionality
	networkURLs := supportedURLs.GetNetworkURLs()
	if len(networkURLs) > 0 {
		fmt.Printf("🌐 Network & Connectivity (%d endpoints):\n", len(networkURLs))

		for _, url := range networkURLs {
			fmt.Printf("    • %s\n", url)
		}

		fmt.Println()
	}

	// Show all supported URLs
	fmt.Printf("📝 Complete Endpoint List:\n")

	allURLs := supportedURLs.GetURLs()
	for i, url := range allURLs {
		fmt.Printf("    %3d. %s\n", i+1, url)
	}
}

// getDeviceAnalysis handles comprehensive device capability analysis
func getDeviceAnalysis(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Analyzing device capabilities", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	supportedURLs, err := client.GetSupportedURLs()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get supported URLs: %v", err))
		return err
	}

	printDeviceAnalysis(supportedURLs)

	return nil
}

// printDeviceAnalysis provides comprehensive device capability analysis
func printDeviceAnalysis(supportedURLs *models.SupportedURLsResponse) {
	fmt.Printf("🔍 Device Capability Analysis:\n")
	fmt.Printf("  Device ID: %s\n", supportedURLs.DeviceID)

	// Overall score
	completeness, supported, total := supportedURLs.GetFeatureCompleteness()
	fmt.Printf("  Feature Coverage: %d%% (%d/%d features)\n", completeness, supported, total)

	// Device classification
	classification := classifyDevice(supportedURLs)
	fmt.Printf("  Device Type: %s\n\n", classification)

	// Essential features check
	missingEssential := supportedURLs.GetMissingEssentialFeatures()
	if len(missingEssential) > 0 {
		fmt.Printf("❌ Missing Essential Features:\n")

		for _, feature := range missingEssential {
			fmt.Printf("    • %s - %s\n", feature.Name, feature.Description)
			fmt.Printf("      Impact: Device may not function properly without this\n")
		}

		fmt.Println()
	} else {
		fmt.Printf("✅ All essential features are supported\n\n")
	}

	// Show what works
	supportedFeatures := supportedURLs.GetSupportedFeatures()
	fmt.Printf("✅ Available Features (%d):\n", len(supportedFeatures))

	categoryCount := make(map[string]int)
	for _, feature := range supportedFeatures {
		categoryCount[feature.Category]++
	}

	for category, count := range categoryCount {
		emoji := getCategoryEmoji(category)
		fmt.Printf("    %s %s: %d features\n", emoji, category, count)
	}

	fmt.Println()

	// Show what's missing
	unsupportedFeatures := supportedURLs.GetUnsupportedFeatures()
	if len(unsupportedFeatures) > 0 {
		fmt.Printf("❌ Unsupported Features (%d):\n", len(unsupportedFeatures))

		for _, feature := range unsupportedFeatures {
			fmt.Printf("    • %s - %s\n", feature.Name, feature.Description)
		}

		fmt.Println()
	}

	// Partial implementations
	partial := supportedURLs.GetPartiallyImplementedFeatures()
	if len(partial) > 0 {
		fmt.Printf("⚠️ Partially Supported Features (%d):\n", len(partial))

		for _, feature := range partial {
			supportedCount := 0

			for _, endpoint := range feature.Endpoints {
				if supportedURLs.HasURL(endpoint) {
					supportedCount++
				}
			}

			fmt.Printf("    • %s (%d/%d endpoints)\n", feature.Name, supportedCount, len(feature.Endpoints))
		}

		fmt.Println()
	}

	// Recommendations
	printRecommendations(supportedURLs)

	// CLI usage suggestions
	printCLIUsageSuggestions(supportedURLs)
}

// classifyDevice determines the device type based on supported features
func classifyDevice(supportedURLs *models.SupportedURLsResponse) string {
	if supportedURLs.HasMultiroomSupport() && supportedURLs.HasAdvancedAudioSupport() {
		return "Premium SoundTouch Speaker (Full Feature Set)"
	}

	if supportedURLs.HasMultiroomSupport() {
		return "Standard SoundTouch Speaker (Multiroom Capable)"
	}

	if supportedURLs.HasStreamingSupport() && supportedURLs.HasPresetSupport() {
		return "Basic SoundTouch Speaker"
	}

	if supportedURLs.HasCorePlaybackSupport() {
		return "Essential SoundTouch Device"
	}

	return "Limited SoundTouch Device"
}

// printRecommendations provides usage recommendations based on device capabilities
func printRecommendations(supportedURLs *models.SupportedURLsResponse) {
	fmt.Printf("💡 Recommendations:\n")

	if supportedURLs.HasMultiroomSupport() {
		fmt.Printf("    🏠 This device supports multiroom - you can create speaker groups\n")
		fmt.Printf("       Try: soundtouch-cli zone create --master <this-device> --members <other-devices>\n")
	}

	if supportedURLs.HasPresetSupport() {
		fmt.Printf("    ⭐ Save your favorite content as presets for quick access\n")
		fmt.Printf("       Try: soundtouch-cli preset store-current --slot 1\n")
	}

	if supportedURLs.HasStreamingSupport() {
		fmt.Printf("    📻 Browse and discover new content from streaming services\n")
		fmt.Printf("       Try: soundtouch-cli browse tunein, station search-tunein --query jazz\n")
	}

	if supportedURLs.HasAdvancedAudioSupport() {
		fmt.Printf("    🔧 Fine-tune your audio with advanced controls\n")
		fmt.Printf("       Try: soundtouch-cli audio dsp get, audio tone get\n")
	}

	if !supportedURLs.HasURL("/bassCapabilities") {
		fmt.Printf("    ⚠️  Device may have limited bass control options\n")
	}

	if !supportedURLs.HasURL("/balance") {
		fmt.Printf("    ⚠️  No balance control available on this device\n")
	}

	fmt.Println()
}

// printCLIUsageSuggestions shows common CLI commands for this device
func printCLIUsageSuggestions(supportedURLs *models.SupportedURLsResponse) {
	fmt.Printf("🚀 Common Commands for This Device:\n")

	// Always available
	fmt.Printf("    • Get device info: soundtouch-cli info get\n")
	fmt.Printf("    • Control volume: soundtouch-cli volume set --level 50\n")

	if supportedURLs.HasURL("/nowPlaying") {
		fmt.Printf("    • Check what's playing: soundtouch-cli play now\n")
	}

	if supportedURLs.HasURL("/sources") {
		fmt.Printf("    • List audio sources: soundtouch-cli source list\n")
	}

	if supportedURLs.HasURL("/presets") {
		fmt.Printf("    • Manage presets: soundtouch-cli preset list\n")
	}

	if supportedURLs.HasURL("/bass") {
		fmt.Printf("    • Adjust bass: soundtouch-cli bass set --level 5\n")
	}

	if supportedURLs.HasURL("/setZone") {
		fmt.Printf("    • Create speaker group: soundtouch-cli zone create\n")
	}

	if supportedURLs.HasURL("/search") {
		fmt.Printf("    • Search content: soundtouch-cli station search-tunein --query \"classic rock\"\n")
	}

	fmt.Println()
}

// getCategoryEmoji returns emoji for feature categories
func getCategoryEmoji(category string) string {
	emojis := map[string]string{
		"Core":      "⚡",
		"Audio":     "🔊",
		"Playback":  "▶️",
		"Sources":   "📱",
		"Content":   "📻",
		"Presets":   "⭐",
		"Multiroom": "🏠",
		"Network":   "🌐",
		"System":    "⚙️",
	}
	if emoji, exists := emojis[category]; exists {
		return emoji
	}

	return "📋"
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
