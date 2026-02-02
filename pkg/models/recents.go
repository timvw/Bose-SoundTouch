package models

import "encoding/xml"

// RecentsResponse represents the response from the /recents endpoint
type RecentsResponse struct {
	XMLName xml.Name              `xml:"recents"`
	Items   []RecentsResponseItem `xml:"recent"`
}

// RecentsResponseItem represents a recently played item from the /recents API endpoint
type RecentsResponseItem struct {
	XMLName     xml.Name     `xml:"recent"`
	DeviceID    string       `xml:"deviceID,attr"`
	UTCTime     int64        `xml:"utcTime,attr"`
	ID          string       `xml:"id,attr,omitempty"`
	ContentItem *ContentItem `xml:"contentItem"`
}

// GetItemCount returns the number of recent items
func (r *RecentsResponse) GetItemCount() int {
	return len(r.Items)
}

// IsEmpty returns true if there are no recent items
func (r *RecentsResponse) IsEmpty() bool {
	return len(r.Items) == 0
}

// GetMostRecent returns the most recently played item (first in the list)
func (r *RecentsResponse) GetMostRecent() *RecentsResponseItem {
	if len(r.Items) == 0 {
		return nil
	}

	return &r.Items[0]
}

// GetItemsBySource returns recent items filtered by source type
func (r *RecentsResponse) GetItemsBySource(source string) []RecentsResponseItem {
	var filtered []RecentsResponseItem

	for _, item := range r.Items {
		if item.ContentItem != nil && item.ContentItem.Source == source {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// GetSpotifyItems returns only Spotify recent items
func (r *RecentsResponse) GetSpotifyItems() []RecentsResponseItem {
	return r.GetItemsBySource("SPOTIFY")
}

// GetLocalMusicItems returns only local music recent items
func (r *RecentsResponse) GetLocalMusicItems() []RecentsResponseItem {
	return r.GetItemsBySource("LOCAL_MUSIC")
}

// GetStoredMusicItems returns only stored music recent items
func (r *RecentsResponse) GetStoredMusicItems() []RecentsResponseItem {
	return r.GetItemsBySource("STORED_MUSIC")
}

// GetTuneInItems returns only TuneIn radio recent items
func (r *RecentsResponse) GetTuneInItems() []RecentsResponseItem {
	return r.GetItemsBySource("TUNEIN")
}

// GetPandoraItems returns only Pandora recent items
func (r *RecentsResponse) GetPandoraItems() []RecentsResponseItem {
	return r.GetItemsBySource("PANDORA")
}

// GetPresetableItems returns recent items that can be saved as presets
func (r *RecentsResponse) GetPresetableItems() []RecentsResponseItem {
	var presetable []RecentsResponseItem

	for _, item := range r.Items {
		if item.ContentItem != nil && item.ContentItem.IsPresetable {
			presetable = append(presetable, item)
		}
	}

	return presetable
}

// GetItemsByType returns recent items filtered by content type
func (r *RecentsResponse) GetItemsByType(contentType string) []RecentsResponseItem {
	var filtered []RecentsResponseItem

	for _, item := range r.Items {
		if item.ContentItem != nil && item.ContentItem.Type == contentType {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// GetTracks returns only track-type recent items
func (r *RecentsResponse) GetTracks() []RecentsResponseItem {
	return r.GetItemsByType("track")
}

// GetStations returns only station-type recent items
func (r *RecentsResponse) GetStations() []RecentsResponseItem {
	return r.GetItemsByType("stationurl")
}

// GetPlaylistsAndAlbums returns playlist and album-type recent items
func (r *RecentsResponse) GetPlaylistsAndAlbums() []RecentsResponseItem {
	var items []RecentsResponseItem

	for _, item := range r.Items {
		if item.ContentItem != nil {
			contentType := item.ContentItem.Type
			if contentType == "playlist" || contentType == "album" || contentType == "container" {
				items = append(items, item)
			}
		}
	}

	return items
}

// HasContent returns true if the recent item has content information
func (ri *RecentsResponseItem) HasContent() bool {
	return ri.ContentItem != nil
}

// GetDisplayName returns the display name for the recent item
func (ri *RecentsResponseItem) GetDisplayName() string {
	if ri.ContentItem != nil && ri.ContentItem.ItemName != "" {
		return ri.ContentItem.ItemName
	}

	return "Unknown Item"
}

// GetSource returns the content source
func (ri *RecentsResponseItem) GetSource() string {
	if ri.ContentItem != nil {
		return ri.ContentItem.Source
	}

	return ""
}

// GetSourceAccount returns the source account
func (ri *RecentsResponseItem) GetSourceAccount() string {
	if ri.ContentItem != nil {
		return ri.ContentItem.SourceAccount
	}

	return ""
}

// GetLocation returns the content location
func (ri *RecentsResponseItem) GetLocation() string {
	if ri.ContentItem != nil {
		return ri.ContentItem.Location
	}

	return ""
}

// GetContentType returns the content type
func (ri *RecentsResponseItem) GetContentType() string {
	if ri.ContentItem != nil {
		return ri.ContentItem.Type
	}

	return ""
}

// IsPresetable returns true if the item can be saved as a preset
func (ri *RecentsResponseItem) IsPresetable() bool {
	return ri.ContentItem != nil && ri.ContentItem.IsPresetable
}

// IsTrack returns true if the recent item is a track
func (ri *RecentsResponseItem) IsTrack() bool {
	return ri.GetContentType() == "track"
}

// IsStation returns true if the recent item is a radio station
func (ri *RecentsResponseItem) IsStation() bool {
	return ri.GetContentType() == "stationurl"
}

// IsPlaylist returns true if the recent item is a playlist
func (ri *RecentsResponseItem) IsPlaylist() bool {
	return ri.GetContentType() == "playlist"
}

// IsAlbum returns true if the recent item is an album
func (ri *RecentsResponseItem) IsAlbum() bool {
	return ri.GetContentType() == "album"
}

// IsContainer returns true if the recent item is a container (folder/collection)
func (ri *RecentsResponseItem) IsContainer() bool {
	contentType := ri.GetContentType()
	return contentType == "container" || contentType == "dir"
}

// IsSpotifyContent returns true if the recent item is from Spotify
func (ri *RecentsResponseItem) IsSpotifyContent() bool {
	return ri.GetSource() == "SPOTIFY"
}

// IsLocalContent returns true if the recent item is from local sources
func (ri *RecentsResponseItem) IsLocalContent() bool {
	source := ri.GetSource()
	return source == "LOCAL_MUSIC" || source == "STORED_MUSIC"
}

// IsStreamingContent returns true if the recent item is from streaming services
func (ri *RecentsResponseItem) IsStreamingContent() bool {
	source := ri.GetSource()

	return source == "SPOTIFY" || source == "PANDORA" || source == "TUNEIN" ||
		source == "AMAZON" || source == "DEEZER" || source == "IHEART"
}

// GetArtwork returns the artwork URL if available
func (ri *RecentsResponseItem) GetArtwork() string {
	if ri.ContentItem != nil {
		return ri.ContentItem.ContainerArt
	}

	return ""
}

// HasArtwork returns true if artwork is available
func (ri *RecentsResponseItem) HasArtwork() bool {
	return ri.GetArtwork() != ""
}

// GetUTCTime returns the UTC timestamp when the item was played
func (ri *RecentsResponseItem) GetUTCTime() int64 {
	return ri.UTCTime
}

// HasID returns true if the recent item has an ID
func (ri *RecentsResponseItem) HasID() bool {
	return ri.ID != ""
}

// GetID returns the recent item ID
func (ri *RecentsResponseItem) GetID() string {
	return ri.ID
}
