# SoundTouch CLI Reference

**Complete command reference for the soundtouch-cli tool**

This document provides comprehensive documentation for all available commands and options in the `soundtouch-cli` tool.

## Overview

The SoundTouch CLI uses a hierarchical command structure with subcommands for different operations:

```bash
soundtouch-cli [global-flags] <command> [command-flags] [subcommand] [subcommand-flags]
```

## Global Flags

These flags can be used with any command:

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--host` | `-h` | Device IP address or hostname | Required for most commands |
| `--port` | `-p` | Device port number | `8090` |
| `--timeout` | `-t` | Request timeout duration | `10s` |
| `--help` | | Show command help | |
| `--version` | `-v` | Show CLI version | |

## Commands

### Discovery

Discover SoundTouch devices on the network.

#### `discover devices`

Discover and list all SoundTouch devices.

```bash
soundtouch-cli discover devices [flags]
```

**Flags:**
- `--all`, `-a`: Show detailed information for all devices
- `--timeout`: Discovery timeout (default: 10s)

**Examples:**
```bash
# Basic discovery
soundtouch-cli discover devices

# Show detailed info for all discovered devices
soundtouch-cli discover devices --all

# Discovery with custom timeout
soundtouch-cli discover devices --timeout 15s
```

### Device Information

Get information about your SoundTouch device.

#### `info`

Get basic device information.

```bash
soundtouch-cli --host <device> info
```

**Example:**
```bash
soundtouch-cli --host 192.168.1.10 info
```

#### `name get|set`

Get or set the device name.

```bash
# Get current name
soundtouch-cli --host <device> name get

# Set new name
soundtouch-cli --host <device> name set --value "My SoundTouch"
```

#### `capabilities`

Get device capabilities and features.

```bash
soundtouch-cli --host <device> capabilities
```

### Preset Management

Manage device presets (favorite content shortcuts).

#### `preset <subcommand>`

Preset management commands.

```bash
# List all presets
soundtouch-cli --host <device> preset list

# Store currently playing content as preset
soundtouch-cli --host <device> preset store-current --slot <1-6>

# Store specific content as preset
soundtouch-cli --host <device> preset store --slot <1-6> --source <SOURCE> --location <LOCATION> [options]

# Select and play a preset
soundtouch-cli --host <device> preset select --slot <1-6>

# Remove a preset
soundtouch-cli --host <device> preset remove --slot <1-6>
```

**Store Current Content Examples:**
```bash
# Store what's currently playing as preset 1
soundtouch-cli --host 192.168.1.10 preset store-current --slot 1

# Store current Spotify track as preset 3
soundtouch-cli --host 192.168.1.10 preset store-current --slot 3
```

**Store Specific Content Examples:**
```bash
# Store Spotify playlist
soundtouch-cli --host 192.168.1.10 preset store \
  --slot 1 \
  --source SPOTIFY \
  --location "spotify:playlist:37i9dQZF1DXcBWIGoYBM5M" \
  --source-account "your_username" \
  --name "Today's Top Hits"

# Store radio station
soundtouch-cli --host 192.168.1.10 preset store \
  --slot 2 \
  --source TUNEIN \
  --location "/v1/playbook/station/s33828" \
  --name "K-LOVE Radio"

# Store internet radio
soundtouch-cli --host 192.168.1.10 preset store \
  --slot 3 \
  --source LOCAL_INTERNET_RADIO \
  --location "https://stream.example.com/jazz" \
  --name "Jazz Radio Stream"
```

**Selection and Management Examples:**
```bash
# List all presets
soundtouch-cli --host 192.168.1.10 preset list

# Select preset 1
soundtouch-cli --host 192.168.1.10 preset select --slot 1

# Remove preset 6
soundtouch-cli --host 192.168.1.10 preset remove --slot 6
```

**Getting Content Locations:**

To find content locations for the `--location` parameter:

```bash
# Show current content details (includes location for all sources)
soundtouch-cli --host 192.168.1.10 play now

