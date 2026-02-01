package main

import (
	"fmt"

	"github.com/gesellix/bose-soundtouch/pkg/models"
	"github.com/urfave/cli/v2"
)

// storeCurrentPreset handles storing currently playing content as preset
func storeCurrentPreset(c *cli.Context) error {
	slot := c.Int("slot")
	clientConfig := GetClientConfig(c)

	PrintDeviceHeader(fmt.Sprintf("Storing current content as preset %d", slot), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	// Check what's currently playing
	nowPlaying, err := client.GetNowPlaying()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get current content: %v", err))
		return err
	}

	if nowPlaying.IsEmpty() {
		PrintError("No content currently playing")
		return fmt.Errorf("no content currently playing")
	}

	if nowPlaying.ContentItem == nil {
		PrintError("Current content has no preset information")
		return fmt.Errorf("current content cannot be saved as preset")
	}

	if !nowPlaying.ContentItem.IsPresetable {
		PrintError("Current content cannot be saved as preset")
		fmt.Printf("  Content: %s\n", nowPlaying.Track)
		fmt.Printf("  Source: %s\n", nowPlaying.Source)
		return fmt.Errorf("current content cannot be preset")
	}

	// Show what we're about to store
	fmt.Printf("Current Content:\n")
	fmt.Printf("  Track: %s\n", nowPlaying.Track)
	if nowPlaying.Artist != "" {
		fmt.Printf("  Artist: %s\n", nowPlaying.Artist)
	}
	if nowPlaying.Album != "" {
		fmt.Printf("  Album: %s\n", nowPlaying.Album)
	}
	fmt.Printf("  Source: %s\n", nowPlaying.Source)
	if nowPlaying.ContentItem.Location != "" {
		fmt.Printf("  Location: %s\n", nowPlaying.ContentItem.Location)
	}

	// Store as preset
	err = client.StoreCurrentAsPreset(slot)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to store preset: %v", err))
		return err
	}

	PrintSuccess(fmt.Sprintf("Stored current content as preset %d", slot))

	return nil
}

// storePreset handles storing specific content as preset
func storePreset(c *cli.Context) error {
	slot := c.Int("slot")
	source := c.String("source")
	location := c.String("location")
	sourceAccount := c.String("source-account")
	name := c.String("name")
	itemType := c.String("type")
	artwork := c.String("artwork")

	// Resolve location and source from potential URLs
	resolvedSource, resolvedLocation := resolveLocation(source, location)
	if resolvedLocation != location && (source == "" || source == "TUNEIN") {
		// If location was a TuneIn URL, fetch metadata if name or artwork is missing
		if name == "" || artwork == "" {
			metadata, err := fetchTuneInMetadata(location)
			if err == nil && metadata != nil {
				if name == "" {
					name = metadata.Name
				}
				if artwork == "" {
					artwork = metadata.Artwork
				}
			}
		}
	}
	source = resolvedSource
	location = resolvedLocation

	clientConfig := GetClientConfig(c)

	if source == "" {
		return fmt.Errorf("source is required (use --source)")
	}

	if location == "" {
		return fmt.Errorf("location is required (use --location)")
	}

	PrintDeviceHeader(fmt.Sprintf("Storing %s content as preset %d", source, slot), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	// Create content item
	contentItem := &models.ContentItem{
		Source:        source,
		Type:          itemType,
		Location:      location,
		SourceAccount: sourceAccount,
		IsPresetable:  true,
		ItemName:      name,
		ContainerArt:  artwork,
	}

	// Set default type if not specified
	if itemType == "" {
		switch source {
		case "SPOTIFY":
			contentItem.Type = "uri"
		case "TUNEIN", "LOCAL_INTERNET_RADIO":
			contentItem.Type = "stationurl"
		default:
			contentItem.Type = ""
		}
	}

	// Show what we're storing
	fmt.Printf("Content to store:\n")
	fmt.Printf("  Name: %s\n", name)
	fmt.Printf("  Source: %s\n", source)
	fmt.Printf("  Location: %s\n", location)
	if sourceAccount != "" {
		fmt.Printf("  Source Account: %s\n", sourceAccount)
	}
	if itemType != "" {
		fmt.Printf("  Type: %s\n", itemType)
	}

	// Store preset
	err = client.StorePreset(slot, contentItem)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to store preset: %v", err))
		return err
	}

	PrintSuccess(fmt.Sprintf("Stored content as preset %d", slot))

	return nil
}

// removePreset handles removing a preset
func removePreset(c *cli.Context) error {
	slot := c.Int("slot")
	clientConfig := GetClientConfig(c)

	PrintDeviceHeader(fmt.Sprintf("Removing preset %d", slot), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	// Check if preset exists first
	presets, err := client.GetPresets()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get presets: %v", err))
		return err
	}

	preset := presets.GetPresetByID(slot)
	if preset == nil || preset.IsEmpty() {
		PrintError(fmt.Sprintf("Preset %d is already empty", slot))
		return fmt.Errorf("preset %d does not exist", slot)
	}

	// Show what we're removing
	fmt.Printf("Removing preset %d:\n", slot)
	fmt.Printf("  Name: %s\n", preset.GetDisplayName())
	fmt.Printf("  Source: %s\n", preset.GetSource())

	// Remove preset
	err = client.RemovePreset(slot)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to remove preset: %v", err))
		return err
	}

	PrintSuccess(fmt.Sprintf("Removed preset %d", slot))

	return nil
}

// selectPresetNew handles selecting a preset (new version that works with subcommands)
func selectPresetNew(c *cli.Context) error {
	slot := c.Int("slot")
	clientConfig := GetClientConfig(c)

	PrintDeviceHeader(fmt.Sprintf("Selecting preset %d", slot), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.SelectPreset(slot)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to select preset: %v", err))
		return err
	}

	PrintSuccess(fmt.Sprintf("Preset %d selected", slot))

	return nil
}

// listPresets handles listing all presets (alias for existing getPresets command)
func listPresets(c *cli.Context) error {
	return getPresets(c)
}
