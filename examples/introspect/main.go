// Package main demonstrates introspect functionality for Bose SoundTouch devices.
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gesellix/bose-soundtouch/pkg/client"
	"github.com/gesellix/bose-soundtouch/pkg/models"
)

// displayBasicInfo prints basic service information
func displayBasicInfo(source string, response *models.IntrospectResponse) {
	fmt.Printf("\n=== %s Service Introspect Data ===\n", source)
	fmt.Printf("State: %s\n", response.State)

	if response.HasUser() {
		fmt.Printf("User: %s\n", response.User)
	}

	fmt.Printf("Currently Playing: %t\n", response.IsPlaying)

	if response.HasCurrentContent() {
		fmt.Printf("Current Content: %s\n", response.CurrentURI)
	}

	fmt.Printf("Shuffle Mode: %s\n", response.ShuffleMode)

	if response.HasSubscription() {
		fmt.Printf("Subscription Type: %s\n", response.SubscriptionType)
	}
}

// displayServiceState prints service state information
func displayServiceState(response *models.IntrospectResponse) {
	fmt.Printf("\n=== Service State ===\n")

	if response.IsActive() {
		fmt.Println("✅ Service is ACTIVE")
	} else if response.IsInactive() {
		fmt.Println("❌ Service is INACTIVE")
	}
}

// displayCapabilities prints service capabilities
func displayCapabilities(response *models.IntrospectResponse) {
	fmt.Printf("\n=== Service Capabilities ===\n")

	if response.SupportsSkipPrevious() {
		fmt.Println("✅ Skip Previous supported")
	} else {
		fmt.Println("❌ Skip Previous not supported")
	}

	if response.SupportsSeek() {
		fmt.Println("✅ Seek supported")
	} else {
		fmt.Println("❌ Seek not supported")
	}

	if response.SupportsResume() {
		fmt.Println("✅ Resume supported")
	} else {
		fmt.Println("❌ Resume not supported")
	}

	if response.CollectsData() {
		fmt.Println("📊 Data collection enabled")
	} else {
		fmt.Println("🚫 Data collection disabled")
	}
}

// displayHistoryInfo prints content history information
func displayHistoryInfo(response *models.IntrospectResponse) {
	historySize := response.GetMaxHistorySize()
	if historySize > 0 {
		fmt.Printf("\n=== Content History ===\n")
		fmt.Printf("Max History Size: %d items\n", historySize)
	}
}

// displayTechnicalDetails prints technical service details
func displayTechnicalDetails(response *models.IntrospectResponse) {
	if response.TokenLastChangedTimeSeconds > 0 {
		fmt.Printf("\n=== Technical Details ===\n")
		fmt.Printf("Token Last Changed: %d seconds\n", response.TokenLastChangedTimeSeconds)

		if response.TokenLastChangedTimeMicroseconds > 0 {
			fmt.Printf("Token Microseconds: %d\n", response.TokenLastChangedTimeMicroseconds)
		}

		fmt.Printf("Play Status State: %s\n", response.PlayStatusState)
		fmt.Printf("Received Playback Request: %t\n", response.ReceivedPlaybackRequest)
	}
}

// displayServiceAvailability shows service availability for comparison
func displayServiceAvailability(soundTouchClient *client.Client, source string) {
	fmt.Printf("\n=== Service Availability Check ===\n")

	availability, err := soundTouchClient.GetServiceAvailability()
	if err != nil {
		fmt.Printf("Could not check service availability: %v\n", err)
		return
	}

	switch source {
	case "SPOTIFY":
		if availability.HasSpotify() {
			fmt.Println("✅ Spotify is available on this device")
		} else {
			fmt.Println("❌ Spotify is not available on this device")
		}
	case "PANDORA":
		if availability.HasPandora() {
			fmt.Println("✅ Pandora is available on this device")
		} else {
			fmt.Println("❌ Pandora is not available on this device")
		}
	case "TUNEIN":
		if availability.HasTuneIn() {
			fmt.Println("✅ TuneIn is available on this device")
		} else {
			fmt.Println("❌ TuneIn is not available on this device")
		}
	default:
		fmt.Printf("Service availability check not implemented for %s\n", source)
	}
}

func main() {
	var (
		host          = flag.String("host", "", "SoundTouch device IP address")
		source        = flag.String("source", "SPOTIFY", "Music service source (SPOTIFY, PANDORA, TUNEIN)")
		sourceAccount = flag.String("account", "", "Source account name (optional)")
		timeout       = flag.Duration("timeout", 10*time.Second, "Request timeout")
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

	fmt.Printf("Getting introspect data for %s", *source)

	if *sourceAccount != "" {
		fmt.Printf(" (account: %s)", *sourceAccount)
	}

	fmt.Println()

	// Get introspect data
	response, err := soundTouchClient.Introspect(*source, *sourceAccount)
	if err != nil {
		log.Fatalf("Failed to get introspect data: %v", err)
	}

	// Display all information using helper functions
	displayBasicInfo(*source, response)
	displayServiceState(response)
	displayCapabilities(response)
	displayHistoryInfo(response)
	displayTechnicalDetails(response)
	displayServiceAvailability(soundTouchClient, *source)

	fmt.Println("\nDone!")
}
