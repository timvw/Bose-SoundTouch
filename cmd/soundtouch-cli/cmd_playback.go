package main

import (
	"fmt"
	"strings"

	"github.com/user_account/bose-soundtouch/pkg/models"
	"github.com/urfave/cli/v2"
)

// getNowPlaying handles getting the current playback status
func getNowPlaying(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Getting current playback status", clientConfig.Host, clientConfig.Port)

	nowPlaying, err := client.GetNowPlaying()
	if err != nil {
		return fmt.Errorf("failed to get now playing: %w", err)
	}

	fmt.Printf("Now Playing:\n")
	fmt.Printf("  Device ID: %s\n", nowPlaying.DeviceID)

	if nowPlaying.IsEmpty() {
		fmt.Printf("  Status: No content playing\n")
		return nil
	}

	fmt.Printf("  Source: %s\n", nowPlaying.Source)
	fmt.Printf("  Status: %s\n", nowPlaying.PlayStatus.String())

	if nowPlaying.Track != "" {
		fmt.Printf("  Track: %s\n", nowPlaying.Track)
	}

	if nowPlaying.Artist != "" {
		fmt.Printf("  Artist: %s\n", nowPlaying.Artist)
	}

	if nowPlaying.Album != "" {
		fmt.Printf("  Album: %s\n", nowPlaying.Album)
	}

	if nowPlaying.HasTimeInfo() {
		fmt.Printf("  Duration: %s\n", nowPlaying.FormatDuration())

		if nowPlaying.Position != nil {
			fmt.Printf("  Position: %s\n", nowPlaying.FormatPosition())
		}
	}

	if nowPlaying.StreamType != "" {
		fmt.Printf("  Stream Type: %s\n", nowPlaying.StreamType)
	}

	if nowPlaying.PlayStatus == models.PlayStatusBuffering {
		fmt.Printf("  Note: Content is buffering\n")
	}

	return nil
}

// playCommand handles play command
func playCommand(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Sending play command", clientConfig.Host, clientConfig.Port)

	err = client.SendKeyPressOnly(models.KeyPlay)
	if err != nil {
		return fmt.Errorf("failed to send play command: %w", err)
	}

	PrintSuccess("Play command sent")

	return nil
}

// pauseCommand handles pause command
func pauseCommand(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Sending pause command", clientConfig.Host, clientConfig.Port)

	err = client.SendKeyPressOnly(models.KeyPause)
	if err != nil {
		return fmt.Errorf("failed to send pause command: %w", err)
	}

	PrintSuccess("Pause command sent")

	return nil
}

// stopCommand handles stop command
func stopCommand(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Sending stop command", clientConfig.Host, clientConfig.Port)

	err = client.SendKeyPressOnly(models.KeyStop)
	if err != nil {
		return fmt.Errorf("failed to send stop command: %w", err)
	}

	PrintSuccess("Stop command sent")

	return nil
}

// nextCommand handles next track command
func nextCommand(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Sending next track command", clientConfig.Host, clientConfig.Port)

	err = client.SendKeyPressOnly(models.KeyNextTrack)
	if err != nil {
		return fmt.Errorf("failed to send next track command: %w", err)
	}

	PrintSuccess("Next track command sent")

	return nil
}

// prevCommand handles previous track command
func prevCommand(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Sending previous track command", clientConfig.Host, clientConfig.Port)

	err = client.SendKeyPressOnly(models.KeyPrevTrack)
	if err != nil {
		return fmt.Errorf("failed to send previous track command: %w", err)
	}

	PrintSuccess("Previous track command sent")

	return nil
}

// sendKey sends a generic key command
func sendKey(c *cli.Context) error {
	key := c.String("key")
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader(fmt.Sprintf("Sending %s key command", key), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.SendKey(strings.ToUpper(key))
	if err != nil {
		PrintError(fmt.Sprintf("Failed to send key command: %v", err))
		return err
	}

	PrintSuccess(fmt.Sprintf("%s key command sent", key))

	return nil
}

// powerCommand sends the POWER key command
func powerCommand(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Sending POWER key command", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.SendKey(models.KeyPower)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to send power command: %v", err))
		return err
	}

	PrintSuccess("Power command sent")

	return nil
}

// muteCommand sends the MUTE key command
func muteCommand(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Sending MUTE key command", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.SendKey(models.KeyMute)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to send mute command: %v", err))
		return err
	}

	PrintSuccess("Mute command sent")

	return nil
}

// thumbsUpCommand sends the THUMBS_UP key command
func thumbsUpCommand(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Sending THUMBS_UP key command", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.SendKey(models.KeyThumbsUp)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to send thumbs up command: %v", err))
		return err
	}

	PrintSuccess("Thumbs up command sent")

	return nil
}

// thumbsDownCommand sends the THUMBS_DOWN key command
func thumbsDownCommand(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Sending THUMBS_DOWN key command", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.SendKey(models.KeyThumbsDown)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to send thumbs down command: %v", err))
		return err
	}

	PrintSuccess("Thumbs down command sent")

	return nil
}

// volumeUpKey sends the VOLUME_UP key command
func volumeUpKey(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Sending VOLUME_UP key command", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.VolumeUp()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to send volume up command: %v", err))
		return err
	}

	PrintSuccess("Volume up command sent")

	return nil
}

// volumeDownKey sends the VOLUME_DOWN key command
func volumeDownKey(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Sending VOLUME_DOWN key command", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.VolumeDown()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to send volume down command: %v", err))
		return err
	}

	PrintSuccess("Volume down command sent")

	return nil
}
