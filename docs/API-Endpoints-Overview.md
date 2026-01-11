# Bose SoundTouch Web API - Endpoints Overview

This document provides a comprehensive overview of the available API endpoints verified against the official Bose SoundTouch Web API v1.0 specification (January 7, 2026).

## Implementation Status Legend
- ✅ **Implemented** - Fully implemented with tests and real device validation
- 🔍 **Extra** - Implemented but not in official API v1.0 (may be newer version or undocumented)
- ⚠️ **Different** - Implemented with different approach than official API
- ℹ️ **N/A** - Documented but officially unsupported or non-functional on real hardware

## API Basics

- **Protocol**: HTTP REST-like
- **Data Format**: XML Request/Response
- **Standard Port**: 8090
- **Base URL**: `http://<device-ip>:8090/`
- **Authentication**: No complex authentication required
- **Real-time Updates**: WebSocket connection available

## Device Information

### GET /info ✅ **Implemented**
Retrieves basic device information.

**Response XML Structure:**
```xml
<info deviceID="..." type="..." name="..." ...>
  <name>Device Name</name>
  <type>Device Type</type>
  <margeAccountUUID>UUID</margeAccountUUID>
  <components>...</components>
</info>
```

## Playback Control

### GET /now_playing ✅ **Implemented**
Retrieves information about the currently playing music.

**Response XML Structure:**
```xml
<nowPlaying deviceID="..." source="...">
  <ContentItem source="..." type="..." location="..." sourceAccount="...">
    <itemName>Track Name</itemName>
    <containerArt>Album Art URL</containerArt>
  </ContentItem>
  <track>Track Name</track>
  <artist>Artist Name</artist>
  <album>Album Name</album>
  <stationName>Station Name</stationName>
  <art artImageStatus="...">Art URL</art>
  <playStatus>PLAY_STATE</playStatus>
  <shuffleSetting>...</shuffleSetting>
  <repeatSetting>...</repeatSetting>
</nowPlaying>
```

### POST /key ✅ **Implemented**
Sends key commands to the device.

**Important**: Proper key simulation requires sending both press and release states:

**Request XML (Press + Release):**
```xml
<key state="press" sender="Gabbo">KEY_NAME</key>
<key state="release" sender="Gabbo">KEY_NAME</key>
```

**Available Keys:**

**Playback Controls:**
- `PLAY` - Start playback
- `PAUSE` - Pause current playback
- `STOP` - Stop current playback
- `PREV_TRACK` - Go to previous track
- `NEXT_TRACK` - Go to next track

**Rating and Bookmark Controls:**
- `THUMBS_UP` - Rate current content positively (Pandora, etc.)
- `THUMBS_DOWN` - Rate current content negatively
- `BOOKMARK` - Bookmark current content

**Power and System Controls:**
- `POWER` - Toggle device power state
- `MUTE` - Toggle mute state

**Volume Controls:**
- `VOLUME_UP` - Increase volume
- `VOLUME_DOWN` - Decrease volume

**Preset Controls:**
- `PRESET_1` to `PRESET_6` - Select preset 1-6

**Input Controls:**
- `AUX_INPUT` - Switch to auxiliary input

**Shuffle Controls:**
- `SHUFFLE_OFF` - Turn shuffle mode off
- `SHUFFLE_ON` - Turn shuffle mode on

**Repeat Controls:**
- `REPEAT_OFF` - Turn repeat mode off
- `REPEAT_ONE` - Repeat current track
- `REPEAT_ALL` - Repeat all tracks in playlist

## Volume Control

### GET /volume ✅ **Implemented**
Retrieves the current volume.

**Response XML:**
```xml
<volume deviceID="...">
  <targetvolume>50</targetvolume>
  <actualvolume>50</actualvolume>
  <muteenabled>false</muteenabled>
</volume>
```

### POST /volume ✅ **Implemented**
Sets the volume.

**Request XML:**
```xml
<volume>50</volume>
```

## Bass Settings

### GET /bass ✅ **Implemented**
Retrieves the current bass settings.

**Response XML:**
```xml
<bass deviceID="...">
  <targetbass>0</targetbass>
  <actualbass>0</actualbass>
</bass>
```

### POST /bass ✅ **Implemented**
Sets the bass settings (-9 to +9).

**Request XML:**
```xml
<bass>0</bass>
```

## Source Management

### GET /sources ✅ **Implemented**
Retrieves the available audio sources.

