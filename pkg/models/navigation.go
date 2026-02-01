package models

import "encoding/xml"

// NavigateRequest represents a request to navigate content sources
type NavigateRequest struct {
	XMLName       xml.Name      `xml:"navigate"`
	Source        string        `xml:"source,attr"`
	SourceAccount string        `xml:"sourceAccount,attr,omitempty"`
	Menu          string        `xml:"menu,attr,omitempty"`
	Sort          string        `xml:"sort,attr,omitempty"`
	StartItem     int           `xml:"startItem"`
	NumItems      int           `xml:"numItems"`
	Item          *NavigateItem `xml:"item,omitempty"`
}

// NavigateResponse represents the response from a navigate request
type NavigateResponse struct {
	XMLName       xml.Name       `xml:"navigateResponse"`
	Source        string         `xml:"source,attr"`
	SourceAccount string         `xml:"sourceAccount,attr,omitempty"`
	TotalItems    int            `xml:"totalItems"`
	Items         []NavigateItem `xml:"items>item"`
}

// NavigateItem represents a single item in a navigate response
type NavigateItem struct {
	XMLName            xml.Name            `xml:"item"`
	Playable           int                 `xml:"Playable,attr,omitempty"`
	Name               string              `xml:"name"`
	Type               string              `xml:"type"`
	ContentItem        *ContentItem        `xml:"ContentItem,omitempty"`
	MediaItemContainer *MediaItemContainer `xml:"mediaItemContainer,omitempty"`
	ArtistName         string              `xml:"artistName,omitempty"`
	AlbumName          string              `xml:"albumName,omitempty"`
}

// MediaItemContainer represents a media container within a navigate item
type MediaItemContainer struct {
	XMLName     xml.Name     `xml:"mediaItemContainer"`
	Offset      int          `xml:"offset,attr"`
	ContentItem *ContentItem `xml:"ContentItem,omitempty"`
}

// AddStationRequest represents a request to add a station
type AddStationRequest struct {
	XMLName       xml.Name `xml:"addStation"`
	Source        string   `xml:"source,attr"`
	SourceAccount string   `xml:"sourceAccount,attr,omitempty"`
	Token         string   `xml:"token,attr,omitempty"`
	Name          string   `xml:"name"`
}

// RemoveStationRequest represents a request to remove a station
// Note: For removeStation, we send the ContentItem directly, not wrapped
type RemoveStationRequest = ContentItem

// StationResponse represents the response from add/remove station operations
type StationResponse struct {
	XMLName xml.Name `xml:"status"`
	Status  string   `xml:",chardata"`
}

// SearchStationRequest represents a request to search for stations
type SearchStationRequest struct {
	XMLName       xml.Name `xml:"search"`
	Source        string   `xml:"source,attr"`
	SourceAccount string   `xml:"sourceAccount,attr,omitempty"`
	SearchTerm    string   `xml:",chardata"`
}

// SearchStationResponse represents the response from a station search
type SearchStationResponse struct {
	XMLName       xml.Name       `xml:"results"`
	DeviceID      string         `xml:"deviceID,attr"`
	Source        string         `xml:"source,attr"`
	SourceAccount string         `xml:"sourceAccount,attr,omitempty"`
	Songs         []SearchResult `xml:"songs>searchResult"`
	Artists       []SearchResult `xml:"artists>searchResult"`
	Stations      []SearchResult `xml:"stations>searchResult"`
}

// SearchResult represents a single search result (song, artist, or station)
type SearchResult struct {
	XMLName       xml.Name `xml:"searchResult"`
	Source        string   `xml:"source,attr"`
	SourceAccount string   `xml:"sourceAccount,attr,omitempty"`
	Token         string   `xml:"token,attr"`
	Name          string   `xml:"name"`
	Artist        string   `xml:"artist,omitempty"`
	Album         string   `xml:"album,omitempty"`
	Logo          string   `xml:"logo,omitempty"`
	Description   string   `xml:"description,omitempty"`
}

// NewNavigateRequest creates a new navigate request for browsing content
func NewNavigateRequest(source, sourceAccount string, startItem, numItems int) *NavigateRequest {
	return &NavigateRequest{
		Source:        source,
		SourceAccount: sourceAccount,
		StartItem:     startItem,
		NumItems:      numItems,
	}
}

// NewNavigateRequestWithMenu creates a navigate request with menu and sort parameters
func NewNavigateRequestWithMenu(source, sourceAccount, menu, sort string, startItem, numItems int) *NavigateRequest {
	return &NavigateRequest{
		Source:        source,
		SourceAccount: sourceAccount,
		Menu:          menu,
		Sort:          sort,
		StartItem:     startItem,
		NumItems:      numItems,
	}
}

// NewNavigateRequestWithItem creates a navigate request to browse a specific container item
func NewNavigateRequestWithItem(source, sourceAccount string, startItem, numItems int, item *ContentItem) *NavigateRequest {
	navigateItem := &NavigateItem{
		Playable:    1,
		Name:        item.ItemName,
		Type:        "dir",
		ContentItem: item,
	}

	return &NavigateRequest{
		Source:        source,
		SourceAccount: sourceAccount,
		StartItem:     startItem,
		NumItems:      numItems,
		Item:          navigateItem,
	}
}

// NewAddStationRequest creates a new add station request
func NewAddStationRequest(source, sourceAccount, token, name string) *AddStationRequest {
	return &AddStationRequest{
		Source:        source,
		SourceAccount: sourceAccount,
		Token:         token,
		Name:          name,
	}
}

