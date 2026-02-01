package models

import "encoding/xml"

// SupportedURLsResponse represents the response from the /supportedURLs endpoint
type SupportedURLsResponse struct {
	XMLName  xml.Name `xml:"supportedURLs"`
	DeviceID string   `xml:"deviceID,attr"`
	URLs     []URL    `xml:"URL"`
}

// URL represents a single supported URL endpoint
type URL struct {
	Location string `xml:"location,attr"`
}

// GetURLs returns a slice of all supported URL locations
func (s *SupportedURLsResponse) GetURLs() []string {
	urls := make([]string, len(s.URLs))
	for i, url := range s.URLs {
		urls[i] = url.Location
	}
	return urls
}

// HasURL checks if a specific URL location is supported
func (s *SupportedURLsResponse) HasURL(location string) bool {
	for _, url := range s.URLs {
		if url.Location == location {
			return true
		}
	}
	return false
}

// GetURLCount returns the total number of supported URLs
func (s *SupportedURLsResponse) GetURLCount() int {
	return len(s.URLs)
}

// GetCoreURLs returns URLs for core device functionality
func (s *SupportedURLsResponse) GetCoreURLs() []string {
	coreEndpoints := []string{
		"/info", "/capabilities", "/name", "/sources", "/volume",
		"/bass", "/balance", "/presets", "/nowPlaying", "/clock",
		"/key", "/powerManagement",
	}

	var available []string
	for _, endpoint := range coreEndpoints {
		if s.HasURL(endpoint) {
			available = append(available, endpoint)
		}
	}
	return available
}

// GetStreamingURLs returns URLs for streaming service functionality
func (s *SupportedURLsResponse) GetStreamingURLs() []string {
	streamingEndpoints := []string{
		"/navigate", "/search", "/addStation", "/removeStation",
		"/serviceAvailability", "/sources", "/select",
	}

	var available []string
	for _, endpoint := range streamingEndpoints {
		if s.HasURL(endpoint) {
			available = append(available, endpoint)
		}
	}
	return available
}

// GetAdvancedURLs returns URLs for advanced audio and system functionality
func (s *SupportedURLsResponse) GetAdvancedURLs() []string {
	advancedEndpoints := []string{
		"/audiodspcontrols", "/audioproducttonecontrols",
		"/audioproductlevelcontrols", "/videoSyncAudioDelay",
		"/group", "/setZone", "/addZoneSlave", "/removeZoneSlave",
		"/bassCapabilities", "/recents", "/trackInfo",
	}

	var available []string
	for _, endpoint := range advancedEndpoints {
		if s.HasURL(endpoint) {
			available = append(available, endpoint)
		}
	}
	return available
}

// GetNetworkURLs returns URLs for network and connectivity functionality
func (s *SupportedURLsResponse) GetNetworkURLs() []string {
	networkEndpoints := []string{
		"/networkInfo", "/netStats", "/wifiProfile", "/bluetoothInfo",
		"/airplay", "/wirelessProfile",
	}

	var available []string
	for _, endpoint := range networkEndpoints {
		if s.HasURL(endpoint) {
			available = append(available, endpoint)
		}
	}
	return available
}

// HasCorePlaybackSupport checks if device supports basic playback functionality
func (s *SupportedURLsResponse) HasCorePlaybackSupport() bool {
	required := []string{"/nowPlaying", "/key", "/volume"}
	for _, endpoint := range required {
		if !s.HasURL(endpoint) {
			return false
		}
	}
	return true
}

// HasPresetSupport checks if device supports preset functionality
func (s *SupportedURLsResponse) HasPresetSupport() bool {
	return s.HasURL("/presets")
}

// HasMultiroomSupport checks if device supports multiroom/zone functionality
func (s *SupportedURLsResponse) HasMultiroomSupport() bool {
	return s.HasURL("/setZone") || s.HasURL("/addZoneSlave")
}

// HasAdvancedAudioSupport checks if device supports advanced audio controls
func (s *SupportedURLsResponse) HasAdvancedAudioSupport() bool {
	return s.HasURL("/audiodspcontrols") || s.HasURL("/audioproducttonecontrols")
}

// HasStreamingSupport checks if device supports streaming service navigation
func (s *SupportedURLsResponse) HasStreamingSupport() bool {
	return s.HasURL("/navigate") || s.HasURL("/search")
}

// GetUnsupportedURLs returns a list of common URLs that this device doesn't support
func (s *SupportedURLsResponse) GetUnsupportedURLs(checkList []string) []string {
	var unsupported []string
	for _, endpoint := range checkList {
		if !s.HasURL(endpoint) {
			unsupported = append(unsupported, endpoint)
		}
	}
	return unsupported
}

