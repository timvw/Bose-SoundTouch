# Introspect CLI Commands Demo

This document demonstrates the usage and output of the new introspect CLI commands added to the soundtouch-cli tool.

## Available Commands

The introspect functionality is available through three commands in the `source` command group:

1. `source introspect` - Get introspect data for any supported service
2. `source introspect-spotify` - Convenience command specifically for Spotify
3. `source introspect-all` - Get introspect data for all available services

## Command Examples and Expected Output

### 1. Basic Spotify Introspect

```bash
$ soundtouch-cli --host 192.168.1.100 source introspect --source SPOTIFY
```

**Expected Output:**
```
⠎⠕⠥⠝⠙⠤⠞⠕⠥⠉⠓ SoundTouch CLI v1.0.0
🔗 Connecting to SoundTouch device at 192.168.1.100:8090

Getting introspect data for SPOTIFY

=== SPOTIFY Service Introspect Data ===
State: InactiveUnselected
User: SpotifyConnectUserName
Currently Playing: ❌ No
Current Content: 
Shuffle Mode: OFF
Subscription Type: 

=== Service State ===
❌ Service is INACTIVE (Never been used)
⏸️  Not currently playing
➡️  Shuffle mode is OFF

=== Service Capabilities ===
❌ ⏮️ Skip Previous
❌ 🎯 Seek within tracks
✅ ▶️ Resume playback
✅ 📊 Data collection: ENABLED

=== Spotify Content History ===
Max History Size: 10 items

=== Technical Details ===
Token Last Changed: 2023-12-14 10:48:15 MST
Token Timestamp: 1702566495 seconds since Unix epoch
Token Microseconds: 427884
Play Status State: 2
Received Playback Request: ❌ No
```

### 2. Spotify Introspect with Account

```bash
$ soundtouch-cli --host 192.168.1.100 source introspect --source SPOTIFY --account my_spotify_user
```

**Expected Output:**
```
⠎⠕⠥⠝⠙⠤⠞⠕⠥⠉⠓ SoundTouch CLI v1.0.0
🔗 Connecting to SoundTouch device at 192.168.1.100:8090

Getting introspect data for SPOTIFY
Source Account: my_spotify_user

=== SPOTIFY Service Introspect Data ===
State: Active
User: my_spotify_user
Currently Playing: ✅ Yes
Current Content: spotify://track/4iV5W9uYEdYUVa79Axb7Rh
Shuffle Mode: ON
Subscription Type: Premium

=== Service State ===
✅ Service is ACTIVE
🎵 Currently playing content
🔀 Shuffle mode is ON

=== Service Capabilities ===
✅ ⏮️ Skip Previous
✅ 🎯 Seek within tracks
✅ ▶️ Resume playback
🚫 Data collection: DISABLED

=== Spotify Content History ===
Max History Size: 15 items

=== Technical Details ===
Token Last Changed: 2023-12-14 15:30:22 MST
Token Timestamp: 1702583422 seconds since Unix epoch
Token Microseconds: 123456
Play Status State: 1
Received Playback Request: ✅ Yes
```

### 3. Spotify Convenience Command

```bash
$ soundtouch-cli --host 192.168.1.100 source introspect-spotify
```

**Expected Output:**
```
⠎⠕⠥⠝⠙⠤⠞⠕⠥⠉⠓ SoundTouch CLI v1.0.0
🔗 Connecting to SoundTouch device at 192.168.1.100:8090

Getting Spotify introspect data

=== Spotify Service Introspect Data ===
State: Active
User: premium_user
Currently Playing: ✅ Yes
Current Content: spotify://playlist/37i9dQZF1DXcBWIGoYBM5M
Shuffle Mode: ON
Subscription Type: Premium

=== Spotify Service State ===
✅ Service is ACTIVE
🎵 Currently playing content
🔀 Shuffle mode is ON

=== Spotify Service Capabilities ===
✅ ⏮️ Skip Previous
✅ 🎯 Seek within tracks
✅ ▶️ Resume playback
🚫 Data collection: DISABLED

💡 Spotify Setup Recommendations:
   (None - service is properly configured and active)

=== Spotify Content History ===
Max History Size: 20 items

=== Technical Details ===
Token Last Changed: 2023-12-14 16:45:10 MST
Token Timestamp: 1702587910 seconds since Unix epoch
Token Microseconds: 789012
Play Status State: 1
Received Playback Request: ✅ Yes
```

