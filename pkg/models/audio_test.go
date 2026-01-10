package models

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestAudioDSPControls_GetSupportedAudioModes(t *testing.T) {
	tests := []struct {
		name           string
		supportedModes string
		expected       []string
	}{
		{
			name:           "multiple modes",
			supportedModes: "NORMAL|DIALOG|SURROUND|MUSIC",
			expected:       []string{"NORMAL", "DIALOG", "SURROUND", "MUSIC"},
		},
		{
			name:           "single mode",
			supportedModes: "NORMAL",
			expected:       []string{"NORMAL"},
		},
		{
			name:           "empty modes",
			supportedModes: "",
			expected:       []string{},
		},
		{
			name:           "modes with spaces",
			supportedModes: "NORMAL|DIALOG CLEAR|MUSIC",
			expected:       []string{"NORMAL", "DIALOG CLEAR", "MUSIC"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsp := AudioDSPControls{
				SupportedAudioModes: tt.supportedModes,
			}

			result := dsp.GetSupportedAudioModes()

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d modes, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("Expected mode %d to be '%s', got '%s'", i, expected, result[i])
				}
			}
		})
	}
}

func TestAudioDSPControls_IsAudioModeSupported(t *testing.T) {
	dsp := AudioDSPControls{
		SupportedAudioModes: "NORMAL|DIALOG|SURROUND|MUSIC",
	}

	tests := []struct {
		mode     string
		expected bool
	}{
		{"NORMAL", true},
		{"DIALOG", true},
		{"SURROUND", true},
		{"MUSIC", true},
		{"MOVIE", false},
		{"INVALID", false},
		{"", false},
		{"normal", false}, // Case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.mode, func(t *testing.T) {
			result := dsp.IsAudioModeSupported(tt.mode)
			if result != tt.expected {
				t.Errorf("Expected IsAudioModeSupported('%s') to be %v, got %v", tt.mode, tt.expected, result)
			}
		})
	}
}

func TestAudioDSPControls_String(t *testing.T) {
	dsp := AudioDSPControls{
		AudioMode:           "MUSIC",
		VideoSyncAudioDelay: 50,
		SupportedAudioModes: "NORMAL|DIALOG|MUSIC",
	}

	result := dsp.String()
	expected := "Audio Mode: MUSIC, Video Sync Delay: 50 ms, Supported Modes: [NORMAL, DIALOG, MUSIC]"

	if result != expected {
		t.Errorf("Expected string representation '%s', got '%s'", expected, result)
	}
}

func TestAudioDSPControlsRequest_Validate(t *testing.T) {
	capabilities := &AudioDSPControls{
		SupportedAudioModes: "NORMAL|DIALOG|MUSIC",
	}

	tests := []struct {
		name        string
		request     *AudioDSPControlsRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid audio mode",
			request: &AudioDSPControlsRequest{
				AudioMode: "MUSIC",
			},
			expectError: false,
		},
		{
			name: "invalid audio mode",
			request: &AudioDSPControlsRequest{
				AudioMode: "INVALID",
			},
			expectError: true,
			errorMsg:    "audio mode 'INVALID' is not supported",
		},
		{
			name: "negative video sync delay",
			request: &AudioDSPControlsRequest{
				VideoSyncAudioDelay: -10,
			},
			expectError: true,
			errorMsg:    "video sync audio delay cannot be negative",
		},
		{
			name: "valid video sync delay",
			request: &AudioDSPControlsRequest{
				VideoSyncAudioDelay: 100,
			},
			expectError: false,
		},
		{
			name: "valid combined request",
			request: &AudioDSPControlsRequest{
				AudioMode:           "DIALOG",
				VideoSyncAudioDelay: 25,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate(capabilities)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestToneControlSetting_ValidateBass(t *testing.T) {
	setting := BassControlSetting{
		MinValue: -10,
		MaxValue: 10,
	}

	tests := []struct {
		value       int
		expectError bool
	}{
		{0, false},
		{-10, false},
		{10, false},
		{5, false},
		{-5, false},
		{-11, true},
		{11, true},
		{100, true},
		{-100, true},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.value)), func(t *testing.T) {
			err := setting.ValidateBass(tt.value)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for value %d but got none", tt.value)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for value %d but got: %v", tt.value, err)
				}
			}
		})
	}
}

