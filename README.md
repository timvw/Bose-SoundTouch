# Bose SoundTouch API Client

A modern Go library and CLI tool for interacting with Bose SoundTouch devices via their Web API.

## Features

### ✅ Implemented (85% Complete - 16/19 endpoints)
- **HTTP Client with XML Support**: Complete client for SoundTouch Web API
- **Device Information**: Get detailed device info via `/info` endpoint
- **Device Name**: Get device name via `/name` endpoint  
- **Device Capabilities**: Get device capabilities via `/capabilities` endpoint
- **Configured Presets**: Get preset configurations via `/presets` endpoint
- **Now Playing Status**: Get current playback information via `/now_playing` endpoint
- **Audio Sources**: Get available sources via `/sources` endpoint
- **Media Controls**: Play, pause, stop, track navigation via `/key` endpoint
- **Volume Management**: Get/set volume, incremental control via `/volume` endpoint
- **Bass Control**: Get/set bass levels (-9 to +9 range) via `/bass` endpoint
- **Balance Control**: Get/set balance (-50 to +50 range) via `/balance` endpoint
- **Clock/Time Management**: Get/set device time via `/clockTime` and `/clockDisplay` endpoints
- **Network Information**: Get network details via `/networkInfo` endpoint
- **Real-time WebSocket Events**: Live monitoring of device state changes
- **UPnP/SSDP Discovery**: Automatic device discovery using Universal Plug and Play
- **mDNS/Bonjour Discovery**: Multicast DNS device discovery support
- **Cross-Platform**: Works on Windows, macOS, Linux, and WASM
- **CLI Tool**: Command-line interface for testing and control operations
- **Flexible Configuration**: Support for .env files and environment variables
- **Unified Discovery**: Combines UPnP, mDNS, and configured device lists
- **Safety Features**: Volume warnings, increment limits, error validation

### 🔄 Remaining High Priority (15% - 3/19 endpoints)
- **Device System**: POST /reboot for device restart
- **Multiroom Support**: GET/POST /getZone, /setZone (if supported by device)

### ❌ Not Supported by API
- **Preset Creation**: POST /presets (officially not supported by SoundTouch API)

## Recent Additions - WebSocket Events ⚡

**NEW**: Real-time WebSocket support has been added! Monitor device state changes in real-time with comprehensive event handling.

### Key Features:
- 🎵 **Live Now Playing Updates**: Track changes, playback status, shuffle/repeat
- 🔊 **Real-time Volume Changes**: Volume levels and mute status  
- 🌐 **Connection Monitoring**: Network connectivity and signal strength
- 📻 **Preset Notifications**: Preset updates and selections
- 🏠 **Multiroom Events**: Zone membership changes
- 🎚️ **Audio Settings**: Bass level adjustments
- 🔄 **Auto-Reconnection**: Robust connection management
- 🎛️ **Event Filtering**: Subscribe to specific event types
- 📊 **Comprehensive Logging**: Debug and monitoring capabilities

### CLI Demo:
```bash
# Quick start - auto-discover and monitor all events
go run ./cmd/websocket-demo -discover

# Monitor specific device with event filtering  
go run ./cmd/websocket-demo -host 192.168.1.10 -filter nowPlaying,volume -verbose
```

## Installation

### Using Go
```bash
go install github.com/user_account/bose-soundtouch/cmd/soundtouch-cli@latest
```

### From Source
```bash
git clone https://github.com/user_account/bose-soundtouch.git
cd bose-soundtouch
make build
```

## Quick Start

### Configuration

Create a `.env` file in your working directory to configure preferred devices:

```bash
# Copy the example file
cp .env.example .env
```

Example `.env` configuration:
```bash
# Discovery Settings
DISCOVERY_TIMEOUT=5s
UPNP_ENABLED=true
MDNS_ENABLED=true

# Preferred Devices (alternative to UPnP)
# Format: name@host:port;name@host:port;...
PREFERRED_DEVICES="Living Room@192.168.1.10;Kitchen@192.168.1.11;192.168.1.12:8091"

# HTTP Client Settings
HTTP_TIMEOUT=10s
USER_AGENT="Bose-SoundTouch-Go-Client/1.0"
```

### CLI Usage

#### Device Discovery

