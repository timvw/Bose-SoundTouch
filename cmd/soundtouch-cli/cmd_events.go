// Package main provides the soundtouch-cli events command for WebSocket event monitoring.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gesellix/bose-soundtouch/pkg/client"
	"github.com/gesellix/bose-soundtouch/pkg/models"
	"github.com/urfave/cli/v2"
)

// eventSubscribe handles the events subscribe command
func eventSubscribe(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	// Parse filters
	filterStr := c.String("filter")
	filters := parseEventFilters(filterStr)

	// Parse duration
	duration := c.Duration("duration")
	verbose := c.Bool("verbose")
	reconnect := !c.Bool("no-reconnect")

	PrintDeviceHeader("Starting WebSocket event monitoring", clientConfig.Host, clientConfig.Port)

	// Create SoundTouch client
	soundTouchClient, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	// Test basic connectivity
	fmt.Println("Testing device connectivity...")

	deviceInfo, err := soundTouchClient.GetDeviceInfo()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to connect to device: %v", err))
		return err
	}

	macAddress := ""
	if len(deviceInfo.NetworkInfo) > 0 {
		macAddress = deviceInfo.NetworkInfo[0].MacAddress
	}

	fmt.Printf("✅ Connected to: %s (Type: %s, MAC: %s)\n",
		deviceInfo.Name, deviceInfo.Type, macAddress)

	// Create WebSocket client
	wsClient := setupWebSocketClient(soundTouchClient, reconnect, verbose)

	// Set up event handlers
	setupEventHandlers(wsClient, filters, verbose)

	// Connect to WebSocket
	fmt.Println("🔌 Connecting to WebSocket...")

	err = wsClient.Connect()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to connect to WebSocket: %v", err))
		return err
	}

	fmt.Println("✅ Connected! Listening for events...")

	if len(filters) > 0 {
		fmt.Printf("📋 Filtering events: %s\n", strings.Join(getFilterKeys(filters), ", "))
	}

	if duration > 0 {
		fmt.Printf("⏰ Will listen for %v\n", duration)
	} else {
		fmt.Println("⏸️  Press Ctrl+C to stop")
	}

	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle duration limit
	if duration > 0 {
		go func() {
			select {
			case <-time.After(duration):
				fmt.Println("\n⏰ Duration limit reached, shutting down...")
				cancel()
			case <-ctx.Done():
				return
			}
		}()
	}

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-sigChan:
			fmt.Printf("\n🛑 Received signal %v, shutting down...\n", sig)
			cancel()
		case <-ctx.Done():
			return
		}
	}()

	// Wait for shutdown
	<-ctx.Done()

	// Disconnect WebSocket
	fmt.Println("🔌 Disconnecting...")

	if err := wsClient.Disconnect(); err != nil {
		PrintError(fmt.Sprintf("Error during disconnect: %v", err))
	}

	fmt.Println("✅ Disconnected successfully")

	return nil
}

// parseEventFilters validates and parses the filter string
func parseEventFilters(eventFilter string) map[string]bool {
	validFilters := map[string]bool{
		"nowPlaying": true, "volume": true, "connection": true,
		"preset": true, "zone": true, "bass": true,
		"sdkInfo": true, "userActivity": true,
	}

	if eventFilter == "" {
		return nil
	}

	filters := make(map[string]bool)
	filterList := strings.Split(eventFilter, ",")

	for _, f := range filterList {
		f = strings.TrimSpace(f)
		if !validFilters[f] {
			PrintError(fmt.Sprintf("Invalid filter '%s'. Valid filters: %s",
				f, strings.Join(getFilterKeys(validFilters), ", ")))
			os.Exit(1)
		}

		filters[f] = true
	}

	return filters
}

// setupWebSocketClient creates and configures the WebSocket client
func setupWebSocketClient(soundTouchClient *client.Client, reconnect, verbose bool) *client.WebSocketClient {
	wsConfig := &client.WebSocketConfig{
		ReconnectInterval:    5 * time.Second,
		MaxReconnectAttempts: 0, // Unlimited if reconnect enabled
		PingInterval:         30 * time.Second,
		PongTimeout:          10 * time.Second,
		ReadBufferSize:       2048,
		WriteBufferSize:      2048,
	}

	if verbose {
		wsConfig.Logger = &VerboseLogger{}
	} else {
		wsConfig.Logger = &SilentLogger{}
	}

	if !reconnect {
		wsConfig.MaxReconnectAttempts = 1
	}

	return soundTouchClient.NewWebSocketClient(wsConfig)
}