# Show detailed content information
soundtouch-cli --host 192.168.1.10 play now --verbose
```

#### `presets` (Legacy)

Get configured presets (legacy command for backward compatibility).

```bash
soundtouch-cli --host <device> presets
```

### Playback Control

Control music playback on your device.

#### `play <subcommand>`

Playback control commands.

```bash
# Get current playback status
soundtouch-cli --host <device> play now

# Start playback
soundtouch-cli --host <device> play start

# Pause playback
soundtouch-cli --host <device> play pause

# Stop playback
soundtouch-cli --host <device> play stop

# Next track
soundtouch-cli --host <device> play next

# Previous track
soundtouch-cli --host <device> play prev
```

#### `preset`

Select a preset by number.

```bash
soundtouch-cli --host <device> preset --preset <1-6>
```

**Examples:**
```bash
# Select preset 1
soundtouch-cli --host 192.168.1.10 preset --preset 1

# Select preset 6
soundtouch-cli --host 192.168.1.10 preset --preset 6
```

#### `track`

Get current track information.

```bash
soundtouch-cli --host <device> track
```

### Key Commands

Send key commands to the device (simulates remote control).

#### `key <subcommand>`

Send various key commands.

```bash
# Send generic key command
soundtouch-cli --host <device> key send --key <KEY_NAME>

# Specific key commands
soundtouch-cli --host <device> key power
soundtouch-cli --host <device> key mute
soundtouch-cli --host <device> key thumbs-up
soundtouch-cli --host <device> key thumbs-down
soundtouch-cli --host <device> key volume-up
soundtouch-cli --host <device> key volume-down
```

**Available Key Names:**
- `PLAY`, `PAUSE`, `STOP`
- `POWER`, `MUTE`
- `VOLUME_UP`, `VOLUME_DOWN`
- `PRESET_1` through `PRESET_6`
- `NEXT_TRACK`, `PREV_TRACK`
- `THUMBS_UP`, `THUMBS_DOWN`
- `SHUFFLE_ON`, `SHUFFLE_OFF`
- `REPEAT_ON`, `REPEAT_OFF`

### Volume Control

Manage device volume.

#### `volume <subcommand>`

Volume control commands.

```bash
# Get current volume
soundtouch-cli --host <device> volume get

# Set specific volume level (0-100)
soundtouch-cli --host <device> volume set --level <0-100>

# Increase volume
soundtouch-cli --host <device> volume up [--amount <1-10>]

# Decrease volume
soundtouch-cli --host <device> volume down [--amount <1-10>]
```

**Examples:**
```bash
# Get volume
soundtouch-cli --host 192.168.1.10 volume get

# Set volume to 50
soundtouch-cli --host 192.168.1.10 volume set --level 50

# Increase volume by 5
soundtouch-cli --host 192.168.1.10 volume up --amount 5

# Decrease volume by 3 (default amount is 2)
soundtouch-cli --host 192.168.1.10 volume down --amount 3
```

### Audio Sources

Manage audio input sources.

#### `source <subcommand>`

Audio source commands.

```bash
# List available sources
soundtouch-cli --host <device> source list

# Select specific source
soundtouch-cli --host <device> source select --source <SOURCE> [--account <ACCOUNT>]

# Quick source selection
soundtouch-cli --host <device> source spotify
soundtouch-cli --host <device> source bluetooth
soundtouch-cli --host <device> source aux
```

**Source Names:**
- `SPOTIFY` - Spotify streaming
- `BLUETOOTH` - Bluetooth input
- `AUX` - AUX input
- `AIRPLAY` - AirPlay
- `STORED_MUSIC` - Local music library
- `INTERNET_RADIO` - Internet radio
- `PRODUCT` - Product-specific sources

**Examples:**
```bash
# List all sources
soundtouch-cli --host 192.168.1.10 source list

# Select Spotify
soundtouch-cli --host 192.168.1.10 source spotify

# Select Spotify with specific account
soundtouch-cli --host 192.168.1.10 source select --source SPOTIFY --account user@example.com

# Select Bluetooth
soundtouch-cli --host 192.168.1.10 source bluetooth
```

### Bass Control

Adjust bass levels (equalizer).

#### `bass <subcommand>`

Bass control commands.

```bash
# Get current bass level
soundtouch-cli --host <device> bass get