The library supports multiple discovery methods automatically:
- **Configuration**: Manually specified devices in `.env` file (fastest, most reliable)
- **UPnP/SSDP**: Universal Plug and Play discovery (widely supported)
- **mDNS/Bonjour**: Multicast DNS discovery (Apple ecosystem friendly)

See [docs/DISCOVERY.md](docs/DISCOVERY.md) for detailed information.

```bash
# Discover SoundTouch devices (combines UPnP, mDNS + configured devices)
soundtouch-cli -discover

# Discover and show detailed info for all devices
soundtouch-cli -discover-all

# Discover with custom timeout
soundtouch-cli -discover -timeout 10s
```

#### Real-time WebSocket Events

Monitor device state changes in real-time using WebSocket connections:

```bash
# Auto-discover device and monitor all events
go run ./cmd/websocket-demo -discover

# Connect to specific device and monitor all events
go run ./cmd/websocket-demo -host 192.168.1.10

# Monitor only volume and now playing events
go run ./cmd/websocket-demo -host 192.168.1.10 -filter volume,nowPlaying

# Monitor for 5 minutes with verbose output
go run ./cmd/websocket-demo -host 192.168.1.10 -duration 5m -verbose

# Available event types for filtering:
# nowPlaying, volume, connection, preset, zone, bass
```

**Supported WebSocket Events:**
- 🎵 **Now Playing**: Track changes, playback status, shuffle/repeat settings
- 🔊 **Volume**: Volume level and mute status changes
- 🌐 **Connection**: Network connectivity and signal strength
- 📻 **Preset**: Preset configuration updates
- 🏠 **Zone**: Multiroom zone membership changes
- 🎚️ **Bass**: Bass equalizer level adjustments

See [docs/websocket-events.md](docs/websocket-events.md) for complete WebSocket documentation.

#### Device Information
```bash
# Get device information by IP address
soundtouch-cli -host 192.168.1.10 -info

# With custom port and timeout
soundtouch-cli -host 192.168.1.10 -port 8090 -timeout 15s -info
```

#### Now Playing Status
```bash
# Get current playback information
soundtouch-cli -host 192.168.1.10 -nowplaying

# Example output:
# Now Playing:
#   Device ID: ABCD1234EFGH
#   Source: SPOTIFY
#   Status: Playing
#   Title: In Between Breaths - Paris Unplugged
#   Artist: SYML
#   Album: Paris Unplugged
#   Duration: 2:32 / 3:30
#   Shuffle: Off
#   Repeat: Off
#   Artwork: https://i.scdn.co/image/...
#   Capabilities: Skip, Skip Previous, Seek, Favorite

#### Audio Sources
```bash
# Get available audio sources
soundtouch-cli -host 192.168.1.10:8090 -sources

# Example output:
# Audio Sources:
#   Device ID: ABCD1234EFGH
#   Total Sources: 14
#   Ready Sources: 5
#
# Ready Sources:
#   • AUX IN [Local, Multiroom]
#   • user+spotify@example.com (user) [Remote, Multiroom, Streaming]
#   • Alexa [Remote, Multiroom]
#   • Tunein [Remote, Multiroom, Streaming]
#   • Local_internet_radio [Remote, Multiroom, Streaming]
#
# Categories:
#   Spotify: 1 account(s) ready
#   AUX Input: Ready
#   Streaming Services: 3 ready
```

#### Media Controls
```bash
# Basic playback controls
soundtouch-cli -host 192.168.1.10:8090 -play
soundtouch-cli -host 192.168.1.10:8090 -pause
soundtouch-cli -host 192.168.1.10:8090 -stop

# Track navigation
soundtouch-cli -host 192.168.1.10:8090 -next
soundtouch-cli -host 192.168.1.10:8090 -prev

# Volume controls (key-based)
soundtouch-cli -host 192.168.1.10:8090 -volume-up
soundtouch-cli -host 192.168.1.10:8090 -volume-down

# Preset selection
soundtouch-cli -host 192.168.1.10:8090 -preset 1
soundtouch-cli -host 192.168.1.10:8090 -preset 6

# Generic key command
soundtouch-cli -host 192.168.1.10:8090 -key STOP
```

#### Volume Management
```bash
# Get current volume
soundtouch-cli -host 192.168.1.10:8090 -volume