### 4. Inactive Service Example

```bash
$ soundtouch-cli --host 192.168.1.100 source introspect-spotify
```

**Expected Output (when Spotify is not set up):**
```
⠎⠕⠥⠝⠙⠤⠞⠕⠥⠉⠓ SoundTouch CLI v1.0.0
🔗 Connecting to SoundTouch device at 192.168.1.100:8090

Getting Spotify introspect data

=== Spotify Service Introspect Data ===
State: InactiveUnselected
User: 
Currently Playing: ❌ No
Current Content: 
Shuffle Mode: OFF
Subscription Type: 

=== Spotify Service State ===
❌ Service is INACTIVE (Never been used)
⏸️  Not currently playing
➡️  Shuffle mode is OFF

=== Spotify Service Capabilities ===
❌ ⏮️ Skip Previous
❌ 🎯 Seek within tracks
✅ ▶️ Resume playback
✅ 📊 Data collection: ENABLED

💡 Spotify Setup Recommendations:
   • Sign in to your Spotify account on the device
   • Use 'soundtouch-cli source select --source SPOTIFY' to activate Spotify
   • Ensure you have Spotify Premium for full functionality
```

### 5. All Services Introspect

```bash
$ soundtouch-cli --host 192.168.1.100 source introspect-all
```

**Expected Output:**
```
⠎⠕⠥⠝⠙⠤⠞⠕⠥⠉⠓ SoundTouch CLI v1.0.0
🔗 Connecting to SoundTouch device at 192.168.1.100:8090

Getting introspect data for all services

🔍 Getting introspect data for SPOTIFY...
✅ SPOTIFY: Successfully retrieved introspect data
   State: Active (User: spotify_user)
   Playing: ✅ Yes | Content: spotify://track/4iV5W9uYEdYUVa79Axb7Rh
   Capabilities: Skip, Seek, Resume

──────────────────────────────────────────────────
🔍 Getting introspect data for PANDORA...
❌ PANDORA: Service not available on this device

──────────────────────────────────────────────────
🔍 Getting introspect data for TUNEIN...
✅ TUNEIN: Successfully retrieved introspect data
   State: Inactive
   Playing: ❌ No
   Capabilities: Resume

──────────────────────────────────────────────────
🔍 Getting introspect data for AMAZON...
❌ AMAZON: Failed to get introspect data - service not configured

──────────────────────────────────────────────────
🔍 Getting introspect data for DEEZER...
❌ DEEZER: Service not available on this device

══════════════════════════════════════════════════
📊 Introspect Summary:
   ✅ Successful: 2 services
   ❌ Failed: 3 services
   📡 Total checked: 5 services

✅ Successfully retrieved introspect data for 2 services
```

### 6. Error Handling Examples

#### Missing Source Parameter
```bash
$ soundtouch-cli --host 192.168.1.100 source introspect
```

**Output:**
```
NAME:
   soundtouch-cli source introspect - Get introspect data for a music service

USAGE:
   soundtouch-cli source introspect [command options]

OPTIONS:
   --account value, -a value  Source account name (optional)
   --source value, -s value   Music service source (SPOTIFY, PANDORA, TUNEIN, etc.)
   --help, -h                 show help

Required flag "source" not set
```

#### Missing Host Parameter
```bash
$ soundtouch-cli source introspect --source SPOTIFY
```

**Output:**
```
host is required. Use --host flag or set SOUNDTOUCH_HOST environment variable
```