**Response XML:**
```xml
<sources deviceID="...">
  <sourceItem source="SPOTIFY" sourceAccount="..." status="READY" multiroomallowed="true">
    <itemName>Spotify</itemName>
  </sourceItem>
  <sourceItem source="BLUETOOTH" status="READY" multiroomallowed="false">
    <itemName>Bluetooth</itemName>
  </sourceItem>
  <!-- Additional sources -->
</sources>
```

**Typical Sources:**
- `SPOTIFY`
- `AMAZON`
- `PANDORA`
- `IHEARTRADIO`
- `TUNEIN`
- `BLUETOOTH`
- `AUX`
- `STORED_MUSIC`

### POST /select ✅ **Implemented**
Selects an audio source.

**Request XML:**
```xml
<ContentItem source="SPOTIFY" sourceAccount="...">
  <itemName>Spotify</itemName>
</ContentItem>
```

## Preset Management

### GET /presets ✅ **Implemented**
Retrieves the configured presets.

**Response XML:**
```xml
<presets deviceID="...">
  <preset id="1" createdOn="..." updatedOn="...">
    <ContentItem source="..." sourceAccount="..." location="...">
      <itemName>Preset Name</itemName>
      <containerArt>Art URL</containerArt>
    </ContentItem>
  </preset>
  <!-- Additional presets -->
</presets>
```

### POST /presets ℹ️ **N/A**
Creates or updates a preset.

**Status**: According to the official Bose SoundTouch API documentation, POST operations on `/presets` are marked as "N/A" - this endpoint officially does not support preset creation or modification via any API client.

**Alternative Methods**:
- Use the official Bose SoundTouch mobile app
- Use physical preset buttons on the device (long-press while content is playing)
- Changes made via these methods will be visible through the GET endpoint

## Advanced Features

### GET /getZone ✅ **Implemented**
Retrieves multiroom zone information.

### POST /setZone ✅ **Implemented**
Configures multiroom zones.

### GET /balance ✅ **Implemented**
Retrieves balance settings (stereo devices).

### POST /balance ✅ **Implemented**
Sets balance settings.

### GET /clockTime ✅ **Implemented**
Retrieves the device time.

### POST /clockTime ✅ **Implemented**
Sets the device time.

### GET /clockDisplay ✅ **Implemented**
Retrieves clock display settings.

### POST /clockDisplay ✅ **Implemented**
Configures the clock display.

## WebSocket Connection

### WebSocket / ✅ **Implemented**
Establishes a persistent connection for live updates.

**Event Types:**
- `nowPlayingUpdated`
- `volumeUpdated`
- `connectionStateUpdated`
- `presetUpdated`

## Network and System

### GET /networkInfo ✅ **Implemented**
Retrieves network information.

### GET /capabilities ✅ **Implemented**
Retrieves device capabilities.

### GET /name 🔍 **Extra** 
Retrieves the device name.

**Note**: Official API only documents `POST /name` for setting device name. Our GET implementation appears to be an undocumented extension.

### POST /name ✅ **Implemented**
Sets the device name via `SetName()` method.

**Official Request Format:**
```xml
<name>$STRING</name>
```

### GET /bassCapabilities ✅ **Implemented**
Checks if bass customization is supported on the device.

**Official Response Format:**
```xml
<bassCapabilities deviceID="$MACADDR">
    <bassAvailable>$BOOL</bassAvailable>
    <bassMin>$INT</bassMin>
    <bassMax>$INT</bassMax>
    <bassDefault>$INT</bassDefault>
</bassCapabilities>
```

### GET /trackInfo ✅ **Implemented**
Gets track information (duplicate of `/now_playing` per official API).

**Status**: Fully implemented but times out on SoundTouch 10 & 20 test devices (AllegroWebserver timeout). May work on other SoundTouch models or firmware versions. Use `/now_playing` endpoint as reliable alternative.

**Implementation**: Available via `GetTrackInfo()` method. Consider using `GetNowPlaying()` method for guaranteed compatibility.

### Zone Slave Management ✅ **Implemented**
Both official low-level endpoints and high-level zone management are available:

#### POST /addZoneSlave ✅ **Implemented**
Add individual device to existing zone using official API format.

**Implementation**: Available via `AddZoneSlave()` and `AddZoneSlaveByDeviceID()` methods

#### POST /removeZoneSlave ✅ **Implemented** 
Remove individual device from existing zone using official API format.

**Implementation**: Available via `RemoveZoneSlave()` and `RemoveZoneSlaveByDeviceID()` methods