# Example output:
# Current Volume:
#   Device ID: ABCD1234EFGH
#   Current Level: 50 (Medium)
#   Target Level: 50
#   Muted: false

# Set specific volume (0-100, shows warning for >30)
soundtouch-cli -host 192.168.1.10:8090 -set-volume 25
soundtouch-cli -host 192.168.1.10:8090 -set-volume 0  # Mute

# Incremental volume control
soundtouch-cli -host 192.168.1.10:8090 -inc-volume 3
soundtouch-cli -host 192.168.1.10:8090 -dec-volume 5
```

#### Device Name
```bash
# Get device name
soundtouch-cli -host 192.168.1.10 -name

# Example output:
# Device Name: My SoundTouch
```

#### Device Capabilities
```bash
# Get device capabilities
soundtouch-cli -host 192.168.1.10 -capabilities

# Example output:
# Device Capabilities:
#   Device ID: ABCD1234EFGH
#
# System Features:
#   • Power Saving Disabled
#
# Audio Features:
#   • L/R Stereo
#
# Network Features:
#   • Dual Mode
#   • WSAPI Proxy
#
# Extended Capabilities:
#   • systemtimeout (/systemtimeout)
#   • rebroadcastlatencymode (/rebroadcastlatencymode)
```

#### Configured Presets
```bash
# Get configured presets
soundtouch-cli -host 192.168.1.10 -presets

# Example output:
# Configured Presets:
#   Used Slots: 6/6
#   Spotify Presets: 6
#
# Preset 1: My Playlist
#   Source: SPOTIFY (user@example.com)
#   Type: tracklisturl
#   Created: 2024-06-23 09:40:36
#   Updated: 2024-10-12 15:39:42
#   Artwork: https://i.scdn.co/image/...
#
# Most Recent: Preset 4 (Movie Soundtrack)
```

#### System Information
```bash
# Get device clock time
soundtouch-cli -host 192.168.1.10 -clock-time

# Set device time to current system time
soundtouch-cli -host 192.168.1.10 -set-clock-time now

# Get clock display settings
soundtouch-cli -host 192.168.1.10 -clock-display

# Configure clock display
soundtouch-cli -host 192.168.1.10 -enable-clock
soundtouch-cli -host 192.168.1.10 -clock-format 24
soundtouch-cli -host 192.168.1.10 -clock-brightness 75

# Get network information
soundtouch-cli -host 192.168.1.10 -network-info
```

### Go Library Usage

#### Basic HTTP Client Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/user_account/bose-soundtouch/pkg/client"
)

func main() {
    // Create client
    soundTouchClient := client.NewClientFromHost("192.168.1.10")
    
    // Get device information
    deviceInfo, err := soundTouchClient.GetDeviceInfo()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Device: %s (%s)\n", deviceInfo.Name, deviceInfo.Type)
    
    // Get now playing
    nowPlaying, err := soundTouchClient.GetNowPlaying()
    if err != nil {
        log.Fatal(err)
    }
    
    if !nowPlaying.IsEmpty() {
        fmt.Printf("Now Playing: %s by %s\n", nowPlaying.Track, nowPlaying.Artist)
        fmt.Printf("Status: %s\n", nowPlaying.PlayStatus.String())
    }
    
    // Volume control
    volume, err := soundTouchClient.GetVolume()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Volume: %d\n", volume.ActualVolume)
    
    // Set volume safely (with warnings)
    err = soundTouchClient.SetVolumeSafe(25)
    if err != nil {
        log.Fatal(err)
    }
    
    // Media controls
    soundTouchClient.Play()
    soundTouchClient.Pause()
    soundTouchClient.VolumeUp()
    
    // Source selection
    soundTouchClient.SelectSpotify()
    soundTouchClient.SelectPreset(1)
}
```

#### Real-time WebSocket Events

