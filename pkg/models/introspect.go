package models

import "encoding/xml"

// IntrospectRequest represents a request to get introspect data for a music service
type IntrospectRequest struct {
	XMLName       xml.Name `xml:"introspect"`
	Source        string   `xml:"source,attr"`
	SourceAccount string   `xml:"sourceAccount,attr,omitempty"`
}

// IntrospectResponse represents a generic introspect response
// The actual XML name will vary based on the source (e.g., spotifyAccountIntrospectResponse)
type IntrospectResponse struct {
	XMLName                          xml.Name               `xml:""`
	State                            string                 `xml:"state,attr,omitempty"`
	User                             string                 `xml:"user,attr,omitempty"`
	IsPlaying                        bool                   `xml:"isPlaying,attr,omitempty"`
	TokenLastChangedTimeSeconds      int64                  `xml:"tokenLastChangedTimeSeconds,attr,omitempty"`
	TokenLastChangedTimeMicroseconds int64                  `xml:"tokenLastChangedTimeMicroseconds,attr,omitempty"`
	ShuffleMode                      string                 `xml:"shuffleMode,attr,omitempty"`
	PlayStatusState                  string                 `xml:"playStatusState,attr,omitempty"`
	CurrentURI                       string                 `xml:"currentUri,attr,omitempty"`
	ReceivedPlaybackRequest          bool                   `xml:"receivedPlaybackRequest,attr,omitempty"`
	SubscriptionType                 string                 `xml:"subscriptionType,attr,omitempty"`
	CachedPlaybackRequest            *CachedPlaybackRequest `xml:"cachedPlaybackRequest,omitempty"`
	NowPlaying                       *IntrospectNowPlaying  `xml:"nowPlaying,omitempty"`
	ContentItemHistory               *ContentItemHistory    `xml:"contentItemHistory,omitempty"`
}

// SpotifyIntrospectResponse represents a Spotify-specific introspect response
type SpotifyIntrospectResponse struct {
	XMLName                          xml.Name               `xml:"spotifyAccountIntrospectResponse"`
	State                            string                 `xml:"state,attr"`
	User                             string                 `xml:"user,attr"`
	IsPlaying                        bool                   `xml:"isPlaying,attr"`
	TokenLastChangedTimeSeconds      int64                  `xml:"tokenLastChangedTimeSeconds,attr"`
	TokenLastChangedTimeMicroseconds int64                  `xml:"tokenLastChangedTimeMicroseconds,attr"`
	ShuffleMode                      string                 `xml:"shuffleMode,attr"`
	PlayStatusState                  string                 `xml:"playStatusState,attr"`
	CurrentURI                       string                 `xml:"currentUri,attr"`
	ReceivedPlaybackRequest          bool                   `xml:"receivedPlaybackRequest,attr"`
	SubscriptionType                 string                 `xml:"subscriptionType,attr"`
	CachedPlaybackRequest            *CachedPlaybackRequest `xml:"cachedPlaybackRequest"`
	NowPlaying                       *IntrospectNowPlaying  `xml:"nowPlaying"`
	ContentItemHistory               *ContentItemHistory    `xml:"contentItemHistory"`
}

// CachedPlaybackRequest represents cached playback request information
type CachedPlaybackRequest struct {
	XMLName xml.Name `xml:"cachedPlaybackRequest"`
	// Add fields as discovered from actual responses
}

// IntrospectNowPlaying represents now playing information in introspect response
type IntrospectNowPlaying struct {
	XMLName               xml.Name `xml:"nowPlaying"`
	SkipPreviousSupported bool     `xml:"skipPreviousSupported,attr"`
	SeekSupported         bool     `xml:"seekSupported,attr"`
	ResumeSupported       bool     `xml:"resumeSupported,attr"`
	CollectData           bool     `xml:"collectData,attr"`
}

// ContentItemHistory represents the content item history
type ContentItemHistory struct {
	XMLName xml.Name `xml:"contentItemHistory"`
	MaxSize int      `xml:"maxSize,attr"`
	// Add items as discovered from actual responses
}

// IntrospectState represents possible introspect states
type IntrospectState string

const (
	// IntrospectStateInactiveUnselected indicates the service is inactive and unselected
	IntrospectStateInactiveUnselected IntrospectState = "InactiveUnselected"
	// IntrospectStateActive indicates the service is active
	IntrospectStateActive IntrospectState = "Active"
	// IntrospectStateInactive indicates the service is inactive
	IntrospectStateInactive IntrospectState = "Inactive"
)

// ShuffleMode represents possible shuffle modes
type ShuffleMode string

const (
	// ShuffleModeOff indicates shuffle is disabled
	ShuffleModeOff ShuffleMode = "OFF"
	// ShuffleModeOn indicates shuffle is enabled
	ShuffleModeOn ShuffleMode = "ON"
)

// NewIntrospectRequest creates a new introspect request
func NewIntrospectRequest(source, sourceAccount string) *IntrospectRequest {
	return &IntrospectRequest{
		Source:        source,
		SourceAccount: sourceAccount,
	}
}

// GetState returns the introspect state as a typed value
func (ir *IntrospectResponse) GetState() IntrospectState {
	return IntrospectState(ir.State)
}

// GetShuffleMode returns the shuffle mode as a typed value
func (ir *IntrospectResponse) GetShuffleMode() ShuffleMode {
	return ShuffleMode(ir.ShuffleMode)
}

