package models

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// AudioDSPControls represents the response from GET /audiodspcontrols endpoint
type AudioDSPControls struct {
	XMLName             xml.Name `xml:"audiodspcontrols"`
	AudioMode           string   `xml:"audiomode,attr"`
	VideoSyncAudioDelay int      `xml:"videosyncaudiodelay,attr"`
	SupportedAudioModes string   `xml:"supportedaudiomodes,attr"`
}

// AudioDSPControlsRequest represents the request for POST /audiodspcontrols endpoint
type AudioDSPControlsRequest struct {
	XMLName             xml.Name `xml:"audiodspcontrols"`
	AudioMode           string   `xml:"audiomode,attr,omitempty"`
	VideoSyncAudioDelay int      `xml:"videosyncaudiodelay,attr,omitempty"`
}

// AudioProductToneControls represents the response from GET /audioproducttonecontrols endpoint
type AudioProductToneControls struct {
	XMLName xml.Name             `xml:"audioproducttonecontrols"`
	Bass    BassControlSetting   `xml:"bass"`
	Treble  TrebleControlSetting `xml:"treble"`
}

// AudioProductToneControlsRequest represents the request for POST /audioproducttonecontrols endpoint
type AudioProductToneControlsRequest struct {
	XMLName xml.Name            `xml:"audioproducttonecontrols"`
	Bass    *BassControlValue   `xml:"bass,omitempty"`
	Treble  *TrebleControlValue `xml:"treble,omitempty"`
}

// BassControlSetting represents a bass control setting with constraints
type BassControlSetting struct {
	XMLName  xml.Name `xml:"bass"`
	Value    int      `xml:"value,attr"`
	MinValue int      `xml:"minValue,attr"`
	MaxValue int      `xml:"maxValue,attr"`
	Step     int      `xml:"step,attr"`
}

// TrebleControlSetting represents a treble control setting with constraints
type TrebleControlSetting struct {
	XMLName  xml.Name `xml:"treble"`
	Value    int      `xml:"value,attr"`
	MinValue int      `xml:"minValue,attr"`
	MaxValue int      `xml:"maxValue,attr"`
	Step     int      `xml:"step,attr"`
}

// BassControlValue represents a bass control value for requests
type BassControlValue struct {
	XMLName xml.Name `xml:"bass"`
	Value   int      `xml:"value,attr"`
}

// TrebleControlValue represents a treble control value for requests
type TrebleControlValue struct {
	XMLName xml.Name `xml:"treble"`
	Value   int      `xml:"value,attr"`
}

// AudioProductLevelControls represents the response from GET /audioproductlevelcontrols endpoint
type AudioProductLevelControls struct {
	XMLName                   xml.Name                 `xml:"audioproductlevelcontrols"`
	FrontCenterSpeakerLevel   FrontCenterLevelSetting  `xml:"frontCenterSpeakerLevel"`
	RearSurroundSpeakersLevel RearSurroundLevelSetting `xml:"rearSurroundSpeakersLevel"`
}

// AudioProductLevelControlsRequest represents the request for POST /audioproductlevelcontrols endpoint
type AudioProductLevelControlsRequest struct {
	XMLName                   xml.Name                  `xml:"audioproductlevelcontrols"`
	FrontCenterSpeakerLevel   *FrontCenterControlValue  `xml:"frontCenterSpeakerLevel,omitempty"`
	RearSurroundSpeakersLevel *RearSurroundControlValue `xml:"rearSurroundSpeakersLevel,omitempty"`
}

// FrontCenterLevelSetting represents a front-center speaker level control setting with constraints
type FrontCenterLevelSetting struct {
	XMLName  xml.Name `xml:"frontCenterSpeakerLevel"`
	Value    int      `xml:"value,attr"`
	MinValue int      `xml:"minValue,attr"`
	MaxValue int      `xml:"maxValue,attr"`
	Step     int      `xml:"step,attr"`
}

// RearSurroundLevelSetting represents a rear-surround speakers level control setting with constraints
type RearSurroundLevelSetting struct {
	XMLName  xml.Name `xml:"rearSurroundSpeakersLevel"`
	Value    int      `xml:"value,attr"`
	MinValue int      `xml:"minValue,attr"`
	MaxValue int      `xml:"maxValue,attr"`
	Step     int      `xml:"step,attr"`
}

// FrontCenterControlValue represents a front-center speaker level control value for requests
type FrontCenterControlValue struct {
	XMLName xml.Name `xml:"frontCenterSpeakerLevel"`
	Value   int      `xml:"value,attr"`
}

