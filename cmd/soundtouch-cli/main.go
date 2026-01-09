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
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Enable verbose output",
			},
		},
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
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "all",
								Aliases: []string{"a"},
								Usage:   "Show detailed information for all devices",
							},
							&cli.DurationFlag{
								Name:    "timeout",
								Aliases: []string{"t"},
								Usage:   "Discovery timeout",
								Value:   10000000000, // 10 seconds in nanoseconds
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
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
