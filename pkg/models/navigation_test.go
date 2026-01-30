package models

import (
	"encoding/xml"
	"testing"
)

func TestNavigateRequest_NewNavigateRequest(t *testing.T) {
	req := NewNavigateRequest("SPOTIFY", "user@example.com", 1, 50)

	if req.Source != "SPOTIFY" {
		t.Errorf("Expected source SPOTIFY, got %s", req.Source)
	}
	if req.SourceAccount != "user@example.com" {
		t.Errorf("Expected sourceAccount user@example.com, got %s", req.SourceAccount)
	}
	if req.StartItem != 1 {
		t.Errorf("Expected startItem 1, got %d", req.StartItem)
	}
	if req.NumItems != 50 {
		t.Errorf("Expected numItems 50, got %d", req.NumItems)
	}
}

func TestNavigateRequest_NewNavigateRequestWithMenu(t *testing.T) {
	req := NewNavigateRequestWithMenu("PANDORA", "user123", "radioStations", "dateCreated", 1, 100)

	if req.Source != "PANDORA" {
		t.Errorf("Expected source PANDORA, got %s", req.Source)
	}
	if req.Menu != "radioStations" {
		t.Errorf("Expected menu radioStations, got %s", req.Menu)
	}
	if req.Sort != "dateCreated" {
		t.Errorf("Expected sort dateCreated, got %s", req.Sort)
	}
}

