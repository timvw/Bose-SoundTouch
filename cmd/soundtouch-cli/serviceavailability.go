package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gesellix/bose-soundtouch/pkg/client"
	"github.com/gesellix/bose-soundtouch/pkg/models"
)

// ServiceAvailabilityChecker provides service availability validation for CLI commands
type ServiceAvailabilityChecker struct {
	client                *client.Client
	serviceAvailability   *models.ServiceAvailability
	skipAvailabilityCheck bool
	cached                bool
}

// NewServiceAvailabilityChecker creates a new service availability checker
func NewServiceAvailabilityChecker(client *client.Client) *ServiceAvailabilityChecker {
	skipCheck := os.Getenv("SOUNDTOUCH_SKIP_AVAILABILITY_CHECK") == "true" ||
		os.Getenv("SOUNDTOUCH_SKIP_AVAILABILITY_CHECK") == "1"

	return &ServiceAvailabilityChecker{
		client:                client,
		skipAvailabilityCheck: skipCheck,
		cached:                false,
	}
}

// loadServiceAvailability loads service availability data (cached after first call)
func (sac *ServiceAvailabilityChecker) loadServiceAvailability() {
	if sac.cached {
		return
	}

	if sac.skipAvailabilityCheck {
		// Create a mock availability that allows everything
		sac.serviceAvailability = &models.ServiceAvailability{}
		sac.cached = true
		return
	}

	serviceAvailability, err := sac.client.GetServiceAvailability()
	if err != nil {
		// If availability check fails, warn but don't fail the command
		PrintWarning(fmt.Sprintf("Could not check service availability: %v", err))
		if !sac.skipAvailabilityCheck {
			PrintWarning("Command will proceed without availability validation")
			PrintWarning("Set SOUNDTOUCH_SKIP_AVAILABILITY_CHECK=true to disable these checks")
		}
		// Create empty availability to prevent further errors
		sac.serviceAvailability = &models.ServiceAvailability{}
		sac.cached = true
		return
	}

	sac.serviceAvailability = serviceAvailability
	sac.cached = true
}

// CheckServiceAvailable validates if a service is available and provides user feedback
func (sac *ServiceAvailabilityChecker) CheckServiceAvailable(serviceType models.ServiceType, actionDescription string) bool {
	if sac.skipAvailabilityCheck {
		return true
	}

	sac.loadServiceAvailability()

	// If we couldn't load availability data, allow the operation
	if sac.serviceAvailability == nil || sac.serviceAvailability.Services == nil {
		return true
	}

	if sac.serviceAvailability.IsServiceAvailable(serviceType) {
		return true
	}

	// Service is not available - provide helpful feedback
	serviceName := formatServiceTypeForDisplay(serviceType)
	PrintError(fmt.Sprintf("Cannot %s: %s service is not available", actionDescription, serviceName))

	// Get specific reason if available
	service := sac.serviceAvailability.GetServiceByType(serviceType)
	if service != nil && service.Reason != "" {
		PrintError(fmt.Sprintf("Reason: %s", service.Reason))
	}

	// Provide troubleshooting hints
	sac.provideTroubleshootingHints(serviceType)

	// Suggest alternatives
	sac.suggestAlternatives(serviceType, actionDescription)

	PrintWarning("To bypass this check, set SOUNDTOUCH_SKIP_AVAILABILITY_CHECK=true")
	return false
}

// CheckSourceAvailable validates if a source string corresponds to an available service
func (sac *ServiceAvailabilityChecker) CheckSourceAvailable(source, actionDescription string) bool {
	if sac.skipAvailabilityCheck {
		return true
	}

	serviceType := sourceToServiceType(source)
	if serviceType == "" {
		// Unknown source type, allow it (might be a valid source not in our list)
		return true
	}

	return sac.CheckServiceAvailable(serviceType, actionDescription)
}

// ValidateSpotifyAvailable checks Spotify availability for Spotify-specific operations
func (sac *ServiceAvailabilityChecker) ValidateSpotifyAvailable(actionDescription string) bool {
	return sac.CheckServiceAvailable(models.ServiceTypeSpotify, actionDescription)
}

// ValidateBluetoothAvailable checks Bluetooth availability for Bluetooth operations
func (sac *ServiceAvailabilityChecker) ValidateBluetoothAvailable(actionDescription string) bool {
	return sac.CheckServiceAvailable(models.ServiceTypeBluetooth, actionDescription)
}

// ValidateTuneInAvailable checks TuneIn availability for radio operations
func (sac *ServiceAvailabilityChecker) ValidateTuneInAvailable(actionDescription string) bool {
	return sac.CheckServiceAvailable(models.ServiceTypeTuneIn, actionDescription)
}

