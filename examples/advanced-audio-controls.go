package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gesellix/bose-soundtouch/pkg/client"
)

func main() {
	// Configure your device
	deviceIP := "192.168.1.100" // Replace with your SoundTouch device IP

	// Create client
	soundtouchClient := client.NewClientFromHost(deviceIP)

	fmt.Println("🎵 Bose SoundTouch Advanced Audio Controls Example")
	fmt.Println("=================================================")

	// Example 1: Check device capabilities first
	fmt.Println("\n1. Checking device capabilities...")
	capabilities, err := soundtouchClient.GetCapabilities()
	if err != nil {
		log.Printf("❌ Failed to get capabilities: %v", err)
		return
	}

	fmt.Printf("📋 Device: %s\n", capabilities.DeviceID)
	fmt.Printf("   Type: %s\n", capabilities.Type)

	// Look for advanced audio capabilities in the response
	// (Note: Advanced audio controls are only available on professional/high-end devices)
	fmt.Println("   Advanced Audio Features:")
	fmt.Println("   - DSP Controls: Check device response for 'audiodspcontrols'")
	fmt.Println("   - Tone Controls: Check device response for 'audioproducttonecontrols'")
	fmt.Println("   - Level Controls: Check device response for 'audioproductlevelcontrols'")

	// Example 2: DSP Audio Controls
	fmt.Println("\n2. DSP Audio Controls...")

	dspControls, err := soundtouchClient.GetAudioDSPControls()
	if err != nil {
		log.Printf("⚠️  DSP controls not available on this device: %v", err)
		fmt.Println("   This is normal for consumer-grade SoundTouch devices")
	} else {
		fmt.Printf("🎛️  Current DSP Settings: %s\n", dspControls.String())

		// Try setting a different audio mode
		supportedModes := dspControls.GetSupportedAudioModes()
		if len(supportedModes) > 0 {
			newMode := supportedModes[0]
			if newMode != dspControls.AudioMode && newMode != "" {
				fmt.Printf("   Changing audio mode to: %s\n", newMode)

				err = soundtouchClient.SetAudioMode(newMode)
				if err != nil {
					log.Printf("❌ Failed to set audio mode: %v", err)
				} else {
					fmt.Printf("✅ Audio mode changed successfully\n")
				}
			}
		}

		// Demonstrate video sync delay adjustment
		if dspControls.VideoSyncAudioDelay != 50 {
			fmt.Println("   Setting video sync audio delay to 50ms...")
			err = soundtouchClient.SetVideoSyncAudioDelay(50)
			if err != nil {
				log.Printf("❌ Failed to set video sync delay: %v", err)
			} else {
				fmt.Printf("✅ Video sync delay adjusted\n")
			}
		}

		// Combined DSP settings update
		fmt.Println("   Updating DSP controls (mode + delay)...")
		err = soundtouchClient.SetAudioDSPControls("NORMAL", 25)
		if err != nil {
			log.Printf("❌ Failed to set DSP controls: %v", err)
		} else {
			fmt.Printf("✅ DSP controls updated\n")
		}
	}

	time.Sleep(2 * time.Second)

	// Example 3: Advanced Tone Controls (Bass/Treble)
	fmt.Println("\n3. Advanced Tone Controls...")

	toneControls, err := soundtouchClient.GetAudioProductToneControls()
	if err != nil {
		log.Printf("⚠️  Advanced tone controls not available on this device: %v", err)
		fmt.Println("   Use the basic bass control instead (soundtouch-cli bass)")
	} else {
		fmt.Printf("🎚️  Current Tone Settings: %s\n", toneControls.String())

		// Adjust bass only
		newBassLevel := 3
		if toneControls.Bass.Value != newBassLevel {
			fmt.Printf("   Setting advanced bass to %d...\n", newBassLevel)
			err = soundtouchClient.SetAdvancedBass(newBassLevel)
			if err != nil {
				log.Printf("❌ Failed to set advanced bass: %v", err)
			} else {
				fmt.Printf("✅ Advanced bass adjusted\n")
			}
		}

		time.Sleep(1 * time.Second)

		// Adjust treble only
		newTrebleLevel := -1
		if toneControls.Treble.Value != newTrebleLevel {
			fmt.Printf("   Setting advanced treble to %d...\n", newTrebleLevel)
			err = soundtouchClient.SetAdvancedTreble(newTrebleLevel)
			if err != nil {
				log.Printf("❌ Failed to set advanced treble: %v", err)
			} else {
				fmt.Printf("✅ Advanced treble adjusted\n")
			}
		}

		time.Sleep(1 * time.Second)

		// Adjust both bass and treble together
		combinedBass := 2
		combinedTreble := 1
		fmt.Printf("   Setting bass to %d and treble to %d together...\n", combinedBass, combinedTreble)
		err = soundtouchClient.SetAudioProductToneControls(&combinedBass, &combinedTreble)
		if err != nil {
			log.Printf("❌ Failed to set tone controls: %v", err)
		} else {
			fmt.Printf("✅ Both tone controls adjusted\n")
		}
	}

	time.Sleep(2 * time.Second)

	// Example 4: Speaker Level Controls
	fmt.Println("\n4. Speaker Level Controls...")

	levelControls, err := soundtouchClient.GetAudioProductLevelControls()
	if err != nil {
		log.Printf("⚠️  Speaker level controls not available on this device: %v", err)
		fmt.Println("   This feature is only available on surround sound systems")
	} else {
		fmt.Printf("🔊 Current Speaker Levels: %s\n", levelControls.String())

		// Adjust front-center speaker level
		newFrontCenterLevel := 2
		if levelControls.FrontCenterSpeakerLevel.Value != newFrontCenterLevel {
			fmt.Printf("   Setting front-center speaker level to %d...\n", newFrontCenterLevel)
			err = soundtouchClient.SetFrontCenterSpeakerLevel(newFrontCenterLevel)
			if err != nil {
				log.Printf("❌ Failed to set front-center level: %v", err)
			} else {
				fmt.Printf("✅ Front-center speaker level adjusted\n")
			}
		}

		time.Sleep(1 * time.Second)

		// Adjust rear-surround speakers level
		newRearSurroundLevel := -1
		if levelControls.RearSurroundSpeakersLevel.Value != newRearSurroundLevel {
			fmt.Printf("   Setting rear-surround speakers level to %d...\n", newRearSurroundLevel)
			err = soundtouchClient.SetRearSurroundSpeakersLevel(newRearSurroundLevel)
			if err != nil {
				log.Printf("❌ Failed to set rear-surround level: %v", err)
			} else {
				fmt.Printf("✅ Rear-surround speakers level adjusted\n")
			}
		}

		time.Sleep(1 * time.Second)

		// Adjust both speaker levels together
		combinedFrontCenter := 1
		combinedRearSurround := 0
		fmt.Printf("   Setting front-center to %d and rear-surround to %d together...\n",
			combinedFrontCenter, combinedRearSurround)
		err = soundtouchClient.SetAudioProductLevelControls(&combinedFrontCenter, &combinedRearSurround)
		if err != nil {
			log.Printf("❌ Failed to set speaker levels: %v", err)
		} else {
			fmt.Printf("✅ Both speaker levels adjusted\n")
		}
	}

	// Example 5: Compare with basic controls
	fmt.Println("\n5. Comparison with Basic Audio Controls...")
	fmt.Println("   Basic controls available on all devices:")

	// Basic bass control (available on all devices)
	basicBass, err := soundtouchClient.GetBass()
	if err != nil {
		log.Printf("❌ Failed to get basic bass: %v", err)
	} else {
		fmt.Printf("   Basic Bass: %d (range: -9 to +9)\n", basicBass.TargetBass)
	}

	// Basic volume control
	volume, err := soundtouchClient.GetVolume()
	if err != nil {
		log.Printf("❌ Failed to get volume: %v", err)
	} else {
		fmt.Printf("   Volume: %d%%\n", volume.TargetVolume)
	}

	// Balance control (if available)
	balance, err := soundtouchClient.GetBalance()
	if err != nil {
		log.Printf("   Balance: Not available on this device")
	} else {
		fmt.Printf("   Balance: %d (range: -50 to +50)\n", balance.TargetBalance)
	}

	// Example 6: Error handling and validation
	fmt.Println("\n6. Error Handling Examples...")

	// Try to set invalid DSP controls to demonstrate validation
	fmt.Println("   Testing invalid audio mode...")
	err = soundtouchClient.SetAudioMode("INVALID_MODE")
	if err != nil {
		fmt.Printf("⚠️  Expected error for invalid mode: %v\n", err)
	}

	fmt.Println("   Testing negative video sync delay...")
	err = soundtouchClient.SetVideoSyncAudioDelay(-10)
	if err != nil {
		fmt.Printf("⚠️  Expected error for negative delay: %v\n", err)
	}

	// Example 7: CLI command equivalents
	fmt.Println("\n7. CLI Command Equivalents...")
	fmt.Println("   You can also use the CLI for these operations:")
	fmt.Println("   ")
	fmt.Println("   # DSP Controls")
	fmt.Printf("   soundtouch-cli audio dsp get --host %s\n", deviceIP)
	fmt.Printf("   soundtouch-cli audio dsp set --host %s --mode MUSIC --delay 50\n", deviceIP)
	fmt.Printf("   soundtouch-cli audio dsp mode --host %s --mode DIALOG\n", deviceIP)
	fmt.Println("   ")
	fmt.Println("   # Tone Controls")
	fmt.Printf("   soundtouch-cli audio tone get --host %s\n", deviceIP)
	fmt.Printf("   soundtouch-cli audio tone set --host %s --bass 3 --treble -1\n", deviceIP)
	fmt.Printf("   soundtouch-cli audio tone bass --host %s --level 5\n", deviceIP)
	fmt.Println("   ")
	fmt.Println("   # Level Controls")
	fmt.Printf("   soundtouch-cli audio level get --host %s\n", deviceIP)
	fmt.Printf("   soundtouch-cli audio level set --host %s --front-center 2 --rear-surround -1\n", deviceIP)
	fmt.Printf("   soundtouch-cli audio level front-center --host %s --level 3\n", deviceIP)

	fmt.Println("\n🎉 Advanced audio controls example completed!")
	fmt.Println("\nNotes:")
	fmt.Println("• Advanced audio controls are only available on professional/high-end devices")
	fmt.Println("• Consumer SoundTouch devices typically only support basic controls")
	fmt.Println("• Check device capabilities first to see which features are supported")
	fmt.Println("• Use GetCapabilities() to see 'audiodspcontrols', 'audioproducttonecontrols', etc.")
	fmt.Println("• All methods include comprehensive validation and error handling")
	fmt.Println("• Ranges and steps vary by device - check the response for valid values")
}

// Device Compatibility Notes:
//
// Consumer Devices (SoundTouch 10, 20, 30):
// - Basic bass control: ✅ Available
// - Basic volume control: ✅ Available
// - Basic balance control: ✅ Available (some models)
// - Advanced DSP controls: ❌ Not available
// - Advanced tone controls: ❌ Not available
// - Speaker level controls: ❌ Not available
//
// Professional/High-end Devices:
// - All basic controls: ✅ Available
// - DSP audio modes: ✅ Available
// - Video sync delay: ✅ Available
// - Advanced bass/treble: ✅ Available
// - Speaker level controls: ✅ Available (surround systems)
//
// API Endpoints Implemented:
// - GET/POST /audiodspcontrols - DSP settings and audio modes
// - GET/POST /audioproducttonecontrols - Advanced bass/treble
// - GET/POST /audioproductlevelcontrols - Speaker level controls
//
// These complement the existing basic audio controls:
// - GET/POST /bass - Basic bass control (-9 to +9)
// - GET/POST /volume - Volume and mute control
// - GET/POST /balance - Stereo balance control (-50 to +50)