# Set bass level (-9 to 9)
soundtouch-cli --host <device> bass set --level <-9 to 9>

# Increase bass
soundtouch-cli --host <device> bass up [--amount <1-5>]

# Decrease bass
soundtouch-cli --host <device> bass down [--amount <1-5>]

# Get bass capabilities
soundtouch-cli --host <device> bass capabilities
```

**Examples:**
```bash
# Get current bass
soundtouch-cli --host 192.168.1.10 bass get

# Set bass to +3
soundtouch-cli --host 192.168.1.10 bass set --level 3

# Increase bass by 2
soundtouch-cli --host 192.168.1.10 bass up --amount 2

# Decrease bass by 1 (default)
soundtouch-cli --host 192.168.1.10 bass down
```

### Balance Control

Adjust left/right balance.

#### `balance <subcommand>`

Balance control commands.

```bash
# Get current balance
soundtouch-cli --host <device> balance get

# Set balance (-50 to 50, negative=left, positive=right)
soundtouch-cli --host <device> balance set --level <-50 to 50>

# Shift balance left
soundtouch-cli --host <device> balance left [--amount <1-10>]

# Shift balance right
soundtouch-cli --host <device> balance right [--amount <1-10>]

# Center balance
soundtouch-cli --host <device> balance center
```

**Examples:**
```bash
# Get balance
soundtouch-cli --host 192.168.1.10 balance get

# Set balance 10 units to the right
soundtouch-cli --host 192.168.1.10 balance set --level 10

# Shift left by 5 units (default)
soundtouch-cli --host 192.168.1.10 balance left

# Center the balance
soundtouch-cli --host 192.168.1.10 balance center
```

### Clock and Time

Manage device clock settings.

#### `clock <subcommand>`

Clock control commands.

```bash
# Get current time
soundtouch-cli --host <device> clock get

# Set time manually (HH:MM format)
soundtouch-cli --host <device> clock set --time "14:30"

# Set to current system time
soundtouch-cli --host <device> clock now

# Display settings
soundtouch-cli --host <device> clock display get
soundtouch-cli --host <device> clock display enable
soundtouch-cli --host <device> clock display disable
soundtouch-cli --host <device> clock display brightness --brightness <low|medium|high|off>
soundtouch-cli --host <device> clock display format --format <12|24>
```

**Examples:**
```bash
# Get current time
soundtouch-cli --host 192.168.1.10 clock get

# Set time to 2:30 PM
soundtouch-cli --host 192.168.1.10 clock set --time "14:30"

# Sync with system time
soundtouch-cli --host 192.168.1.10 clock now

# Enable clock display
soundtouch-cli --host 192.168.1.10 clock display enable

# Set 24-hour format
soundtouch-cli --host 192.168.1.10 clock display format --format 24

# Set high brightness
soundtouch-cli --host 192.168.1.10 clock display brightness --brightness high
```

### Network Information

Get network and connectivity information.

#### `network <subcommand>`

Network information commands.

```bash
# Get network information
soundtouch-cli --host <device> network info

# Ping the device
soundtouch-cli --host <device> network ping

# Get device base URL
soundtouch-cli --host <device> network url
```

### Zone Management

Manage multi-room zones (multiple speakers playing together).

#### `zone <subcommand>`

Zone management commands.

```bash
# Get current zone configuration
soundtouch-cli --host <device> zone get

# Get zone status
soundtouch-cli --host <device> zone status

# List zone members
soundtouch-cli --host <device> zone members

# Create new zone
soundtouch-cli --host <device> zone create --members <ip1,ip2,ip3>

# Add device to zone
soundtouch-cli --host <device> zone add --member <ip>

# Remove device from zone
soundtouch-cli --host <device> zone remove --member <ip>

# Dissolve current zone
soundtouch-cli --host <device> zone dissolve

# Set zone configuration
soundtouch-cli --host <device> zone set --master <ip> --members <ip1,ip2>
```

**Examples:**
```bash
# Get current zone info
soundtouch-cli --host 192.168.1.10 zone get

# Create zone with three speakers
soundtouch-cli --host 192.168.1.10 zone create --members 192.168.1.11,192.168.1.12

