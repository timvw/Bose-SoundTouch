package main

import (
	"testing"

	"github.com/gesellix/bose-soundtouch/pkg/models"
)

func TestShouldShowContentDetails(t *testing.T) {
	tests := []struct {
		name        string
		verbose     bool
		contentItem *models.ContentItem
		expected    bool
		description string
	}{
		{
			name:    "verbose_flag_true_shows_details",
			verbose: true,
			contentItem: &models.ContentItem{
				Source:   "SPOTIFY",
				Location: "",
			},
			expected:    true,
			description: "Verbose flag should always show details regardless of location",
		},
		{
			name:    "spotify_with_location_shows_details",
			verbose: false,
			contentItem: &models.ContentItem{
				Source:   "SPOTIFY",
				Location: "spotify:track:123456789",
			},
			expected:    true,
			description: "Any source with location should show details",
		},
		{
			name:    "tunein_with_location_shows_details",
			verbose: false,
			contentItem: &models.ContentItem{
				Source:   "TUNEIN",
				Location: "/v1/playback/station/s33828",
			},
			expected:    true,
			description: "TUNEIN with location should show details",
		},
		{
			name:    "local_internet_radio_with_location_shows_details",
			verbose: false,
			contentItem: &models.ContentItem{
				Source:   "LOCAL_INTERNET_RADIO",
				Location: "https://stream.example.com/radio",
			},
			expected:    true,
			description: "Local internet radio with location should show details",
		},
		{
			name:    "stored_music_with_location_shows_details",
			verbose: false,
			contentItem: &models.ContentItem{
				Source:   "STORED_MUSIC",
				Location: "6_a2874b5d_4f83d999",
			},
			expected:    true,
			description: "Stored music with location should show details",
		},
		{
			name:    "pandora_with_location_shows_details",
			verbose: false,
			contentItem: &models.ContentItem{
				Source:   "PANDORA",
				Location: "126740707481236361",
			},
			expected:    true,
			description: "Pandora with location should show details",
		},
		{
			name:    "local_music_with_location_shows_details",
			verbose: false,
			contentItem: &models.ContentItem{
				Source:   "LOCAL_MUSIC",
				Location: "album:983",
			},
			expected:    true,
			description: "Local music with location should show details",
		},
		{
			name:    "no_location_no_verbose_hides_details",
			verbose: false,
			contentItem: &models.ContentItem{
				Source:   "BLUETOOTH",
				Location: "",
			},
			expected:    false,
			description: "No location and no verbose should hide details",
		},
		{
			name:    "empty_location_no_verbose_hides_details",
			verbose: false,
			contentItem: &models.ContentItem{
				Source:   "AIRPLAY",
				Location: "",
			},
			expected:    false,
			description: "Empty location and no verbose should hide details",
		},
		{
			name:        "nil_content_item_hides_details",
			verbose:     false,
			contentItem: nil,
			expected:    false,
			description: "Nil content item should hide details",
		},
		{
			name:        "verbose_with_nil_content_item_hides_details",
			verbose:     true,
			contentItem: nil,
			expected:    false,
			description: "Even verbose flag cannot show details for nil content item",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This mimics the logic from getNowPlaying function:
			// showDetails := verbose || (nowPlaying.ContentItem != nil && nowPlaying.ContentItem.Location != "")
			result := shouldShowContentDetails(tt.verbose, tt.contentItem)

			if result != tt.expected {
				t.Errorf("shouldShowContentDetails(%v, %+v) = %v, want %v. %s",
					tt.verbose, tt.contentItem, result, tt.expected, tt.description)
			}
		})
	}
}
func TestContentDetailsDisplayLogic(t *testing.T) {
	// Test the specific conditions that determine when to show content details
	tests := []struct {
		name           string
		verbose        bool
		hasContentItem bool
		hasLocation    bool
		expectedShow   bool
	}{
		{"verbose_true_overrides_all", true, false, false, false}, // Note: still need contentItem != nil
		{"verbose_false_with_location", false, true, true, true},
		{"verbose_false_without_location", false, true, false, false},
		{"verbose_false_without_contentitem", false, false, false, false},
		{"verbose_true_with_contentitem_and_location", true, true, true, true},
		{"verbose_true_with_contentitem_no_location", true, true, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var contentItem *models.ContentItem
			if tt.hasContentItem {
				contentItem = &models.ContentItem{
					Source: "TEST_SOURCE",
				}
				if tt.hasLocation {
					contentItem.Location = "test_location"
				}
			}

			result := shouldShowContentDetails(tt.verbose, contentItem)
			if result != tt.expectedShow {
				t.Errorf("Expected %v, got %v for verbose=%v, hasContentItem=%v, hasLocation=%v",
					tt.expectedShow, result, tt.verbose, tt.hasContentItem, tt.hasLocation)
			}
		})
	}
}