```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/user_account/bose-soundtouch/pkg/client"
    "github.com/user_account/bose-soundtouch/pkg/models"
)

func main() {
    // Create SoundTouch client
    soundTouchClient := client.NewClientFromHost("192.168.1.10")
    
    // Create WebSocket client
    wsClient := soundTouchClient.NewWebSocketClient(nil)
    
    // Set up event handlers
    wsClient.OnNowPlaying(func(event *models.NowPlayingUpdatedEvent) {
        np := &event.NowPlaying
        log.Printf("🎵 Now Playing: %s by %s", np.Track, np.Artist)
        log.Printf("   Status: %s, Source: %s", np.PlayStatus.String(), np.Source)
        
        if np.HasTimeInfo() {
            log.Printf("   Duration: %s", np.FormatDuration())
        }
    })
    
    wsClient.OnVolumeUpdated(func(event *models.VolumeUpdatedEvent) {
        vol := &event.Volume
        if vol.IsMuted() {
            log.Println("🔇 Volume: Muted")
        } else {
            log.Printf("🔊 Volume: %d (%s)", vol.ActualVolume, 
                models.GetVolumeLevelName(vol.ActualVolume))
        }
    })
    
    wsClient.OnConnectionState(func(event *models.ConnectionStateUpdatedEvent) {
        cs := &event.ConnectionState
        if cs.IsConnected() {
            log.Printf("✅ Connected (Signal: %s)", cs.GetSignalStrength())
        } else {
            log.Printf("❌ Connection: %s", cs.State)
        }
    })
    
    wsClient.OnBassUpdated(func(event *models.BassUpdatedEvent) {
        bass := &event.Bass
        log.Printf("🎚️ Bass: %d", bass.ActualBass)
    })
    
    // Handle unknown events for debugging
    wsClient.OnUnknownEvent(func(event *models.WebSocketEvent) {
        log.Printf("❓ Unknown event types: %v", event.GetEventTypes())
    })
    
    // Connect to WebSocket
    if err := wsClient.Connect(); err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    
    log.Println("Connected! Listening for events... (Press Ctrl+C to stop)")
    
    // Set up graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    // Wait for shutdown signal
    <-sigChan
    log.Println("Shutting down...")
    
    // Disconnect
    if err := wsClient.Disconnect(); err != nil {
        log.Printf("Error during disconnect: %v", err)
    }
    
    log.Println("Disconnected successfully")
}
```

#### Advanced WebSocket Configuration

```go
package main

import (
    "log"
    "time"
    
    "github.com/user_account/bose-soundtouch/pkg/client"
    "github.com/user_account/bose-soundtouch/pkg/models"
)

// Custom logger for WebSocket events
type CustomLogger struct{}

func (c *CustomLogger) Printf(format string, v ...interface{}) {
    timestamp := time.Now().Format("15:04:05.000")
    log.Printf("[%s] [WebSocket] %s", timestamp, fmt.Sprintf(format, v...))
}

func main() {
    soundTouchClient := client.NewClientFromHost("192.168.1.10")
    
    // Custom WebSocket configuration
    config := &client.WebSocketConfig{
        ReconnectInterval:    3 * time.Second,  // Reconnect every 3 seconds
        MaxReconnectAttempts: 5,                // Try 5 times before giving up
        PingInterval:         15 * time.Second, // Ping every 15 seconds
        PongTimeout:          5 * time.Second,  // Wait 5 seconds for pong
        ReadBufferSize:       4096,             // 4KB read buffer
        WriteBufferSize:      4096,             // 4KB write buffer
        Logger:               &CustomLogger{},  // Custom logger
    }
    
    wsClient := soundTouchClient.NewWebSocketClient(config)
    
    // Set up handlers for specific events only
    wsClient.OnNowPlaying(func(event *models.NowPlayingUpdatedEvent) {
        // Handle only now playing events
        log.Printf("Track changed: %s", event.NowPlaying.GetDisplayTitle())
    })
    
    // Connect with custom config
    if err := wsClient.ConnectWithConfig(config); err != nil {
        log.Fatal(err)
    }
    
    // Keep running
    wsClient.Wait()
}
```

#### Device Discovery with WebSocket

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/user_account/bose-soundtouch/pkg/client"
    "github.com/user_account/bose-soundtouch/pkg/config"
    "github.com/user_account/bose-soundtouch/pkg/discovery"
    "github.com/user_account/bose-soundtouch/pkg/models"
)