# Add speaker to existing zone
soundtouch-cli --host 192.168.1.10 zone add --member 192.168.1.13

# Remove speaker from zone
soundtouch-cli --host 192.168.1.10 zone remove --member 192.168.1.12

# Dissolve the zone (make all speakers independent)
soundtouch-cli --host 192.168.1.10 zone dissolve
```

### Browse and Navigation

Browse and navigate content sources on your device.

#### `browse <subcommand>`

Browse content from different sources.

```bash
# Browse TuneIn stations
soundtouch-cli --host <device> browse tunein

# Browse Pandora stations (requires account)
soundtouch-cli --host <device> browse pandora --source-account <pandora_account>

# Browse stored music library (requires device ID)
soundtouch-cli --host <device> browse stored-music --source-account <device_id>

# Browse any content source with pagination
soundtouch-cli --host <device> browse content --source <SOURCE> [--start <num>] [--limit <num>]

# Browse with menu navigation (for sources that support it)
soundtouch-cli --host <device> browse menu --source <SOURCE> --menu <MENU_TYPE> [--sort <SORT_ORDER>]

# Browse into a container/directory
soundtouch-cli --host <device> browse container --source <SOURCE> --location <LOCATION> [--type <TYPE>]
```

**Examples:**
```bash
# Browse TuneIn stations
soundtouch-cli --host 192.168.1.10 browse tunein

# Browse first 50 TuneIn stations
soundtouch-cli --host 192.168.1.10 browse tunein --limit 50

# Browse Pandora radio stations
soundtouch-cli --host 192.168.1.10 browse pandora --source-account myuser123

# Browse Pandora with menu navigation
soundtouch-cli --host 192.168.1.10 browse menu --source PANDORA --source-account myuser123 --menu radioStations --sort dateCreated

# Browse stored music library
soundtouch-cli --host 192.168.1.10 browse stored-music --source-account device_12345

# Browse into a music album container
soundtouch-cli --host 192.168.1.10 browse container --source STORED_MUSIC --location "album:983" --type dir
```

### Station Search and Management

Search for and manage radio stations and streaming content.

#### `station <subcommand>`

Search and manage stations.

```bash
# Search across any source
soundtouch-cli --host <device> station search --source <SOURCE> --query <SEARCH_TERM>

# Search TuneIn specifically
soundtouch-cli --host <device> station search-tunein --query <SEARCH_TERM>

# Search Pandora specifically (requires account)
soundtouch-cli --host <device> station search-pandora --source-account <ACCOUNT> --query <SEARCH_TERM>

# Search Spotify specifically (requires account)
soundtouch-cli --host <device> station search-spotify --source-account <ACCOUNT> --query <SEARCH_TERM>

# Add station and play immediately
soundtouch-cli --host <device> station add --source <SOURCE> --token <TOKEN> --name <NAME>

# Remove station from collection
soundtouch-cli --host <device> station remove --source <SOURCE> --location <LOCATION>
```

**Search Examples:**
```bash
# Search TuneIn for jazz stations
soundtouch-cli --host 192.168.1.10 station search-tunein --query "jazz"

# Search Pandora for Taylor Swift
soundtouch-cli --host 192.168.1.10 station search-pandora --source-account myuser123 --query "Taylor Swift"

# Search Spotify for workout playlists
soundtouch-cli --host 192.168.1.10 station search-spotify --source-account spotify_user --query "workout playlist"

# General search across any source
soundtouch-cli --host 192.168.1.10 station search --source TUNEIN --query "classic rock"
```

**Station Management Examples:**
```bash
# Add a station found from search results (use token from search output)
soundtouch-cli --host 192.168.1.10 station add \
  --source TUNEIN \
  --token "c121508" \
  --name "Classic Rock Radio"

# Add Pandora station with account
soundtouch-cli --host 192.168.1.10 station add \
  --source PANDORA \
  --source-account myuser123 \
  --token "TR:12345" \
  --name "My Custom Station"

# Remove a station (use location from browse/search results)
soundtouch-cli --host 192.168.1.10 station remove \
  --source TUNEIN \
  --location "/v1/playbook/station/s33828"