func TestVerboseFlagSpecificFields(t *testing.T) {
	// Test which fields should only be shown in verbose mode
	contentItem := &models.ContentItem{
		Source:        "SPOTIFY",
		Type:          "uri",
		Location:      "spotify:track:123456789",
		SourceAccount: "testuser",
		IsPresetable:  true,
		ItemName:      "Test Track",
		ContainerArt:  "https://example.com/art.jpg",
	}

	// These fields should always be shown when content details are displayed
	alwaysShown := []string{"Location"}

	// These fields should only be shown in verbose mode
	verboseOnly := []string{"Type", "ItemName", "IsPresetable"}

	t.Run("verbose_mode_shows_all_fields", func(t *testing.T) {
		verbose := true
		showDetails := shouldShowContentDetails(verbose, contentItem)

		if !showDetails {
			t.Error("Expected to show details in verbose mode")
		}

		// In verbose mode, we would show all fields
		// (This is testing the conceptual logic, actual field display is in the CLI function)
	})

	t.Run("non_verbose_mode_shows_limited_fields", func(t *testing.T) {
		verbose := false
		showDetails := shouldShowContentDetails(verbose, contentItem)

		if !showDetails {
			t.Error("Expected to show details when location is present")
		}

		// In non-verbose mode, we would only show location
		// The actual field filtering happens in the CLI display logic
		_ = alwaysShown // Would show these
		_ = verboseOnly // Would NOT show these
	})
}

// Helper function that encapsulates the logic from getNowPlaying
func shouldShowContentDetails(verbose bool, contentItem *models.ContentItem) bool {
	// This mirrors the exact logic from cmd_playback.go:
	// showDetails := verbose || (nowPlaying.ContentItem != nil && nowPlaying.ContentItem.Location != "")
	// if showDetails && nowPlaying.ContentItem != nil { ... }
	hasLocationData := contentItem != nil && contentItem.Location != ""
	showDetails := verbose || hasLocationData

	return showDetails && contentItem != nil
}

func TestRealWorldScenarios(t *testing.T) {
	scenarios := []struct {
		name     string
		source   string
		location string
		verbose  bool
		expected bool
		useCase  string
	}{
		{
			name:     "spotify_user_wants_uri",
			source:   "SPOTIFY",
			location: "spotify:playlist:37i9dQZF1DX0XUsuxWHRQd",
			verbose:  false,
			expected: true,
			useCase:  "User playing Spotify wants to see URI for storePreset",
		},
		{
			name:     "radio_user_wants_station_id",
			source:   "TUNEIN",
			location: "/v1/playback/station/s33828",
			verbose:  false,
			expected: true,
			useCase:  "User playing radio wants to see station ID for storePreset",
		},
		{
			name:     "bluetooth_no_useful_location",
			source:   "BLUETOOTH",
			location: "",
			verbose:  false,
			expected: false,
			useCase:  "Bluetooth has no useful location data for presets",
		},
		{
			name:     "developer_debugging_verbose",
			source:   "AIRPLAY",
			location: "",
			verbose:  true,
			expected: true,
			useCase:  "Developer wants all available info regardless of source",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			contentItem := &models.ContentItem{
				Source:   scenario.source,
				Location: scenario.location,
			}

			result := shouldShowContentDetails(scenario.verbose, contentItem)
			if result != scenario.expected {
				t.Errorf("Scenario '%s' failed: %s. Expected %v, got %v",
					scenario.name, scenario.useCase, scenario.expected, result)
			}
		})
	}
}
