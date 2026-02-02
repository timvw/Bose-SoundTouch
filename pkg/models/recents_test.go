package models

import (
	"encoding/xml"
	"testing"
)

func TestRecentsResponse_Unmarshal(t *testing.T) {
	tests := []struct {
		name        string
		xmlData     string
		expected    *RecentsResponse
		expectError bool
	}{
		{
			name: "complete recents response",
			xmlData: `<recents>
  <recent deviceID="1004567890AA" utcTime="1701202831">
    <contentItem source="STORED_MUSIC" location="6_a2874b5d_4f83d999" sourceAccount="d09708a1-5953-44bc-a413-123456789012/0" isPresetable="true">
      <itemName>MercyMe, It's Christmas!</itemName>
    </contentItem>
  </recent>
  <recent deviceID="1004567890AA" utcTime="1700232917" id="2487503626">
    <contentItem source="LOCAL_MUSIC" type="track" location="track:2590" sourceAccount="3f205110-4a57-4e91-810a-123456789012" isPresetable="true">
      <itemName>Baby It's Cold Outside - ANNE MURRAY</itemName>
    </contentItem>
  </recent>
</recents>`,
			expected: &RecentsResponse{
				Items: []RecentsResponseItem{
					{
						DeviceID: "1004567890AA",
						UTCTime:  1701202831,
						ContentItem: &ContentItem{
							Source:        "STORED_MUSIC",
							Location:      "6_a2874b5d_4f83d999",
							SourceAccount: "d09708a1-5953-44bc-a413-123456789012/0",
							IsPresetable:  true,
							ItemName:      "MercyMe, It's Christmas!",
						},
					},
					{
						DeviceID: "1004567890AA",
						UTCTime:  1700232917,
						ID:       "2487503626",
						ContentItem: &ContentItem{
							Source:        "LOCAL_MUSIC",
							Type:          "track",
							Location:      "track:2590",
							SourceAccount: "3f205110-4a57-4e91-810a-123456789012",
							IsPresetable:  true,
							ItemName:      "Baby It's Cold Outside - ANNE MURRAY",
						},
					},
				},
			},
		},
		{
			name: "spotify recent item",
			xmlData: `<recents>
  <recent deviceID="1004567890AA" utcTime="1701300000" id="spotify123">
    <contentItem source="SPOTIFY" type="track" location="spotify:track:4iV5W9uYEdYUVa79Axb7Rh" sourceAccount="spotify_user" isPresetable="true">
      <itemName>Shape of You - Ed Sheeran</itemName>
      <containerArt>https://i.scdn.co/image/ab67616d0000b273ba5db46f4b838ef6027e6f96</containerArt>
    </contentItem>
  </recent>
</recents>`,
			expected: &RecentsResponse{
				Items: []RecentsResponseItem{
					{
						DeviceID: "1004567890AA",
						UTCTime:  1701300000,
						ID:       "spotify123",
						ContentItem: &ContentItem{
							Source:        "SPOTIFY",
							Type:          "track",
							Location:      "spotify:track:4iV5W9uYEdYUVa79Axb7Rh",
							SourceAccount: "spotify_user",
							IsPresetable:  true,
							ItemName:      "Shape of You - Ed Sheeran",
							ContainerArt:  "https://i.scdn.co/image/ab67616d0000b273ba5db46f4b838ef6027e6f96",
						},
					},
				},
			},
		},
		{
			name: "empty recents",
			xmlData: `<recents>
</recents>`,
			expected: &RecentsResponse{
				Items: []RecentsResponseItem{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var response RecentsResponse
			err := xml.Unmarshal([]byte(tt.xmlData), &response)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			// Compare basic structure
			if len(response.Items) != len(tt.expected.Items) {
				t.Errorf("expected %d items, got %d", len(tt.expected.Items), len(response.Items))
			}

			// Compare each item
			for i, expectedItem := range tt.expected.Items {
				if i >= len(response.Items) {
					break
				}

				actualItem := response.Items[i]

				if actualItem.DeviceID != expectedItem.DeviceID {
					t.Errorf("item %d: expected deviceID %s, got %s", i, expectedItem.DeviceID, actualItem.DeviceID)
				}
				if actualItem.UTCTime != expectedItem.UTCTime {
					t.Errorf("item %d: expected utcTime %d, got %d", i, expectedItem.UTCTime, actualItem.UTCTime)
				}
				if actualItem.ID != expectedItem.ID {
					t.Errorf("item %d: expected id %s, got %s", i, expectedItem.ID, actualItem.ID)
				}

				// Compare ContentItem
				if expectedItem.ContentItem != nil {
					if actualItem.ContentItem == nil {
						t.Errorf("item %d: expected contentItem, got nil", i)
						continue
					}

					if actualItem.ContentItem.Source != expectedItem.ContentItem.Source {
						t.Errorf("item %d: expected source %s, got %s", i, expectedItem.ContentItem.Source, actualItem.ContentItem.Source)
					}
					if actualItem.ContentItem.Type != expectedItem.ContentItem.Type {
						t.Errorf("item %d: expected type %s, got %s", i, expectedItem.ContentItem.Type, actualItem.ContentItem.Type)
					}
					if actualItem.ContentItem.Location != expectedItem.ContentItem.Location {
						t.Errorf("item %d: expected location %s, got %s", i, expectedItem.ContentItem.Location, actualItem.ContentItem.Location)
					}
					if actualItem.ContentItem.ItemName != expectedItem.ContentItem.ItemName {
						t.Errorf("item %d: expected itemName %s, got %s", i, expectedItem.ContentItem.ItemName, actualItem.ContentItem.ItemName)
					}
					if actualItem.ContentItem.IsPresetable != expectedItem.ContentItem.IsPresetable {
						t.Errorf("item %d: expected isPresetable %t, got %t", i, expectedItem.ContentItem.IsPresetable, actualItem.ContentItem.IsPresetable)
					}
				} else if actualItem.ContentItem != nil {
					t.Errorf("item %d: expected nil contentItem, got non-nil", i)
				}
			}
		})
	}
}

func TestRecentsResponse_Methods(t *testing.T) {
	response := &RecentsResponse{
		Items: []RecentsResponseItem{
			{
				DeviceID: "device1",
				UTCTime:  1701200000,
				ContentItem: &ContentItem{
					Source:       "SPOTIFY",
					Type:         "track",
					ItemName:     "Song 1",
					IsPresetable: true,
				},
			},
			{
				DeviceID: "device1",
				UTCTime:  1701100000,
				ContentItem: &ContentItem{
					Source:       "LOCAL_MUSIC",
					Type:         "track",
					ItemName:     "Song 2",
					IsPresetable: false,
				},
			},
			{
				DeviceID: "device1",
				UTCTime:  1701000000,
				ContentItem: &ContentItem{
					Source:       "TUNEIN",
					Type:         "stationurl",
					ItemName:     "Radio Station",
					IsPresetable: true,
				},
			},
		},
	}

	// Test GetItemCount
	if response.GetItemCount() != 3 {
		t.Errorf("expected item count 3, got %d", response.GetItemCount())
	}

	// Test IsEmpty
	if response.IsEmpty() {
		t.Error("expected IsEmpty() to return false")
	}

	// Test GetMostRecent
	mostRecent := response.GetMostRecent()
	if mostRecent == nil {
		t.Error("expected most recent item, got nil")
	} else if mostRecent.UTCTime != 1701200000 {
		t.Errorf("expected most recent UTCTime 1701200000, got %d", mostRecent.UTCTime)
	}

	// Test GetItemsBySource
	spotifyItems := response.GetItemsBySource("SPOTIFY")
	if len(spotifyItems) != 1 {
		t.Errorf("expected 1 Spotify item, got %d", len(spotifyItems))
	}

	localItems := response.GetItemsBySource("LOCAL_MUSIC")
	if len(localItems) != 1 {
		t.Errorf("expected 1 LOCAL_MUSIC item, got %d", len(localItems))
	}

	// Test GetSpotifyItems
	spotifyItems2 := response.GetSpotifyItems()
	if len(spotifyItems2) != 1 {
		t.Errorf("expected 1 Spotify item, got %d", len(spotifyItems2))
	}

	// Test GetPresetableItems
	presetableItems := response.GetPresetableItems()
	if len(presetableItems) != 2 {
		t.Errorf("expected 2 presetable items, got %d", len(presetableItems))
	}

	// Test GetTracks
	tracks := response.GetTracks()
	if len(tracks) != 2 {
		t.Errorf("expected 2 track items, got %d", len(tracks))
	}

	// Test GetStations
	stations := response.GetStations()
	if len(stations) != 1 {
		t.Errorf("expected 1 station item, got %d", len(stations))
	}
}

func TestRecentsResponse_EmptyResponse(t *testing.T) {
	response := &RecentsResponse{
		Items: []RecentsResponseItem{},
	}

	// Test empty response methods
	if response.GetItemCount() != 0 {
		t.Errorf("expected item count 0, got %d", response.GetItemCount())
	}

	if !response.IsEmpty() {
		t.Error("expected IsEmpty() to return true")
	}

	if response.GetMostRecent() != nil {
		t.Error("expected GetMostRecent() to return nil")
	}

	if len(response.GetSpotifyItems()) != 0 {
		t.Errorf("expected 0 Spotify items, got %d", len(response.GetSpotifyItems()))
	}
}

func TestRecentItem_Methods(t *testing.T) {
	tests := []struct {
		name string
		item RecentsResponseItem
		test func(t *testing.T, item *RecentsResponseItem)
	}{
		{
			name: "spotify track item",
			item: RecentsResponseItem{
				DeviceID: "device1",
				UTCTime:  1701200000,
				ID:       "spotify123",
				ContentItem: &ContentItem{
					Source:        "SPOTIFY",
					Type:          "track",
					Location:      "spotify:track:123",
					SourceAccount: "user@spotify.com",
					IsPresetable:  true,
					ItemName:      "Test Song",
					ContainerArt:  "https://example.com/art.jpg",
				},
			},
			test: func(t *testing.T, item *RecentsResponseItem) {
				if !item.HasContent() {
					t.Error("expected HasContent() to return true")
				}
				if item.GetDisplayName() != "Test Song" {
					t.Errorf("expected display name 'Test Song', got %s", item.GetDisplayName())
				}
				if item.GetSource() != "SPOTIFY" {
					t.Errorf("expected source 'SPOTIFY', got %s", item.GetSource())
				}
				if !item.IsTrack() {
					t.Error("expected IsTrack() to return true")
				}
				if !item.IsSpotifyContent() {
					t.Error("expected IsSpotifyContent() to return true")
				}
				if !item.IsStreamingContent() {
					t.Error("expected IsStreamingContent() to return true")
				}
				if item.IsLocalContent() {
					t.Error("expected IsLocalContent() to return false")
				}
				if !item.IsPresetable() {
					t.Error("expected IsPresetable() to return true")
				}
				if !item.HasArtwork() {
					t.Error("expected HasArtwork() to return true")
				}
				if item.GetArtwork() != "https://example.com/art.jpg" {
					t.Errorf("expected artwork URL, got %s", item.GetArtwork())
				}
				if item.GetUTCTime() != 1701200000 {
					t.Errorf("expected UTC time 1701200000, got %d", item.GetUTCTime())
				}
				if !item.HasID() {
					t.Error("expected HasID() to return true")
				}
				if item.GetID() != "spotify123" {
					t.Errorf("expected ID 'spotify123', got %s", item.GetID())
				}
			},
		},
		{
			name: "local music item",
			item: RecentsResponseItem{
				DeviceID: "device1",
				UTCTime:  1701100000,
				ContentItem: &ContentItem{
					Source:       "LOCAL_MUSIC",
					Type:         "track",
					Location:     "/music/song.mp3",
					IsPresetable: false,
					ItemName:     "Local Song",
				},
			},
			test: func(t *testing.T, item *RecentsResponseItem) {
				if !item.IsLocalContent() {
					t.Error("expected IsLocalContent() to return true")
				}
				if item.IsStreamingContent() {
					t.Error("expected IsStreamingContent() to return false")
				}
				if item.IsSpotifyContent() {
					t.Error("expected IsSpotifyContent() to return false")
				}
				if item.HasArtwork() {
					t.Error("expected HasArtwork() to return false")
				}
				if item.GetArtwork() != "" {
					t.Errorf("expected empty artwork, got %s", item.GetArtwork())
				}
			},
		},
		{
			name: "radio station item",
			item: RecentsResponseItem{
				DeviceID: "device1",
				UTCTime:  1701000000,
				ContentItem: &ContentItem{
					Source:       "TUNEIN",
					Type:         "stationurl",
					Location:     "tunein:station:123",
					IsPresetable: true,
					ItemName:     "Rock FM",
				},
			},
			test: func(t *testing.T, item *RecentsResponseItem) {
				if !item.IsStation() {
					t.Error("expected IsStation() to return true")
				}
				if item.IsTrack() {
					t.Error("expected IsTrack() to return false")
				}
				if !item.IsStreamingContent() {
					t.Error("expected IsStreamingContent() to return true")
				}
			},
		},
		{
			name: "empty content item",
			item: RecentsResponseItem{
				DeviceID: "device1",
				UTCTime:  1701000000,
			},
			test: func(t *testing.T, item *RecentsResponseItem) {
				if item.HasContent() {
					t.Error("expected HasContent() to return false")
				}
				if item.GetDisplayName() != "Unknown Item" {
					t.Errorf("expected display name 'Unknown Item', got %s", item.GetDisplayName())
				}
				if item.GetSource() != "" {
					t.Errorf("expected empty source, got %s", item.GetSource())
				}
				if item.IsTrack() {
					t.Error("expected IsTrack() to return false")
				}
				if item.IsPresetable() {
					t.Error("expected IsPresetable() to return false")
				}
				if item.HasID() {
					t.Error("expected HasID() to return false")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t, &tt.item)
		})
	}
}

func TestRecentItem_ContentTypes(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		expected    map[string]bool
	}{
		{
			name:        "track type",
			contentType: "track",
			expected: map[string]bool{
				"IsTrack":     true,
				"IsStation":   false,
				"IsPlaylist":  false,
				"IsAlbum":     false,
				"IsContainer": false,
			},
		},
		{
			name:        "station type",
			contentType: "stationurl",
			expected: map[string]bool{
				"IsTrack":     false,
				"IsStation":   true,
				"IsPlaylist":  false,
				"IsAlbum":     false,
				"IsContainer": false,
			},
		},
		{
			name:        "playlist type",
			contentType: "playlist",
			expected: map[string]bool{
				"IsTrack":     false,
				"IsStation":   false,
				"IsPlaylist":  true,
				"IsAlbum":     false,
				"IsContainer": false,
			},
		},
		{
			name:        "album type",
			contentType: "album",
			expected: map[string]bool{
				"IsTrack":     false,
				"IsStation":   false,
				"IsPlaylist":  false,
				"IsAlbum":     true,
				"IsContainer": false,
			},
		},
		{
			name:        "container type",
			contentType: "container",
			expected: map[string]bool{
				"IsTrack":     false,
				"IsStation":   false,
				"IsPlaylist":  false,
				"IsAlbum":     false,
				"IsContainer": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := RecentsResponseItem{
				ContentItem: &ContentItem{
					Type: tt.contentType,
				},
			}

			results := map[string]bool{
				"IsTrack":     item.IsTrack(),
				"IsStation":   item.IsStation(),
				"IsPlaylist":  item.IsPlaylist(),
				"IsAlbum":     item.IsAlbum(),
				"IsContainer": item.IsContainer(),
			}

			for method, expected := range tt.expected {
				if results[method] != expected {
					t.Errorf("expected %s() to return %t, got %t", method, expected, results[method])
				}
			}
		})
	}
}