// EndpointFeature represents a feature that maps to one or more endpoints
type EndpointFeature struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Endpoints   []string `json:"endpoints"`
	Category    string   `json:"category"`
	Essential   bool     `json:"essential"`   // Required for basic operation
	CLICommand  string   `json:"cli_command"` // Corresponding CLI command
}

// GetEndpointFeatureMap returns a comprehensive mapping of endpoints to implemented features
func GetEndpointFeatureMap() []EndpointFeature {
	return []EndpointFeature{
		// Core Device Information
		{
			Name:        "Device Information",
			Description: "Basic device details, name, and identification",
			Endpoints:   []string{"/info", "/name"},
			Category:    "Core",
			Essential:   true,
			CLICommand:  "info get, name get/set",
		},
		{
			Name:        "Device Capabilities",
			Description: "Supported device features and endpoints discovery",
			Endpoints:   []string{"/capabilities", "/supportedURLs"},
			Category:    "Core",
			Essential:   true,
			CLICommand:  "capabilities, supported-urls",
		},
		{
			Name:        "Network Information",
			Description: "Network configuration and connectivity status",
			Endpoints:   []string{"/networkInfo", "/netStats", "/wifiProfile"},
			Category:    "Network",
			Essential:   false,
			CLICommand:  "network info",
		},

		// Audio Control
		{
			Name:        "Volume Control",
			Description: "Audio volume management and adjustment",
			Endpoints:   []string{"/volume"},
			Category:    "Audio",
			Essential:   true,
			CLICommand:  "volume get/set/up/down",
		},
		{
			Name:        "Bass Control",
			Description: "Bass level adjustment and capabilities",
			Endpoints:   []string{"/bass", "/bassCapabilities"},
			Category:    "Audio",
			Essential:   false,
			CLICommand:  "bass get/set/up/down",
		},
		{
			Name:        "Balance Control",
			Description: "Left/right audio balance adjustment",
			Endpoints:   []string{"/balance"},
			Category:    "Audio",
			Essential:   false,
			CLICommand:  "balance get/set/left/right",
		},
		{
			Name:        "Advanced Audio Controls",
			Description: "DSP controls, tone controls, and audio processing",
			Endpoints:   []string{"/audiodspcontrols", "/audioproducttonecontrols", "/audioproductlevelcontrols"},
			Category:    "Audio",
			Essential:   false,
			CLICommand:  "audio dsp/tone/level",
		},

		// Playback Control
		{
			Name:        "Playback Control",
			Description: "Play, pause, stop, and track navigation",
			Endpoints:   []string{"/key", "/nowPlaying"},
			Category:    "Playback",
			Essential:   true,
			CLICommand:  "play start/stop/pause, key send",
		},
		{
			Name:        "Track Information",
			Description: "Currently playing track details and metadata",
			Endpoints:   []string{"/trackInfo", "/recents"},
			Category:    "Playback",
			Essential:   false,
			CLICommand:  "play now, track get",
		},

		// Source Management
		{
			Name:        "Audio Sources",
			Description: "Available audio sources and source selection",
			Endpoints:   []string{"/sources", "/select"},
			Category:    "Sources",
			Essential:   true,
			CLICommand:  "source list/select/spotify/bluetooth/aux",
		},
		{
			Name:        "Service Availability",
			Description: "Streaming service availability and status",
			Endpoints:   []string{"/serviceAvailability"},
			Category:    "Sources",
			Essential:   false,
			CLICommand:  "source availability/compare",
		},

		// Content Navigation
		{
			Name:        "Content Navigation",
			Description: "Browse music libraries and streaming services",
			Endpoints:   []string{"/navigate", "/search"},
			Category:    "Content",
			Essential:   false,
			CLICommand:  "browse content/tunein/pandora/spotify",
		},
		{
			Name:        "Station Management",
			Description: "Add, remove, and manage radio stations",
			Endpoints:   []string{"/addStation", "/removeStation"},
			Category:    "Content",
			Essential:   false,
			CLICommand:  "station add/remove/search/list",
		},

		// Presets
		{
			Name:        "Preset Management",
			Description: "Store and recall favorite content as presets",
			Endpoints:   []string{"/presets"},
			Category:    "Presets",
			Essential:   false,
			CLICommand:  "preset list/set/select/remove",
		},

		// Multiroom/Zone
		{
			Name:        "Multiroom Zones",
			Description: "Create and manage speaker groups",
			Endpoints:   []string{"/setZone", "/getZone", "/addZoneSlave", "/removeZoneSlave"},
			Category:    "Multiroom",
			Essential:   false,
			CLICommand:  "zone create/add/remove/list",
		},

		// System Features
		{
			Name:        "Clock and Time",
			Description: "Device clock settings and time display",
			Endpoints:   []string{"/clock"},
			Category:    "System",
			Essential:   false,
			CLICommand:  "clock get/set",
		},
		{
			Name:        "Power Management",
			Description: "Device power state and standby control",
			Endpoints:   []string{"/powerManagement"},
			Category:    "System",
			Essential:   false,
			CLICommand:  "key power",
		},
		{
			Name:        "Bluetooth Connectivity",
			Description: "Bluetooth pairing and device management",
			Endpoints:   []string{"/bluetoothInfo", "/bluetoothPair"},
			Category:    "Network",
			Essential:   false,
			CLICommand:  "source bluetooth",
		},
		{
			Name:        "AirPlay Support",
			Description: "Apple AirPlay streaming capability",
			Endpoints:   []string{"/airplay"},
			Category:    "Network",
			Essential:   false,
			CLICommand:  "source select --source AIRPLAY",
		},
	}
}