// setupEventHandlers configures all event handlers
func setupEventHandlers(wsClient *client.WebSocketClient, filters map[string]bool, verbose bool) {
	// Now Playing events
	if filters == nil || filters["nowPlaying"] {
		wsClient.OnNowPlaying(func(event *models.NowPlayingUpdatedEvent) {
			handleNowPlayingEvent(event, verbose)
		})
	}

	// Volume events
	if filters == nil || filters["volume"] {
		wsClient.OnVolumeUpdated(func(event *models.VolumeUpdatedEvent) {
			handleVolumeEvent(event, verbose)
		})
	}

	// Connection state events
	if filters == nil || filters["connection"] {
		wsClient.OnConnectionState(func(event *models.ConnectionStateUpdatedEvent) {
			handleConnectionEvent(event)
		})
	}

	// Preset events
	if filters == nil || filters["preset"] {
		wsClient.OnPresetUpdated(func(event *models.PresetUpdatedEvent) {
			handlePresetEvent(event, verbose)
		})
	}

	// Zone/Multiroom events
	if filters == nil || filters["zone"] {
		wsClient.OnZoneUpdated(func(event *models.ZoneUpdatedEvent) {
			handleZoneEvent(event)
		})
	}

	// Bass events
	if filters == nil || filters["bass"] {
		wsClient.OnBassUpdated(func(event *models.BassUpdatedEvent) {
			handleBassEvent(event)
		})
	}

	// Special message handler
	wsClient.OnSpecialMessage(func(message *models.SpecialMessage) {
		handleSpecialMessage(message, filters, verbose)
	})

	// Unknown events (always enabled for debugging)
	wsClient.OnUnknownEvent(func(event *models.WebSocketEvent) {
		handleUnknownEvent(event, verbose)
	})
}

// Event handlers
func handleNowPlayingEvent(event *models.NowPlayingUpdatedEvent, verbose bool) {
	fmt.Printf("\n🎵 Now Playing Update [%s]:\n", event.DeviceID)
	np := &event.NowPlaying

	if np.IsEmpty() {
		fmt.Println("  ⏹️  Nothing playing")
		return
	}

	fmt.Printf("  🎵 %s\n", np.GetDisplayTitle())

	if artist := np.GetDisplayArtist(); artist != "" {
		fmt.Printf("  👤 %s\n", artist)
	}

	if np.Album != "" {
		fmt.Printf("  💿 %s\n", np.Album)
	}

	fmt.Printf("  📻 Source: %s\n", np.Source)
	fmt.Printf("  ▶️  Status: %s\n", np.PlayStatus.String())

	if np.HasTimeInfo() {
		fmt.Printf("  ⏱️  Duration: %s\n", np.FormatDuration())
	}

	if np.ShuffleSetting != "" {
		fmt.Printf("  🔀 Shuffle: %s\n", np.ShuffleSetting.String())
	}

	if np.RepeatSetting != "" {
		fmt.Printf("  🔁 Repeat: %s\n", np.RepeatSetting.String())
	}

	if verbose {
		fmt.Printf("  📱 Raw Source: %s, Account: %s\n", np.Source, np.SourceAccount)

		if np.Art != nil && np.Art.URL != "" {
			fmt.Printf("  🖼️  Artwork: %s\n", np.Art.URL)
		}
	}
}

func handleVolumeEvent(event *models.VolumeUpdatedEvent, verbose bool) {
	vol := &event.Volume
	fmt.Printf("\n🔊 Volume Update [%s]:\n", event.DeviceID)

	if vol.IsMuted() {
		fmt.Println("  🔇 Muted")
	} else {
		fmt.Printf("  🔊 Level: %d\n", vol.ActualVolume)

		if vol.TargetVolume != vol.ActualVolume {
			fmt.Printf("  🎯 Target: %d\n", vol.TargetVolume)
		}

		fmt.Printf("  📊 %s\n", models.GetVolumeLevelName(vol.ActualVolume))
	}

	if verbose {
		fmt.Printf("  📱 Sync: %v\n", vol.IsVolumeSync())
	}
}

func handleConnectionEvent(event *models.ConnectionStateUpdatedEvent) {
	cs := &event.ConnectionState
	fmt.Printf("\n🌐 Connection Update [%s]:\n", event.DeviceID)

	if cs.IsConnected() {
		fmt.Println("  ✅ Connected")
	} else {
		fmt.Printf("  ❌ State: %s\n", cs.State)
	}

	if cs.Signal != "" {
		fmt.Printf("  📶 Signal: %s\n", cs.GetSignalStrength())
	}
}

