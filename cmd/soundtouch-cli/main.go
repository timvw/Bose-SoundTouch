package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

// Build-time variables injected via ldflags
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	app := &cli.App{
		Name:  "soundtouch-cli",
		Usage: "Command-line interface for controlling Bose SoundTouch devices",
		Description: `A comprehensive CLI tool for interacting with Bose SoundTouch devices.
   Supports device discovery, playback control, volume/bass/balance adjustment,
   source selection, zone management, and more.`,
		Version: version,
		Authors: []*cli.Author{
			{
				Name:  "SoundTouch CLI Contributors",
				Email: "info@example.com",
			},
		},
		Flags: CommonFlags,
		Commands: []*cli.Command{
			// Version commands
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Show detailed version information",
				Action:  showVersionInfo,
			},
			// Discovery commands
			{
				Name:    "discover",
				Aliases: []string{"d"},
				Usage:   "Discover SoundTouch devices on the network",
				Subcommands: []*cli.Command{
					{
						Name:   "devices",
						Usage:  "Discover and list SoundTouch devices",
						Action: discoverDevices,
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "all",
								Aliases: []string{"a"},
								Usage:   "Show detailed information for all devices",
							},
						},
					},
				},
			},
			// Device information commands
			{
				Name:    "info",
				Aliases: []string{"i"},
				Usage:   "Get device information",
				Action:  getDeviceInfo,
				Before:  RequireHost,
			},
			{
				Name:   "name",
				Usage:  "Get or set device name",
				Before: RequireHost,
				Subcommands: []*cli.Command{
					{
						Name:   "get",
						Usage:  "Get device name",
						Action: getDeviceName,
						Before: RequireHost,
					},
					{
						Name:   "set",
						Usage:  "Set device name",
						Action: setDeviceName,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "value",
								Aliases:  []string{"n"},
								Usage:    "New device name",
								Required: true,
							},
						},
						Before: RequireHost,
					},
				},
			},
			{
				Name:   "capabilities",
				Usage:  "Get device capabilities",
				Action: getCapabilities,
				Before: RequireHost,
			},
			{
				Name:   "presets",
				Usage:  "Get configured presets",
				Action: getPresets,
				Before: RequireHost,
			},
			// Playback commands
			{
				Name:    "play",
				Aliases: []string{"p"},
				Usage:   "Playback control commands",
				Subcommands: []*cli.Command{
					{
						Name:   "now",
						Usage:  "Get current playback status",
						Action: getNowPlaying,
						Before: RequireHost,
					},
					{
						Name:   "start",
						Usage:  "Start playback",
						Action: playCommand,
						Before: RequireHost,
					},
					{
						Name:   "pause",
						Usage:  "Pause playback",
						Action: pauseCommand,
						Before: RequireHost,
					},
					{
						Name:   "stop",
						Usage:  "Stop playback",
						Action: stopCommand,
						Before: RequireHost,
					},
					{
						Name:   "next",
						Usage:  "Next track",
						Action: nextCommand,
						Before: RequireHost,
					},
					{
						Name:   "prev",
						Usage:  "Previous track",
						Action: prevCommand,
						Before: RequireHost,
					},
				},
			},
			// Preset commands
			{
				Name:   "preset",
				Usage:  "Select preset by number",
				Action: selectPreset,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "preset",
						Usage:    "Preset number (1-6)",
						Required: true,
					},
				},
				Before: RequireHost,
			},
			// Key commands
			{
				Name:    "key",
				Aliases: []string{"k"},
				Usage:   "Send key commands",
				Subcommands: []*cli.Command{
					{
						Name:   "send",
						Usage:  "Send generic key command",
						Action: sendKey,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "key",
								Aliases:  []string{"k"},
								Usage:    "Key name (PLAY, PAUSE, STOP, POWER, MUTE, etc.)",
								Required: true,
							},
						},
						Before: RequireHost,
					},
					{
						Name:   "power",
						Usage:  "Send POWER key command",
						Action: powerCommand,
						Before: RequireHost,
					},
					{
						Name:   "mute",
						Usage:  "Send MUTE key command",
						Action: muteCommand,
						Before: RequireHost,
					},
					{
						Name:   "thumbs-up",
						Usage:  "Send THUMBS_UP key command",
						Action: thumbsUpCommand,
						Before: RequireHost,
					},
					{
						Name:   "thumbs-down",
						Usage:  "Send THUMBS_DOWN key command",
						Action: thumbsDownCommand,
						Before: RequireHost,
					},
					{
						Name:   "volume-up",
						Usage:  "Send VOLUME_UP key command",
						Action: volumeUpKey,
						Before: RequireHost,
					},
					{
						Name:   "volume-down",
						Usage:  "Send VOLUME_DOWN key command",
						Action: volumeDownKey,
						Before: RequireHost,
					},
				},
			},
			// Track info
			{
				Name:   "track",
				Usage:  "Get track information (WARNING: times out on real devices, use 'now' instead)",
				Action: getTrackInfo,
				Before: RequireHost,
			},
			// Volume commands
			{
				Name:    "volume",
				Aliases: []string{"vol"},
				Usage:   "Volume control commands",
				Subcommands: []*cli.Command{
					{
						Name:   "get",
						Usage:  "Get current volume level",
						Action: getVolume,
						Before: RequireHost,
					},
					{
						Name:   "set",
						Usage:  "Set volume level",
						Action: setVolume,
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:     "level",
								Aliases:  []string{"l"},
								Usage:    "Volume level (0-100)",
								Required: true,
							},
						},
						Before: RequireHost,
					},
					{
						Name:   "up",
						Usage:  "Increase volume",
						Action: volumeUp,
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:    "amount",
								Aliases: []string{"a"},
								Usage:   "Amount to increase (1-10)",
								Value:   2,
							},
						},
						Before: RequireHost,
					},
					{
						Name:   "down",
						Usage:  "Decrease volume",
						Action: volumeDown,
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:    "amount",
								Aliases: []string{"a"},
								Usage:   "Amount to decrease (1-10)",
								Value:   2,
							},
						},
						Before: RequireHost,
					},
				},
			},
			// Source commands
			{
				Name:    "source",
				Aliases: []string{"src"},
				Usage:   "Audio source commands",
				Subcommands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "List available audio sources",
						Action: listSources,
						Before: RequireHost,
					},
					{
						Name:   "select",
						Usage:  "Select an audio source",
						Action: selectSource,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "source",
								Aliases:  []string{"s"},
								Usage:    "Source to select (SPOTIFY, BLUETOOTH, AUX, etc.)",
								Required: true,
							},
							&cli.StringFlag{
								Name:    "account",
								Aliases: []string{"a"},
								Usage:   "Source account for streaming services (optional)",
							},
						},
						Before: RequireHost,
					},
					{
						Name:   "spotify",
						Usage:  "Select Spotify source",
						Action: selectSpotify,
						Before: RequireHost,
					},
					{
						Name:   "bluetooth",
						Usage:  "Select Bluetooth source",
						Action: selectBluetooth,
						Before: RequireHost,
					},
					{
						Name:   "aux",
						Usage:  "Select AUX input source",
						Action: selectAux,
						Before: RequireHost,
					},
				},
			},
			// Bass commands
			{
				Name:    "bass",
				Aliases: []string{"b"},
				Usage:   "Bass control commands",
				Subcommands: []*cli.Command{
					{
						Name:   "get",
						Usage:  "Get current bass level",
						Action: getBass,
						Before: RequireHost,
					},
					{
						Name:   "set",
						Usage:  "Set bass level",
						Action: setBass,
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:     "level",
								Aliases:  []string{"l"},
								Usage:    "Bass level (-9 to 9)",
								Required: true,
							},
						},
						Before: RequireHost,
					},
					{
						Name:   "up",
						Usage:  "Increase bass",
						Action: bassUp,
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:    "amount",
								Aliases: []string{"a"},
								Usage:   "Amount to increase (1-5)",
								Value:   1,
							},
						},
						Before: RequireHost,
					},
					{
						Name:   "down",
						Usage:  "Decrease bass",
						Action: bassDown,
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:    "amount",
								Aliases: []string{"a"},
								Usage:   "Amount to decrease (1-5)",
								Value:   1,
							},
						},
						Before: RequireHost,
					},
					{
						Name:   "capabilities",
						Usage:  "Get bass capabilities",
						Action: getBassCapabilities,
						Before: RequireHost,
					},
				},
			},
			// Balance commands
			{
				Name:    "balance",
				Aliases: []string{"bal"},
				Usage:   "Balance control commands",
				Subcommands: []*cli.Command{
					{
						Name:   "get",
						Usage:  "Get current balance level",
						Action: getBalance,
						Before: RequireHost,
					},
					{
						Name:   "set",
						Usage:  "Set balance level",
						Action: setBalance,
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:     "level",
								Aliases:  []string{"l"},
								Usage:    "Balance level (-50 to 50, negative=left, positive=right)",
								Required: true,
							},
						},
						Before: RequireHost,
					},
					{
						Name:   "left",
						Usage:  "Shift balance to the left",
						Action: balanceLeft,
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:    "amount",
								Aliases: []string{"a"},
								Usage:   "Amount to shift left (1-10, default: 5)",
								Value:   5,
							},
						},
						Before: RequireHost,
					},
					{
						Name:   "right",
						Usage:  "Shift balance to the right",
						Action: balanceRight,
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:    "amount",
								Aliases: []string{"a"},
								Usage:   "Amount to shift right (1-10, default: 5)",
								Value:   5,
							},
						},
						Before: RequireHost,
					},
					{
						Name:   "center",
						Usage:  "Center the balance",
						Action: balanceCenter,
						Before: RequireHost,
					},
				},
			},
			// Clock commands
			{
				Name:    "clock",
				Aliases: []string{"time"},
				Usage:   "Clock and time commands",
				Subcommands: []*cli.Command{
					{
						Name:   "get",
						Usage:  "Get current time",
						Action: getClockTime,
						Before: RequireHost,
					},
					{
						Name:   "set",
						Usage:  "Set clock time",
						Action: setClockTime,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "time",
								Aliases:  []string{"t"},
								Usage:    "Time in HH:MM format or 'now' for current time",
								Required: true,
							},
						},
						Before: RequireHost,
					},
					{
						Name:   "now",
						Usage:  "Set clock to current system time",
						Action: setClockTimeNow,
						Before: RequireHost,
					},
					{
						Name:  "display",
						Usage: "Clock display commands",
						Subcommands: []*cli.Command{
							{
								Name:   "get",
								Usage:  "Get display settings",
								Action: getClockDisplay,
								Before: RequireHost,
							},
							{
								Name:   "enable",
								Usage:  "Enable clock display",
								Action: enableClockDisplay,
								Before: RequireHost,
							},
							{
								Name:   "disable",
								Usage:  "Disable clock display",
								Action: disableClockDisplay,
								Before: RequireHost,
							},
							{
								Name:   "brightness",
								Usage:  "Set display brightness",
								Action: setClockDisplayBrightness,
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "brightness",
										Aliases:  []string{"b"},
										Usage:    "Brightness level (low, medium, high, off)",
										Required: true,
									},
								},
								Before: RequireHost,
							},
							{
								Name:   "format",
								Usage:  "Set display format",
								Action: setClockDisplayFormat,
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "format",
										Aliases:  []string{"f"},
										Usage:    "Time format (12 or 24)",
										Required: true,
									},
								},
								Before: RequireHost,
							},
						},
					},
				},
			},
			// Network commands
			{
				Name:    "network",
				Aliases: []string{"net"},
				Usage:   "Network information commands",
				Subcommands: []*cli.Command{
					{
						Name:   "info",
						Usage:  "Get network information",
						Action: getNetworkInfo,
						Before: RequireHost,
					},
					{
						Name:   "ping",
						Usage:  "Ping the device",
						Action: pingDevice,
						Before: RequireHost,
					},
					{
						Name:   "url",
						Usage:  "Get device base URL",
						Action: getDeviceURL,
						Before: RequireHost,
					},
				},
			},
			// Zone commands
			{
				Name:    "zone",
				Aliases: []string{"z"},
				Usage:   "Multi-room zone management commands",
				Subcommands: []*cli.Command{
					{
						Name:   "get",
						Usage:  "Get current zone configuration",
						Action: getZone,
						Before: RequireHost,
					},
					{
						Name:   "status",
						Usage:  "Get zone status",
						Action: getZoneStatus,
						Before: RequireHost,
					},
					{
						Name:   "members",
						Usage:  "List zone members",
						Action: getZoneMembers,
						Before: RequireHost,
					},
					{
						Name:   "create",
						Usage:  "Create a new zone",
						Action: createZone,
						Flags: []cli.Flag{
							&cli.StringSliceFlag{
								Name:     "members",
								Aliases:  []string{"m"},
								Usage:    "Member IP addresses",
								Required: true,
							},
						},
						Before: RequireHost,
					},
					{
						Name:   "add",
						Usage:  "Add device to zone",
						Action: addToZone,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "member",
								Aliases:  []string{"m"},
								Usage:    "Member IP address to add",
								Required: true,
							},
						},
						Before: RequireHost,
					},
					{
						Name:   "remove",
						Usage:  "Remove device from zone",
						Action: removeFromZone,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "member",
								Aliases:  []string{"m"},
								Usage:    "Member IP address to remove",
								Required: true,
							},
						},
						Before: RequireHost,
					},
					{
						Name:   "dissolve",
						Usage:  "Dissolve the current zone",
						Action: dissolveZone,
						Before: RequireHost,
					},
					{
						Name:   "set",
						Usage:  "Set zone configuration",
						Action: setZoneConfig,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "master",
								Usage:    "Master device IP address",
								Required: true,
							},
							&cli.StringSliceFlag{
								Name:    "members",
								Aliases: []string{"m"},
								Usage:   "Member IP addresses",
							},
						},
						Before: RequireHost,
					},
					{
						Name:   "add-slave",
						Usage:  "Add slave to zone (official API)",
						Action: addZoneSlave,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "master",
								Usage:    "Master device ID",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "slave",
								Usage:    "Slave device ID",
								Required: true,
							},
							&cli.StringFlag{
								Name:  "slave-ip",
								Usage: "Slave device IP address (optional)",
							},
						},
						Before: RequireHost,
					},
					{
						Name:   "remove-slave",
						Usage:  "Remove slave from zone (official API)",
						Action: removeZoneSlave,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "master",
								Usage:    "Master device ID",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "slave",
								Usage:    "Slave device ID",
								Required: true,
							},
							&cli.StringFlag{
								Name:  "slave-ip",
								Usage: "Slave device IP address (optional)",
							},
						},
						Before: RequireHost,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
