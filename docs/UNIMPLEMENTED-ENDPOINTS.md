# Unimplemented SoundTouch API Endpoints

This document provides detailed information about SoundTouch API endpoints that are supported by real hardware but not yet implemented in this Go library. The examples and XML structures are based on real device responses and the comprehensive [HomeAssistant SoundTouch Plus documentation](https://github.com/thlucas1/homeassistantcomponent_soundtouchplus/wiki/SoundTouch-WebServices-API).

## High Priority Implementation Candidates

### Music Service Management

#### POST /setMusicServiceAccount
Adds a music service account to the sources list.

**Request Examples:**

Pandora:
```xml
<credentials source="PANDORA" displayName="Pandora Music Service">
  <user>YourPandoraUserId</user>
  <pass>YourPandoraPassword$1pd</pass>
</credentials>
```

NAS Music Library:
```xml
<credentials source="STORED_MUSIC" displayName="My NAS Media Library:">
  <user>d09708a1-5953-44bc-a413-123456789012/0</user>
  <pass />
</credentials>
```

**Response:**
```xml
<status>/setMusicServiceAccount</status>
```

**Notes:**
- UPnP media servers must be detected first (check `/listMediaServers`)
- Note the `/0` suffix for STORED_MUSIC user names

#### POST /removeMusicServiceAccount
Removes an existing music service account from the sources list.

**Request Examples:**

Remove Pandora:
```xml
<credentials source="PANDORA" displayName="Pandora Music Service">
  <user>YourPandoraUserId</user>
  <pass />
</credentials>
```

Remove NAS Library:
```xml
<credentials source="STORED_MUSIC" displayName="My NAS Media Library:">
  <user>d09708a1-5953-44bc-a413-123456789012/0</user>
  <pass />
</credentials>
```

### Enhanced Preset Management

#### POST /storePreset
Stores a preset to the device (maximum 6 presets).

**Request XML:**
```xml
<preset id="3" createdOn="1696133215" updatedOn="1696133215">
  <ContentItem source="TUNEIN" type="stationurl" location="/v1/playback/station/s309605" sourceAccount="" isPresetable="true">
    <itemName>K-LOVE 90s</itemName>
    <containerArt>http://cdn-profiles.tunein.com/s309605/images/logog.png</containerArt>
  </ContentItem>
</preset>
```

**Response:** Returns updated presets list

**Behavior:**
- If preset ID exists, overlay existing preset
- If content matches existing preset, move to specified slot
- Generates `presetsUpdated` WebSocket event

#### POST /removePreset
Removes an existing preset from the device.

**Request XML:**
```xml
<preset id="4"></preset>
```

**Response:** Returns updated presets list  
**WebSocket Event:** `presetsUpdated`

#### GET /selectPreset
Selects a preset by ID for playback.

**Usage:** Send preset ID to immediately play stored preset content.

### Station Management (Pandora Tested)

#### POST /searchStation
Searches music service for stations that can be added.

**Request XML:**
```xml
<search source="PANDORA" sourceAccount="yourUserId">
  Zach Williams
</search>
```

**Response XML:**
```xml
<results deviceID="..." source="PANDORA" sourceAccount="yourUserId">
  <songs>
    <searchResult source="PANDORA" sourceAccount="yourUserId" token="S88589974">
      <name>Cornerstone (Radio Edit) (feat. Zach Williams)</name>
      <artist>TobyMac</artist>
      <logo>http://mediaserver-cont-usc-mp1-1-v4v6.pandora.com/images/.../1080W_1080H.jpg</logo>
    </searchResult>
    <!-- More songs -->
  </songs>
  <artists>
    <searchResult source="PANDORA" sourceAccount="yourUserId" token="R324771">
      <name>Zach Williams</name>
      <logo>http://mediaserver-cont-dc6-2-v4v6.pandora.com/images/.../1080W_1080H.jpg</logo>
    </searchResult>
    <!-- More artists -->
  </artists>
</results>
```

#### POST /addStation
Adds a station to music service collection.

**Request XML:**
```xml
<addStation source="PANDORA" sourceAccount="yourUserId" token="R4328162">
  <name>Zach Williams &amp; Essential Worship</name>
</addStation>
```

**Response:**
```xml
<status>/addStation</status>
```

**Notes:**
- Station is immediately selected for playing
- Use token from `/searchStation` results

#### POST /removeStation
Removes a station from music service collection.

**Request XML:**
```xml
<ContentItem source="PANDORA" location="126740707481236361" sourceAccount="yourUserId" isPresetable="true">
  <itemName>Zach Williams Radio</itemName>
</ContentItem>
```

**Response:**
```xml
<status>/removeStation</status>
```

**Behavior:**
- If removed station is currently playing, playback stops and source becomes "INVALID_SOURCE"

### Enhanced User Controls

#### POST /userPlayControl
Sends user play control commands.

**Request XML:**
```xml
<PlayControl>PLAY_CONTROL</PlayControl>
```

**Valid Control Values:**
- `PAUSE_CONTROL` - Pause currently playing content
- `PLAY_CONTROL` - Play content that is paused/stopped  
- `PLAY_PAUSE_CONTROL` - Toggle play/pause state
- `STOP_CONTROL` - Stop currently playing content

**Response:**
```xml
<status>/userPlayControl</status>
```

#### POST /userRating
Rates currently playing media (Pandora support confirmed).

**Request XML:**
```xml
<Rating>UP</Rating>
```

**Valid Rating Values:**
- `UP` - Thumbs up rating
- `DOWN` - Thumbs down rating (stops current track, advances to next)

**Response:**
```xml
<status>/userRating</status>
```

**Notes:**
- Ratings stored in artist profile under "My Collection"
- Currently only works with Pandora

## Medium Priority Implementation Candidates

### Content Discovery and Navigation

#### POST /navigate
Returns child container items from music library containers.

**Request XML (Root Container):**
```xml
<navigate source="STORED_MUSIC" sourceAccount="d09708a1-5953-44bc-a413-123456789012/0">
  <startItem>1</startItem>
  <numItems>1000</numItems>
</navigate>
```

**Response XML:**
```xml
<navigateResponse source="STORED_MUSIC" sourceAccount="...">
  <totalItems>4</totalItems>
  <items>
    <item Playable="1">
      <name>Music</name>
      <type>dir</type>
      <ContentItem source="STORED_MUSIC" location="1" sourceAccount="..." isPresetable="true">
        <itemName>Music</itemName>
      </ContentItem>
    </item>
    <item Playable="1">
      <name>Playlists</name>
      <type>dir</type>
      <ContentItem source="STORED_MUSIC" location="12" sourceAccount="..." isPresetable="true">
        <itemName>Playlists</itemName>
      </ContentItem>
    </item>
  </items>
</navigateResponse>
```

**Navigate Specific Container:**
```xml
<navigate source="STORED_MUSIC" sourceAccount="d09708a1-5953-44bc-a413-123456789012/0">
  <startItem>1</startItem>
  <numItems>1000</numItems>
  <item Playable="1">
    <name>Welcome to the New</name>
    <type>dir</type>
    <ContentItem source="STORED_MUSIC" location="7_114e8de9" sourceAccount="..." isPresetable="true">
      <itemName>Welcome to the New</itemName>
    </ContentItem>
  </item>
</navigate>
```

#### POST /search
Searches specified music library container.

**Request XML:**
```xml
<search source="STORED_MUSIC" sourceAccount="d09708a1-5953-44bc-a413-123456789012/0">
  <startItem>1</startItem>
  <numItems>1000</numItems>
  <searchTerm filter="track">baby</searchTerm>
  <item>
    <name>Music Playlists</name>
    <type>dir</type>
    <ContentItem source="STORED_MUSIC" location="F" sourceAccount="..." isPresetable="true" />
  </item>
</search>
```

**Response XML:**
```xml
<searchResponse source="STORED_MUSIC" sourceAccount="...">
  <totalItems>2</totalItems>
  <items>
    <item Playable="1">
      <name>Baby, It's Cold Outside</name>
      <type>track</type>
      <ContentItem source="STORED_MUSIC" location="F_a62301e2-7788 TRACK" sourceAccount="..." isPresetable="true">
        <itemName>Baby, It's Cold Outside</itemName>
      </ContentItem>
      <artistName>Anne Murray</artistName>
      <albumName>Christmas Album</albumName>
    </item>
  </items>
</searchResponse>
```

**Valid Filters:**
- `track` - Search track names
- `artist` - Search artist names  
- `album` - Search album names

### Power Management

#### GET /powerManagement
Returns power state and battery capability information.

**Response XML:**
```xml
<powerManagementResponse>
  <powerState>FullPower</powerState>
  <battery>
    <capable>false</capable>
  </battery>
</powerManagementResponse>
```

#### GET /standby
Places device into standby (power-saving) mode.

**Response:**
```xml
<status>/standby</status>
```

**WebSocket Event:**
```xml
<updates deviceID="...">
  <nowPlayingUpdated>
    <nowPlaying deviceID="..." source="STANDBY">
      <ContentItem source="STANDBY" isPresetable="false" />
    </nowPlaying>
  </nowPlayingUpdated>
</updates>
```

#### GET /lowPowerStandby
Places device into low-power mode (device becomes unresponsive until physical power button pressed).

**Response:**
```xml
<status>/lowPowerStandby</status>
```

**Warning:** Device will not respond to any commands after this until physically powered on.

### Language and Configuration

#### GET /language
Returns current language configuration.

**Response XML:**
```xml
<sysLanguage>3</sysLanguage>
```

#### POST /language
Sets device language.

**Request XML:**
```xml
<sysLanguage>3</sysLanguage>
```

**Language Codes:**
- 1 = DANISH
- 2 = GERMAN  
- 3 = ENGLISH
- 4 = SPANISH
- 5 = FRENCH
- 6 = ITALIAN
- 7 = DUTCH
- 8 = SWEDISH
- 9 = JAPANESE
- 10 = SIMPLIFIED_CHINESE
- 11 = TRADITIONAL_CHINESE
- 12 = KOREAN
- 13 = THAI
- 15 = CZECH
- 16 = FINNISH
- 17 = GREEK
- 18 = NORWEGIAN
- 19 = POLISH
- 20 = PORTUGUESE
- 21 = ROMANIAN
- 22 = RUSSIAN
- 23 = SLOVENIAN
- 24 = TURKISH
- 25 = HUNGARIAN

### System Information and Status

#### GET /soundTouchConfigurationStatus
Returns current SoundTouch configuration status.

**Response XML:**
```xml
<SoundTouchConfigurationStatus status="SOUNDTOUCH_CONFIGURED" />
```

**Valid Status Values:**
- `SOUNDTOUCH_CONFIGURED` - Device configuration complete
- `SOUNDTOUCH_NOT_CONFIGURED` - Device not configured
- `SOUNDTOUCH_CONFIGURING` - Configuration in progress

#### GET /serviceAvailability
Returns information about which source services are currently available.

**Response XML:**
```xml
<serviceAvailability>
  <services>
    <service type="AIRPLAY" isAvailable="true" />
    <service type="ALEXA" isAvailable="false" />
    <service type="AMAZON" isAvailable="true" />
    <service type="BLUETOOTH" isAvailable="false" reason="INVALID_SOURCE_TYPE" />
    <service type="SPOTIFY" isAvailable="true" />
    <service type="TUNEIN" isAvailable="true" />
    <!-- More services -->
  </services>
</serviceAvailability>
```

#### GET /listMediaServers
Returns information about detected UPnP/DLNA media servers.

**Response XML:**
```xml
<ListMediaServersResponse>
  <media_server id="d09708a1-5953-44bc-a413-123456789012" mac="S-1-5-21-240303764-901663538-1234567890-1001" ip="192.168.1.5" manufacturer="Microsoft Corporation" model_name="Windows Media Player Sharing" friendly_name="My NAS Media Library" model_description="" location="http://192.168.1.5:2869/upnphost/udhisapi.dll?..." />
</ListMediaServersResponse>
```

#### GET /requestToken ✅ **Implemented**
Returns a new bearer token generated by the device.

**Response XML:**
```xml
<bearertoken value="Bearer vUApzBVT6Lh0nw1xVu/plr1UDRNdMYMEpe0cStm4wCH5mWSjrrtORnGGirMn3pspkJ8mNR1MFh/J4OcsbEikMplcDGJVeuZOnDPAskQALvDBCF0PW74qXRms2k1AfLJ/" />
```

**Implementation**: Available via `RequestToken()` client method and `soundtouch-cli token request` command.

### Software Updates

#### GET /swUpdateCheck
Gets latest available software update release information.

**Response XML:**
```xml
<swUpdateCheckResponse deviceID="..." indexFileUrl="https://worldwide.bose.com/updates/soundtouch">
  <release revision="27.0.6.46330.5043500" />
</swUpdateCheckResponse>
```

#### GET /swUpdateQuery
Gets status of a SoundTouch software update.

**Response XML:**
```xml
<swUpdateQueryResponse deviceID="...">
  <state>IDLE</state>
  <percentComplete>0</percentComplete>
  <canAbort>false</canAbort>
</swUpdateQueryResponse>
```

## Low Priority / Specialized Endpoints

### Notification System (ST-10 Series Only)

#### POST /playNotification
Plays a notification beep on the device.

**Response:**
```xml
<status>/playNotification</status>
```

**Behavior:**
- Pauses current media
- Emits double beep sound
- Resumes media playback

**Note:** Only works on ST-10 series. ST-300 and other models do not support this.

#### POST /speaker
Plays TTS messages or URL content (ST-10 Series Only).

**TTS Example:**
```xml
<play_info>
  <url>http://translate.google.com/translate_tts?ie=UTF-8&amp;tl=EN&amp;client=tw-ob&amp;q=a.There%20is%20activity%20at%20the%20front%20door.</url>
  <app_key>YourAppKey</app_key>
  <service>TTS Notification</service>
  <message>Google TTS</message>
  <reason>a.There is activity at the front door.</reason>
  <volume>70</volume>
</play_info>
```

**URL Content Example:**
```xml
<play_info>
  <url>https://freetestdata.com/wp-content/uploads/2021/09/Free_Test_Data_1MB_MP3.mp3</url>
  <app_key>YourAppKey</app_key>
  <service>FreeTestData.com</service>
  <message>MP3 Test Data</message>
  <reason>Free_Test_Data_1MB_MP3</reason>
</play_info>
```

**Notes:**
- Only ST-10 series supported
- Pauses current content during notification
- Volume restored after notification completes
- SoundTouch device limits volume range 10-70

### WiFi Management

#### POST /performWirelessSiteSurvey
Gets list of wireless networks detected by device.

**Response XML:**
```xml
<PerformWirelessSiteSurveyResponse error="none">
  <items>
    <item ssid="my_wireless_ssid" signalStrength="-58" secure="true">
      <securityTypes>
        <type>wpa_or_wpa2</type>
      </securityTypes>
    </item>
    <item ssid="NETGEAR07" signalStrength="-83" secure="true">
      <securityTypes>
        <type>wpa_or_wpa2</type>
      </securityTypes>
    </item>
  </items>
</PerformWirelessSiteSurveyResponse>
```

#### POST /addWirelessProfile
Adds wireless profile configuration to device.

**Request XML:**
```xml
<addWirelessProfile timeout="30">
  <profile ssid="YourSSIDName" password="YourSSIDPassword" securityType="wpa_or_wpa2"></profile>
</addWirelessProfile>
```

**Security Types:**
- `none` - No security
- `wep` - WEP
- `wpatkip` - WPA/TKIP
- `wpaaes` - WPA/AES
- `wpa2tkip` - WPA2/TKIP
- `wpa2aes` - WPA2/AES
- `wpa_or_wpa2` - WPA/WPA2 (recommended)

**Setup Process:**
1. Connect to device WiFi (default IP: 192.0.2.1)
2. Add wireless profile
3. End setup: POST to `/setup` with `<setupState state="SETUP_WIFI_LEAVE" />`

#### GET /getActiveWirelessProfile
Gets current wireless profile configuration.

**Response XML:**
```xml
<GetActiveWirelessProfileResponse>
  <ssid>my_wireless_ssid</ssid>
</GetActiveWirelessProfileResponse>
```

### Bluetooth Management

#### POST /enterBluetoothPairing
Enters Bluetooth pairing mode and waits for device to pair.

**Response:**
```xml
<status>/enterBluetoothPairing</status>
```

**Behavior:**
- Device enters pairing mode
- Bluetooth indicator turns blue
- Emits ascending tone when pairing complete
- Source immediately switches to BLUETOOTH

#### POST /clearBluetoothPaired
Clears all existing Bluetooth pairings.

**Response:**
```xml
<BluetoothInfo BluetoothMACAddress="34:15:13:45:2f:93" />
```

**Behavior:**
- All existing pairings removed
- Devices need to re-pair
- Emits descending tone

### Source Selection Shortcuts

#### GET /selectLastSource
Selects the last source that was selected.

**Response:**
```xml
<status>/selectLastSource</status>
```

#### GET /selectLastSoundTouchSource
Selects the last SoundTouch source that was selected.

**Response:**
```xml
<status>/selectLastSoundTouchSource</status>
```

#### GET /selectLastWiFiSource
Selects the last WiFi source that was selected.

**Response:**
```xml
<status>/selectLastWiFiSource</status>
```

#### GET /selectLocalSource
Selects the LOCAL source (for devices where this is the only way to select LOCAL).

**Response:**
```xml
<status>/selectLocalSource</status>
```

### Group Management (ST-10 Stereo Pairs Only)

The ST-10 is the only SoundTouch product that supports stereo pair groups (different from zones).

#### GET /getGroup
Gets current left/right stereo pair configuration.

**Response XML:**
```xml
<group id="1115893">
  <name>Bose-ST10-1 + Bose-ST10-4</name>
  <masterDeviceId>9070658C9D4A</masterDeviceId>
  <roles>
    <groupRole>
      <deviceId>9070658C9D4A</deviceId>
      <role>LEFT</role>
      <ipAddress>192.168.1.131</ipAddress>
    </groupRole>
    <groupRole>
      <deviceId>F45EAB3115DA</deviceId>
      <role>RIGHT</role>
      <ipAddress>192.168.1.134</ipAddress>
    </groupRole>
  </roles>
  <senderIPAddress>192.168.1.131</senderIPAddress>
  <status>GROUP_OK</status>
</group>
```

#### POST /addGroup
Creates new left/right stereo pair speaker group.

**Request XML:**
```xml
<group>
  <name>Bose-ST10-1 + Bose-ST10-4</name>
  <masterDeviceId>9070658C9D4A</masterDeviceId>
  <roles>
    <groupRole>
      <deviceId>9070658C9D4A</deviceId>
      <role>LEFT</role>
      <ipAddress>192.168.1.131</ipAddress>
    </groupRole>
    <groupRole>
      <deviceId>F45EAB3115DA</deviceId>
      <role>RIGHT</role>
      <ipAddress>192.168.1.134</ipAddress>
    </groupRole>
  </roles>
</group>
```

**WebSocket Event:** `groupUpdated` sent to both devices

#### GET /removeGroup
Removes existing stereo pair group.

**Response:**
```xml
<group />
```

#### POST /updateGroup
Updates name of stereo pair group.

**Request XML:**
```xml
<group id="1116267">
  <name>Updated Group Name</name>
  <masterDeviceId>9070658C9D4A</masterDeviceId>
  <roles>
    <groupRole>
      <deviceId>9070658C9D4A</deviceId>
      <role>LEFT</role>
      <ipAddress>192.168.1.131</ipAddress>
    </groupRole>
    <groupRole>
      <deviceId>F45EAB3115DA</deviceId>
      <role>RIGHT</role>
      <ipAddress>192.168.1.134</ipAddress>
    </groupRole>
  </roles>
</group>
```

### Advanced System Configuration

#### GET /systemtimeout
Gets current system timeout configuration.

**Response XML:**
```xml
<systemtimeout>
  <powersaving_enabled>true</powersaving_enabled>
</systemtimeout>
```

#### GET /rebroadcastlatencymode
Gets current rebroadcast latency mode configuration.

**Response XML:**
```xml
<rebroadcastlatencymode mode="SYNC_TO_ZONE" controllable="true" />
```

#### GET /DSPMonoStereo
Gets current digital signal processor configuration.

**Response XML:**
```xml
<DSPMonoStereo deviceID="...">
  <mono enable="false" />
</DSPMonoStereo>
```

## Implementation Notes

### Device Compatibility
- Many endpoints work on specific device models only
- Always check `/supportedURLs` before implementing
- Test with real hardware when possible

### Error Handling
- Services may return timeout errors on unsupported devices
- Some endpoints appear in `/supportedURLs` but still don't work
- Graceful degradation recommended

### WebSocket Events
Many POST operations generate corresponding WebSocket events:
- `presetsUpdated` - Preset changes
- `groupUpdated` - Group changes  
- `volumeUpdated` - Volume changes
- `nowPlayingUpdated` - Source/playback changes
- `zoneUpdated` - Zone changes

### Security Considerations
- `/speaker` endpoint requires app_key parameter
- Token-based authentication available via `/requestToken` ✅ **Implemented**
- Some operations require device to be in specific states

### Music Service Specifics
- Pandora: Confirmed working for station management, ratings
- Spotify: Requires PREMIUM account for most operations
- STORED_MUSIC: Requires UPnP/DLNA server setup
- LOCAL_MUSIC: Requires SoundTouch App Media Server running

This documentation provides the foundation for implementing these endpoints in the Go library, with real-world examples and detailed XML structures verified against actual SoundTouch hardware.