func TestNavigateRequest_XMLMarshal(t *testing.T) {
	tests := []struct {
		name     string
		request  *NavigateRequest
		expected string
	}{
		{
			name:     "Basic navigate request",
			request:  NewNavigateRequest("TUNEIN", "", 1, 25),
			expected: `<navigate source="TUNEIN"><startItem>1</startItem><numItems>25</numItems></navigate>`,
		},
		{
			name:     "Navigate with source account",
			request:  NewNavigateRequest("SPOTIFY", "user@example.com", 10, 50),
			expected: `<navigate source="SPOTIFY" sourceAccount="user@example.com"><startItem>10</startItem><numItems>50</numItems></navigate>`,
		},
		{
			name:     "Navigate with menu and sort",
			request:  NewNavigateRequestWithMenu("PANDORA", "user123", "radioStations", "dateCreated", 1, 100),
			expected: `<navigate source="PANDORA" sourceAccount="user123" menu="radioStations" sort="dateCreated"><startItem>1</startItem><numItems>100</numItems></navigate>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xmlData, err := xml.Marshal(tt.request)
			if err != nil {
				t.Fatalf("Failed to marshal XML: %v", err)
			}

			actual := string(xmlData)
			if actual != tt.expected {
				t.Errorf("XML mismatch:\nExpected: %s\nActual:   %s", tt.expected, actual)
			}
		})
	}
}

func TestNavigateRequest_XMLMarshalWithItem(t *testing.T) {
	containerItem := &ContentItem{
		Source:        "STORED_MUSIC",
		Location:      "1",
		SourceAccount: "device123/0",
		IsPresetable:  true,
		ItemName:      "Music",
	}

	request := NewNavigateRequestWithItem("STORED_MUSIC", "device123/0", 1, 1000, containerItem)

	xmlData, err := xml.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal XML: %v", err)
	}

	// Debug: print the actual XML to see what's generated
	xmlStr := string(xmlData)
	t.Logf("Generated XML: %s", xmlStr)

	// Check that the XML contains expected elements
	if !contains(xmlStr, `source="STORED_MUSIC"`) {
		t.Error("XML should contain source attribute")
	}
	if !contains(xmlStr, `<startItem>1</startItem>`) {
		t.Error("XML should contain startItem element")
	}
	if !contains(xmlStr, `<numItems>1000</numItems>`) {
		t.Error("XML should contain numItems element")
	}
	if !contains(xmlStr, `<item`) {
		t.Error("XML should contain item element")
	}
	if !contains(xmlStr, `<ContentItem`) {
		t.Error("XML should contain ContentItem element")
	}
}

func TestNavigateResponse_XMLUnmarshal(t *testing.T) {
	xmlData := `
		<navigateResponse source="STORED_MUSIC" sourceAccount="device123/0">
			<totalItems>2</totalItems>
			<items>
				<item Playable="1">
					<name>Album Artists</name>
					<type>dir</type>
					<ContentItem source="STORED_MUSIC" location="107" sourceAccount="device123/0" isPresetable="true">
						<itemName>Album Artists</itemName>
					</ContentItem>
				</item>
				<item Playable="1">
					<name>Test Track</name>
					<type>track</type>
					<ContentItem source="STORED_MUSIC" location="track123" sourceAccount="device123/0" isPresetable="true">
						<itemName>Test Track</itemName>
					</ContentItem>
					<artistName>Test Artist</artistName>
					<albumName>Test Album</albumName>
				</item>
			</items>
		</navigateResponse>`

	var response NavigateResponse
	err := xml.Unmarshal([]byte(xmlData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	if response.Source != "STORED_MUSIC" {
		t.Errorf("Expected source STORED_MUSIC, got %s", response.Source)
	}
	if response.TotalItems != 2 {
		t.Errorf("Expected totalItems 2, got %d", response.TotalItems)
	}
	if len(response.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(response.Items))
	}

	// Test first item (directory)
	firstItem := response.Items[0]
	if firstItem.Name != "Album Artists" {
		t.Errorf("Expected first item name 'Album Artists', got %s", firstItem.Name)
	}
	if firstItem.Type != "dir" {
		t.Errorf("Expected first item type 'dir', got %s", firstItem.Type)
	}
	if !firstItem.IsPlayable() {
		t.Error("Expected first item to be playable")
	}
	if !firstItem.IsDirectory() {
		t.Error("Expected first item to be a directory")
	}

	// Test second item (track)
	secondItem := response.Items[1]
	if secondItem.Name != "Test Track" {
		t.Errorf("Expected second item name 'Test Track', got %s", secondItem.Name)
	}
	if secondItem.ArtistName != "Test Artist" {
		t.Errorf("Expected artist name 'Test Artist', got %s", secondItem.ArtistName)
	}
	if !secondItem.IsTrack() {
		t.Error("Expected second item to be a track")
	}
}

func TestNavigateResponse_FilterMethods(t *testing.T) {
	response := &NavigateResponse{
		TotalItems: 3,
		Items: []NavigateItem{
			{
				Name:     "Directory 1",
				Type:     "dir",
				Playable: 1,
			},
			{
				Name:     "Track 1",
				Type:     "track",
				Playable: 1,
			},
			{
				Name:     "Station 1",
				Type:     "stationurl",
				Playable: 1,
			},
		},
	}

	// Test GetPlayableItems
	playable := response.GetPlayableItems()
	if len(playable) != 3 {
		t.Errorf("Expected 3 playable items, got %d", len(playable))
	}

	// Test GetDirectories
	directories := response.GetDirectories()
	if len(directories) != 1 {
		t.Errorf("Expected 1 directory, got %d", len(directories))
	}
	if directories[0].Name != "Directory 1" {
		t.Errorf("Expected directory name 'Directory 1', got %s", directories[0].Name)
	}

	// Test GetTracks
	tracks := response.GetTracks()
	if len(tracks) != 1 {
		t.Errorf("Expected 1 track, got %d", len(tracks))
	}
	if tracks[0].Name != "Track 1" {
		t.Errorf("Expected track name 'Track 1', got %s", tracks[0].Name)
	}

	// Test GetStations
	stations := response.GetStations()
	if len(stations) != 1 {
		t.Errorf("Expected 1 station, got %d", len(stations))
	}
	if stations[0].Name != "Station 1" {
		t.Errorf("Expected station name 'Station 1', got %s", stations[0].Name)
	}
}

func TestNavigateResponse_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		response *NavigateResponse
		expected bool
	}{
		{
			name:     "Empty response - no items",
			response: &NavigateResponse{TotalItems: 0, Items: []NavigateItem{}},
			expected: true,
		},
		{
			name:     "Empty response - zero total",
			response: &NavigateResponse{TotalItems: 0, Items: []NavigateItem{{Name: "test"}}},
			expected: true,
		},
		{
			name:     "Empty response - nil items",
			response: &NavigateResponse{TotalItems: 1, Items: nil},
			expected: true,
		},
		{
			name:     "Non-empty response",
			response: &NavigateResponse{TotalItems: 1, Items: []NavigateItem{{Name: "test"}}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual := tt.response.IsEmpty(); actual != tt.expected {
				t.Errorf("IsEmpty() = %v, expected %v", actual, tt.expected)
			}
		})
	}
}

func TestNavigateItem_GetDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		item     *NavigateItem
		expected string
	}{
		{
			name:     "Item with name",
			item:     &NavigateItem{Name: "Test Item"},
			expected: "Test Item",
		},
		{
			name:     "Item with ContentItem itemName",
			item:     &NavigateItem{ContentItem: &ContentItem{ItemName: "Content Name"}},
			expected: "Content Name",
		},
		{
			name:     "Item with both names - prefers item name",
			item:     &NavigateItem{Name: "Item Name", ContentItem: &ContentItem{ItemName: "Content Name"}},
			expected: "Item Name",
		},
		{
			name:     "Item with no names",
			item:     &NavigateItem{},
			expected: "Unknown Item",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual := tt.item.GetDisplayName(); actual != tt.expected {
				t.Errorf("GetDisplayName() = %s, expected %s", actual, tt.expected)
			}
		})
	}
}

func TestAddStationRequest_XMLMarshal(t *testing.T) {
	req := NewAddStationRequest("PANDORA", "user123", "R4328162", "Test Station")

	xmlData, err := xml.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal XML: %v", err)
	}

	expected := `<addStation source="PANDORA" sourceAccount="user123" token="R4328162"><name>Test Station</name></addStation>`
	actual := string(xmlData)

	if actual != expected {
		t.Errorf("XML mismatch:\nExpected: %s\nActual:   %s", expected, actual)
	}
}

func TestAddStationRequest_XMLMarshalWithoutAccount(t *testing.T) {
	req := NewAddStationRequest("TUNEIN", "", "station123", "Radio Station")

	xmlData, err := xml.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal XML: %v", err)
	}

	expected := `<addStation source="TUNEIN" token="station123"><name>Radio Station</name></addStation>`
	actual := string(xmlData)

	if actual != expected {
		t.Errorf("XML mismatch:\nExpected: %s\nActual:   %s", expected, actual)
	}
}

func TestStationResponse_XMLUnmarshal(t *testing.T) {
	xmlData := `<status>/addStation</status>`

	var response StationResponse
	err := xml.Unmarshal([]byte(xmlData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	if response.Status != "/addStation" {
		t.Errorf("Expected status '/addStation', got %s", response.Status)
	}
}

func TestRemoveStationRequest(t *testing.T) {
	contentItem := &ContentItem{
		Source:        "PANDORA",
		Location:      "126740707481236361",
		SourceAccount: "user123",
		IsPresetable:  true,
		ItemName:      "Test Station",
	}

	req := NewRemoveStationRequest(contentItem)

	// Since RemoveStationRequest is just ContentItem, test that it's the same
	if req != contentItem {
		t.Error("NewRemoveStationRequest should return the same ContentItem")
	}

	// Test XML marshaling
	xmlData, err := xml.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal XML: %v", err)
	}

	xmlStr := string(xmlData)
	if !contains(xmlStr, `source="PANDORA"`) {
		t.Error("XML should contain source attribute")
	}
	if !contains(xmlStr, `location="126740707481236361"`) {
		t.Error("XML should contain location attribute")
	}
	if !contains(xmlStr, `<itemName>Test Station</itemName>`) {
		t.Error("XML should contain itemName element")
	}
}

func TestSearchStationRequest_NewSearchStationRequest(t *testing.T) {
	req := NewSearchStationRequest("PANDORA", "user123", "Zach Williams")

	if req.Source != "PANDORA" {
		t.Errorf("Expected source PANDORA, got %s", req.Source)
	}
	if req.SourceAccount != "user123" {
		t.Errorf("Expected sourceAccount user123, got %s", req.SourceAccount)
	}
	if req.SearchTerm != "Zach Williams" {
		t.Errorf("Expected searchTerm 'Zach Williams', got %s", req.SearchTerm)
	}
}

func TestSearchStationRequest_XMLMarshal(t *testing.T) {
	tests := []struct {
		name     string
		request  *SearchStationRequest
		expected string
	}{
		{
			name:     "Basic search request",
			request:  NewSearchStationRequest("PANDORA", "user123", "Classic Rock"),
			expected: `<search source="PANDORA" sourceAccount="user123">Classic Rock</search>`,
		},
		{
			name:     "Search without source account",
			request:  NewSearchStationRequest("TUNEIN", "", "Jazz"),
			expected: `<search source="TUNEIN">Jazz</search>`,
		},
		{
			name:     "Search with special characters",
			request:  NewSearchStationRequest("SPOTIFY", "user@example.com", "Rock & Roll"),
			expected: `<search source="SPOTIFY" sourceAccount="user@example.com">Rock &amp; Roll</search>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xmlData, err := xml.Marshal(tt.request)
			if err != nil {
				t.Fatalf("Failed to marshal XML: %v", err)
			}

			actual := string(xmlData)
			if actual != tt.expected {
				t.Errorf("XML mismatch:\nExpected: %s\nActual:   %s", tt.expected, actual)
			}
		})
	}
}