#### Invalid Service
```bash
$ soundtouch-cli --host 192.168.1.100 source introspect --source INVALID_SERVICE
```

**Expected Output:**
```
⠎⠕⠥⠝⠙⠤⠞⠕⠥⠉⠓ SoundTouch CLI v1.0.0
🔗 Connecting to SoundTouch device at 192.168.1.100:8090

⚠️  Service INVALID_SERVICE may not be available, but continuing with introspect request...

Getting introspect data for INVALID_SERVICE

❌ Error: failed to get introspect data: HTTP 404: endpoint not found or service not supported
```

## Integration with Other Commands

The introspect commands work well with other CLI commands:

### 1. Check Availability First
```bash
# Check what services are available
$ soundtouch-cli --host 192.168.1.100 source availability

# Then introspect specific services
$ soundtouch-cli --host 192.168.1.100 source introspect --source SPOTIFY
```

### 2. Activate Service After Introspect
```bash
# Check service status
$ soundtouch-cli --host 192.168.1.100 source introspect-spotify

# If inactive, activate it
$ soundtouch-cli --host 192.168.1.100 source select --source SPOTIFY
```

### 3. Compare Sources and Introspect Data
```bash
# Compare configured sources vs available services
$ soundtouch-cli --host 192.168.1.100 source compare

# Get detailed introspect data for specific services
$ soundtouch-cli --host 192.168.1.100 source introspect-all
```

## Environment Variables

The introspect commands respect the same environment variables as other CLI commands:

- `SOUNDTOUCH_HOST` - Default device IP address
- `SOUNDTOUCH_SKIP_AVAILABILITY_CHECK` - Skip service availability validation
- `SOUNDTOUCH_TIMEOUT` - Request timeout duration

**Example:**
```bash
export SOUNDTOUCH_HOST=192.168.1.100
soundtouch-cli source introspect-spotify
```

## Use Cases

### 1. Service Setup Verification
Check if streaming services are properly configured and authenticated:
```bash
soundtouch-cli --host $DEVICE source introspect-spotify
soundtouch-cli --host $DEVICE source introspect --source PANDORA
```

### 2. Troubleshooting Playback Issues
Understand why certain playback controls aren't working:
```bash
# Check if seek is supported
soundtouch-cli --host $DEVICE source introspect --source SPOTIFY | grep -i seek

# Check current playback state
soundtouch-cli --host $DEVICE source introspect-spotify | grep -i playing
```

### 3. Service Health Monitoring
Monitor the health and status of streaming services:
```bash
# Quick health check for all services
soundtouch-cli --host $DEVICE source introspect-all

# Detailed status for critical service
soundtouch-cli --host $DEVICE source introspect-spotify
```

### 4. Account Management
Verify which accounts are associated with services:
```bash
# Check current Spotify account
soundtouch-cli --host $DEVICE source introspect-spotify | grep -i user

# Check with specific account parameter
soundtouch-cli --host $DEVICE source introspect --source SPOTIFY --account specific_user
```

## Tips

1. **Use with grep**: Pipe output to `grep` to filter specific information:
   ```bash
   soundtouch-cli --host $DEVICE source introspect-spotify | grep -E "(State|User|Playing)"
   ```

2. **JSON output**: While not currently implemented, future versions may support JSON output for scripting:
   ```bash
   # Future feature
   soundtouch-cli --host $DEVICE source introspect-spotify --format json
   ```

3. **Batch operations**: Use shell scripting to check multiple devices:
   ```bash
   for device in 192.168.1.100 192.168.1.101; do
     echo "=== Device $device ==="
     soundtouch-cli --host $device source introspect-spotify
   done
   ```

4. **Environment setup**: Set up your environment for easier usage:
   ```bash
   export SOUNDTOUCH_HOST=192.168.1.100
   alias st='soundtouch-cli'
   st source introspect-spotify
   ```
