package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

// getAudioDSPControls gets the current DSP audio controls
func getAudioDSPControls(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Getting DSP audio controls", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	dspControls, err := client.GetAudioDSPControls()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get DSP controls: %v", err))
		return err
	}

	fmt.Println("DSP Audio Controls:")
	fmt.Printf("  Audio Mode: %s\n", dspControls.AudioMode)
	fmt.Printf("  Video Sync Audio Delay: %d ms\n", dspControls.VideoSyncAudioDelay)

	supportedModes := dspControls.GetSupportedAudioModes()
	if len(supportedModes) > 0 {
		fmt.Printf("  Supported Audio Modes: %s\n", strings.Join(supportedModes, ", "))
	}

	return nil
}

// setAudioDSPControls sets the DSP audio controls
func setAudioDSPControls(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	audioMode := c.String("mode")
	videoSyncDelay := c.Int("delay")

	if audioMode == "" && videoSyncDelay == 0 {
		return fmt.Errorf("at least one of --mode or --delay must be specified")
	}

	PrintDeviceHeader("Setting DSP audio controls", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.SetAudioDSPControls(audioMode, videoSyncDelay)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to set DSP controls: %v", err))
		return err
	}

	fmt.Println("✅ DSP controls updated successfully")

	if audioMode != "" {
		fmt.Printf("   Audio Mode: %s\n", audioMode)
	}

	if videoSyncDelay != 0 {
		fmt.Printf("   Video Sync Delay: %d ms\n", videoSyncDelay)
	}

	return nil
}

// setAudioMode sets only the audio mode
func setAudioMode(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	audioMode := c.String("mode")

	if audioMode == "" {
		return fmt.Errorf("audio mode is required (use --mode)")
	}

	PrintDeviceHeader(fmt.Sprintf("Setting audio mode to '%s'", audioMode), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.SetAudioMode(audioMode)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to set audio mode: %v", err))
		return err
	}

	fmt.Printf("✅ Audio mode set to '%s'\n", audioMode)

	return nil
}

// setVideoSyncDelay sets only the video sync audio delay
func setVideoSyncDelay(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	delay := c.Int("delay")

	PrintDeviceHeader(fmt.Sprintf("Setting video sync audio delay to %d ms", delay), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.SetVideoSyncAudioDelay(delay)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to set video sync delay: %v", err))
		return err
	}

	fmt.Printf("✅ Video sync audio delay set to %d ms\n", delay)

	return nil
}

// getAudioToneControls gets the current advanced tone controls
func getAudioToneControls(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Getting advanced tone controls", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	toneControls, err := client.GetAudioProductToneControls()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get tone controls: %v", err))
		return err
	}

	fmt.Println("Advanced Tone Controls:")
	fmt.Printf("  Bass: %d (range: %d to %d, step: %d)\n",
		toneControls.Bass.Value, toneControls.Bass.MinValue, toneControls.Bass.MaxValue, toneControls.Bass.Step)
	fmt.Printf("  Treble: %d (range: %d to %d, step: %d)\n",
		toneControls.Treble.Value, toneControls.Treble.MinValue, toneControls.Treble.MaxValue, toneControls.Treble.Step)

	return nil
}

// setAudioToneControls sets the advanced tone controls
func setAudioToneControls(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	bassStr := c.String("bass")
	trebleStr := c.String("treble")

	if bassStr == "" && trebleStr == "" {
		return fmt.Errorf("at least one of --bass or --treble must be specified")
	}

	var bass, treble *int

	var err error

	if bassStr != "" {
		bassVal, errVal := strconv.Atoi(bassStr)
		if errVal != nil {
			return fmt.Errorf("invalid bass value: %s", bassStr)
		}

		bass = &bassVal
	}

	if trebleStr != "" {
		trebleVal, errVal := strconv.Atoi(trebleStr)
		if errVal != nil {
			return fmt.Errorf("invalid treble value: %s", trebleStr)
		}

		treble = &trebleVal
	}

	PrintDeviceHeader("Setting advanced tone controls", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.SetAudioProductToneControls(bass, treble)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to set tone controls: %v", err))
		return err
	}

	fmt.Println("✅ Advanced tone controls updated successfully")

	if bass != nil {
		fmt.Printf("   Bass: %d\n", *bass)
	}

	if treble != nil {
		fmt.Printf("   Treble: %d\n", *treble)
	}

	return nil
}