// RearSurroundControlValue represents a rear-surround speakers level control value for requests
type RearSurroundControlValue struct {
	XMLName xml.Name `xml:"rearSurroundSpeakersLevel"`
	Value   int      `xml:"value,attr"`
}

// Audio mode constants
const (
	AudioModeNormal   = "NORMAL"
	AudioModeDialog   = "DIALOG"
	AudioModeSurround = "SURROUND"
	AudioModeMusic    = "MUSIC"
	AudioModeMovie    = "MOVIE"
	AudioModeSport    = "SPORT"
	AudioModeNight    = "NIGHT"
	AudioModeStandard = "STANDARD"
	AudioModeVivid    = "VIVID"
	AudioModeWarm     = "WARM"
	AudioModeBright   = "BRIGHT"
)

// GetSupportedAudioModes returns a slice of supported audio modes
func (adsp *AudioDSPControls) GetSupportedAudioModes() []string {
	if adsp.SupportedAudioModes == "" {
		return []string{}
	}

	return strings.Split(adsp.SupportedAudioModes, "|")
}

// IsAudioModeSupported checks if the given audio mode is supported
func (adsp *AudioDSPControls) IsAudioModeSupported(mode string) bool {
	supportedModes := adsp.GetSupportedAudioModes()
	for _, supportedMode := range supportedModes {
		if supportedMode == mode {
			return true
		}
	}

	return false
}

// String returns a human-readable string representation of DSP controls
func (adsp *AudioDSPControls) String() string {
	supportedModes := strings.Join(adsp.GetSupportedAudioModes(), ", ")

	return fmt.Sprintf("Audio Mode: %s, Video Sync Delay: %d ms, Supported Modes: [%s]",
		adsp.AudioMode, adsp.VideoSyncAudioDelay, supportedModes)
}

// Validate validates the DSP controls request
func (req *AudioDSPControlsRequest) Validate(capabilities *AudioDSPControls) error {
	if req.AudioMode != "" && capabilities != nil {
		if !capabilities.IsAudioModeSupported(req.AudioMode) {
			return fmt.Errorf("audio mode '%s' is not supported. Supported modes: %s",
				req.AudioMode, strings.Join(capabilities.GetSupportedAudioModes(), ", "))
		}
	}

	if req.VideoSyncAudioDelay < 0 {
		return fmt.Errorf("video sync audio delay cannot be negative: %d", req.VideoSyncAudioDelay)
	}

	return nil
}

// ValidateBass validates the bass value within constraints
func (bc *BassControlSetting) ValidateBass(value int) error {
	if value < bc.MinValue || value > bc.MaxValue {
		return fmt.Errorf("bass value %d is outside valid range [%d, %d]", value, bc.MinValue, bc.MaxValue)
	}

	return nil
}

// ClampValue clamps a value to the valid range
func (bc *BassControlSetting) ClampValue(value int) int {
	if value < bc.MinValue {
		return bc.MinValue
	}

	if value > bc.MaxValue {
		return bc.MaxValue
	}

	return value
}

// ValidateTreble validates the treble value within constraints
func (tc *TrebleControlSetting) ValidateTreble(value int) error {
	if value < tc.MinValue || value > tc.MaxValue {
		return fmt.Errorf("treble value %d is outside valid range [%d, %d]", value, tc.MinValue, tc.MaxValue)
	}

	return nil
}

// ClampValue clamps a value to the valid range
func (tc *TrebleControlSetting) ClampValue(value int) int {
	if value < tc.MinValue {
		return tc.MinValue
	}

	if value > tc.MaxValue {
		return tc.MaxValue
	}

	return value
}

// String returns a human-readable string representation of tone controls
func (atc *AudioProductToneControls) String() string {
	return fmt.Sprintf("Bass: %d [%d-%d], Treble: %d [%d-%d]",
		atc.Bass.Value, atc.Bass.MinValue, atc.Bass.MaxValue,
		atc.Treble.Value, atc.Treble.MinValue, atc.Treble.MaxValue)
}

// Validate validates the tone controls request
func (req *AudioProductToneControlsRequest) Validate(capabilities *AudioProductToneControls) error {
	if req.Bass != nil && capabilities != nil {
		if err := capabilities.Bass.ValidateBass(req.Bass.Value); err != nil {
			return err
		}
	}

	if req.Treble != nil && capabilities != nil {
		if err := capabilities.Treble.ValidateTreble(req.Treble.Value); err != nil {
			return err
		}
	}

	return nil
}

// NewBassControlValue creates a new bass control value for requests
func NewBassControlValue(value int) *BassControlValue {
	return &BassControlValue{
		XMLName: xml.Name{Local: "bass"},
		Value:   value,
	}
}