// NewRemoveStationRequest creates a new remove station request
func NewRemoveStationRequest(contentItem *ContentItem) *RemoveStationRequest {
	return contentItem
}

// NewSearchStationRequest creates a new search station request
func NewSearchStationRequest(source, sourceAccount, searchTerm string) *SearchStationRequest {
	return &SearchStationRequest{
		Source:        source,
		SourceAccount: sourceAccount,
		SearchTerm:    searchTerm,
	}
}

// GetPlayableItems returns only the playable items from the navigate response
func (nr *NavigateResponse) GetPlayableItems() []NavigateItem {
	var playable []NavigateItem

	for _, item := range nr.Items {
		if item.Playable == 1 {
			playable = append(playable, item)
		}
	}

	return playable
}

// GetDirectories returns only the directory items from the navigate response
func (nr *NavigateResponse) GetDirectories() []NavigateItem {
	var directories []NavigateItem

	for _, item := range nr.Items {
		if item.Type == "dir" {
			directories = append(directories, item)
		}
	}

	return directories
}

// GetTracks returns only the track items from the navigate response
func (nr *NavigateResponse) GetTracks() []NavigateItem {
	var tracks []NavigateItem

	for _, item := range nr.Items {
		if item.Type == "track" {
			tracks = append(tracks, item)
		}
	}

	return tracks
}

// GetStations returns only the station items from the navigate response
func (nr *NavigateResponse) GetStations() []NavigateItem {
	var stations []NavigateItem

	for _, item := range nr.Items {
		if item.Type == "stationurl" || (item.ContentItem != nil && item.ContentItem.Type == "stationurl") {
			stations = append(stations, item)
		}
	}

	return stations
}

// IsEmpty returns true if the navigate response contains no items
func (nr *NavigateResponse) IsEmpty() bool {
	return nr.TotalItems == 0 || len(nr.Items) == 0
}

// GetDisplayName returns the display name for a navigate item
func (ni *NavigateItem) GetDisplayName() string {
	if ni.Name != "" {
		return ni.Name
	}

	if ni.ContentItem != nil && ni.ContentItem.ItemName != "" {
		return ni.ContentItem.ItemName
	}

	return "Unknown Item"
}

// IsPlayable returns true if the navigate item can be played directly
func (ni *NavigateItem) IsPlayable() bool {
	return ni.Playable == 1
}

// IsDirectory returns true if the navigate item is a directory/container
func (ni *NavigateItem) IsDirectory() bool {
	return ni.Type == "dir"
}

// IsTrack returns true if the navigate item is a track
func (ni *NavigateItem) IsTrack() bool {
	return ni.Type == "track"
}

// IsStation returns true if the navigate item is a radio station
func (ni *NavigateItem) IsStation() bool {
	return ni.Type == "stationurl" || (ni.ContentItem != nil && ni.ContentItem.Type == "stationurl")
}

// GetContentItem returns the ContentItem for this navigate item
func (ni *NavigateItem) GetContentItem() *ContentItem {
	return ni.ContentItem
}

// GetArtwork returns the artwork URL if available
func (ni *NavigateItem) GetArtwork() string {
	if ni.ContentItem != nil && ni.ContentItem.ContainerArt != "" {
		return ni.ContentItem.ContainerArt
	}

	return ""
}

// GetAllResults returns all search results regardless of type
func (sr *SearchStationResponse) GetAllResults() []SearchResult {
	var allResults []SearchResult

	allResults = append(allResults, sr.Songs...)
	allResults = append(allResults, sr.Artists...)
	allResults = append(allResults, sr.Stations...)
	return allResults
}

// GetSongs returns only song results
func (sr *SearchStationResponse) GetSongs() []SearchResult {
	return sr.Songs
}

// GetArtists returns only artist results
func (sr *SearchStationResponse) GetArtists() []SearchResult {
	return sr.Artists
}

// GetStations returns only station results
func (sr *SearchStationResponse) GetStations() []SearchResult {
	return sr.Stations
}

// HasResults returns true if the search response contains any results
func (sr *SearchStationResponse) HasResults() bool {
	return len(sr.Songs) > 0 || len(sr.Artists) > 0 || len(sr.Stations) > 0
}

// GetResultCount returns the total number of results
func (sr *SearchStationResponse) GetResultCount() int {
	return len(sr.Songs) + len(sr.Artists) + len(sr.Stations)
}

// IsEmpty returns true if the search response contains no results
func (sr *SearchStationResponse) IsEmpty() bool {
	return !sr.HasResults()
}

// GetDisplayName returns the display name for a search result
func (sr *SearchResult) GetDisplayName() string {
	if sr.Name != "" {
		return sr.Name
	}
	return "Unknown"
}

// GetArtworkURL returns the logo/artwork URL if available
func (sr *SearchResult) GetArtworkURL() string {
	return sr.Logo
}

// IsArtist returns true if this is an artist result
func (sr *SearchResult) IsArtist() bool {
	// Artist result: no artist field (name is the artist) and no description
	return sr.Artist == "" && sr.Album == "" && sr.Description == ""
}

// IsSong returns true if this is a song result
func (sr *SearchResult) IsSong() bool {
	// Song result: has artist field populated
	return sr.Artist != ""
}

// IsStation returns true if this is a station result
func (sr *SearchResult) IsStation() bool {
	// Station result: has description or is neither artist nor song
	return sr.Description != "" || (!sr.IsSong() && !sr.IsArtist())
}

// GetFullTitle returns a full title with artist if available
func (sr *SearchResult) GetFullTitle() string {
	if sr.Artist != "" {
		return sr.Name + " - " + sr.Artist
	}
	return sr.Name
}