func TestSearchStationResponse_XMLUnmarshal(t *testing.T) {
	xmlData := `
		<results deviceID="1004567890AA" source="PANDORA" sourceAccount="user123">
			<songs>
				<searchResult source="PANDORA" sourceAccount="user123" token="S10657777">
					<name>Old Church Choir</name>
					<artist>Zach Williams</artist>
					<logo>http://example.com/song.jpg</logo>
				</searchResult>
			</songs>
			<artists>
				<searchResult source="PANDORA" sourceAccount="user123" token="R324771">
					<name>Zach Williams</name>
					<logo>http://example.com/artist.jpg</logo>
				</searchResult>
			</artists>
			<stations>
				<searchResult source="PANDORA" sourceAccount="user123" token="R123456">
					<name>Classic Rock Station</name>
					<description>The best classic rock hits</description>
					<logo>http://example.com/station.jpg</logo>
				</searchResult>
			</stations>
		</results>`

	var response SearchStationResponse
	err := xml.Unmarshal([]byte(xmlData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	if response.DeviceID != "1004567890AA" {
		t.Errorf("Expected deviceID '1004567890AA', got %s", response.DeviceID)
	}
	if response.Source != "PANDORA" {
		t.Errorf("Expected source PANDORA, got %s", response.Source)
	}
	if len(response.Songs) != 1 {
		t.Errorf("Expected 1 song result, got %d", len(response.Songs))
	}
	if len(response.Artists) != 1 {
		t.Errorf("Expected 1 artist result, got %d", len(response.Artists))
	}
	if len(response.Stations) != 1 {
		t.Errorf("Expected 1 station result, got %d", len(response.Stations))
	}

	// Test song result
	song := response.Songs[0]
	if song.Name != "Old Church Choir" {
		t.Errorf("Expected song name 'Old Church Choir', got %s", song.Name)
	}
	if song.Artist != "Zach Williams" {
		t.Errorf("Expected artist 'Zach Williams', got %s", song.Artist)
	}
	if song.Token != "S10657777" {
		t.Errorf("Expected token 'S10657777', got %s", song.Token)
	}

	// Test artist result
	artist := response.Artists[0]
	if artist.Name != "Zach Williams" {
		t.Errorf("Expected artist name 'Zach Williams', got %s", artist.Name)
	}
	if !artist.IsArtist() {
		t.Error("Expected result to be identified as artist")
	}

	// Test station result
	station := response.Stations[0]
	if station.Name != "Classic Rock Station" {
		t.Errorf("Expected station name 'Classic Rock Station', got %s", station.Name)
	}
	if !station.IsStation() {
		t.Error("Expected result to be identified as station")
	}
}

func TestSearchStationResponse_HelperMethods(t *testing.T) {
	response := &SearchStationResponse{
		Songs: []SearchResult{
			{Name: "Song 1", Artist: "Artist 1", Token: "S1"},
			{Name: "Song 2", Artist: "Artist 2", Token: "S2"},
		},
		Artists: []SearchResult{
			{Name: "Artist 3", Token: "A1"},
		},
		Stations: []SearchResult{
			{Name: "Station 1", Description: "Great music", Token: "R1"},
		},
	}

	// Test GetAllResults
	allResults := response.GetAllResults()
	if len(allResults) != 4 {
		t.Errorf("Expected 4 total results, got %d", len(allResults))
	}

	// Test GetSongs
	songs := response.GetSongs()
	if len(songs) != 2 {
		t.Errorf("Expected 2 songs, got %d", len(songs))
	}

	// Test GetArtists
	artists := response.GetArtists()
	if len(artists) != 1 {
		t.Errorf("Expected 1 artist, got %d", len(artists))
	}

	// Test GetStations
	stations := response.GetStations()
	if len(stations) != 1 {
		t.Errorf("Expected 1 station, got %d", len(stations))
	}

	// Test HasResults
	if !response.HasResults() {
		t.Error("Expected response to have results")
	}

	// Test GetResultCount
	if response.GetResultCount() != 4 {
		t.Errorf("Expected result count 4, got %d", response.GetResultCount())
	}

	// Test IsEmpty
	if response.IsEmpty() {
		t.Error("Expected response not to be empty")
	}

	// Test empty response
	emptyResponse := &SearchStationResponse{}
	if !emptyResponse.IsEmpty() {
		t.Error("Expected empty response to be empty")
	}
	if emptyResponse.HasResults() {
		t.Error("Expected empty response to have no results")
	}
}

func TestSearchResult_HelperMethods(t *testing.T) {
	tests := []struct {
		name         string
		result       SearchResult
		expectedType string
		fullTitle    string
	}{
		{
			name:         "Song result",
			result:       SearchResult{Name: "Test Song", Artist: "Test Artist"},
			expectedType: "song",
			fullTitle:    "Test Song - Test Artist",
		},
		{
			name:         "Artist result",
			result:       SearchResult{Name: "Test Artist"},
			expectedType: "artist",
			fullTitle:    "Test Artist",
		},
		{
			name:         "Station result",
			result:       SearchResult{Name: "Test Station", Description: "Great station"},
			expectedType: "station",
			fullTitle:    "Test Station",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test GetDisplayName
			if tt.result.GetDisplayName() != tt.result.Name {
				t.Errorf("Expected display name %s, got %s", tt.result.Name, tt.result.GetDisplayName())
			}

			// Test type detection
			switch tt.expectedType {
			case "song":
				if !tt.result.IsSong() {
					t.Error("Expected result to be identified as song")
				}
				if tt.result.IsArtist() || tt.result.IsStation() {
					t.Error("Result incorrectly identified as artist or station")
				}
			case "artist":
				if !tt.result.IsArtist() {
					t.Error("Expected result to be identified as artist")
				}
				if tt.result.IsSong() || tt.result.IsStation() {
					t.Error("Result incorrectly identified as song or station")
				}
			case "station":
				if !tt.result.IsStation() {
					t.Error("Expected result to be identified as station")
				}
				if tt.result.IsSong() || tt.result.IsArtist() {
					t.Error("Result incorrectly identified as song or artist")
				}
			}

			// Test GetFullTitle
			if tt.result.GetFullTitle() != tt.fullTitle {
				t.Errorf("Expected full title %s, got %s", tt.fullTitle, tt.result.GetFullTitle())
			}
		})
	}
}

func TestSearchResult_GetArtworkURL(t *testing.T) {
	result := SearchResult{
		Name: "Test",
		Logo: "http://example.com/artwork.jpg",
	}

	if result.GetArtworkURL() != "http://example.com/artwork.jpg" {
		t.Errorf("Expected artwork URL 'http://example.com/artwork.jpg', got %s", result.GetArtworkURL())
	}

	// Test empty logo
	emptyResult := SearchResult{Name: "Test"}
	if emptyResult.GetArtworkURL() != "" {
		t.Errorf("Expected empty artwork URL, got %s", emptyResult.GetArtworkURL())
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
