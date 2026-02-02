package models

import (
	"encoding/xml"
	"testing"
)

func TestNewMusicServiceCredentials(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		displayName string
		user        string
		pass        string
	}{
		{
			name:        "Valid credentials",
			source:      "SPOTIFY",
			displayName: "Spotify Premium",
			user:        "user@spotify.com",
			pass:        "password123",
		},
		{
			name:        "Empty display name",
			source:      "PANDORA",
			displayName: "",
			user:        "pandora_user",
			pass:        "pandora_pass",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cred := NewMusicServiceCredentials(tt.source, tt.displayName, tt.user, tt.pass)

			if cred.Source != tt.source {
				t.Errorf("Expected source %s, got %s", tt.source, cred.Source)
			}

			if cred.DisplayName != tt.displayName {
				t.Errorf("Expected displayName %s, got %s", tt.displayName, cred.DisplayName)
			}

			if cred.User != tt.user {
				t.Errorf("Expected user %s, got %s", tt.user, cred.User)
			}

			if cred.Pass != tt.pass {
				t.Errorf("Expected pass %s, got %s", tt.pass, cred.Pass)
			}
		})
	}
}

func TestNewSpotifyCredentials(t *testing.T) {
	cred := NewSpotifyCredentials("user@spotify.com", "mypassword")

	if cred.Source != "SPOTIFY" {
		t.Errorf("Expected source SPOTIFY, got %s", cred.Source)
	}

	if cred.DisplayName != "Spotify Premium" {
		t.Errorf("Expected displayName 'Spotify Premium', got %s", cred.DisplayName)
	}

	if cred.User != "user@spotify.com" {
		t.Errorf("Expected user 'user@spotify.com', got %s", cred.User)
	}

	if cred.Pass != "mypassword" {
		t.Errorf("Expected pass 'mypassword', got %s", cred.Pass)
	}
}

func TestNewPandoraCredentials(t *testing.T) {
	cred := NewPandoraCredentials("pandora_user", "pandora_pass")

	if cred.Source != "PANDORA" {
		t.Errorf("Expected source PANDORA, got %s", cred.Source)
	}

	if cred.DisplayName != "Pandora Music Service" {
		t.Errorf("Expected displayName 'Pandora Music Service', got %s", cred.DisplayName)
	}

	if cred.User != "pandora_user" {
		t.Errorf("Expected user 'pandora_user', got %s", cred.User)
	}

	if cred.Pass != "pandora_pass" {
		t.Errorf("Expected pass 'pandora_pass', got %s", cred.Pass)
	}
}

func TestNewStoredMusicCredentials(t *testing.T) {
	cred := NewStoredMusicCredentials("d09708a1-5953-44bc-a413-123456789012/0", "My NAS Library")

	if cred.Source != "STORED_MUSIC" {
		t.Errorf("Expected source STORED_MUSIC, got %s", cred.Source)
	}

	if cred.DisplayName != "My NAS Library" {
		t.Errorf("Expected displayName 'My NAS Library', got %s", cred.DisplayName)
	}

	if cred.User != "d09708a1-5953-44bc-a413-123456789012/0" {
		t.Errorf("Expected user 'd09708a1-5953-44bc-a413-123456789012/0', got %s", cred.User)
	}

	if cred.Pass != "" {
		t.Errorf("Expected empty pass for STORED_MUSIC, got %s", cred.Pass)
	}
}

func TestNewAmazonMusicCredentials(t *testing.T) {
	cred := NewAmazonMusicCredentials("amazon_user", "amazon_pass")

	if cred.Source != "AMAZON" {
		t.Errorf("Expected source AMAZON, got %s", cred.Source)
	}

	if cred.DisplayName != "Amazon Music" {
		t.Errorf("Expected displayName 'Amazon Music', got %s", cred.DisplayName)
	}

	if cred.User != "amazon_user" {
		t.Errorf("Expected user 'amazon_user', got %s", cred.User)
	}

	if cred.Pass != "amazon_pass" {
		t.Errorf("Expected pass 'amazon_pass', got %s", cred.Pass)
	}
}

