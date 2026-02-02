// Package main demonstrates music service account management functionality for Bose SoundTouch devices.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gesellix/bose-soundtouch/pkg/client"
	"github.com/gesellix/bose-soundtouch/pkg/models"
)

func main() {
	// Configure the SoundTouch client
	config := &client.Config{
		Host:    "192.168.1.100", // Replace with your device IP
		Port:    8090,
		Timeout: 10 * time.Second,
	}

	// Create client
	soundtouchClient := client.NewClient(config)

	fmt.Printf("🎵 SoundTouch Music Service Account Management Example\n")
	fmt.Printf("Device: %s:%d\n\n", config.Host, config.Port)

	// Example 1: Add a Spotify account using convenience method
	fmt.Println("📱 Adding Spotify Premium account...")

	err := soundtouchClient.AddSpotifyAccount("user@spotify.com", "your_password")
	if err != nil {
		log.Printf("Failed to add Spotify account: %v", err)
	} else {
		fmt.Println("✅ Spotify account added successfully")
	}

	// Example 2: Add a Pandora account
	fmt.Println("\n📻 Adding Pandora account...")

	err = soundtouchClient.AddPandoraAccount("pandora_username", "pandora_password")
	if err != nil {
		log.Printf("Failed to add Pandora account: %v", err)
	} else {
		fmt.Println("✅ Pandora account added successfully")
	}

	// Example 3: Add Amazon Music account
	fmt.Println("\n🛒 Adding Amazon Music account...")

	err = soundtouchClient.AddAmazonMusicAccount("amazon_user", "amazon_password")
	if err != nil {
		log.Printf("Failed to add Amazon Music account: %v", err)
	} else {
		fmt.Println("✅ Amazon Music account added successfully")
	}

	// Example 4: Add a network music library (NAS/UPnP)
	fmt.Println("\n🏠 Adding network music library...")

	nasGUID := "d09708a1-5953-44bc-a413-123456789012/0" // Example UPnP server GUID

	err = soundtouchClient.AddStoredMusicAccount(nasGUID, "My Home Music Server")
	if err != nil {
		log.Printf("Failed to add network music library: %v", err)
	} else {
		fmt.Println("✅ Network music library added successfully")
	}

	// Example 5: Add account using generic method with custom credentials
	fmt.Println("\n🎧 Adding Deezer account using generic method...")

	deezerCredentials := models.NewDeezerCredentials("deezer_user", "deezer_password")

	err = soundtouchClient.SetMusicServiceAccount(deezerCredentials)
	if err != nil {
		log.Printf("Failed to add Deezer account: %v", err)
	} else {
		fmt.Println("✅ Deezer account added successfully")
	}

	// Example 6: Add a custom/unknown service
	fmt.Println("\n🎶 Adding custom music service...")

	customCredentials := models.NewMusicServiceCredentials("TIDAL", "Tidal HiFi", "tidal_user", "tidal_password")

	err = soundtouchClient.SetMusicServiceAccount(customCredentials)
	if err != nil {
		log.Printf("Failed to add custom music service: %v", err)
	} else {
		fmt.Println("✅ Custom music service added successfully")
	}

	// Example 7: List current sources to see added accounts
	fmt.Println("\n📋 Checking available sources...")

	sources, err := soundtouchClient.GetSources()
	if err != nil {
		log.Printf("Failed to get sources: %v", err)
	} else {
		fmt.Printf("Available sources (%d total):\n", len(sources.SourceItem))

		for _, source := range sources.SourceItem {
			status := "🔴 Unavailable"
			if source.Status == models.SourceStatusReady {
				status = "🟢 Ready"
			}

			accountInfo := ""
			if source.SourceAccount != "" && source.SourceAccount != source.Source {
				accountInfo = fmt.Sprintf(" (%s)", source.SourceAccount)
			}

			fmt.Printf("  %s %s%s\n", status, source.GetDisplayName(), accountInfo)
		}
	}

	// Example 8: Remove accounts
	fmt.Println("\n🗑️  Removing accounts...")

	// Remove Spotify account
	err = soundtouchClient.RemoveSpotifyAccount("user@spotify.com")
	if err != nil {
		log.Printf("Failed to remove Spotify account: %v", err)
	} else {
		fmt.Println("✅ Spotify account removed successfully")
	}

	// Remove Deezer account using generic method
	deezerRemovalCredentials := models.NewDeezerCredentials("deezer_user", "")

	err = soundtouchClient.RemoveMusicServiceAccount(deezerRemovalCredentials)
	if err != nil {
		log.Printf("Failed to remove Deezer account: %v", err)
	} else {
		fmt.Println("✅ Deezer account removed successfully")
	}

	// Remove network music library
	err = soundtouchClient.RemoveStoredMusicAccount(nasGUID, "My Home Music Server")
	if err != nil {
		log.Printf("Failed to remove network music library: %v", err)
	} else {
		fmt.Println("✅ Network music library removed successfully")
	}

	fmt.Println("\n🎉 Account management example completed!")
	fmt.Println("\n💡 Tips:")
	fmt.Println("   • Use 'account list' to see which services are configured")
	fmt.Println("   • After adding accounts, use 'source list' to verify availability")
	fmt.Println("   • Network libraries (NAS/UPnP) don't require passwords")
	fmt.Println("   • Some services may need additional authentication via their mobile apps")
	fmt.Println("   • Account credentials are stored securely on the SoundTouch device")
}
