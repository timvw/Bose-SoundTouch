package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gesellix/bose-soundtouch/pkg/client"
	"github.com/gesellix/bose-soundtouch/pkg/models"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Get device IP from command line
	deviceIP := os.Args[1]

	// Create client
	config := &client.Config{
		Host:    deviceIP,
		Port:    8090,
		Timeout: 10 * time.Second,
	}
	c := client.NewClient(config)

	fmt.Printf("🎵 SoundTouch Navigation & Station Management Demo\n")
	fmt.Printf("📱 Device: %s:%d\n\n", config.Host, config.Port)

	// Demonstrate navigation and station management
	if err := demonstrateNavigationAndStations(c); err != nil {
		log.Fatalf("Demo failed: %v", err)
	}

	fmt.Println("\n✅ Navigation and station management demo completed!")
}

func demonstrateNavigationAndStations(c *client.Client) error {
	// 1. Browse TuneIn content
	fmt.Println("📻 Step 1: Browsing TuneIn stations...")
	if err := browseTuneInStations(c); err != nil {
		return fmt.Errorf("failed to browse TuneIn: %w", err)
	}

	// 2. Search for specific content
	fmt.Println("\n🔍 Step 2: Searching for jazz stations...")
	searchResults, err := searchForJazzStations(c)
	if err != nil {
		return fmt.Errorf("failed to search stations: %w", err)
	}

	// 3. Add and play a station
	fmt.Println("\n➕ Step 3: Adding and playing a station...")
	if err := addAndPlayStation(c, searchResults); err != nil {
		fmt.Printf("⚠️  Could not add station: %v\n", err)
		// Continue with demo even if this fails
	}

	// 4. Demonstrate Pandora search (if account available)
	fmt.Println("\n🎵 Step 4: Demonstrating Pandora search...")
	if err := demonstratePandoraSearch(c); err != nil {
		fmt.Printf("⚠️  Pandora search not available: %v\n", err)
		// Continue with demo
	}

	// 5. Browse stored music (if available)
	fmt.Println("\n💿 Step 5: Browsing stored music...")
	if err := browseStoredMusic(c); err != nil {
		fmt.Printf("⚠️  Stored music not available: %v\n", err)
		// Continue with demo
	}

	// 6. Search Spotify content (if account available)
	fmt.Println("\n🎧 Step 6: Demonstrating Spotify search...")
	if err := demonstrateSpotifySearch(c); err != nil {
		fmt.Printf("⚠️  Spotify search not available: %v\n", err)
		// Continue with demo
	}

	return nil
}

func browseTuneInStations(c *client.Client) error {
	fmt.Printf("  📡 Getting TuneIn stations (first 10)...\n")

	response, err := c.Navigate("TUNEIN", "", 1, 10)
	if err != nil {
		return err
	}

	fmt.Printf("  📻 Found %d total TuneIn stations\n", response.TotalItems)

	if len(response.Items) > 0 {
		fmt.Printf("  🎵 Sample stations:\n")
		for i, item := range response.Items[:min(5, len(response.Items))] {
			fmt.Printf("    %d. %s\n", i+1, item.GetDisplayName())
			if item.IsPlayable() {
				fmt.Printf("       ▶️  Playable\n")
			} else if item.IsDirectory() {
				fmt.Printf("       📁 Directory\n")
			}
		}
	}

	return nil
}

func searchForJazzStations(c *client.Client) (*models.SearchStationResponse, error) {
	fmt.Printf("  🎷 Searching TuneIn for 'jazz'...\n")

	searchResults, err := c.SearchTuneInStations("jazz")
	if err != nil {
		return nil, err
	}

	fmt.Printf("  📊 Search results: %d total\n", searchResults.GetResultCount())

	songs := searchResults.GetSongs()
	artists := searchResults.GetArtists()
	stations := searchResults.GetStations()

	if len(songs) > 0 {
		fmt.Printf("  🎵 Songs (%d): %s\n", len(songs), songs[0].GetDisplayName())
	}
	if len(artists) > 0 {
		fmt.Printf("  🎤 Artists (%d): %s\n", len(artists), artists[0].GetDisplayName())
	}
	if len(stations) > 0 {
		fmt.Printf("  📻 Stations (%d):\n", len(stations))
		for i, station := range stations[:min(3, len(stations))] {
			fmt.Printf("    %d. %s (Token: %s)\n", i+1, station.GetDisplayName(), station.Token)
		}
	}

	return searchResults, nil
}