func TestNewDeezerCredentials(t *testing.T) {
	cred := NewDeezerCredentials("deezer_user", "deezer_pass")

	if cred.Source != "DEEZER" {
		t.Errorf("Expected source DEEZER, got %s", cred.Source)
	}

	if cred.DisplayName != "Deezer Premium" {
		t.Errorf("Expected displayName 'Deezer Premium', got %s", cred.DisplayName)
	}

	if cred.User != "deezer_user" {
		t.Errorf("Expected user 'deezer_user', got %s", cred.User)
	}

	if cred.Pass != "deezer_pass" {
		t.Errorf("Expected pass 'deezer_pass', got %s", cred.Pass)
	}
}

func TestNewIHeartRadioCredentials(t *testing.T) {
	cred := NewIHeartRadioCredentials("iheart_user", "iheart_pass")

	if cred.Source != "IHEART" {
		t.Errorf("Expected source IHEART, got %s", cred.Source)
	}

	if cred.DisplayName != "iHeartRadio" {
		t.Errorf("Expected displayName 'iHeartRadio', got %s", cred.DisplayName)
	}

	if cred.User != "iheart_user" {
		t.Errorf("Expected user 'iheart_user', got %s", cred.User)
	}

	if cred.Pass != "iheart_pass" {
		t.Errorf("Expected pass 'iheart_pass', got %s", cred.Pass)
	}
}

func TestMusicServiceCredentials_Validate(t *testing.T) {
	tests := []struct {
		name        string
		credentials *MusicServiceCredentials
		wantError   bool
		errorMsg    string
	}{
		{
			name: "Valid Spotify credentials",
			credentials: &MusicServiceCredentials{
				Source: "SPOTIFY",
				User:   "user@spotify.com",
				Pass:   "password",
			},
			wantError: false,
		},
		{
			name: "Valid STORED_MUSIC credentials (no password required)",
			credentials: &MusicServiceCredentials{
				Source: "STORED_MUSIC",
				User:   "guid/0",
				Pass:   "",
			},
			wantError: false,
		},
		{
			name: "Empty source",
			credentials: &MusicServiceCredentials{
				Source: "",
				User:   "user",
				Pass:   "pass",
			},
			wantError: true,
			errorMsg:  "source cannot be empty",
		},
		{
			name: "Empty user",
			credentials: &MusicServiceCredentials{
				Source: "SPOTIFY",
				User:   "",
				Pass:   "pass",
			},
			wantError: true,
			errorMsg:  "user cannot be empty",
		},
		{
			name: "Empty password for non-STORED_MUSIC",
			credentials: &MusicServiceCredentials{
				Source: "SPOTIFY",
				User:   "user",
				Pass:   "",
			},
			wantError: true,
			errorMsg:  "password cannot be empty for SPOTIFY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.credentials.Validate()

			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Expected error message %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestMusicServiceCredentials_IsForRemoval(t *testing.T) {
	tests := []struct {
		name        string
		credentials *MusicServiceCredentials
		expected    bool
	}{
		{
			name: "Has password - not for removal",
			credentials: &MusicServiceCredentials{
				Pass: "password",
			},
			expected: false,
		},
		{
			name: "Empty password - for removal",
			credentials: &MusicServiceCredentials{
				Pass: "",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.credentials.IsForRemoval()
			if result != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, result)
			}
		})
	}
}

func TestMusicServiceCredentials_HasPassword(t *testing.T) {
	tests := []struct {
		name        string
		credentials *MusicServiceCredentials
		expected    bool
	}{
		{
			name: "Has password",
			credentials: &MusicServiceCredentials{
				Pass: "password",
			},
			expected: true,
		},
		{
			name: "No password",
			credentials: &MusicServiceCredentials{
				Pass: "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.credentials.HasPassword()
			if result != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, result)
			}
		})
	}
}