func TestToneControlSetting_ClampValue(t *testing.T) {
	setting := TrebleControlSetting{
		MinValue: -5,
		MaxValue: 5,
	}

	tests := []struct {
		input    int
		expected int
	}{
		{0, 0},
		{3, 3},
		{-3, -3},
		{5, 5},
		{-5, -5},
		{10, 5},
		{-10, -5},
		{100, 5},
		{-100, -5},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.input)), func(t *testing.T) {
			result := setting.ClampValue(tt.input)
			if result != tt.expected {
				t.Errorf("Expected ClampValue(%d) to be %d, got %d", tt.input, tt.expected, result)
			}
		})
	}
}

func TestAudioProductToneControls_String(t *testing.T) {
	controls := AudioProductToneControls{
		Bass: BassControlSetting{
			Value:    3,
			MinValue: -10,
			MaxValue: 10,
		},
		Treble: TrebleControlSetting{
			Value:    -2,
			MinValue: -10,
			MaxValue: 10,
		},
	}

	result := controls.String()
	expected := "Bass: 3 [-10-10], Treble: -2 [-10-10]"

	if result != expected {
		t.Errorf("Expected string representation '%s', got '%s'", expected, result)
	}
}

func TestAudioProductToneControlsRequest_Validate(t *testing.T) {
	capabilities := &AudioProductToneControls{
		Bass: BassControlSetting{
			MinValue: -10,
			MaxValue: 10,
		},
		Treble: TrebleControlSetting{
			MinValue: -5,
			MaxValue: 5,
		},
	}

	tests := []struct {
		name        string
		request     *AudioProductToneControlsRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid bass only",
			request: &AudioProductToneControlsRequest{
				Bass: NewBassControlValue(5),
			},
			expectError: false,
		},
		{
			name: "valid treble only",
			request: &AudioProductToneControlsRequest{
				Treble: NewTrebleControlValue(3),
			},
			expectError: false,
		},
		{
			name: "invalid bass value",
			request: &AudioProductToneControlsRequest{
				Bass: NewBassControlValue(15),
			},
			expectError: true,
			errorMsg:    "bass value 15 is outside valid range",
		},
		{
			name: "invalid treble value",
			request: &AudioProductToneControlsRequest{
				Treble: NewTrebleControlValue(-10),
			},
			expectError: true,
			errorMsg:    "treble value -10 is outside valid range",
		},
		{
			name: "valid combined request",
			request: &AudioProductToneControlsRequest{
				Bass:   NewBassControlValue(-5),
				Treble: NewTrebleControlValue(2),
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate(capabilities)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestNewBassControlValue(t *testing.T) {
	value := NewBassControlValue(5)

	if value.Value != 5 {
		t.Errorf("Expected value 5, got %d", value.Value)
	}

	if value.XMLName.Local != "bass" {
		t.Errorf("Expected XMLName.Local to be 'bass', got '%s'", value.XMLName.Local)
	}
}

func TestNewTrebleControlValue(t *testing.T) {
	value := NewTrebleControlValue(-3)

	if value.Value != -3 {
		t.Errorf("Expected value -3, got %d", value.Value)
	}

	if value.XMLName.Local != "treble" {
		t.Errorf("Expected XMLName.Local to be 'treble', got '%s'", value.XMLName.Local)
	}
}

func TestAudioProductLevelControls_String(t *testing.T) {
	controls := AudioProductLevelControls{
		FrontCenterSpeakerLevel: FrontCenterLevelSetting{
			Value:    2,
			MinValue: -10,
			MaxValue: 10,
		},
		RearSurroundSpeakersLevel: RearSurroundLevelSetting{
			Value:    -1,
			MinValue: -10,
			MaxValue: 10,
		},
	}

	result := controls.String()
	expected := "Front-Center: 2 [-10-10], Rear-Surround: -1 [-10-10]"

	if result != expected {
		t.Errorf("Expected string representation '%s', got '%s'", expected, result)
	}
}

func TestAudioProductLevelControlsRequest_Validate(t *testing.T) {
	capabilities := &AudioProductLevelControls{
		FrontCenterSpeakerLevel: FrontCenterLevelSetting{
			MinValue: -5,
			MaxValue: 5,
		},
		RearSurroundSpeakersLevel: RearSurroundLevelSetting{
			MinValue: -8,
			MaxValue: 8,
		},
	}

	tests := []struct {
		name        string
		request     *AudioProductLevelControlsRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid front center only",
			request: &AudioProductLevelControlsRequest{
				FrontCenterSpeakerLevel: NewFrontCenterLevelValue(3),
			},
			expectError: false,
		},
		{
			name: "valid rear surround only",
			request: &AudioProductLevelControlsRequest{
				RearSurroundSpeakersLevel: NewRearSurroundLevelValue(-4),
			},
			expectError: false,
		},
		{
			name: "invalid front center value",
			request: &AudioProductLevelControlsRequest{
				FrontCenterSpeakerLevel: NewFrontCenterLevelValue(10),
			},
			expectError: true,
			errorMsg:    "speaker level 10 is outside valid range",
		},
		{
			name: "invalid rear surround value",
			request: &AudioProductLevelControlsRequest{
				RearSurroundSpeakersLevel: NewRearSurroundLevelValue(-15),
			},
			expectError: true,
			errorMsg:    "speaker level -15 is outside valid range",
		},
		{
			name: "valid combined request",
			request: &AudioProductLevelControlsRequest{
				FrontCenterSpeakerLevel:   NewFrontCenterLevelValue(-2),
				RearSurroundSpeakersLevel: NewRearSurroundLevelValue(5),
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate(capabilities)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestNewFrontCenterLevelValue(t *testing.T) {
	value := NewFrontCenterLevelValue(3)

	if value.Value != 3 {
		t.Errorf("Expected value 3, got %d", value.Value)
	}

	if value.XMLName.Local != "frontCenterSpeakerLevel" {
		t.Errorf("Expected XMLName.Local to be 'frontCenterSpeakerLevel', got '%s'", value.XMLName.Local)
	}
}

func TestNewRearSurroundLevelValue(t *testing.T) {
	value := NewRearSurroundLevelValue(-2)

	if value.Value != -2 {
		t.Errorf("Expected value -2, got %d", value.Value)
	}

	if value.XMLName.Local != "rearSurroundSpeakersLevel" {
		t.Errorf("Expected XMLName.Local to be 'rearSurroundSpeakersLevel', got '%s'", value.XMLName.Local)
	}
}

func TestAudioCapabilities_HasAdvancedAudioControls(t *testing.T) {
	tests := []struct {
		name         string
		capabilities AudioCapabilities
		expected     bool
	}{
		{
			name: "no controls",
			capabilities: AudioCapabilities{
				DSPControls:          false,
				ProductToneControls:  false,
				ProductLevelControls: false,
			},
			expected: false,
		},
		{
			name: "dsp controls only",
			capabilities: AudioCapabilities{
				DSPControls:          true,
				ProductToneControls:  false,
				ProductLevelControls: false,
			},
			expected: true,
		},
		{
			name: "tone controls only",
			capabilities: AudioCapabilities{
				DSPControls:          false,
				ProductToneControls:  true,
				ProductLevelControls: false,
			},
			expected: true,
		},
		{
			name: "level controls only",
			capabilities: AudioCapabilities{
				DSPControls:          false,
				ProductToneControls:  false,
				ProductLevelControls: true,
			},
			expected: true,
		},
		{
			name: "all controls",
			capabilities: AudioCapabilities{
				DSPControls:          true,
				ProductToneControls:  true,
				ProductLevelControls: true,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.capabilities.HasAdvancedAudioControls()
			if result != tt.expected {
				t.Errorf("Expected HasAdvancedAudioControls() to be %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestAudioCapabilities_GetAvailableControls(t *testing.T) {
	tests := []struct {
		name         string
		capabilities AudioCapabilities
		expected     []string
	}{
		{
			name: "no controls",
			capabilities: AudioCapabilities{
				DSPControls:          false,
				ProductToneControls:  false,
				ProductLevelControls: false,
			},
			expected: []string{},
		},
		{
			name: "dsp controls only",
			capabilities: AudioCapabilities{
				DSPControls:          true,
				ProductToneControls:  false,
				ProductLevelControls: false,
			},
			expected: []string{"DSP Controls"},
		},
		{
			name: "all controls",
			capabilities: AudioCapabilities{
				DSPControls:          true,
				ProductToneControls:  true,
				ProductLevelControls: true,
			},
			expected: []string{"DSP Controls", "Tone Controls", "Level Controls"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.capabilities.GetAvailableControls()

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d controls, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("Expected control %d to be '%s', got '%s'", i, expected, result[i])
				}
			}
		})
	}
}

func TestAudioCapabilities_String(t *testing.T) {
	tests := []struct {
		name         string
		capabilities AudioCapabilities
		expected     string
	}{
		{
			name: "no controls",
			capabilities: AudioCapabilities{
				DSPControls:          false,
				ProductToneControls:  false,
				ProductLevelControls: false,
			},
			expected: "No advanced audio controls available",
		},
		{
			name: "single control",
			capabilities: AudioCapabilities{
				DSPControls:          true,
				ProductToneControls:  false,
				ProductLevelControls: false,
			},
			expected: "Available controls: DSP Controls",
		},
		{
			name: "multiple controls",
			capabilities: AudioCapabilities{
				DSPControls:          true,
				ProductToneControls:  true,
				ProductLevelControls: false,
			},
			expected: "Available controls: DSP Controls, Tone Controls",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.capabilities.String()
			if result != tt.expected {
				t.Errorf("Expected string representation '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestAudioDSPControls_XMLMarshaling(t *testing.T) {
	controls := AudioDSPControls{
		AudioMode:           "MUSIC",
		VideoSyncAudioDelay: 50,
		SupportedAudioModes: "NORMAL|DIALOG|MUSIC",
	}

	xmlData, err := xml.Marshal(controls)
	if err != nil {
		t.Fatalf("Failed to marshal XML: %v", err)
	}

	xmlStr := string(xmlData)

	if !strings.Contains(xmlStr, `audiomode="MUSIC"`) {
		t.Error("Expected XML to contain audiomode attribute")
	}

	if !strings.Contains(xmlStr, `videosyncaudiodelay="50"`) {
		t.Error("Expected XML to contain videosyncaudiodelay attribute")
	}

	if !strings.Contains(xmlStr, `supportedaudiomodes="NORMAL|DIALOG|MUSIC"`) {
		t.Error("Expected XML to contain supportedaudiomodes attribute")
	}
}

func TestAudioProductToneControls_XMLMarshaling(t *testing.T) {
	controls := AudioProductToneControls{
		Bass: BassControlSetting{
			XMLName:  xml.Name{Local: "bass"},
			Value:    3,
			MinValue: -10,
			MaxValue: 10,
			Step:     1,
		},
		Treble: TrebleControlSetting{
			XMLName:  xml.Name{Local: "treble"},
			Value:    -2,
			MinValue: -5,
			MaxValue: 5,
			Step:     1,
		},
	}

	xmlData, err := xml.Marshal(controls)
	if err != nil {
		t.Fatalf("Failed to marshal XML: %v", err)
	}

	xmlStr := string(xmlData)

	if !strings.Contains(xmlStr, `<bass value="3" minValue="-10" maxValue="10" step="1">`) {
		t.Error("Expected XML to contain bass element with correct attributes")
	}

	if !strings.Contains(xmlStr, `<treble value="-2" minValue="-5" maxValue="5" step="1">`) {
		t.Error("Expected XML to contain treble element with correct attributes")
	}
}

func TestAudioProductLevelControls_XMLMarshaling(t *testing.T) {
	controls := AudioProductLevelControls{
		FrontCenterSpeakerLevel: FrontCenterLevelSetting{
			XMLName:  xml.Name{Local: "frontCenterSpeakerLevel"},
			Value:    2,
			MinValue: -10,
			MaxValue: 10,
			Step:     1,
		},
		RearSurroundSpeakersLevel: RearSurroundLevelSetting{
			XMLName:  xml.Name{Local: "rearSurroundSpeakersLevel"},
			Value:    -1,
			MinValue: -8,
			MaxValue: 8,
			Step:     1,
		},
	}

	xmlData, err := xml.Marshal(controls)
	if err != nil {
		t.Fatalf("Failed to marshal XML: %v", err)
	}

	xmlStr := string(xmlData)

	if !strings.Contains(xmlStr, `<frontCenterSpeakerLevel value="2" minValue="-10" maxValue="10" step="1">`) {
		t.Error("Expected XML to contain frontCenterSpeakerLevel element with correct attributes")
	}

	if !strings.Contains(xmlStr, `<rearSurroundSpeakersLevel value="-1" minValue="-8" maxValue="8" step="1">`) {
		t.Error("Expected XML to contain rearSurroundSpeakersLevel element with correct attributes")
	}
}