// IsActive returns true if the service is in an active state
func (ir *IntrospectResponse) IsActive() bool {
	return ir.GetState() == IntrospectStateActive
}

// IsInactive returns true if the service is in an inactive state
func (ir *IntrospectResponse) IsInactive() bool {
	state := ir.GetState()
	return state == IntrospectStateInactive || state == IntrospectStateInactiveUnselected
}

// HasUser returns true if a user is associated with the service
func (ir *IntrospectResponse) HasUser() bool {
	return ir.User != ""
}

// IsShuffleEnabled returns true if shuffle mode is enabled
func (ir *IntrospectResponse) IsShuffleEnabled() bool {
	return ir.GetShuffleMode() == ShuffleModeOn
}

// HasCurrentContent returns true if there is current content playing
func (ir *IntrospectResponse) HasCurrentContent() bool {
	return ir.CurrentURI != ""
}

// SupportsSkipPrevious returns true if the service supports skipping to previous track
func (ir *IntrospectResponse) SupportsSkipPrevious() bool {
	return ir.NowPlaying != nil && ir.NowPlaying.SkipPreviousSupported
}

// SupportsSeek returns true if the service supports seeking within tracks
func (ir *IntrospectResponse) SupportsSeek() bool {
	return ir.NowPlaying != nil && ir.NowPlaying.SeekSupported
}

// SupportsResume returns true if the service supports resuming playback
func (ir *IntrospectResponse) SupportsResume() bool {
	return ir.NowPlaying != nil && ir.NowPlaying.ResumeSupported
}

// CollectsData returns true if the service collects usage data
func (ir *IntrospectResponse) CollectsData() bool {
	return ir.NowPlaying != nil && ir.NowPlaying.CollectData
}

// GetMaxHistorySize returns the maximum size of the content item history
func (ir *IntrospectResponse) GetMaxHistorySize() int {
	if ir.ContentItemHistory != nil {
		return ir.ContentItemHistory.MaxSize
	}

	return 0
}

// HasSubscription returns true if the user has a subscription
func (ir *IntrospectResponse) HasSubscription() bool {
	return ir.SubscriptionType != ""
}

// GetTokenAge returns the age of the token in seconds since last change
func (ir *IntrospectResponse) GetTokenAge() int64 {
	// This would need current time to calculate actual age
	// For now, just return the timestamp
	return ir.TokenLastChangedTimeSeconds
}

// Spotify-specific methods for SpotifyIntrospectResponse

// GetState returns the introspect state as a typed value
func (sir *SpotifyIntrospectResponse) GetState() IntrospectState {
	return IntrospectState(sir.State)
}

// GetShuffleMode returns the shuffle mode as a typed value
func (sir *SpotifyIntrospectResponse) GetShuffleMode() ShuffleMode {
	return ShuffleMode(sir.ShuffleMode)
}

// IsActive returns true if the service is in an active state
func (sir *SpotifyIntrospectResponse) IsActive() bool {
	return sir.GetState() == IntrospectStateActive
}

// IsInactive returns true if the service is in an inactive state
func (sir *SpotifyIntrospectResponse) IsInactive() bool {
	state := sir.GetState()
	return state == IntrospectStateInactive || state == IntrospectStateInactiveUnselected
}

// HasUser returns true if a user is associated with the service
func (sir *SpotifyIntrospectResponse) HasUser() bool {
	return sir.User != ""
}

// IsShuffleEnabled returns true if shuffle mode is enabled
func (sir *SpotifyIntrospectResponse) IsShuffleEnabled() bool {
	return sir.GetShuffleMode() == ShuffleModeOn
}

// HasCurrentContent returns true if there is current content playing
func (sir *SpotifyIntrospectResponse) HasCurrentContent() bool {
	return sir.CurrentURI != ""
}

// SupportsSkipPrevious returns true if the service supports skipping to previous track
func (sir *SpotifyIntrospectResponse) SupportsSkipPrevious() bool {
	return sir.NowPlaying != nil && sir.NowPlaying.SkipPreviousSupported
}

// SupportsSeek returns true if the service supports seeking within tracks
func (sir *SpotifyIntrospectResponse) SupportsSeek() bool {
	return sir.NowPlaying != nil && sir.NowPlaying.SeekSupported
}

// SupportsResume returns true if the service supports resuming playback
func (sir *SpotifyIntrospectResponse) SupportsResume() bool {
	return sir.NowPlaying != nil && sir.NowPlaying.ResumeSupported
}

// CollectsData returns true if the service collects usage data
func (sir *SpotifyIntrospectResponse) CollectsData() bool {
	return sir.NowPlaying != nil && sir.NowPlaying.CollectData
}

// GetMaxHistorySize returns the maximum size of the content item history
func (sir *SpotifyIntrospectResponse) GetMaxHistorySize() int {
	if sir.ContentItemHistory != nil {
		return sir.ContentItemHistory.MaxSize
	}

	return 0
}

// HasSubscription returns true if the user has a subscription
func (sir *SpotifyIntrospectResponse) HasSubscription() bool {
	return sir.SubscriptionType != ""
}

// GetTokenAge returns the age of the token in seconds since last change
func (sir *SpotifyIntrospectResponse) GetTokenAge() int64 {
	return sir.TokenLastChangedTimeSeconds
}
