package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gesellix/bose-soundtouch/pkg/models"
	"github.com/urfave/cli/v2"
)

// getRecents handles getting recently played content
func getRecents(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Getting recently played content", clientConfig.Host, clientConfig.Port)

	response, err := client.GetRecents()
	if err != nil {
		return fmt.Errorf("failed to get recent items: %w", err)
	}

	if response.IsEmpty() {
		fmt.Printf("📭 No recent items found\n")
		fmt.Printf("💡 Play some content to populate the recent items list\n")
		return nil
	}

	// Display summary
	fmt.Printf("📊 Recent Items Summary:\n")
	fmt.Printf("   Total Items: %d\n", response.GetItemCount())

	// Show source breakdown
	sources := map[string]int{
		"Spotify":      len(response.GetSpotifyItems()),
		"Local Music":  len(response.GetLocalMusicItems()),
		"Stored Music": len(response.GetStoredMusicItems()),
		"TuneIn":       len(response.GetTuneInItems()),
		"Pandora":      len(response.GetPandoraItems()),
	}

	fmt.Printf("   By Source:\n")
	for source, count := range sources {
		if count > 0 {
			fmt.Printf("     • %s: %d items\n", source, count)
		}
	}

	// Show type breakdown
	tracks := len(response.GetTracks())
	stations := len(response.GetStations())
	playlists := len(response.GetPlaylistsAndAlbums())
	presetable := len(response.GetPresetableItems())

	fmt.Printf("   By Type:\n")
	if tracks > 0 {
		fmt.Printf("     • 🎵 Tracks: %d\n", tracks)
	}
	if stations > 0 {
		fmt.Printf("     • 📻 Stations: %d\n", stations)
	}
	if playlists > 0 {
		fmt.Printf("     • 📋 Playlists/Albums: %d\n", playlists)
	}
	if presetable > 0 {
		fmt.Printf("     • ⭐ Presetable: %d\n", presetable)
	}

	fmt.Printf("\n=== Recent Items ===\n")

	// Display items with details
	maxItems := c.Int("limit")
	if maxItems <= 0 || maxItems > len(response.Items) {
		maxItems = len(response.Items)
	}

	for i, item := range response.Items[:maxItems] {
		printRecentItem(i+1, &item, c.Bool("detailed"))
	}

	if len(response.Items) > maxItems {
		fmt.Printf("\n... and %d more items (use --limit to show more)\n", len(response.Items)-maxItems)
	}

	return nil
}

// getRecentsFiltered handles getting filtered recent content
func getRecentsFiltered(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	source := strings.ToUpper(c.String("source"))
	contentType := strings.ToLower(c.String("type"))

	filterDesc := ""
	if source != "" && contentType != "" {
		filterDesc = fmt.Sprintf(" (filtered by source: %s, type: %s)", source, contentType)
	} else if source != "" {
		filterDesc = fmt.Sprintf(" (filtered by source: %s)", source)
	} else if contentType != "" {
		filterDesc = fmt.Sprintf(" (filtered by type: %s)", contentType)
	}

	PrintDeviceHeader("Getting filtered recent content"+filterDesc, clientConfig.Host, clientConfig.Port)

	response, err := client.GetRecents()
	if err != nil {
		return fmt.Errorf("failed to get recent items: %w", err)
	}

	if response.IsEmpty() {
		fmt.Printf("📭 No recent items found\n")
		return nil
	}

	// Apply filters
	var filteredItems []models.RecentsResponseItem

	if source != "" {
		filteredItems = response.GetItemsBySource(source)
	} else {
		filteredItems = response.Items
	}

	// Apply type filter
	if contentType != "" {
		var typeFiltered []models.RecentsResponseItem
		for _, item := range filteredItems {
			switch contentType {
			case "track", "tracks":
				if item.IsTrack() {
					typeFiltered = append(typeFiltered, item)
				}
			case "station", "stations":
				if item.IsStation() {
					typeFiltered = append(typeFiltered, item)
				}
			case "playlist", "playlists":
				if item.IsPlaylist() {
					typeFiltered = append(typeFiltered, item)
				}
			case "album", "albums":
				if item.IsAlbum() {
					typeFiltered = append(typeFiltered, item)
				}
			case "container", "containers":
				if item.IsContainer() {
					typeFiltered = append(typeFiltered, item)
				}
			case "presetable":
				if item.IsPresetable() {
					typeFiltered = append(typeFiltered, item)
				}
			}
		}
		filteredItems = typeFiltered
	}

	if len(filteredItems) == 0 {
		fmt.Printf("📭 No items match the specified filters\n")
		fmt.Printf("💡 Try different filter criteria or check available content\n")
		return nil
	}

	fmt.Printf("📊 Filtered Results: %d items\n\n", len(filteredItems))

	// Display filtered items
	maxItems := c.Int("limit")
	if maxItems <= 0 || maxItems > len(filteredItems) {
		maxItems = len(filteredItems)
	}

	for i, item := range filteredItems[:maxItems] {
		printRecentItem(i+1, &item, c.Bool("detailed"))
	}

	if len(filteredItems) > maxItems {
		fmt.Printf("\n... and %d more items (use --limit to show more)\n", len(filteredItems)-maxItems)
	}

	return nil
}