// ValidatePandoraAvailable checks Pandora availability for Pandora operations
func (sac *ServiceAvailabilityChecker) ValidatePandoraAvailable(actionDescription string) bool {
	return sac.CheckServiceAvailable(models.ServiceTypePandora, actionDescription)
}

// GetAvailableStreamingServices returns a list of available streaming services for user feedback
func (sac *ServiceAvailabilityChecker) GetAvailableStreamingServices() []string {
	if sac.skipAvailabilityCheck {
		return []string{"All services (availability check disabled)"}
	}

	sac.loadServiceAvailability()

	if sac.serviceAvailability == nil || sac.serviceAvailability.Services == nil {
		return []string{"Unable to determine available services"}
	}

	streamingServices := sac.serviceAvailability.GetStreamingServices()

	var available []string

	for _, service := range streamingServices {
		if service.IsAvailable {
			available = append(available, formatServiceTypeForDisplay(models.ServiceType(service.Type)))
		}
	}

	if len(available) == 0 {
		return []string{"No streaming services currently available"}
	}

	return available
}

// GetAvailableLocalServices returns a list of available local input services
func (sac *ServiceAvailabilityChecker) GetAvailableLocalServices() []string {
	if sac.skipAvailabilityCheck {
		return []string{"All services (availability check disabled)"}
	}

	sac.loadServiceAvailability()

	if sac.serviceAvailability == nil || sac.serviceAvailability.Services == nil {
		return []string{"Unable to determine available services"}
	}

	localServices := sac.serviceAvailability.GetLocalServices()

	var available []string

	for _, service := range localServices {
		if service.IsAvailable {
			available = append(available, formatServiceTypeForDisplay(models.ServiceType(service.Type)))
		}
	}

	if len(available) == 0 {
		return []string{"No local input services currently available"}
	}

	return available
}

// provideTroubleshootingHints provides specific troubleshooting advice based on service type
func (sac *ServiceAvailabilityChecker) provideTroubleshootingHints(serviceType models.ServiceType) {
	switch serviceType {
	case models.ServiceTypeBluetooth:
		PrintWarning("💡 Bluetooth troubleshooting:")
		PrintWarning("   • Check if your device supports Bluetooth audio input")
		PrintWarning("   • Ensure Bluetooth is enabled on the SoundTouch device")
		PrintWarning("   • Try restarting the device")

	case models.ServiceTypeSpotify:
		PrintWarning("💡 Spotify troubleshooting:")
		PrintWarning("   • Ensure you have a Spotify Premium account")
		PrintWarning("   • Check if you're logged in to Spotify on the device")
		PrintWarning("   • Verify your network connection")

	case models.ServiceTypeAirPlay:
		PrintWarning("💡 AirPlay troubleshooting:")
		PrintWarning("   • Ensure your Apple device and SoundTouch are on the same network")
		PrintWarning("   • Check that AirPlay is enabled in device settings")
		PrintWarning("   • Verify network connectivity")

	case models.ServiceTypeAlexa:
		PrintWarning("💡 Alexa troubleshooting:")
		PrintWarning("   • Check if Amazon Alexa is properly configured")
		PrintWarning("   • Ensure the device is connected to your Amazon account")
		PrintWarning("   • Verify internet connectivity")

	case models.ServiceTypeTuneIn:
		PrintWarning("💡 TuneIn troubleshooting:")
		PrintWarning("   • Check internet connectivity")
		PrintWarning("   • Verify the device can access external streaming services")

	case models.ServiceTypePandora:
		PrintWarning("💡 Pandora troubleshooting:")
		PrintWarning("   • Ensure you have a valid Pandora account")
		PrintWarning("   • Check if you're logged in to Pandora on the device")
		PrintWarning("   • Verify internet connectivity")
	}
}

// suggestAlternatives suggests alternative services when the requested one is unavailable
func (sac *ServiceAvailabilityChecker) suggestAlternatives(serviceType models.ServiceType, _ string) {
	if sac.serviceAvailability == nil {
		return
	}

	switch serviceType {
	case models.ServiceTypeSpotify:
		if sac.serviceAvailability.HasTuneIn() {
			PrintWarning("💡 Alternative: TuneIn Radio is available for music streaming")
		}
		if sac.serviceAvailability.HasPandora() {
			PrintWarning("💡 Alternative: Pandora is available for music streaming")
		}

	case models.ServiceTypeBluetooth:
		if sac.serviceAvailability.HasAirPlay() {
			PrintWarning("💡 Alternative: AirPlay is available for wireless audio")
		}
		if sac.serviceAvailability.HasLocalMusic() {
			PrintWarning("💡 Alternative: Local Music Library is available")
		}

	case models.ServiceTypeTuneIn:
		if sac.serviceAvailability.HasSpotify() {
			PrintWarning("💡 Alternative: Spotify is available for music streaming")
		}
		if sac.serviceAvailability.HasPandora() {
			PrintWarning("💡 Alternative: Pandora is available for music streaming")
		}
	}

	// Show all available streaming services as suggestions
	available := sac.GetAvailableStreamingServices()
	if len(available) > 0 && available[0] != "No streaming services currently available" {
		PrintWarning(fmt.Sprintf("💡 Available streaming services: %s", strings.Join(available, ", ")))
	}
}

