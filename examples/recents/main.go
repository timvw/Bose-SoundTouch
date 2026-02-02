// Package main demonstrates recent content functionality for Bose SoundTouch devices.
package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/gesellix/bose-soundtouch/pkg/client"
	"github.com/gesellix/bose-soundtouch/pkg/models"
)

// applyFilters applies source and type filters to the items
func applyFilters(response *models.RecentsResponse, source, itemType string) ([]models.RecentsResponseItem, error) {
	items := response.Items

	// Apply source filter
	if source != "" {
		items = response.GetItemsBySource(strings.ToUpper(source))
		if len(items) == 0 {
			fmt.Printf("📭 No items found for source: %s\n", source)
			fmt.Println("💡 Available sources:", getAvailableSources(response))

			return nil, fmt.Errorf("no items found for source")
		}
	}

	// Apply type filter
	if itemType != "" {
		filteredItems, err := filterItemsByType(items, itemType)
		if err != nil {
			return nil, err
		}

		items = filteredItems

		if len(items) == 0 {
			fmt.Printf("📭 No items found for type: %s\n", itemType)
			return nil, fmt.Errorf("no items found for type")
		}
	}

	return items, nil
}

// filterItemsByType filters items by content type
func filterItemsByType(items []models.RecentsResponseItem, itemType string) ([]models.RecentsResponseItem, error) {
	// Define type predicates
	predicates := map[string]func(*models.RecentsResponseItem) bool{
		"track":      (*models.RecentsResponseItem).IsTrack,
		"tracks":     (*models.RecentsResponseItem).IsTrack,
		"station":    (*models.RecentsResponseItem).IsStation,
		"stations":   (*models.RecentsResponseItem).IsStation,
		"playlist":   (*models.RecentsResponseItem).IsPlaylist,
		"playlists":  (*models.RecentsResponseItem).IsPlaylist,
		"album":      (*models.RecentsResponseItem).IsAlbum,
		"albums":     (*models.RecentsResponseItem).IsAlbum,
		"presetable": (*models.RecentsResponseItem).IsPresetable,
	}

	predicate, exists := predicates[strings.ToLower(itemType)]
	if !exists {
		fmt.Printf("❌ Unknown type filter: %s\n", itemType)
		fmt.Println("💡 Available types: track, station, playlist, album, presetable")

		return nil, fmt.Errorf("unknown type filter")
	}

	var filteredItems []models.RecentsResponseItem

	for _, item := range items {
		if predicate(&item) {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems, nil
}

func main() {
	var (
		host     = flag.String("host", "", "SoundTouch device IP address")
		timeout  = flag.Duration("timeout", 10*time.Second, "Request timeout")
		detailed = flag.Bool("detailed", false, "Show detailed information for each item")
		limit    = flag.Int("limit", 10, "Maximum number of items to display (0 for all)")
		source   = flag.String("source", "", "Filter by source (SPOTIFY, LOCAL_MUSIC, etc.)")
		itemType = flag.String("type", "", "Filter by type (track, station, playlist, presetable)")
		stats    = flag.Bool("stats", false, "Show statistics only")
	)

	flag.Parse()

	if *host == "" {
		log.Fatal("Please provide a SoundTouch device IP address with -host flag")
	}

	// Create client
	config := &client.Config{
		Host:    *host,
		Port:    8090,
		Timeout: *timeout,
	}
	soundTouchClient := client.NewClient(config)

	fmt.Printf("Getting recent items from %s\n", *host)

	// Get recent items
	response, err := soundTouchClient.GetRecents()
	if err != nil {
		log.Fatalf("Failed to get recent items: %v", err)
	}

	if response.IsEmpty() {
		fmt.Println("\n📭 No recent items found")
		fmt.Println("💡 Play some content to populate the recent items list")

		return
	}

	// Show statistics if requested
	if *stats {
		showStatistics(response)
		return
	}

	// Apply filters
	items, err := applyFilters(response, *source, *itemType)
	if err != nil {
		return
	}

	// Apply limit
	if *limit > 0 && *limit < len(items) {
		items = items[:*limit]
	}

	// Display results
	displayResults(response, items, *detailed, *source, *itemType)

	fmt.Println("\nDone!")
}

// sourceCount represents a count for a named category
type sourceCount struct {
	name  string
	count int
}

// printBasicStatistics prints overall statistics
func printBasicStatistics(response *models.RecentsResponse) {
	fmt.Printf("Overall Statistics:\n")
	fmt.Printf("  Total Items: %d\n", response.GetItemCount())

	if !response.IsEmpty() {
		mostRecent := response.GetMostRecent()
		if mostRecent != nil {
			lastPlayTime := time.Unix(mostRecent.GetUTCTime(), 0)
			fmt.Printf("  Last Played: %s\n", lastPlayTime.Format("2006-01-02 15:04:05"))
		}
	}
}

// printSourceStatistics prints statistics by source
func printSourceStatistics(response *models.RecentsResponse) {
	fmt.Printf("\n📍 By Source:\n")

	sourceStats := map[string]int{
		"Spotify":      len(response.GetSpotifyItems()),
		"Pandora":      len(response.GetPandoraItems()),
		"TuneIn":       len(response.GetTuneInItems()),
		"Local Music":  len(response.GetLocalMusicItems()),
		"Stored Music": len(response.GetStoredMusicItems()),
	}

	var sources []sourceCount

	for name, count := range sourceStats {
		if count > 0 {
			sources = append(sources, sourceCount{name, count})
		}
	}

	sort.Slice(sources, func(i, j int) bool {
		return sources[i].count > sources[j].count
	})

	for _, sc := range sources {
		percentage := float64(sc.count) / float64(response.GetItemCount()) * 100
		fmt.Printf("  %-15s %3d items (%5.1f%%)\n", sc.name+":", sc.count, percentage)
	}
}

// printContentTypeStatistics prints statistics by content type
func printContentTypeStatistics(response *models.RecentsResponse) {
	fmt.Printf("\n🎼 By Content Type:\n")

	tracks := len(response.GetTracks())
	stations := len(response.GetStations())
	playlists := len(response.GetPlaylistsAndAlbums())

	typeStats := []sourceCount{
		{"Tracks", tracks},
		{"Stations", stations},
		{"Playlists/Albums", playlists},
	}

	for _, ts := range typeStats {
		if ts.count > 0 {
			percentage := float64(ts.count) / float64(response.GetItemCount()) * 100
			fmt.Printf("  %-15s %3d items (%5.1f%%)\n", ts.name+":", ts.count, percentage)
		}
	}
}

// printSpecialCategoryStatistics prints special category statistics
func printSpecialCategoryStatistics(response *models.RecentsResponse) {
	presetable := len(response.GetPresetableItems())
	if presetable > 0 {
		fmt.Printf("\n⭐ Special Categories:\n")

		percentage := float64(presetable) / float64(response.GetItemCount()) * 100
		fmt.Printf("  %-15s %3d items (%5.1f%%)\n", "Presetable:", presetable, percentage)
	}
}

// printSourceAnalysis prints streaming vs local content analysis
func printSourceAnalysis(response *models.RecentsResponse) {
	streamingCount := 0
	localCount := 0

	for _, item := range response.Items {
		if item.IsStreamingContent() {
			streamingCount++
		} else if item.IsLocalContent() {
			localCount++
		}
	}

	if streamingCount > 0 || localCount > 0 {
		fmt.Printf("\n📡 Source Analysis:\n")

		if streamingCount > 0 {
			percentage := float64(streamingCount) / float64(response.GetItemCount()) * 100
			fmt.Printf("  %-15s %3d items (%5.1f%%)\n", "Streaming:", streamingCount, percentage)
		}

		if localCount > 0 {
			percentage := float64(localCount) / float64(response.GetItemCount()) * 100
			fmt.Printf("  %-15s %3d items (%5.1f%%)\n", "Local:", localCount, percentage)
		}
	}
}

// printTimeAnalysis prints when items were played
func printTimeAnalysis(response *models.RecentsResponse) {
	fmt.Printf("\n🕐 Time Analysis:\n")

	now := time.Now()
	today := 0
	yesterday := 0
	thisWeek := 0
	older := 0

	for _, item := range response.Items {
		if item.GetUTCTime() > 0 {
			playTime := time.Unix(item.GetUTCTime(), 0)
			diff := now.Sub(playTime)

			switch {
			case diff < 24*time.Hour:
				today++
			case diff < 48*time.Hour:
				yesterday++
			case diff < 7*24*time.Hour:
				thisWeek++
			default:
				older++
			}
		}
	}

	if today > 0 {
		fmt.Printf("  %-15s %3d items\n", "Today:", today)
	}

	if yesterday > 0 {
		fmt.Printf("  %-15s %3d items\n", "Yesterday:", yesterday)
	}

	if thisWeek > 0 {
		fmt.Printf("  %-15s %3d items\n", "This Week:", thisWeek)
	}

	if older > 0 {
		fmt.Printf("  %-15s %3d items\n", "Older:", older)
	}
}

func showStatistics(response *models.RecentsResponse) {
	fmt.Printf("\n📊 Recent Items Statistics\n\n")

	printBasicStatistics(response)
	printSourceStatistics(response)
	printContentTypeStatistics(response)
	printSpecialCategoryStatistics(response)
	printSourceAnalysis(response)
	printTimeAnalysis(response)
}

func displayResults(response *models.RecentsResponse, items []models.RecentsResponseItem, detailed bool, sourceFilter, typeFilter string) {
	// Build filter description
	var filters []string
	if sourceFilter != "" {
		filters = append(filters, fmt.Sprintf("source: %s", sourceFilter))
	}

	if typeFilter != "" {
		filters = append(filters, fmt.Sprintf("type: %s", typeFilter))
	}

	filterDesc := ""
	if len(filters) > 0 {
		filterDesc = fmt.Sprintf(" (filtered by %s)", strings.Join(filters, ", "))
	}

	// Display header
	fmt.Printf("\n📊 Recent Items Summary%s:\n", filterDesc)
	fmt.Printf("   Showing: %d items", len(items))

	if len(items) < response.GetItemCount() {
		fmt.Printf(" (of %d total)", response.GetItemCount())
	}

	fmt.Println()

	if len(filters) == 0 {
		// Show source breakdown for unfiltered results
		sources := []string{}
		sourceCounts := map[string]int{
			"Spotify": len(response.GetSpotifyItems()),
			"Local":   len(response.GetLocalMusicItems()) + len(response.GetStoredMusicItems()),
			"TuneIn":  len(response.GetTuneInItems()),
			"Pandora": len(response.GetPandoraItems()),
		}

		for source, count := range sourceCounts {
			if count > 0 {
				sources = append(sources, fmt.Sprintf("%s: %d", source, count))
			}
		}

		if len(sources) > 0 {
			fmt.Printf("   By Source: %s\n", strings.Join(sources, ", "))
		}
	}

	fmt.Printf("\n=== Recent Items ===\n")

	// Display items
	for i, item := range items {
		displayItem(i+1, &item, detailed)
	}

	if len(items) < response.GetItemCount() {
		fmt.Printf("\n💡 Showing %d of %d total items\n", len(items), response.GetItemCount())
		fmt.Printf("   Use -limit 0 to show all items\n")
	}
}

func displayItem(index int, item *models.RecentsResponseItem, detailed bool) {
	// Basic information
	displayName := item.GetDisplayName()
	source := formatSource(item.GetSource())
	contentType := item.GetContentType()

	// Content type icon
	icon := getIcon(item)

	fmt.Printf("%d. %s %s\n", index, icon, displayName)
	fmt.Printf("   Source: %s", source)

	if contentType != "" {
		fmt.Printf(" | Type: %s", contentType)
	}

	fmt.Println()

	// Time information
	if item.GetUTCTime() > 0 {
		playTime := time.Unix(item.GetUTCTime(), 0)
		timeAgo := time.Since(playTime)
		fmt.Printf("   Played: %s", playTime.Format("2006-01-02 15:04:05"))
		fmt.Printf(" (%s ago)\n", formatDuration(timeAgo))
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
			fmt.Printf("   🎨 Has artwork\n")
		}

		location := item.GetLocation()
		if location != "" {
			fmt.Printf("   📍 Location: %s\n", truncateString(location, 60))
		}

		sourceAccount := item.GetSourceAccount()
		if sourceAccount != "" && sourceAccount != item.GetSource() {
			fmt.Printf("   👤 Account: %s\n", truncateString(sourceAccount, 40))
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
			fmt.Printf("   🏷️  Type: %s\n", strings.Join(classifications, ", "))
		}
	}

	fmt.Println()
}

func getIcon(item *models.RecentsResponseItem) string {
	switch {
	case item.IsTrack():
		return "🎵"
	case item.IsStation():
		return "📻"
	case item.IsPlaylist():
		return "📋"
	case item.IsAlbum():
		return "💿"
	case item.IsContainer():
		return "📁"
	default:
		return "🎶"
	}
}

func formatSource(source string) string {
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

func formatDuration(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "< 1 minute"
	case d < time.Hour:
		minutes := int(d.Minutes())
		return fmt.Sprintf("%d minute%s", minutes, pluralize(minutes))
	case d < 24*time.Hour:
		hours := int(d.Hours())
		return fmt.Sprintf("%d hour%s", hours, pluralize(hours))
	default:
		days := int(d.Hours() / 24)
		return fmt.Sprintf("%d day%s", days, pluralize(days))
	}
}

func pluralize(count int) string {
	if count == 1 {
		return ""
	}

	return "s"
}

func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}

	if maxLength <= 3 {
		return "..."
	}

	return s[:maxLength-3] + "..."
}

func getAvailableSources(response *models.RecentsResponse) string {
	sourceMap := make(map[string]bool)

	for _, item := range response.Items {
		if source := item.GetSource(); source != "" {
			sourceMap[source] = true
		}
	}

	var sources []string
	for source := range sourceMap {
		sources = append(sources, source)
	}

	sort.Strings(sources)

	if len(sources) == 0 {
		return "none"
	}

	return strings.Join(sources, ", ")
}