// getRecentsMostRecent shows only the most recent item
func getRecentsMostRecent(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Getting most recent item", clientConfig.Host, clientConfig.Port)

	response, err := client.GetRecents()
	if err != nil {
		return fmt.Errorf("failed to get recent items: %w", err)
	}

	mostRecent := response.GetMostRecent()
	if mostRecent == nil {
		fmt.Printf("📭 No recent items found\n")
		return nil
	}

	fmt.Printf("🕒 Most Recent Item:\n\n")
	printRecentItem(1, mostRecent, true)

	return nil
}

// printRecentItem prints details about a recent item
func printRecentItem(index int, item *models.RecentsResponseItem, detailed bool) {
	// Basic information
	displayName := item.GetDisplayName()
	source := item.GetSource()
	contentType := item.GetContentType()

	// Format source display
	sourceDisplay := formatSourceForDisplay(source)

	// Content type icon
	typeIcon := getContentTypeIcon(item)

	fmt.Printf("%d. %s %s\n", index, typeIcon, displayName)
	fmt.Printf("   Source: %s", sourceDisplay)

	if contentType != "" {
		fmt.Printf(" | Type: %s", strings.Title(contentType))
	}
	fmt.Printf("\n")

	// Time information
	if item.GetUTCTime() > 0 {
		playTime := time.Unix(item.GetUTCTime(), 0)
		fmt.Printf("   Played: %s\n", playTime.Format("2006-01-02 15:04:05"))
	}

	// Additional details if requested
	if detailed {
		if item.HasID() {
			fmt.Printf("   ID: %s\n", item.GetID())
		}

		if item.IsPresetable() {
			fmt.Printf("   ⭐ Can be saved as preset\n")
		}

		if item.HasArtwork() {
			fmt.Printf("   🎨 Has artwork: %s\n", truncateString(item.GetArtwork(), 50))
		}

		location := item.GetLocation()
		if location != "" {
			fmt.Printf("   📍 Location: %s\n", truncateString(location, 50))
		}

		sourceAccount := item.GetSourceAccount()
		if sourceAccount != "" && sourceAccount != source {
			fmt.Printf("   👤 Account: %s\n", truncateString(sourceAccount, 30))
		}

		// Content classification
		var classifications []string
		if item.IsStreamingContent() {
			classifications = append(classifications, "Streaming")
		}
		if item.IsLocalContent() {
			classifications = append(classifications, "Local")
		}
		if len(classifications) > 0 {
			fmt.Printf("   🏷️  Classification: %s\n", strings.Join(classifications, ", "))
		}
	}

	fmt.Println()
}

// getContentTypeIcon returns an emoji icon for the content type
func getContentTypeIcon(item *models.RecentsResponseItem) string {
	if item.IsTrack() {
		return "🎵"
	} else if item.IsStation() {
		return "📻"
	} else if item.IsPlaylist() {
		return "📋"
	} else if item.IsAlbum() {
		return "💿"
	} else if item.IsContainer() {
		return "📁"
	}
	return "🎼"
}

// formatSourceForDisplay formats source names for user-friendly display
func formatSourceForDisplay(source string) string {
	switch source {
	case "SPOTIFY":
		return "Spotify"
	case "LOCAL_MUSIC":
		return "Local Music"
	case "STORED_MUSIC":
		return "Stored Music"
	case "TUNEIN":
		return "TuneIn Radio"
	case "PANDORA":
		return "Pandora"
	case "AMAZON":
		return "Amazon Music"
	case "DEEZER":
		return "Deezer"
	case "IHEART":
		return "iHeartRadio"
	case "BLUETOOTH":
		return "Bluetooth"
	case "AUX":
		return "AUX Input"
	case "AIRPLAY":
		return "AirPlay"
	default:
		return source
	}
}

