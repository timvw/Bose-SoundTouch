# Bose SoundTouch API Client

A comprehensive Go library and CLI tool for controlling Bose SoundTouch devices via their Web API.

[![Go Reference](https://pkg.go.dev/badge/github.com/gesellix/bose-soundtouch.svg)](https://pkg.go.dev/github.com/gesellix/bose-soundtouch)
[![Go Report Card](https://goreportcard.com/badge/github.com/gesellix/bose-soundtouch)](https://goreportcard.com/report/github.com/gesellix/bose-soundtouch)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **Note**: This is an independent project based on the [official Bose SoundTouch Web API documentation](https://assets.bosecreative.com/m/496577402d128874/original/SoundTouch-Web-API.pdf). Not affiliated with or endorsed by Bose Corporation.

## Features

- ✅ **Complete API Coverage**: All available SoundTouch Web API endpoints implemented
- 🎵 **Media Control**: Play, pause, stop, volume, bass, balance, source selection
- 🏠 **Multiroom Support**: Create and manage zones across multiple speakers
- ⚡ **Real-time Events**: WebSocket connection for live device state monitoring
- 🔍 **Device Discovery**: Automatic discovery via UPnP/SSDP and mDNS
- 🖥️ **CLI Tool**: Comprehensive command-line interface
- 🔒 **Production Ready**: Extensive testing with real SoundTouch hardware
- 🌐 **Cross-Platform**: Windows, macOS, Linux support

## Quick Start

### Installation

#### Install CLI Tool
```bash
go install github.com/gesellix/bose-soundtouch/cmd/soundtouch-cli@latest
```

#### Add Library to Your Project
```bash
go get github.com/gesellix/bose-soundtouch
```

### CLI Usage

#### Discover Devices
```bash
# Find SoundTouch devices on your network
soundtouch-cli discover devices
```

#### Control a Device
```bash
# Basic device information
soundtouch-cli --host 192.168.1.100 info get

# Media controls
soundtouch-cli --host 192.168.1.100 play start
soundtouch-cli --host 192.168.1.100 volume set --level 50
soundtouch-cli --host 192.168.1.100 source select --source SPOTIFY

# Real-time monitoring
soundtouch-cli --host 192.168.1.100 events subscribe
```

### Library Usage

#### Basic Control
```go
package main

import (
    "fmt"
    "log"
    
    "github.com/gesellix/bose-soundtouch/pkg/client"
)

func main() {
    // Connect to your SoundTouch device
    c := client.NewClient(&client.Config{
        Host: "192.168.1.100",
        Port: 8090,
    })
    
    // Get device information
    info, err := c.GetDeviceInfo()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Device: %s\n", info.Name)
    
    // Control playback
    err = c.Play()
    if err != nil {
        log.Fatal(err)
    }
    
    // Set volume
    err = c.SetVolume(50)
    if err != nil {
        log.Fatal(err)
    }
}
```

#### Device Discovery
```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/gesellix/bose-soundtouch/pkg/discovery"
)

func main() {
    // Discover SoundTouch devices
    service := discovery.NewService(5 * time.Second)
    devices, err := service.DiscoverDevices(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    
    for _, device := range devices {
        fmt.Printf("Found: %s at %s:%d\n", 
            device.Name, device.Host, device.Port)
    }
}
```

#### Real-time Events
```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/gesellix/bose-soundtouch/pkg/client"
    "github.com/gesellix/bose-soundtouch/pkg/models"
)

func main() {
    c := client.NewClient(&client.Config{
        Host: "192.168.1.100",
        Port: 8090,
    })
    
    // Subscribe to device events
    events, err := c.SubscribeToEvents(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    
    for event := range events {
        switch e := event.(type) {
        case *models.NowPlayingUpdated:
            fmt.Printf("Now playing: %s by %s\n", e.Track, e.Artist)
        case *models.VolumeUpdated:
            fmt.Printf("Volume changed to: %d\n", e.ActualVolume)
        case *models.ConnectionStateUpdated:
            fmt.Printf("Connection state: %s\n", e.State)
        }
    }
}
```

#### Multiroom Zones
```go
package main

import (
    "log"
    
    "github.com/gesellix/bose-soundtouch/pkg/client"
    "github.com/gesellix/bose-soundtouch/pkg/models"
)

func main() {
    master := client.NewClient(&client.Config{
        Host: "192.168.1.100", // Master speaker
        Port: 8090,
    })
    
    // Create a multiroom zone
    zone := &models.Zone{
        Master: "192.168.1.100",
        Members: []models.ZoneMember{
            {IPAddress: "192.168.1.101"}, // Living room
            {IPAddress: "192.168.1.102"}, // Kitchen
        },
    }
    
    err := master.SetZone(zone)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Multiroom zone created!")
}
```

## Supported Devices

This library supports all Bose SoundTouch-compatible devices, including:

- SoundTouch 10, 20, 30 series
- SoundTouch Portable  
- Wave SoundTouch music system
- SoundTouch-enabled Bose speakers

**Tested Hardware**:
- ✅ SoundTouch 10
- ✅ SoundTouch 20

## API Coverage

| Feature | Status | Description |
|---------|--------|-------------|
| Device Info | ✅ Complete | Device details, name, capabilities |
| Media Control | ✅ Complete | Play/pause/stop, track navigation |
| Volume & Audio | ✅ Complete | Volume, bass, balance control |
| Source Selection | ✅ Complete | Spotify, Bluetooth, AUX, etc. |
| Preset Management | ✅ Complete | Read preset configurations |
| Real-time Events | ✅ Complete | WebSocket event streaming |
| Multiroom Zones | ✅ Complete | Zone creation and management |
| System Settings | ✅ Complete | Clock, display, network info |
| Advanced Audio | ✅ Complete | DSP controls, tone controls |

**API Limitations**: Preset creation is not supported by the SoundTouch API itself.

## Documentation

- 📖 [Contributing Guide](CONTRIBUTING.md) - How to contribute to the project
- 📚 [API Reference](docs/API-Endpoints-Overview.md) - Complete endpoint documentation
- 🔧 [CLI Reference](docs/CLI-REFERENCE.md) - Command-line tool guide
- 🎯 [Getting Started](docs/GETTING-STARTED.md) - Detailed setup and usage
- ⚙️ [Advanced Features](docs/SYSTEM-ENDPOINTS.md) - Advanced functionality
- 🏠 [Multiroom Setup](docs/zone-management.md) - Zone configuration guide
- ⚡ [WebSocket Events](docs/websocket-events.md) - Real-time event handling
- 🔍 [Device Discovery](docs/DISCOVERY.md) - Discovery configuration
- 🛠️ [Troubleshooting](docs/TROUBLESHOOTING.md) - Common issues and solutions

## Development

### Prerequisites
- Go 1.21 or later
- Optional: SoundTouch device for testing

### Building from Source
```bash
# Clone the repository
git clone https://github.com/gesellix/bose-soundtouch.git
cd Bose-SoundTouch

# Install dependencies
go mod download

# Build CLI tool
make build

# Run tests
make test

# Install CLI locally
go install ./cmd/soundtouch-cli
```

### Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details on:

- Setting up your development environment
- Coding guidelines and best practices
- Testing with real devices
- Submitting pull requests

## Examples

Check out the [examples/](examples/) directory for more usage patterns:

- **Basic HTTP Client**: Simple device control
- **WebSocket Events**: Real-time monitoring
- **Device Discovery**: Finding devices on your network
- **Multiroom Management**: Zone operations
- **Advanced Audio**: DSP and tone controls

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This is an independent project created by reverse-engineering and documenting the Bose SoundTouch Web API. It is not affiliated with, endorsed by, or supported by Bose Corporation. Use at your own risk.

SoundTouch is a trademark of Bose Corporation.

## Support

- 🐛 **Bug Reports**: [Create an issue](https://github.com/gesellix/bose-soundtouch/issues/new)
- 💡 **Feature Requests**: [Start a discussion](https://github.com/gesellix/bose-soundtouch/discussions)
- ❓ **Questions**: Check [existing discussions](https://github.com/gesellix/bose-soundtouch/discussions)
- 📖 **Documentation**: Browse the [docs/](docs/) directory

---

**Star this project** ⭐ if you find it useful!