func handlePresetEvent(event *models.PresetUpdatedEvent, verbose bool) {
	presets := &event.Presets

	deviceHeader := "\n📻 Presets Update"
	if event.DeviceID != "" {
		deviceHeader += fmt.Sprintf(" [%s]", event.DeviceID)
	}

	fmt.Printf("%s:\n", deviceHeader)
	fmt.Printf("  📻 Total presets: %d\n", len(presets.Preset))

	for _, preset := range presets.Preset {
		fmt.Printf("  📻 Preset %d:", preset.ID)

		if preset.ContentItem != nil {
			fmt.Printf(" %s", preset.ContentItem.ItemName)
			fmt.Printf(" (%s)", preset.ContentItem.Source)
		}

		fmt.Println()
	}

	if verbose {
		fmt.Printf("  📱 Raw presets data: %d total presets\n", len(presets.Preset))
	}
}

func handleZoneEvent(event *models.ZoneUpdatedEvent) {
	zone := &event.Zone
	fmt.Printf("\n🏠 Zone Update [%s]:\n", event.DeviceID)
	fmt.Printf("  👑 Master: %s\n", zone.Master)

	if len(zone.Members) > 0 {
		fmt.Printf("  👥 Members (%d):\n", len(zone.Members))

		for i, member := range zone.Members {
			fmt.Printf("    %d. %s (%s)\n", i+1, member.DeviceID, member.IP)
		}
	} else {
		fmt.Println("  👤 Single device (no zone)")
	}
}

func handleBassEvent(event *models.BassUpdatedEvent) {
	bass := &event.Bass
	fmt.Printf("\n🎵 Bass Update [%s]:\n", event.DeviceID)
	fmt.Printf("  🎚️  Level: %d\n", bass.ActualBass)

	if bass.TargetBass != bass.ActualBass {
		fmt.Printf("  🎯 Target: %d\n", bass.TargetBass)
	}

	levelDesc := "Neutral"
	if bass.ActualBass > 0 {
		levelDesc = "Boosted"
	} else if bass.ActualBass < 0 {
		levelDesc = "Reduced"
	}

	fmt.Printf("  📊 %s\n", levelDesc)
}

func handleSpecialMessage(message *models.SpecialMessage, filters map[string]bool, verbose bool) {
	// Check if we should filter this message type
	if filters != nil {
		switch message.Type {
		case models.MessageTypeSdkInfo:
			if !filters["sdkInfo"] {
				return
			}
		case models.MessageTypeUserActivity:
			if !filters["userActivity"] {
				return
			}
		}
	}

	switch message.Type {
	case models.MessageTypeSdkInfo:
		if sdkInfo := message.GetSdkInfo(); sdkInfo != nil {
			fmt.Printf("\n📡 SDK Info:\n")
			fmt.Printf("  📋 Server Version: %s\n", sdkInfo.ServerVersion)
			fmt.Printf("  🔧 Server Build: %s\n", sdkInfo.ServerBuild)
		}
	case models.MessageTypeUserActivity:
		fmt.Printf("\n👤 User Activity [%s]\n", message.DeviceID)

		if verbose {
			fmt.Printf("  ⏰ Timestamp: %s\n", message.Timestamp.Format("15:04:05"))
		}
	default:
		fmt.Printf("\n❓ Unknown Special Message: %s\n", message.String())

		if verbose {
			fmt.Printf("  📱 Raw data: %s\n", string(message.RawData))
		}
	}
}

func handleUnknownEvent(event *models.WebSocketEvent, verbose bool) {
	fmt.Printf("\n❓ Unknown Event [%s]:\n", event.DeviceID)
	types := event.GetEventTypes()

	for _, eventType := range types {
		fmt.Printf("  📝 Type: %s\n", eventType)
	}

	if verbose {
		events := event.GetEvents()
		fmt.Printf("  📱 Event count: %d\n", len(events))
		fmt.Printf("  ⏰ Timestamp: %s\n", event.Timestamp.Format(time.RFC3339))
	}
}

// getFilterKeys extracts keys from filter map
func getFilterKeys(filters map[string]bool) []string {
	var keys []string
	for k := range filters {
		keys = append(keys, k)
	}

	return keys
}

// Logger implementations
type VerboseLogger struct{}

func (v *VerboseLogger) Printf(format string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("[%s] [WebSocket] %s\n", timestamp, fmt.Sprintf(format, args...))
}

type SilentLogger struct{}

func (s *SilentLogger) Printf(_ string, _ ...interface{}) {
	// Do nothing - silent logging
}