// truncateString truncates a string to the specified length with ellipsis
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	if maxLength <= 3 {
		return "..."
	}
	return s[:maxLength-3] + "..."
}

// recentsStats shows statistics about recent items
func recentsStats(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Getting recent items statistics", clientConfig.Host, clientConfig.Port)

	response, err := client.GetRecents()
	if err != nil {
		return fmt.Errorf("failed to get recent items: %w", err)
	}

	if response.IsEmpty() {
		fmt.Printf("📊 Statistics: No recent items found\n")
		return nil
	}

	fmt.Printf("📊 Recent Items Statistics\n\n")

	// Basic stats
	fmt.Printf("Overall Statistics:\n")
	fmt.Printf("  Total Items: %d\n", response.GetItemCount())

	if !response.IsEmpty() {
		mostRecent := response.GetMostRecent()
		if mostRecent != nil {
			lastPlayTime := time.Unix(mostRecent.GetUTCTime(), 0)
			fmt.Printf("  Last Played: %s\n", lastPlayTime.Format("2006-01-02 15:04:05"))
		}
	}

	// Source breakdown
	fmt.Printf("\nBy Source:\n")
	sourceStats := map[string]int{
		"Spotify":      len(response.GetSpotifyItems()),
		"Pandora":      len(response.GetPandoraItems()),
		"TuneIn":       len(response.GetTuneInItems()),
		"Local Music":  len(response.GetLocalMusicItems()),
		"Stored Music": len(response.GetStoredMusicItems()),
	}

	// Add other sources if they exist
	otherSources := make(map[string]int)
	for _, item := range response.Items {
		source := item.GetSource()
		found := false
		for knownSource := range sourceStats {
			if strings.Contains(strings.ToLower(knownSource), strings.ToLower(source)) ||
				strings.Contains(strings.ToLower(source), strings.ToLower(knownSource)) {
				found = true
				break
			}
		}
		if !found && source != "" {
			otherSources[formatSourceForDisplay(source)]++
		}
	}

	// Merge other sources
	for source, count := range otherSources {
		sourceStats[source] = count
	}

	for source, count := range sourceStats {
		if count > 0 {
			percentage := float64(count) / float64(response.GetItemCount()) * 100
			fmt.Printf("  %-15s %3d items (%5.1f%%)\n", source+":", count, percentage)
		}
	}

	// Content type breakdown
	fmt.Printf("\nBy Content Type:\n")
	tracks := len(response.GetTracks())
	stations := len(response.GetStations())
	playlists := len(response.GetPlaylistsAndAlbums())

	if tracks > 0 {
		percentage := float64(tracks) / float64(response.GetItemCount()) * 100
		fmt.Printf("  %-15s %3d items (%5.1f%%)\n", "Tracks:", tracks, percentage)
	}
	if stations > 0 {
		percentage := float64(stations) / float64(response.GetItemCount()) * 100
		fmt.Printf("  %-15s %3d items (%5.1f%%)\n", "Stations:", stations, percentage)
	}
	if playlists > 0 {
		percentage := float64(playlists) / float64(response.GetItemCount()) * 100
		fmt.Printf("  %-15s %3d items (%5.1f%%)\n", "Playlists/Albums:", playlists, percentage)
	}

	// Special categories
	presetable := len(response.GetPresetableItems())
	if presetable > 0 {
		fmt.Printf("\nSpecial Categories:\n")
		percentage := float64(presetable) / float64(response.GetItemCount()) * 100
		fmt.Printf("  %-15s %3d items (%5.1f%%)\n", "Presetable:", presetable, percentage)
	}

	// Content source analysis
	streamingCount := 0
	localCount := 0
	for _, item := range response.Items {
		if item.IsStreamingContent() {
			streamingCount++
		} else if item.IsLocalContent() {
			localCount++
		}
	}

	fmt.Printf("\nSource Analysis:\n")
	if streamingCount > 0 {
		percentage := float64(streamingCount) / float64(response.GetItemCount()) * 100
		fmt.Printf("  %-15s %3d items (%5.1f%%)\n", "Streaming:", streamingCount, percentage)
	}
	if localCount > 0 {
		percentage := float64(localCount) / float64(response.GetItemCount()) * 100
		fmt.Printf("  %-15s %3d items (%5.1f%%)\n", "Local:", localCount, percentage)
	}

	return nil
}