func addAndPlayStation(c *client.Client, searchResults *models.SearchStationResponse) error {
	stations := searchResults.GetStations()
	if len(stations) == 0 {
		return fmt.Errorf("no stations found to add")
	}

	// Use the first station from search results
	station := stations[0]
	stationName := station.GetDisplayName()

	fmt.Printf("  ➕ Adding station: %s\n", stationName)
	fmt.Printf("  🎯 Token: %s\n", station.Token)

	err := c.AddStation("TUNEIN", station.SourceAccount, station.Token, stationName)
	if err != nil {
		return err
	}

	fmt.Printf("  ✅ Successfully added and started playing: %s\n", stationName)

	// Wait a moment and show what's playing
	time.Sleep(2 * time.Second)
	fmt.Println("  🎵 Checking what's now playing...")

	nowPlaying, err := c.GetNowPlaying()
	if err != nil {
		fmt.Printf("  ⚠️  Could not get now playing: %v\n", err)
		return nil
	}

	if !nowPlaying.IsEmpty() {
		fmt.Printf("      Now Playing: %s\n", nowPlaying.Track)
		fmt.Printf("      Source: %s\n", nowPlaying.Source)
	}

	return nil
}

func demonstratePandoraSearch(c *client.Client) error {
	// Note: This would require a valid Pandora account
	// For demo purposes, we'll show how it would work
	fmt.Printf("  🎵 Pandora search requires a valid source account\n")
	fmt.Printf("  💡 Example usage:\n")
	fmt.Printf("     searchResults, err := client.SearchPandoraStations(\"your_pandora_account\", \"rock\")\n")
	fmt.Printf("     if err == nil {\n")
	fmt.Printf("         // Process Pandora search results\n")
	fmt.Printf("         stations := searchResults.GetStations()\n")
	fmt.Printf("     }\n")

	return nil
}

func browseStoredMusic(c *client.Client) error {
	// Note: This would require a valid device ID for stored music
	fmt.Printf("  💿 Stored music browsing requires device ID\n")
	fmt.Printf("  💡 Example usage:\n")
	fmt.Printf("     musicLibrary, err := client.GetStoredMusicLibrary(\"device_12345\")\n")
	fmt.Printf("     if err == nil {\n")
	fmt.Printf("         // Browse local music library\n")
	fmt.Printf("         directories := musicLibrary.GetDirectories()\n")
	fmt.Printf("         tracks := musicLibrary.GetTracks()\n")
	fmt.Printf("     }\n")

	return nil
}

func demonstrateSpotifySearch(c *client.Client) error {
	// Note: This would require a valid Spotify account
	fmt.Printf("  🎧 Spotify search requires a valid source account\n")
	fmt.Printf("  💡 Example usage:\n")
	fmt.Printf("     searchResults, err := client.SearchSpotifyContent(\"spotify_username\", \"workout\")\n")
	fmt.Printf("     if err == nil {\n")
	fmt.Printf("         // Process Spotify search results\n")
	fmt.Printf("         songs := searchResults.GetSongs()\n")
	fmt.Printf("         artists := searchResults.GetArtists()\n")
	fmt.Printf("     }\n")

	return nil
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func printUsage() {
	fmt.Println("🎵 SoundTouch Navigation & Station Management Demo")
	fmt.Println()
	fmt.Println("This example demonstrates content navigation and station management:")
	fmt.Println("• Browse TuneIn stations")
	fmt.Println("• Search for content across different sources")
	fmt.Println("• Add stations and play them immediately")
	fmt.Println("• Show how to work with Pandora, Spotify, and stored music")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Printf("  %s <device_ip>\n", os.Args[0])
	fmt.Println()
	fmt.Println("Example:")
	fmt.Printf("  %s 192.168.1.100\n", os.Args[0])
	fmt.Println()
	fmt.Println("Prerequisites:")
	fmt.Println("• SoundTouch device on your network")
	fmt.Println("• Device IP address")
	fmt.Println("• Device powered on and connected")
	fmt.Println()
	fmt.Println("CLI Equivalent Commands:")
	fmt.Println("• Browse: soundtouch-cli --host 192.168.1.100 browse tunein")
	fmt.Println("• Search: soundtouch-cli --host 192.168.1.100 station search-tunein --query jazz")
	fmt.Println("• Add:    soundtouch-cli --host 192.168.1.100 station add --source TUNEIN --token <token> --name <name>")
}