// GetFeaturesByCategory returns features grouped by category
func (s *SupportedURLsResponse) GetFeaturesByCategory() map[string][]EndpointFeature {
	features := GetEndpointFeatureMap()
	result := make(map[string][]EndpointFeature)

	for _, feature := range features {
		// Check if device supports this feature (any of its endpoints)
		supported := false
		for _, endpoint := range feature.Endpoints {
			if s.HasURL(endpoint) {
				supported = true
				break
			}
		}

		if supported {
			result[feature.Category] = append(result[feature.Category], feature)
		}
	}

	return result
}

// GetSupportedFeatures returns all features supported by this device
func (s *SupportedURLsResponse) GetSupportedFeatures() []EndpointFeature {
	features := GetEndpointFeatureMap()

	var supported []EndpointFeature

	for _, feature := range features {
		// Check if device supports this feature (any of its endpoints)
		hasSupport := false
		for _, endpoint := range feature.Endpoints {
			if s.HasURL(endpoint) {
				hasSupport = true
				break
			}
		}

		if hasSupport {
			supported = append(supported, feature)
		}
	}

	return supported
}

// GetUnsupportedFeatures returns features not supported by this device
func (s *SupportedURLsResponse) GetUnsupportedFeatures() []EndpointFeature {
	features := GetEndpointFeatureMap()

	var unsupported []EndpointFeature

	for _, feature := range features {
		// Check if device supports this feature (any of its endpoints)
		hasSupport := false
		for _, endpoint := range feature.Endpoints {
			if s.HasURL(endpoint) {
				hasSupport = true
				break
			}
		}

		if !hasSupport {
			unsupported = append(unsupported, feature)
		}
	}

	return unsupported
}

// GetPartiallyImplementedFeatures returns features where only some endpoints are supported
func (s *SupportedURLsResponse) GetPartiallyImplementedFeatures() []EndpointFeature {
	features := GetEndpointFeatureMap()

	var partial []EndpointFeature

	for _, feature := range features {
		if len(feature.Endpoints) <= 1 {
			continue // Skip single-endpoint features
		}

		supportedCount := 0
		for _, endpoint := range feature.Endpoints {
			if s.HasURL(endpoint) {
				supportedCount++
			}
		}

		// Partially implemented if some but not all endpoints are supported
		if supportedCount > 0 && supportedCount < len(feature.Endpoints) {
			partial = append(partial, feature)
		}
	}

	return partial
}

// GetMissingEssentialFeatures returns essential features that are not supported
func (s *SupportedURLsResponse) GetMissingEssentialFeatures() []EndpointFeature {
	features := GetEndpointFeatureMap()

	var missing []EndpointFeature

	for _, feature := range features {
		if !feature.Essential {
			continue
		}

		// Check if device supports this essential feature
		hasSupport := false
		for _, endpoint := range feature.Endpoints {
			if s.HasURL(endpoint) {
				hasSupport = true
				break
			}
		}

		if !hasSupport {
			missing = append(missing, feature)
		}
	}

	return missing
}

// GetFeatureCompleteness returns a completeness score (0-100) based on supported features
func (s *SupportedURLsResponse) GetFeatureCompleteness() (int, int, int) {
	features := GetEndpointFeatureMap()

	var total, supported, essential int

	for _, feature := range features {
		total++

		if feature.Essential {
			essential++
		}

		// Check if device supports this feature
		hasSupport := false
		for _, endpoint := range feature.Endpoints {
			if s.HasURL(endpoint) {
				hasSupport = true
				break
			}
		}

		if hasSupport {
			supported++
		}
	}

	completeness := 0
	if total > 0 {
		completeness = (supported * 100) / total
	}

	return completeness, supported, total
}