// sourceToServiceType maps source strings to service types
func sourceToServiceType(source string) models.ServiceType {
	switch strings.ToUpper(source) {
	case "SPOTIFY":
		return models.ServiceTypeSpotify
	case "BLUETOOTH":
		return models.ServiceTypeBluetooth
	case "AIRPLAY":
		return models.ServiceTypeAirPlay
	case "ALEXA":
		return models.ServiceTypeAlexa
	case "AMAZON":
		return models.ServiceTypeAmazon
	case "PANDORA":
		return models.ServiceTypePandora
	case "TUNEIN":
		return models.ServiceTypeTuneIn
	case "DEEZER":
		return models.ServiceTypeDeezer
	case "IHEART", "IHEARTRADIO":
		return models.ServiceTypeIHeart
	case "LOCAL_INTERNET_RADIO":
		return models.ServiceTypeLocalInternetRadio
	case "LOCAL_MUSIC":
		return models.ServiceTypeLocalMusic
	case "BMX":
		return models.ServiceTypeBMX
	case "NOTIFICATION":
		return models.ServiceTypeNotification
	default:
		return ""
	}
}

// formatServiceTypeForDisplay formats service types for user-friendly display
func formatServiceTypeForDisplay(serviceType models.ServiceType) string {
	switch serviceType {
	case models.ServiceTypeSpotify:
		return "Spotify"
	case models.ServiceTypeBluetooth:
		return "Bluetooth"
	case models.ServiceTypeAirPlay:
		return "AirPlay"
	case models.ServiceTypeAlexa:
		return "Amazon Alexa"
	case models.ServiceTypeAmazon:
		return "Amazon Music"
	case models.ServiceTypePandora:
		return "Pandora"
	case models.ServiceTypeTuneIn:
		return "TuneIn Radio"
	case models.ServiceTypeDeezer:
		return "Deezer"
	case models.ServiceTypeIHeart:
		return "iHeartRadio"
	case models.ServiceTypeLocalInternetRadio:
		return "Internet Radio"
	case models.ServiceTypeLocalMusic:
		return "Local Music Library"
	case models.ServiceTypeBMX:
		return "BMX"
	case models.ServiceTypeNotification:
		return "Notifications"
	default:
		return string(serviceType)
	}
}

// PrintServiceAvailabilitySummary prints a summary of available services
func (sac *ServiceAvailabilityChecker) PrintServiceAvailabilitySummary() {
	if sac.skipAvailabilityCheck {
		PrintWarning("Service availability checking is disabled")
		return
	}

	sac.loadServiceAvailability()

	if sac.serviceAvailability == nil || sac.serviceAvailability.Services == nil {
		PrintWarning("Unable to determine service availability")
		return
	}

	fmt.Printf("📊 Service Availability Summary:\n")
	fmt.Printf("   Total services: %d\n", sac.serviceAvailability.GetServiceCount())
	fmt.Printf("   Available: %d\n", sac.serviceAvailability.GetAvailableServiceCount())
	fmt.Printf("   Unavailable: %d\n", sac.serviceAvailability.GetUnavailableServiceCount())

	// Show quick status for popular services
	fmt.Printf("   Popular services:\n")
	popularChecks := []struct {
		check func() bool
		name  string
	}{
		{sac.serviceAvailability.HasSpotify, "Spotify"},
		{sac.serviceAvailability.HasBluetooth, "Bluetooth"},
		{sac.serviceAvailability.HasAirPlay, "AirPlay"},
		{sac.serviceAvailability.HasTuneIn, "TuneIn Radio"},
		{sac.serviceAvailability.HasPandora, "Pandora"},
	}

	for _, check := range popularChecks {
		status := "❌"
		if check.check() {
			status = "✅"
		}
		fmt.Printf("     %s %s\n", status, check.name)
	}

	fmt.Printf("💡 Use 'soundtouch-cli sources list' to see configured sources\n")
}
