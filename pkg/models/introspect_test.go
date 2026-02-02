package models

import (
	"encoding/xml"
	"testing"
)

func TestIntrospectRequest_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		request  *IntrospectRequest
		expected string
	}{
		{
			name: "with source account",
			request: &IntrospectRequest{
				Source:        "SPOTIFY",
				SourceAccount: "SpotifyConnectUserName",
			},
			expected: `<introspect source="SPOTIFY" sourceAccount="SpotifyConnectUserName"></introspect>`,
		},
		{
			name: "without source account",
			request: &IntrospectRequest{
				Source: "BLUETOOTH",
			},
			expected: `<introspect source="BLUETOOTH"></introspect>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := xml.Marshal(tt.request)
			if err != nil {
				t.Fatalf("failed to marshal request: %v", err)
			}

			if string(data) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(data))
			}
		})
	}
}

func TestIntrospectResponse_Unmarshal(t *testing.T) {
	tests := []struct {
		name        string
		xmlData     string
		expected    *IntrospectResponse
		expectError bool
	}{
		{
			name: "spotify introspect response",
			xmlData: `<spotifyAccountIntrospectResponse state="InactiveUnselected" user="SpotifyConnectUserName" isPlaying="false" tokenLastChangedTimeSeconds="1702566495" tokenLastChangedTimeMicroseconds="427884" shuffleMode="OFF" playStatusState="2" currentUri="" receivedPlaybackRequest="false" subscriptionType="">
  <cachedPlaybackRequest />
  <nowPlaying skipPreviousSupported="false" seekSupported="false" resumeSupported="true" collectData="true" />
  <contentItemHistory maxSize="10" />
</spotifyAccountIntrospectResponse>`,
			expected: &IntrospectResponse{
				State:                            "InactiveUnselected",
				User:                             "SpotifyConnectUserName",
				IsPlaying:                        false,
				TokenLastChangedTimeSeconds:      1702566495,
				TokenLastChangedTimeMicroseconds: 427884,
				ShuffleMode:                      "OFF",
				PlayStatusState:                  "2",
				CurrentURI:                       "",
				ReceivedPlaybackRequest:          false,
				SubscriptionType:                 "",
				CachedPlaybackRequest:            &CachedPlaybackRequest{},
				NowPlaying: &IntrospectNowPlaying{
					SkipPreviousSupported: false,
					SeekSupported:         false,
					ResumeSupported:       true,
					CollectData:           true,
				},
				ContentItemHistory: &ContentItemHistory{
					MaxSize: 10,
				},
			},
		},
		{
			name: "pandora introspect response",
			xmlData: `<pandoraAccountIntrospectResponse state="Active" user="pandora_user" isPlaying="true" shuffleMode="ON" currentUri="pandora://track/123" subscriptionType="Premium">
  <nowPlaying skipPreviousSupported="true" seekSupported="false" resumeSupported="true" collectData="false" />
  <contentItemHistory maxSize="20" />
</pandoraAccountIntrospectResponse>`,
			expected: &IntrospectResponse{
				State:            "Active",
				User:             "pandora_user",
				IsPlaying:        true,
				ShuffleMode:      "ON",
				CurrentURI:       "pandora://track/123",
				SubscriptionType: "Premium",
				NowPlaying: &IntrospectNowPlaying{
					SkipPreviousSupported: true,
					SeekSupported:         false,
					ResumeSupported:       true,
					CollectData:           false,
				},
				ContentItemHistory: &ContentItemHistory{
					MaxSize: 20,
				},
			},
		},
		{
			name: "minimal response",
			xmlData: `<serviceIntrospectResponse state="Inactive">
</serviceIntrospectResponse>`,
			expected: &IntrospectResponse{
				State: "Inactive",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var response IntrospectResponse

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

			// Compare basic fields
			if response.State != tt.expected.State {
				t.Errorf("expected state %s, got %s", tt.expected.State, response.State)
			}

			if response.User != tt.expected.User {
				t.Errorf("expected user %s, got %s", tt.expected.User, response.User)
			}

			if response.IsPlaying != tt.expected.IsPlaying {
				t.Errorf("expected isPlaying %t, got %t", tt.expected.IsPlaying, response.IsPlaying)
			}

			if response.TokenLastChangedTimeSeconds != tt.expected.TokenLastChangedTimeSeconds {
				t.Errorf("expected tokenLastChangedTimeSeconds %d, got %d",
					tt.expected.TokenLastChangedTimeSeconds, response.TokenLastChangedTimeSeconds)
			}

			if response.TokenLastChangedTimeMicroseconds != tt.expected.TokenLastChangedTimeMicroseconds {
				t.Errorf("expected tokenLastChangedTimeMicroseconds %d, got %d",
					tt.expected.TokenLastChangedTimeMicroseconds, response.TokenLastChangedTimeMicroseconds)
			}

			if response.ShuffleMode != tt.expected.ShuffleMode {
				t.Errorf("expected shuffleMode %s, got %s", tt.expected.ShuffleMode, response.ShuffleMode)
			}

			if response.PlayStatusState != tt.expected.PlayStatusState {
				t.Errorf("expected playStatusState %s, got %s", tt.expected.PlayStatusState, response.PlayStatusState)
			}

			if response.CurrentURI != tt.expected.CurrentURI {
				t.Errorf("expected currentUri %s, got %s", tt.expected.CurrentURI, response.CurrentURI)
			}

			if response.ReceivedPlaybackRequest != tt.expected.ReceivedPlaybackRequest {
				t.Errorf("expected receivedPlaybackRequest %t, got %t",
					tt.expected.ReceivedPlaybackRequest, response.ReceivedPlaybackRequest)
			}

			if response.SubscriptionType != tt.expected.SubscriptionType {
				t.Errorf("expected subscriptionType %s, got %s", tt.expected.SubscriptionType, response.SubscriptionType)
			}

			// Compare nested structures
			if tt.expected.CachedPlaybackRequest != nil {
				if response.CachedPlaybackRequest == nil {
					t.Error("expected cachedPlaybackRequest, got nil")
				}
			} else if response.CachedPlaybackRequest != nil {
				t.Error("expected cachedPlaybackRequest to be nil, got non-nil")
			}

			if tt.expected.NowPlaying != nil {
				if response.NowPlaying == nil {
					t.Error("expected nowPlaying, got nil")
				} else {
					if response.NowPlaying.SkipPreviousSupported != tt.expected.NowPlaying.SkipPreviousSupported {
						t.Errorf("expected skipPreviousSupported %t, got %t",
							tt.expected.NowPlaying.SkipPreviousSupported,
							response.NowPlaying.SkipPreviousSupported)
					}

					if response.NowPlaying.SeekSupported != tt.expected.NowPlaying.SeekSupported {
						t.Errorf("expected seekSupported %t, got %t",
							tt.expected.NowPlaying.SeekSupported,
							response.NowPlaying.SeekSupported)
					}

					if response.NowPlaying.ResumeSupported != tt.expected.NowPlaying.ResumeSupported {
						t.Errorf("expected resumeSupported %t, got %t",
							tt.expected.NowPlaying.ResumeSupported,
							response.NowPlaying.ResumeSupported)
					}

					if response.NowPlaying.CollectData != tt.expected.NowPlaying.CollectData {
						t.Errorf("expected collectData %t, got %t",
							tt.expected.NowPlaying.CollectData,
							response.NowPlaying.CollectData)
					}
				}
			} else if response.NowPlaying != nil {
				t.Error("expected nowPlaying to be nil, got non-nil")
			}

			if tt.expected.ContentItemHistory != nil {
				if response.ContentItemHistory == nil {
					t.Error("expected contentItemHistory, got nil")
				} else {
					if response.ContentItemHistory.MaxSize != tt.expected.ContentItemHistory.MaxSize {
						t.Errorf("expected maxSize %d, got %d",
							tt.expected.ContentItemHistory.MaxSize,
							response.ContentItemHistory.MaxSize)
					}
				}
			} else if response.ContentItemHistory != nil {
				t.Error("expected contentItemHistory to be nil, got non-nil")
			}
		})
	}
}

func TestSpotifyIntrospectResponse_Unmarshal(t *testing.T) {
	xmlData := `<spotifyAccountIntrospectResponse state="InactiveUnselected" user="SpotifyConnectUserName" isPlaying="false" tokenLastChangedTimeSeconds="1702566495" tokenLastChangedTimeMicroseconds="427884" shuffleMode="OFF" playStatusState="2" currentUri="" receivedPlaybackRequest="false" subscriptionType="">
  <cachedPlaybackRequest />
  <nowPlaying skipPreviousSupported="false" seekSupported="false" resumeSupported="true" collectData="true" />
  <contentItemHistory maxSize="10" />
</spotifyAccountIntrospectResponse>`

	var response SpotifyIntrospectResponse

	err := xml.Unmarshal([]byte(xmlData), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal spotify response: %v", err)
	}

	if response.State != "InactiveUnselected" {
		t.Errorf("expected state InactiveUnselected, got %s", response.State)
	}

	if response.User != "SpotifyConnectUserName" {
		t.Errorf("expected user SpotifyConnectUserName, got %s", response.User)
	}

	if response.IsPlaying != false {
		t.Errorf("expected isPlaying false, got %t", response.IsPlaying)
	}

	if response.TokenLastChangedTimeSeconds != 1702566495 {
		t.Errorf("expected tokenLastChangedTimeSeconds 1702566495, got %d", response.TokenLastChangedTimeSeconds)
	}

	if response.ShuffleMode != "OFF" {
		t.Errorf("expected shuffleMode OFF, got %s", response.ShuffleMode)
	}
}

func TestIntrospectState_Constants(t *testing.T) {
	tests := []struct {
		name     string
		state    IntrospectState
		expected string
	}{
		{"InactiveUnselected", IntrospectStateInactiveUnselected, "InactiveUnselected"},
		{"Active", IntrospectStateActive, "Active"},
		{"Inactive", IntrospectStateInactive, "Inactive"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.state) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(tt.state))
			}
		})
	}
}

func TestShuffleMode_Constants(t *testing.T) {
	tests := []struct {
		name     string
		mode     ShuffleMode
		expected string
	}{
		{"Off", ShuffleModeOff, "OFF"},
		{"On", ShuffleModeOn, "ON"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.mode) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(tt.mode))
			}
		})
	}
}

func TestNewIntrospectRequest(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		sourceAccount string
	}{
		{
			name:          "with source account",
			source:        "SPOTIFY",
			sourceAccount: "test_user",
		},
		{
			name:          "without source account",
			source:        "BLUETOOTH",
			sourceAccount: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := NewIntrospectRequest(tt.source, tt.sourceAccount)

			if request == nil {
				t.Error("expected request, got nil")
				return
			}

			if request.Source != tt.source {
				t.Errorf("expected source %s, got %s", tt.source, request.Source)
			}

			if request.SourceAccount != tt.sourceAccount {
				t.Errorf("expected sourceAccount %s, got %s", tt.sourceAccount, request.SourceAccount)
			}
		})
	}
}

func TestIntrospectResponse_Methods(t *testing.T) {
	tests := []struct {
		name     string
		response *IntrospectResponse
		testFunc func(t *testing.T, r *IntrospectResponse)
	}{
		{
			name: "active spotify response",
			response: &IntrospectResponse{
				State:            "Active",
				User:             "test_user",
				IsPlaying:        true,
				ShuffleMode:      "ON",
				CurrentURI:       "spotify://track/123",
				SubscriptionType: "Premium",
				NowPlaying: &IntrospectNowPlaying{
					SkipPreviousSupported: true,
					SeekSupported:         true,
					ResumeSupported:       true,
					CollectData:           false,
				},
				ContentItemHistory: &ContentItemHistory{
					MaxSize: 15,
				},
			},
			testFunc: func(t *testing.T, r *IntrospectResponse) {
				t.Helper()

				if !r.IsActive() {
					t.Error("expected IsActive() to return true")
				}

				if r.IsInactive() {
					t.Error("expected IsInactive() to return false")
				}

				if !r.HasUser() {
					t.Error("expected HasUser() to return true")
				}

				if !r.IsShuffleEnabled() {
					t.Error("expected IsShuffleEnabled() to return true")
				}

				if !r.HasCurrentContent() {
					t.Error("expected HasCurrentContent() to return true")
				}

				if !r.SupportsSkipPrevious() {
					t.Error("expected SupportsSkipPrevious() to return true")
				}

				if !r.SupportsSeek() {
					t.Error("expected SupportsSeek() to return true")
				}

				if !r.SupportsResume() {
					t.Error("expected SupportsResume() to return true")
				}

				if r.CollectsData() {
					t.Error("expected CollectsData() to return false")
				}

				if r.GetMaxHistorySize() != 15 {
					t.Errorf("expected GetMaxHistorySize() to return 15, got %d", r.GetMaxHistorySize())
				}

				if !r.HasSubscription() {
					t.Error("expected HasSubscription() to return true")
				}
			},
		},
		{
			name: "inactive response",
			response: &IntrospectResponse{
				State:            "InactiveUnselected",
				User:             "",
				IsPlaying:        false,
				ShuffleMode:      "OFF",
				CurrentURI:       "",
				SubscriptionType: "",
			},
			testFunc: func(t *testing.T, r *IntrospectResponse) {
				t.Helper()

				if r.IsActive() {
					t.Error("expected IsActive() to return false")
				}

				if !r.IsInactive() {
					t.Error("expected IsInactive() to return true")
				}

				if r.HasUser() {
					t.Error("expected HasUser() to return false")
				}

				if r.IsShuffleEnabled() {
					t.Error("expected IsShuffleEnabled() to return false")
				}

				if r.HasCurrentContent() {
					t.Error("expected HasCurrentContent() to return false")
				}

				if r.HasSubscription() {
					t.Error("expected HasSubscription() to return false")
				}

				if r.GetMaxHistorySize() != 0 {
					t.Errorf("expected GetMaxHistorySize() to return 0, got %d", r.GetMaxHistorySize())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, tt.response)
		})
	}
}

func TestSpotifyIntrospectResponse_Methods(t *testing.T) {
	response := &SpotifyIntrospectResponse{
		State:            "Active",
		User:             "test_user",
		IsPlaying:        true,
		ShuffleMode:      "ON",
		CurrentURI:       "spotify://track/123",
		SubscriptionType: "Premium",
		NowPlaying: &IntrospectNowPlaying{
			SkipPreviousSupported: true,
			SeekSupported:         true,
			ResumeSupported:       true,
			CollectData:           false,
		},
		ContentItemHistory: &ContentItemHistory{
			MaxSize: 15,
		},
	}

	// Test that Spotify-specific response has same methods as generic response
	if !response.IsActive() {
		t.Error("expected IsActive() to return true")
	}

	if response.IsInactive() {
		t.Error("expected IsInactive() to return false")
	}

	if !response.HasUser() {
		t.Error("expected HasUser() to return true")
	}

	if !response.IsShuffleEnabled() {
		t.Error("expected IsShuffleEnabled() to return true")
	}

	if !response.HasCurrentContent() {
		t.Error("expected HasCurrentContent() to return true")
	}

	if !response.SupportsSkipPrevious() {
		t.Error("expected SupportsSkipPrevious() to return true")
	}

	if !response.SupportsSeek() {
		t.Error("expected SupportsSeek() to return true")
	}

	if !response.SupportsResume() {
		t.Error("expected SupportsResume() to return true")
	}

	if response.CollectsData() {
		t.Error("expected CollectsData() to return false")
	}

	if response.GetMaxHistorySize() != 15 {
		t.Errorf("expected GetMaxHistorySize() to return 15, got %d", response.GetMaxHistorySize())
	}

	if !response.HasSubscription() {
		t.Error("expected HasSubscription() to return true")
	}
}
