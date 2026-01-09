package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "soundtouch-cli",
		Usage: "Command-line interface for controlling Bose SoundTouch devices",
		Description: `A comprehensive CLI tool for interacting with Bose SoundTouch devices.
   Supports device discovery, playback control, volume/bass/balance adjustment,
   source selection, zone management, and more.`,
		Version: "1.0.0",
		Authors: []*cli.Author{
			{
				Name:  "SoundTouch CLI Contributors",
				Email: "info@example.com",
			},
		},
		Flags: []cli.Flag{},
		Commands: []*cli.Command{
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
						Flags: append(CommonFlags, &cli.BoolFlag{
							Name:    "all",
							Aliases: []string{"a"},
							Usage:   "Show detailed information for all devices",
						}),
					},
				},
			},
			// Device information commands
			{
				Name:    "info",
				Aliases: []string{"i"},
				Usage:   "Get device information",
				Action:  getDeviceInfo,
				Flags:   CommonFlags,
				Before:  RequireHost,
			},
			{
				Name:   "name",
				Usage:  "Get or set device name",
				Flags:  CommonFlags,
				Before: RequireHost,
				Subcommands: []*cli.Command{
					{
						Name:   "get",
						Usage:  "Get device name",
						Action: getDeviceName,
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "set",
						Usage:  "Set device name",
						Action: setDeviceName,
						Flags: append(CommonFlags, &cli.StringFlag{
							Name:     "value",
							Aliases:  []string{"n"},
							Usage:    "New device name",
							Required: true,
						}),
						Before: RequireHost,
					},
				},
			},
			{
				Name:   "capabilities",
				Usage:  "Get device capabilities",
				Action: getCapabilities,
				Flags:  CommonFlags,
				Before: RequireHost,
			},
			{
				Name:   "presets",
				Usage:  "Get configured presets",
				Action: getPresets,
				Flags:  CommonFlags,
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
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "start",
						Usage:  "Start playback",
						Action: playCommand,
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "pause",
						Usage:  "Pause playback",
						Action: pauseCommand,
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "stop",
						Usage:  "Stop playback",
						Action: stopCommand,
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "next",
						Usage:  "Next track",
						Action: nextCommand,
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "prev",
						Usage:  "Previous track",
						Action: prevCommand,
						Flags:  CommonFlags,
						Before: RequireHost,
					},
				},
			},
			// Preset commands
			{
				Name:   "preset",
				Usage:  "Select preset by number",
				Action: selectPreset,
				Flags: append(CommonFlags, &cli.IntFlag{
					Name:     "preset",
					Usage:    "Preset number (1-6)",
					Required: true,
				}),
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
						Flags: append(CommonFlags, &cli.StringFlag{
							Name:     "key",
							Aliases:  []string{"k"},
							Usage:    "Key name (PLAY, PAUSE, STOP, POWER, MUTE, etc.)",
							Required: true,
						}),
						Before: RequireHost,
					},
					{
						Name:   "power",
						Usage:  "Send POWER key command",
						Action: powerCommand,
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "mute",
						Usage:  "Send MUTE key command",
						Action: muteCommand,
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "thumbs-up",
						Usage:  "Send THUMBS_UP key command",
						Action: thumbsUpCommand,
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "thumbs-down",
						Usage:  "Send THUMBS_DOWN key command",
						Action: thumbsDownCommand,
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "volume-up",
						Usage:  "Send VOLUME_UP key command",
						Action: volumeUpKey,
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "volume-down",
						Usage:  "Send VOLUME_DOWN key command",
						Action: volumeDownKey,
						Flags:  CommonFlags,
						Before: RequireHost,
					},
				},
			},
			// Track info
			{
				Name:   "track",
				Usage:  "Get track information",
				Action: getTrackInfo,
				Flags:  CommonFlags,
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
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "set",
						Usage:  "Set volume level",
						Action: setVolume,
						Flags: append(CommonFlags, &cli.IntFlag{
							Name:     "level",
							Aliases:  []string{"l"},
							Usage:    "Volume level (0-100)",
							Required: true,
						}),
						Before: RequireHost,
					},
					{
						Name:   "up",
						Usage:  "Increase volume",
						Action: volumeUp,
						Flags: append(CommonFlags, &cli.IntFlag{
							Name:    "amount",
							Aliases: []string{"a"},
							Usage:   "Amount to increase (1-10)",
							Value:   2,
						}),
						Before: RequireHost,
					},
					{
						Name:   "down",
						Usage:  "Decrease volume",
						Action: volumeDown,
						Flags: append(CommonFlags, &cli.IntFlag{
							Name:    "amount",
							Aliases: []string{"a"},
							Usage:   "Amount to decrease (1-10)",
							Value:   2,
						}),
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
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "select",
						Usage:  "Select an audio source",
						Action: selectSource,
						Flags: append(CommonFlags,
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
						),
						Before: RequireHost,
					},
					{
						Name:   "spotify",
						Usage:  "Select Spotify source",
						Action: selectSpotify,
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "bluetooth",
						Usage:  "Select Bluetooth source",
						Action: selectBluetooth,
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "aux",
						Usage:  "Select AUX input source",
						Action: selectAux,
						Flags:  CommonFlags,
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
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "set",
						Usage:  "Set bass level",
						Action: setBass,
						Flags: append(CommonFlags, &cli.IntFlag{
							Name:     "level",
							Aliases:  []string{"l"},
							Usage:    "Bass level (-9 to 9)",
							Required: true,
						}),
						Before: RequireHost,
					},
					{
						Name:   "up",
						Usage:  "Increase bass",
						Action: bassUp,
						Flags: append(CommonFlags, &cli.IntFlag{
							Name:    "amount",
							Aliases: []string{"a"},
							Usage:   "Amount to increase (1-5)",
							Value:   1,
						}),
						Before: RequireHost,
					},
					{
						Name:   "down",
						Usage:  "Decrease bass",
						Action: bassDown,
						Flags: append(CommonFlags, &cli.IntFlag{
							Name:    "amount",
							Aliases: []string{"a"},
							Usage:   "Amount to decrease (1-5)",
							Value:   1,
						}),
						Before: RequireHost,
					},
					{
						Name:   "capabilities",
						Usage:  "Get bass capabilities",
						Action: getBassCapabilities,
						Flags:  CommonFlags,
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
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "set",
						Usage:  "Set balance level",
						Action: setBalance,
						Flags: append(CommonFlags, &cli.IntFlag{
							Name:     "level",
							Aliases:  []string{"l"},
							Usage:    "Balance level (-50 to 50, negative=left, positive=right)",
							Required: true,
						}),
						Before: RequireHost,
					},
					{
						Name:   "left",
						Usage:  "Shift balance to the left",
						Action: balanceLeft,
						Flags: append(CommonFlags, &cli.IntFlag{
							Name:    "amount",
							Aliases: []string{"a"},
							Usage:   "Amount to shift left (1-10, default: 5)",
							Value:   5,
						}),
						Before: RequireHost,
					},
					{
						Name:   "right",
						Usage:  "Shift balance to the right",
						Action: balanceRight,
						Flags: append(CommonFlags, &cli.IntFlag{
							Name:    "amount",
							Aliases: []string{"a"},
							Usage:   "Amount to shift right (1-10, default: 5)",
							Value:   5,
						}),
						Before: RequireHost,
					},
					{
						Name:   "center",
						Usage:  "Center the balance",
						Action: balanceCenter,
						Flags:  CommonFlags,
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
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "set",
						Usage:  "Set clock time",
						Action: setClockTime,
						Flags: append(CommonFlags, &cli.StringFlag{
							Name:     "time",
							Aliases:  []string{"t"},
							Usage:    "Time in HH:MM format or 'now' for current time",
							Required: true,
						}),
						Before: RequireHost,
					},
					{
						Name:   "now",
						Usage:  "Set clock to current system time",
						Action: setClockTimeNow,
						Flags:  CommonFlags,
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
								Flags:  CommonFlags,
								Before: RequireHost,
							},
							{
								Name:   "enable",
								Usage:  "Enable clock display",
								Action: enableClockDisplay,
								Flags:  CommonFlags,
								Before: RequireHost,
							},
							{
								Name:   "disable",
								Usage:  "Disable clock display",
								Action: disableClockDisplay,
								Flags:  CommonFlags,
								Before: RequireHost,
							},
							{
								Name:   "brightness",
								Usage:  "Set display brightness",
								Action: setClockDisplayBrightness,
								Flags: append(CommonFlags, &cli.StringFlag{
									Name:     "brightness",
									Aliases:  []string{"b"},
									Usage:    "Brightness level (low, medium, high, off)",
									Required: true,
								}),
								Before: RequireHost,
							},
							{
								Name:   "format",
								Usage:  "Set display format",
								Action: setClockDisplayFormat,
								Flags: append(CommonFlags, &cli.StringFlag{
									Name:     "format",
									Aliases:  []string{"f"},
									Usage:    "Time format (12 or 24)",
									Required: true,
								}),
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
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "ping",
						Usage:  "Ping the device",
						Action: pingDevice,
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "url",
						Usage:  "Get device base URL",
						Action: getDeviceURL,
						Flags:  CommonFlags,
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
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "status",
						Usage:  "Get zone status",
						Action: getZoneStatus,
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "members",
						Usage:  "List zone members",
						Action: getZoneMembers,
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "create",
						Usage:  "Create a new zone",
						Action: createZone,
						Flags: append(CommonFlags, &cli.StringSliceFlag{
							Name:     "members",
							Aliases:  []string{"m"},
							Usage:    "Member IP addresses",
							Required: true,
						}),
						Before: RequireHost,
					},
					{
						Name:   "add",
						Usage:  "Add device to zone",
						Action: addToZone,
						Flags: append(CommonFlags, &cli.StringFlag{
							Name:     "member",
							Aliases:  []string{"m"},
							Usage:    "Member IP address to add",
							Required: true,
						}),
						Before: RequireHost,
					},
					{
						Name:   "remove",
						Usage:  "Remove device from zone",
						Action: removeFromZone,
						Flags: append(CommonFlags, &cli.StringFlag{
							Name:     "member",
							Aliases:  []string{"m"},
							Usage:    "Member IP address to remove",
							Required: true,
						}),
						Before: RequireHost,
					},
					{
						Name:   "dissolve",
						Usage:  "Dissolve the current zone",
						Action: dissolveZone,
						Flags:  CommonFlags,
						Before: RequireHost,
					},
					{
						Name:   "set",
						Usage:  "Set zone configuration",
						Action: setZoneConfig,
						Flags: append(CommonFlags,
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
						),
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