```

**Workflow Example - Discover and Play New Content:**
```bash
# 1. Search for content
soundtouch-cli --host 192.168.1.10 station search-tunein --query "smooth jazz"

# 2. Add interesting station from results (copy token from output)
soundtouch-cli --host 192.168.1.10 station add \
  --source TUNEIN \
  --token "c456789" \
  --name "Smooth Jazz 24/7"

# 3. Station is automatically playing! Or browse for more options:
soundtouch-cli --host 192.168.1.10 browse tunein --limit 10
```

## Common Usage Patterns

### Quick Device Setup

```bash
# Discover devices
soundtouch-cli discover devices

# Get device info
soundtouch-cli --host 192.168.1.10 info

# Set comfortable volume and start playing
soundtouch-cli --host 192.168.1.10 volume set --level 30
soundtouch-cli --host 192.168.1.10 source spotify
soundtouch-cli --host 192.168.1.10 play start
```

### Daily Usage

```bash
# Morning routine
soundtouch-cli --host 192.168.1.10 preset --preset 1  # Morning playlist
soundtouch-cli --host 192.168.1.10 volume set --level 25

# Pause for a call
soundtouch-cli --host 192.168.1.10 play pause

# Resume
soundtouch-cli --host 192.168.1.10 play start

# Evening routine
soundtouch-cli --host 192.168.1.10 preset --preset 3  # Evening playlist
soundtouch-cli --host 192.168.1.10 volume set --level 15
```

### Multi-room Setup

```bash
# Create a zone with living room as master
soundtouch-cli --host 192.168.1.10 zone create --members 192.168.1.11,192.168.1.12

# Control the whole zone from master
soundtouch-cli --host 192.168.1.10 volume set --level 40
soundtouch-cli --host 192.168.1.10 source spotify
soundtouch-cli --host 192.168.1.10 preset --preset 2

# Later, dissolve the zone
soundtouch-cli --host 192.168.1.10 zone dissolve
```

### Audio Tuning

```bash
# Get current audio settings
soundtouch-cli --host 192.168.1.10 volume get
soundtouch-cli --host 192.168.1.10 bass get
soundtouch-cli --host 192.168.1.10 balance get

# Adjust for better sound
soundtouch-cli --host 192.168.1.10 bass set --level 2      # Slight bass boost
soundtouch-cli --host 192.168.1.10 balance set --level -5  # Slightly left
soundtouch-cli --host 192.168.1.10 volume set --level 35   # Good listening level
```

## Error Handling

The CLI provides clear error messages for common issues:

### Device Not Found
```
Error: Failed to connect to device: connection refused
```
**Solutions:**
- Check IP address is correct
- Ensure device is powered on
- Verify network connectivity with `soundtouch-cli --host <device> network ping`

### Invalid Commands
```
Error: unknown command "volumee" for "soundtouch-cli"
```
**Solution:** Check command spelling and structure using `--help`

### Missing Required Flags
```
Error: required flag "host" not set
```
**Solution:** Provide required flags: `--host <device>`

## Getting Help

```bash
# General help
soundtouch-cli --help

# Command-specific help
soundtouch-cli volume --help
soundtouch-cli zone --help

# Subcommand help
soundtouch-cli volume set --help
soundtouch-cli zone create --help
```

## Configuration

### Environment Variables

You can set default values using environment variables:

```bash
export SOUNDTOUCH_HOST=192.168.1.10
export SOUNDTOUCH_PORT=8090
export SOUNDTOUCH_TIMEOUT=15s

# Now you can omit these flags
soundtouch-cli info
soundtouch-cli volume get
```

### Configuration File

Create `~/.soundtouch.env`:

```
SOUNDTOUCH_HOST=192.168.1.10
SOUNDTOUCH_PORT=8090
SOUNDTOUCH_TIMEOUT=15s
SOUNDTOUCH_DISCOVERY_TIMEOUT=10s
```

## See Also

- [Getting Started Guide](GETTING-STARTED.md) - Basic setup and usage
- [WebSocket Events](websocket-events.md) - Real-time monitoring
- [Zone Management](zone-management.md) - Multi-room setup
- [API Endpoints](API-Endpoints-Overview.md) - Complete API reference