#### High-Level Zone API ✅ **Enhanced**
- **Enhanced**: `CreateZone()`, `AddToZone()`, `RemoveFromZone()` methods via `/setZone`
- **Status**: Provides both official low-level API and enhanced high-level operations

### Advanced Audio Controls ✅ **Conditionally Available**
Professional/high-end device features (only available on devices that list these capabilities):

#### `/audiodspcontrols` - GET/POST ✅ **Implemented**
Access DSP settings including audio modes and video sync delay.

**Availability**: Only available if `audiodspcontrols` is listed in the reply to `GET /capabilities`

**Implementation**: Available via `GetAudioDSPControls()`, `SetAudioDSPControls()`, `SetAudioMode()`, `SetVideoSyncAudioDelay()` methods with automatic capability checking

#### `/audioproducttonecontrols` - GET/POST ✅ **Implemented**
Advanced bass and treble controls (beyond basic `/bass` endpoint).

**Availability**: Only available if `audioproducttonecontrols` is listed in the reply to `GET /capabilities`

**Implementation**: Available via `GetAudioProductToneControls()`, `SetAudioProductToneControls()`, `SetAdvancedBass()`, `SetAdvancedTreble()` methods with automatic capability checking

#### `/audioproductlevelcontrols` - GET/POST ✅ **Implemented**
Speaker level controls for front-center and rear-surround speakers.

**Availability**: Only available if `audioproductlevelcontrols` is listed in the reply to `GET /capabilities`

**Implementation**: Available via `GetAudioProductLevelControls()`, `SetAudioProductLevelControls()`, `SetFrontCenterSpeakerLevel()`, `SetRearSurroundSpeakersLevel()` methods with automatic capability checking

### Clock and Network Endpoints 🔍 **Extra**
These endpoints work with real hardware but are NOT in official API v1.0:
- `GET/POST /clockTime` ✅ **Implemented** - Device time management
- `GET/POST /clockDisplay` ✅ **Implemented** - Clock display settings  
- `GET /networkInfo` ✅ **Implemented** - Network information

### Balance Control 🔍 **Extra**
- `GET/POST /balance` ✅ **Implemented** - Stereo balance adjustment

**Note**: Not documented in official API v1.0 but works with real devices.

## Coverage Summary

### Official API Coverage: 100%
- **Total Official Endpoints**: 19
- **Implemented**: 19 (100%)
- **Conditionally Available**: 3 (16%) - Advanced audio endpoints require device support
- **Device-Dependent**: 1 (5%) - GET /trackInfo times out on some models
- **Excluded**: 1 endpoint (POST /presets officially N/A)

### Feature Coverage: 100%
- ✅ All available user functionality implemented
- ✅ All functional device operations supported  
- ✅ Complete WebSocket event system
- ✅ Full multiroom capabilities
- ✅ Complete advanced audio controls (where supported by device)
- 🔍 Additional features beyond official specification


## Error Handling

The API uses standard HTTP status codes:
- `200 OK` - Successful request
- `400 Bad Request` - Invalid request
- `404 Not Found` - Endpoint or resource not found
- `500 Internal Server Error` - Internal device error

## Example Implementation

```go
// Example for a GET request
func GetNowPlaying(deviceIP string) (*NowPlaying, error) {
    url := fmt.Sprintf("http://%s:8090/now_playing", deviceIP)
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var nowPlaying NowPlaying
    err = xml.NewDecoder(resp.Body).Decode(&nowPlaying)
    return &nowPlaying, err
}

// Example for a POST request
func SendKey(deviceIP string, key string) error {
    url := fmt.Sprintf("http://%s:8090/key", deviceIP)
    xmlData := fmt.Sprintf(`<key state="press" sender="GoClient">%s</key>`, key)
    
    resp, err := http.Post(url, "application/xml", strings.NewReader(xmlData))
    if err != nil {
        return err
    }
    resp.Body.Close()
    return nil
}
```

## Notes

1. **XML Namespace**: Most responses use no explicit XML namespace
2. **Encoding**: UTF-8 is used for all XML documents
3. **Timeouts**: Recommended timeout for HTTP requests: 10 seconds
4. **Rate Limiting**: No explicit limits documented, but moderate usage recommended
5. **Device Discovery**: Devices can be found via UPnP on the local network

## Reference

Based on the official Bose SoundTouch Web API documentation:
https://assets.bosecreative.com/m/496577402d128874/original/SoundTouch-Web-API.pdf
