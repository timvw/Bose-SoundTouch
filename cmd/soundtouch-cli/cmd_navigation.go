package main

import (
	"fmt"

	"github.com/gesellix/bose-soundtouch/pkg/models"
	"github.com/urfave/cli/v2"
)

// browseContent handles browsing content sources
func browseContent(c *cli.Context) error {
	source := c.String("source")
	sourceAccount := c.String("source-account")
	startItem := c.Int("start")
	numItems := c.Int("limit")

	clientConfig := GetClientConfig(c)
	PrintDeviceHeader(fmt.Sprintf("Browsing %s content", source), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	response, err := client.Navigate(source, sourceAccount, startItem, numItems)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to browse content: %v", err))
		return err
	}

	printNavigationResults(response, "Content")
	return nil
}

// browseWithMenu handles browsing with menu navigation
func browseWithMenu(c *cli.Context) error {
	source := c.String("source")
	sourceAccount := c.String("source-account")
	menu := c.String("menu")
	sort := c.String("sort")
	startItem := c.Int("start")
	numItems := c.Int("limit")

	clientConfig := GetClientConfig(c)
	PrintDeviceHeader(fmt.Sprintf("Browsing %s menu: %s", source, menu), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	response, err := client.NavigateWithMenu(source, sourceAccount, menu, sort, startItem, numItems)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to browse menu: %v", err))
		return err
	}

	printNavigationResults(response, "Menu Items")
	return nil
}

// browseContainer handles browsing into containers/directories
func browseContainer(c *cli.Context) error {
	source := c.String("source")
	sourceAccount := c.String("source-account")
	location := c.String("location")
	itemType := c.String("type")
	startItem := c.Int("start")
	numItems := c.Int("limit")

	clientConfig := GetClientConfig(c)
	PrintDeviceHeader(fmt.Sprintf("Browsing %s container: %s", source, location), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	// Create container content item
	containerItem := &models.ContentItem{
		Source:   source,
		Location: location,
		Type:     itemType,
	}

	response, err := client.NavigateContainer(source, sourceAccount, startItem, numItems, containerItem)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to browse container: %v", err))
		return err
	}

	printNavigationResults(response, "Container Contents")
	return nil
}

// browseTuneIn handles browsing TuneIn content
func browseTuneIn(c *cli.Context) error {
	sourceAccount := c.String("source-account")
	startItem := c.Int("start")
	numItems := c.Int("limit")

	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Browsing TuneIn stations", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	response, err := client.GetTuneInStations(sourceAccount)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get TuneIn stations: %v", err))
		return err
	}

	// Apply pagination if different from defaults
	if startItem != 1 || numItems != 100 {
		response, err = client.Navigate("TUNEIN", sourceAccount, startItem, numItems)
		if err != nil {
			PrintError(fmt.Sprintf("Failed to browse TuneIn with pagination: %v", err))
			return err
		}
	}

	printNavigationResults(response, "TuneIn Stations")
	return nil
}

// browsePandora handles browsing Pandora content
func browsePandora(c *cli.Context) error {
	sourceAccount := c.String("source-account")

	if sourceAccount == "" {
		PrintError("Pandora source account is required")
		return fmt.Errorf("source account required for Pandora")
	}

	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Browsing Pandora stations", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	response, err := client.GetPandoraStations(sourceAccount)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get Pandora stations: %v", err))
		return err
	}

	printNavigationResults(response, "Pandora Stations")
	return nil
}

// browseStoredMusic handles browsing local/stored music
func browseStoredMusic(c *cli.Context) error {
	sourceAccount := c.String("source-account")

	if sourceAccount == "" {
		PrintError("Source account (device ID) is required for stored music")
		return fmt.Errorf("source account required for stored music")
	}

	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Browsing stored music library", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	response, err := client.GetStoredMusicLibrary(sourceAccount)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get stored music library: %v", err))
		return err
	}

	printNavigationResults(response, "Stored Music Library")
	return nil
}

// printNavigationResults formats and displays navigation results
func printNavigationResults(response *models.NavigateResponse, title string) {
	fmt.Printf("%s:\n", title)

	if response.TotalItems == 0 {
		fmt.Printf("  No items found\n")
		return
	}

	fmt.Printf("  Total items: %d\n", response.TotalItems)

	if len(response.Items) == 0 {
		fmt.Printf("  No items in current page\n")
		return
	}

	fmt.Printf("  Items:\n")
	for i, item := range response.Items {
		fmt.Printf("    %d. %s\n", i+1, item.GetDisplayName())

		// Show type and source for identification
		if item.ContentItem != nil {
			if item.ContentItem.Source != "" && item.ContentItem.Source != response.Source {
				fmt.Printf("       Source: %s\n", item.ContentItem.Source)
			}
			if item.Type != "" {
				fmt.Printf("       Type: %s\n", item.Type)
			}
			if item.ContentItem.Location != "" && len(item.ContentItem.Location) < 100 {
				fmt.Printf("       Location: %s\n", item.ContentItem.Location)
			}
		}

		// Show additional metadata
		if item.ArtistName != "" {
			fmt.Printf("       Artist: %s\n", item.ArtistName)
		}
		if item.AlbumName != "" {
			fmt.Printf("       Album: %s\n", item.AlbumName)
		}

		// Show if it's a container that can be browsed further
		if item.IsDirectory() {
			fmt.Printf("       📁 Directory (can browse into)\n")
		} else if item.IsPlayable() {
			fmt.Printf("       ▶️  Playable content\n")
		}

		fmt.Println()
	}

	// Show navigation hints
	directories := response.GetDirectories()
	if len(directories) > 0 {
		fmt.Printf("  💡 To browse into a directory, use: browse container --location <location> --type <type>\n")
	}

	playableItems := response.GetPlayableItems()
	if len(playableItems) > 0 {
		fmt.Printf("  💡 Found %d playable items\n", len(playableItems))
	}
}