func main() {
    // Discover devices
    cfg := &config.Config{
        DiscoveryTimeout: 10 * time.Second,
        CacheEnabled:     false,
    }
    
    discoveryService := discovery.NewUnifiedDiscoveryService(cfg)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    devices, err := discoveryService.DiscoverDevices(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    if len(devices) == 0 {
        log.Fatal("No devices found")
    }
    
    // Connect to first device found
    device := devices[0]
    log.Printf("Connecting to: %s (%s:%d)", device.Name, device.Host, device.Port)
    
    clientConfig := client.ClientConfig{
        Host: device.Host,
        Port: device.Port,
    }
    
    soundTouchClient := client.NewClient(clientConfig)
    
    // Test basic connectivity
    deviceInfo, err := soundTouchClient.GetDeviceInfo()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Connected to: %s", deviceInfo.Name)
    
    // Set up WebSocket monitoring
    wsClient := soundTouchClient.NewWebSocketClient(nil)
    
    wsClient.OnNowPlaying(func(event *models.NowPlayingUpdatedEvent) {
        log.Printf("[%s] Now Playing: %s", 
            deviceInfo.Name, event.NowPlaying.GetDisplayTitle())
    })
    
    if err := wsClient.Connect(); err != nil {
        log.Fatal(err)
    }
    
    log.Println("Monitoring events...")
    wsClient.Wait()
}
```

## Project Structure

```
Bose-SoundTouch/
├── cmd/
│   ├── soundtouch-cli/          # Main CLI tool (fully functional)
│   ├── websocket-demo/          # WebSocket event monitoring demo
│   ├── example-upnp/            # UPnP discovery examples
│   ├── example-mdns/            # mDNS discovery examples
│   └── mdns-scanner/            # Network scanning utility
├── pkg/
│   ├── client/                  # HTTP & WebSocket clients
│   │   ├── client.go            # Main HTTP API client
│   │   ├── websocket.go         # WebSocket event client
│   │   └── *_test.go           # Comprehensive tests
│   ├── models/                  # Typed XML models
│   │   ├── websocket.go         # WebSocket event models
│   │   ├── nowplaying.go        # Now playing models
│   │   ├── volume.go            # Volume control models
│   │   ├── bass.go              # Bass control models
│   │   ├── balance.go           # Balance control models
│   │   └── *.go                 # Other endpoint models
│   ├── discovery/               # Device discovery
│   │   ├── unified.go           # Unified discovery service
│   │   ├── upnp.go              # UPnP/SSDP discovery
│   │   └── mdns.go              # mDNS/Bonjour discovery
│   └── config/                  # Configuration management
└── docs/                        # Comprehensive documentation
    ├── websocket-events.md      # WebSocket API documentation
    ├── DISCOVERY.md             # Device discovery guide
    └── API.md                   # HTTP API reference
```

## API Coverage Status

| Endpoint | Method | Status | Description |
|----------|--------|--------|-------------|
| `/info` | GET | ✅ Complete | Device information and capabilities |
| `/name` | GET | ✅ Complete | Device name |
| `/capabilities` | GET | ✅ Complete | Device feature capabilities |
| `/now_playing` | GET | ✅ Complete | Current playback status |
| `/sources` | GET | ✅ Complete | Available audio sources |
| `/sources` | POST | ✅ Complete | Select audio source |
| `/key` | POST | ✅ Complete | Send key commands (24 commands) |
| `/volume` | GET/POST | ✅ Complete | Volume control with safety features |
| `/bass` | GET/POST | ✅ Complete | Bass control (-9 to +9) |
| `/balance` | GET/POST | ✅ Complete | Balance control (-50 to +50) |
| `/presets` | GET | ✅ Complete | Preset configurations (read-only) |
| `/presets` | POST | ❌ Not Supported | **Officially not supported by SoundTouch API** |
| `/clockTime` | GET/POST | ✅ Complete | Device time management |
| `/clockDisplay` | GET/POST | ✅ Complete | Clock display settings |
| `/networkInfo` | GET | ✅ Complete | Network connectivity information |
| **WebSocket** | `/` | ✅ **NEW** | **Real-time event monitoring** |
| **Discovery** | UPnP/mDNS | ✅ Complete | Device discovery services |
| `/reboot` | POST | 🔄 Planned | Device restart |
| `/getZone` | GET | 🔄 Planned | Multiroom zone info |
| `/setZone` | POST | 🔄 Planned | Multiroom zone configuration |

## Testing Coverage

- **Unit Tests**: 150+ test cases covering all functionality
- **Integration Tests**: Real device testing scenarios  
- **Benchmark Tests**: Performance validation
- **WebSocket Tests**: Comprehensive event handling tests
- **Discovery Tests**: Multi-protocol device discovery tests

```go
// Run all tests
go test ./... -v

// Run specific test suites
go test ./pkg/client -v -run TestWebSocket
go test ./pkg/models -v -run TestWebSocket
go test ./pkg/discovery -v

// Run benchmarks
go test ./pkg/client -bench=. 
go test ./pkg/models -bench=.
```

## Quick Start Examples

### Basic HTTP Client

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/user_account/bose-soundtouch/pkg/client"
    "github.com/user_account/bose-soundtouch/pkg/discovery"
)

func main() {
    // Option 1: Connect to known device
    soundtouchClient := client.NewClientFromHost("192.168.1.10")
    
    deviceInfo, err := soundtouchClient.GetDeviceInfo()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Device: %s (%s)\n", deviceInfo.Name, deviceInfo.Type)

    // Option 2: Discover devices automatically (unified: UPnP + mDNS + config)
    cfg, _ := config.LoadFromEnv()
    discoveryService := discovery.NewUnifiedDiscoveryService(cfg)
    ctx := context.Background()
    
    devices, err := discoveryService.DiscoverDevices(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    for _, device := range devices {
        fmt.Printf("Found: %s at %s:%d\n", device.Name, device.Host, device.Port)
    }
    
    // Get current playback status
    nowPlaying, err := soundtouchClient.GetNowPlaying()
    if err != nil {
        log.Fatal(err)
    }
    
    if !nowPlaying.IsEmpty() {
        fmt.Printf("Now Playing: %s by %s\n", 
            nowPlaying.GetDisplayTitle(), 
            nowPlaying.GetDisplayArtist())
    }
    
    // Get available audio sources
    sources, err := soundtouchClient.GetSources()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Ready Sources: %d/%d\n", 
        sources.GetReadySourceCount(), 
        sources.GetSourceCount())
    
    if sources.HasSpotify() {
        fmt.Println("Spotify is available")
    }
    
    // Get device name
    name, err := soundtouchClient.GetName()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Device: %s\n", name.GetName())
    
    // Get device capabilities
    capabilities, err := soundtouchClient.GetCapabilities()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("L/R Stereo Support: %v\n", capabilities.HasLRStereoCapability())
    
    // Get presets
    presets, err := soundtouchClient.GetPresets()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Presets Used: %d/6\n", len(presets.GetUsedPresetSlots()))
    
    // Media controls
    fmt.Println("Playing music...")
    if err := soundtouchClient.Play(); err != nil {
        log.Printf("Failed to play: %v", err)
    }
    
    // Volume control
    volume, err := soundtouchClient.GetVolume()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Current Volume: %d (%s)\n", 
        volume.GetLevel(), 
        volume.GetVolumeString())
    
    // Set comfortable volume
    if err := soundtouchClient.SetVolume(25); err != nil {
        log.Printf("Failed to set volume: %v", err)
    }
}
```

## Project Structure

```
├── cmd/
│   └── soundtouch-cli/     # CLI application
├── pkg/
│   ├── client/             # HTTP client with XML support
│   ├── discovery/          # Device discovery (UPnP/SSDP + mDNS/Bonjour)
│   └── models/             # XML data models
├── docs/                   # Documentation
└── build/                  # Build artifacts
```

## Development

### Prerequisites
- Go 1.25.5 or later
- Make (optional, for convenience)

### Building
```bash
# Build CLI tool
make build

# Build for all platforms
make build-all

# Build and run tests
make check

# Run tests with coverage
make test-coverage
```

### Testing
```bash
# Run all tests
make test

# Run specific package tests
go test -v ./pkg/client
go test -v ./pkg/discovery

# Test with real devices
make dev-info HOST=192.168.1.10

# Test device discovery
make dev-discover

# Test mDNS discovery example
make dev-mdns
```

### Development Commands
```bash
# Format code
make fmt

# Run linter (requires golangci-lint)
make lint

# Clean build artifacts
make clean

# Show help
make help
```

## API Documentation

The SoundTouch Web API uses HTTP with XML payloads. Key endpoints include:

- `GET /info` - Device information ✅ Implemented
- `GET /name` - Device name ✅ Implemented  
- `GET /capabilities` - Device capabilities ✅ Implemented
- `GET /presets` - Configured presets ✅ Implemented
- `GET /now_playing` - Current playback status ✅ Implemented
- `GET /sources` - Available audio sources ✅ Implemented
- `POST /key` - Send key commands (play, pause, etc.) ✅ Implemented
- `GET/POST /volume` - Volume control ✅ Implemented
- `GET/POST /bass` - Bass control ✅ Implemented
- `GET/POST /balance` - Stereo balance control ✅ Implemented
- `POST /select` - Source selection ✅ Implemented
- `GET/POST /clockTime` - Device time management ✅ Implemented
- `GET/POST /clockDisplay` - Clock display settings ✅ Implemented
- `GET /networkInfo` - Network information ✅ Implemented
- WebSocket `/` - Real-time event stream

For complete API documentation, see:
- [API Endpoints Overview](docs/API-Endpoints-Overview.md)
- [Key Controls Documentation](docs/KEY-CONTROLS.md)
- [Volume Controls Documentation](docs/VOLUME-CONTROLS.md)
- [Host:Port Parsing Feature](docs/HOST-PORT-PARSING.md)

## Configuration Options

The application supports configuration through `.env` files and environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `DISCOVERY_TIMEOUT` | `5s` | Timeout for device discovery |
| `UPNP_ENABLED` | `true` | Enable/disable UPnP/SSDP discovery |
| `MDNS_ENABLED` | `true` | Enable/disable mDNS/Bonjour discovery |
| `PREFERRED_DEVICES` | (empty) | Semicolon-separated list of devices |
| `HTTP_TIMEOUT` | `10s` | HTTP client timeout |
| `CACHE_ENABLED` | `true` | Enable device caching |
| `CACHE_TTL` | `30s` | Cache time-to-live |

### Device Configuration Format

The `PREFERRED_DEVICES` environment variable supports multiple formats:

```bash
# Host only (uses default port 8090)
PREFERRED_DEVICES="192.168.1.10"

# Host with port
PREFERRED_DEVICES="192.168.1.10:8091"

# Named device
PREFERRED_DEVICES="Living Room@192.168.1.10"

# Multiple devices
PREFERRED_DEVICES="Living Room@192.168.1.10;Kitchen@192.168.1.11:8091"
```

## Supported Devices

Tested with:
- Bose SoundTouch 10
- Bose SoundTouch 20

Should work with all SoundTouch series devices that support the Web API.

## Real Device Examples

### SoundTouch 10 Response
```xml
<info deviceID="ABCD1234EFGH">
    <name>My SoundTouch Device</name>
    <type>SoundTouch 10</type>
    <moduleType>sm2</moduleType>
    <variant>rhino</variant>
    <components>
        <component>
            <componentCategory>SCM</componentCategory>
            <softwareVersion>27.0.6.46330.5043500</softwareVersion>
        </component>
    </components>
</info>
```

### SoundTouch 20 Response
```xml
<info deviceID="ABCD1234EFGH">
    <name>My SoundTouch Device</name>
    <type>SoundTouch 20</type>
    <moduleType>scm</moduleType>
    <variant>spotty</variant>
    <components>
        <component>
            <componentCategory>SCM</componentCategory>
            <softwareVersion>27.0.6.46330.5043500</softwareVersion>
        </component>
        <component>
            <componentCategory>Lightswitch</componentCategory>
        </component>
    </components>
</info>
```

## Architecture

This project follows modern Go patterns:

- **Clean Architecture**: Separated concerns with pkg structure
- **Interface-Based Design**: Testable and mockable components
- **Cross-Platform**: Supports Windows, macOS, Linux, and WASM
- **Test-Driven**: Comprehensive unit and integration tests
- **Real Device Integration**: Tested with actual SoundTouch hardware

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass: `make check`
5. Submit a pull request

### Development Guidelines

- **Tests are mandatory**: Every feature needs corresponding tests
- **KISS principle**: Keep implementations simple and readable
- **Small iterations**: Break large features into testable chunks
- **Real device testing**: Validate against actual SoundTouch devices
- **Cross-platform compatibility**: Test on multiple platforms

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## References

- [Official Bose SoundTouch Web API Documentation](docs/2025.12.18%20SoundTouch%20Web%20API.pdf)
- [Project Development Plan](docs/PLAN.md)
- [Development Guidelines](docs/CLAUDE.md)