// NewTrebleControlValue creates a new treble control value for requests
func NewTrebleControlValue(value int) *TrebleControlValue {
	return &TrebleControlValue{
		XMLName: xml.Name{Local: "treble"},
		Value:   value,
	}
}

// ValidateLevel validates the front-center speaker level value within constraints
func (fc *FrontCenterLevelSetting) ValidateLevel(value int) error {
	if value < fc.MinValue || value > fc.MaxValue {
		return fmt.Errorf("front-center speaker level %d is outside valid range [%d, %d]", value, fc.MinValue, fc.MaxValue)
	}

	return nil
}

// ClampLevel clamps a front-center speaker level value to the valid range
func (fc *FrontCenterLevelSetting) ClampLevel(value int) int {
	if value < fc.MinValue {
		return fc.MinValue
	}

	if value > fc.MaxValue {
		return fc.MaxValue
	}

	return value
}

// ValidateLevel validates the rear-surround speaker level value within constraints
func (rs *RearSurroundLevelSetting) ValidateLevel(value int) error {
	if value < rs.MinValue || value > rs.MaxValue {
		return fmt.Errorf("rear-surround speaker level %d is outside valid range [%d, %d]", value, rs.MinValue, rs.MaxValue)
	}

	return nil
}

// ClampLevel clamps a rear-surround speaker level value to the valid range
func (rs *RearSurroundLevelSetting) ClampLevel(value int) int {
	if value < rs.MinValue {
		return rs.MinValue
	}

	if value > rs.MaxValue {
		return rs.MaxValue
	}

	return value
}

// String returns a human-readable string representation of level controls
func (alc *AudioProductLevelControls) String() string {
	return fmt.Sprintf("Front-Center: %d [%d-%d], Rear-Surround: %d [%d-%d]",
		alc.FrontCenterSpeakerLevel.Value, alc.FrontCenterSpeakerLevel.MinValue, alc.FrontCenterSpeakerLevel.MaxValue,
		alc.RearSurroundSpeakersLevel.Value, alc.RearSurroundSpeakersLevel.MinValue, alc.RearSurroundSpeakersLevel.MaxValue)
}

// Validate validates the level controls request
func (req *AudioProductLevelControlsRequest) Validate(capabilities *AudioProductLevelControls) error {
	if req.FrontCenterSpeakerLevel != nil && capabilities != nil {
		if err := capabilities.FrontCenterSpeakerLevel.ValidateLevel(req.FrontCenterSpeakerLevel.Value); err != nil {
			return err
		}
	}

	if req.RearSurroundSpeakersLevel != nil && capabilities != nil {
		if err := capabilities.RearSurroundSpeakersLevel.ValidateLevel(req.RearSurroundSpeakersLevel.Value); err != nil {
			return err
		}
	}

	return nil
}

// NewFrontCenterLevelValue creates a new level control value for front-center speaker
func NewFrontCenterLevelValue(value int) *FrontCenterControlValue {
	return &FrontCenterControlValue{
		XMLName: xml.Name{Local: "frontCenterSpeakerLevel"},
		Value:   value,
	}
}

// NewRearSurroundLevelValue creates a new level control value for rear-surround speakers
func NewRearSurroundLevelValue(value int) *RearSurroundControlValue {
	return &RearSurroundControlValue{
		XMLName: xml.Name{Local: "rearSurroundSpeakersLevel"},
		Value:   value,
	}
}

// AudioCapabilities represents the combined audio capabilities
type AudioCapabilities struct {
	DSPControls          bool `json:"dspControls"`
	ProductToneControls  bool `json:"productToneControls"`
	ProductLevelControls bool `json:"productLevelControls"`
}

// HasAdvancedAudioControls returns true if any advanced audio controls are available
func (ac *AudioCapabilities) HasAdvancedAudioControls() bool {
	return ac.DSPControls || ac.ProductToneControls || ac.ProductLevelControls
}

// GetAvailableControls returns a list of available advanced audio controls
func (ac *AudioCapabilities) GetAvailableControls() []string {
	var controls []string

	if ac.DSPControls {
		controls = append(controls, "DSP Controls")
	}

	if ac.ProductToneControls {
		controls = append(controls, "Tone Controls")
	}

	if ac.ProductLevelControls {
		controls = append(controls, "Level Controls")
	}

	return controls
}

// String returns a human-readable string representation of audio capabilities
func (ac *AudioCapabilities) String() string {
	if !ac.HasAdvancedAudioControls() {
		return "No advanced audio controls available"
	}

	controls := ac.GetAvailableControls()

	return fmt.Sprintf("Available controls: %s", strings.Join(controls, ", "))
}