func TestMusicServiceCredentials_GetDescription(t *testing.T) {
	tests := []struct {
		name        string
		credentials *MusicServiceCredentials
		expected    string
	}{
		{
			name: "Has display name",
			credentials: &MusicServiceCredentials{
				Source:      "SPOTIFY",
				DisplayName: "Custom Spotify Name",
			},
			expected: "Custom Spotify Name",
		},
		{
			name: "Spotify default",
			credentials: &MusicServiceCredentials{
				Source: "SPOTIFY",
			},
			expected: "Spotify Premium",
		},
		{
			name: "Pandora default",
			credentials: &MusicServiceCredentials{
				Source: "PANDORA",
			},
			expected: "Pandora Music Service",
		},
		{
			name: "Amazon default",
			credentials: &MusicServiceCredentials{
				Source: "AMAZON",
			},
			expected: "Amazon Music",
		},
		{
			name: "Deezer default",
			credentials: &MusicServiceCredentials{
				Source: "DEEZER",
			},
			expected: "Deezer Premium",
		},
		{
			name: "iHeartRadio default",
			credentials: &MusicServiceCredentials{
				Source: "IHEART",
			},
			expected: "iHeartRadio",
		},
		{
			name: "STORED_MUSIC default",
			credentials: &MusicServiceCredentials{
				Source: "STORED_MUSIC",
			},
			expected: "Network Music Library",
		},
		{
			name: "LOCAL_MUSIC default",
			credentials: &MusicServiceCredentials{
				Source: "LOCAL_MUSIC",
			},
			expected: "Local Music Server",
		},
		{
			name: "Unknown source",
			credentials: &MusicServiceCredentials{
				Source: "UNKNOWN",
			},
			expected: "UNKNOWN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.credentials.GetDescription()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestMusicServiceAccountResponse_IsSuccess(t *testing.T) {
	tests := []struct {
		name     string
		response *MusicServiceAccountResponse
		expected bool
	}{
		{
			name: "Set account success",
			response: &MusicServiceAccountResponse{
				Status: "/setMusicServiceAccount",
			},
			expected: true,
		},
		{
			name: "Remove account success",
			response: &MusicServiceAccountResponse{
				Status: "/removeMusicServiceAccount",
			},
			expected: true,
		},
		{
			name: "Other status",
			response: &MusicServiceAccountResponse{
				Status: "/someOtherEndpoint",
			},
			expected: false,
		},
		{
			name: "Empty status",
			response: &MusicServiceAccountResponse{
				Status: "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.IsSuccess()
			if result != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, result)
			}
		})
	}
}

func TestMusicServiceCredentials_XMLMarshaling(t *testing.T) {
	cred := &MusicServiceCredentials{
		Source:      "SPOTIFY",
		DisplayName: "Spotify Premium",
		User:        "user@spotify.com",
		Pass:        "mypassword",
	}

	// Test marshaling
	data, err := xml.Marshal(cred)
	if err != nil {
		t.Errorf("Failed to marshal credentials: %v", err)
	}

	expectedXML := `<credentials source="SPOTIFY" displayName="Spotify Premium"><user>user@spotify.com</user><pass>mypassword</pass></credentials>`
	if string(data) != expectedXML {
		t.Errorf("Expected XML %s, got %s", expectedXML, string(data))
	}

	// Test unmarshaling
	var unmarshaledCred MusicServiceCredentials

	err = xml.Unmarshal(data, &unmarshaledCred)
	if err != nil {
		t.Errorf("Failed to unmarshal credentials: %v", err)
	}

	if unmarshaledCred.Source != cred.Source {
		t.Errorf("Expected source %s, got %s", cred.Source, unmarshaledCred.Source)
	}

	if unmarshaledCred.DisplayName != cred.DisplayName {
		t.Errorf("Expected displayName %s, got %s", cred.DisplayName, unmarshaledCred.DisplayName)
	}

	if unmarshaledCred.User != cred.User {
		t.Errorf("Expected user %s, got %s", cred.User, unmarshaledCred.User)
	}

	if unmarshaledCred.Pass != cred.Pass {
		t.Errorf("Expected pass %s, got %s", cred.Pass, unmarshaledCred.Pass)
	}
}

func TestMusicServiceAccountResponse_XMLMarshaling(t *testing.T) {
	response := &MusicServiceAccountResponse{
		Status: "/setMusicServiceAccount",
	}

	// Test marshaling
	data, err := xml.Marshal(response)
	if err != nil {
		t.Errorf("Failed to marshal response: %v", err)
	}

	expectedXML := `<status>/setMusicServiceAccount</status>`
	if string(data) != expectedXML {
		t.Errorf("Expected XML %s, got %s", expectedXML, string(data))
	}

	// Test unmarshaling
	var unmarshaledResponse MusicServiceAccountResponse

	err = xml.Unmarshal(data, &unmarshaledResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if unmarshaledResponse.Status != response.Status {
		t.Errorf("Expected status %s, got %s", response.Status, unmarshaledResponse.Status)
	}
}