func TestRecentsResponse_FilterMethods(t *testing.T) {
	response := &RecentsResponse{
		Items: []RecentsResponseItem{
			{
				ContentItem: &ContentItem{Source: "SPOTIFY", Type: "track"},
			},
			{
				ContentItem: &ContentItem{Source: "PANDORA", Type: "track"},
			},
			{
				ContentItem: &ContentItem{Source: "LOCAL_MUSIC", Type: "track"},
			},
			{
				ContentItem: &ContentItem{Source: "STORED_MUSIC", Type: "track"},
			},
			{
				ContentItem: &ContentItem{Source: "TUNEIN", Type: "stationurl"},
			},
			{
				ContentItem: &ContentItem{Source: "SPOTIFY", Type: "playlist"},
			},
		},
	}

	// Test individual service filters
	if len(response.GetSpotifyItems()) != 2 {
		t.Errorf("expected 2 Spotify items, got %d", len(response.GetSpotifyItems()))
	}

	if len(response.GetPandoraItems()) != 1 {
		t.Errorf("expected 1 Pandora item, got %d", len(response.GetPandoraItems()))
	}

	if len(response.GetLocalMusicItems()) != 1 {
		t.Errorf("expected 1 LOCAL_MUSIC item, got %d", len(response.GetLocalMusicItems()))
	}

	if len(response.GetStoredMusicItems()) != 1 {
		t.Errorf("expected 1 STORED_MUSIC item, got %d", len(response.GetStoredMusicItems()))
	}

	if len(response.GetTuneInItems()) != 1 {
		t.Errorf("expected 1 TuneIn item, got %d", len(response.GetTuneInItems()))
	}

	// Test type filters
	if len(response.GetTracks()) != 4 {
		t.Errorf("expected 4 track items, got %d", len(response.GetTracks()))
	}

	if len(response.GetStations()) != 1 {
		t.Errorf("expected 1 station item, got %d", len(response.GetStations()))
	}

	if len(response.GetPlaylistsAndAlbums()) != 1 {
		t.Errorf("expected 1 playlist/album item, got %d", len(response.GetPlaylistsAndAlbums()))
	}
}
