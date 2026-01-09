package main

import (
	"fmt"
	"time"

	"github.com/user_account/bose-soundtouch/pkg/models"
	"github.com/urfave/cli/v2"
)

// getVolume handles getting the current volume level
func getVolume(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Getting current volume", clientConfig.Host, clientConfig.Port)

	volume, err := client.GetVolume()
	if err != nil {
		return fmt.Errorf("failed to get volume: %w", err)
	}

	fmt.Printf("Current Volume:\n")
	fmt.Printf("  Device ID: %s\n", volume.DeviceID)
	fmt.Printf("  Current Level: %d (%s)\n", volume.GetLevel(), models.GetVolumeLevelName(volume.GetLevel()))
	fmt.Printf("  Target Level: %d\n", volume.GetTargetLevel())
	fmt.Printf("  Muted: %v\n", volume.IsMuted())

	if !volume.IsVolumeSync() {
		fmt.Printf("  Note: Volume is adjusting (target: %d, actual: %d)\n", volume.GetTargetLevel(), volume.GetLevel())
	}

	return nil
}

// setVolume handles setting the volume level
func setVolume(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	level := c.Int("level")
	if level < 0 || level > 100 {
		return fmt.Errorf("volume level must be between 0 and 100, got %d", level)
	}

	// Safety warning for loud volumes
	if level > 30 {
		PrintWarning(fmt.Sprintf("Setting volume to %d (this is quite loud!)", level))
		fmt.Printf("Proceeding in 2 seconds... Press Ctrl+C to cancel\n")
		time.Sleep(2 * time.Second)
	}

	PrintDeviceHeader(fmt.Sprintf("Setting volume to %d", level), clientConfig.Host, clientConfig.Port)

	err = client.SetVolume(level)
	if err != nil {
		return fmt.Errorf("failed to set volume: %w", err)
	}

	// Get updated volume to confirm
	volume, err := client.GetVolume()
	if err != nil {
		PrintSuccess("Volume set successfully")
	} else {
		PrintSuccess(fmt.Sprintf("Volume set to %d (%s)", volume.GetLevel(), models.GetVolumeLevelName(volume.GetLevel())))
	}

	return nil
}

// volumeUp handles increasing the volume
func volumeUp(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	amount := c.Int("amount")
	if amount < 1 || amount > 10 {
		return fmt.Errorf("volume increase amount must be between 1 and 10, got %d", amount)
	}

	PrintDeviceHeader(fmt.Sprintf("Increasing volume by %d", amount), clientConfig.Host, clientConfig.Port)

	volume, err := client.IncreaseVolume(amount)
	if err != nil {
		return fmt.Errorf("failed to increase volume: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Volume increased to %d (%s)", volume.GetLevel(), models.GetVolumeLevelName(volume.GetLevel())))

	return nil
}

// volumeDown handles decreasing the volume
func volumeDown(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	amount := c.Int("amount")
	if amount < 1 || amount > 10 {
		return fmt.Errorf("volume decrease amount must be between 1 and 10, got %d", amount)
	}

	PrintDeviceHeader(fmt.Sprintf("Decreasing volume by %d", amount), clientConfig.Host, clientConfig.Port)

	volume, err := client.DecreaseVolume(amount)
	if err != nil {
		return fmt.Errorf("failed to decrease volume: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Volume decreased to %d (%s)", volume.GetLevel(), models.GetVolumeLevelName(volume.GetLevel())))

	return nil
}