// setAdvancedBass sets only the advanced bass control
func setAdvancedBass(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	level := c.Int("level")

	PrintDeviceHeader(fmt.Sprintf("Setting advanced bass to %d", level), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.SetAdvancedBass(level)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to set advanced bass: %v", err))
		return err
	}

	fmt.Printf("✅ Advanced bass set to %d\n", level)

	return nil
}

// setAdvancedTreble sets only the advanced treble control
func setAdvancedTreble(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	level := c.Int("level")

	PrintDeviceHeader(fmt.Sprintf("Setting advanced treble to %d", level), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.SetAdvancedTreble(level)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to set advanced treble: %v", err))
		return err
	}

	fmt.Printf("✅ Advanced treble set to %d\n", level)

	return nil
}

// getAudioLevelControls gets the current speaker level controls
func getAudioLevelControls(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Getting speaker level controls", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	levelControls, err := client.GetAudioProductLevelControls()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get level controls: %v", err))
		return err
	}

	fmt.Println("Speaker Level Controls:")
	fmt.Printf("  Front-Center Speaker: %d (range: %d to %d, step: %d)\n",
		levelControls.FrontCenterSpeakerLevel.Value,
		levelControls.FrontCenterSpeakerLevel.MinValue,
		levelControls.FrontCenterSpeakerLevel.MaxValue,
		levelControls.FrontCenterSpeakerLevel.Step)
	fmt.Printf("  Rear-Surround Speakers: %d (range: %d to %d, step: %d)\n",
		levelControls.RearSurroundSpeakersLevel.Value,
		levelControls.RearSurroundSpeakersLevel.MinValue,
		levelControls.RearSurroundSpeakersLevel.MaxValue,
		levelControls.RearSurroundSpeakersLevel.Step)

	return nil
}

// setAudioLevelControls sets the speaker level controls
func setAudioLevelControls(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	frontCenterStr := c.String("front-center")
	rearSurroundStr := c.String("rear-surround")

	if frontCenterStr == "" && rearSurroundStr == "" {
		return fmt.Errorf("at least one of --front-center or --rear-surround must be specified")
	}

	var frontCenter, rearSurround *int

	var err error

	if frontCenterStr != "" {
		frontCenterVal, errVal := strconv.Atoi(frontCenterStr)
		if errVal != nil {
			return fmt.Errorf("invalid front-center value: %s", frontCenterStr)
		}

		frontCenter = &frontCenterVal
	}

	if rearSurroundStr != "" {
		rearSurroundVal, errVal := strconv.Atoi(rearSurroundStr)
		if errVal != nil {
			return fmt.Errorf("invalid rear-surround value: %s", rearSurroundStr)
		}

		rearSurround = &rearSurroundVal
	}

	PrintDeviceHeader("Setting speaker level controls", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.SetAudioProductLevelControls(frontCenter, rearSurround)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to set level controls: %v", err))
		return err
	}

	fmt.Println("✅ Speaker level controls updated successfully")

	if frontCenter != nil {
		fmt.Printf("   Front-Center Speaker: %d\n", *frontCenter)
	}

	if rearSurround != nil {
		fmt.Printf("   Rear-Surround Speakers: %d\n", *rearSurround)
	}

	return nil
}

// setFrontCenterLevel sets only the front-center speaker level
func setFrontCenterLevel(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	level := c.Int("level")

	PrintDeviceHeader(fmt.Sprintf("Setting front-center speaker level to %d", level), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.SetFrontCenterSpeakerLevel(level)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to set front-center speaker level: %v", err))
		return err
	}

	fmt.Printf("✅ Front-center speaker level set to %d\n", level)

	return nil
}

// setRearSurroundLevel sets only the rear-surround speakers level
func setRearSurroundLevel(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	level := c.Int("level")

	PrintDeviceHeader(fmt.Sprintf("Setting rear-surround speakers level to %d", level), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.SetRearSurroundSpeakersLevel(level)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to set rear-surround speakers level: %v", err))
		return err
	}

	fmt.Printf("✅ Rear-surround speakers level set to %d\n", level)

	return nil
}
