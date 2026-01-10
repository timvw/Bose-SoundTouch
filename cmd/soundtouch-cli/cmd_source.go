package main

import (
	"fmt"
	"strings"

	"github.com/gesellix/bose-soundtouch/pkg/models"
	"github.com/urfave/cli/v2"
)

func printSource(source models.SourceItem) {
	fmt.Printf("    • %s", source.GetDisplayName())

	if source.SourceAccount != "" && source.SourceAccount != source.Source {
		fmt.Printf(" (%s)", source.SourceAccount)
	}

	var attributes []string
	if source.IsLocalSource() {
		attributes = append(attributes, "Local")
		attributes = append(attributes, "Available")
	}

	if len(attributes) > 0 {
		fmt.Printf(" [%s]", strings.Join(attributes, ", "))
	}

	fmt.Println()
}

// listSources handles listing available audio sources
func listSources(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Getting available sources", clientConfig.Host, clientConfig.Port)

	sources, err := client.GetSources()
	if err != nil {
		return fmt.Errorf("failed to get sources: %w", err)
	}

	fmt.Printf("Available Audio Sources:\n")
	fmt.Printf("  Device ID: %s\n", sources.DeviceID)

	// Show ready sources first
	availableSources := sources.GetAvailableSources()
	if len(availableSources) > 0 {
		fmt.Printf("  Ready Sources:\n")

		for _, source := range availableSources {
			printSource(source)
		}
	}

	// Show all configured sources
	fmt.Printf("  All Sources:\n")

	for _, source := range sources.SourceItem {
		status := "Available"
		if !source.IsLocalSource() {
			status = "Remote"
		}

		fmt.Printf("    • %s (%s)\n", source.GetDisplayName(), status)

		if source.SourceAccount != "" && source.SourceAccount != source.Source {
			fmt.Printf("      Account: %s\n", source.SourceAccount)
		}
	}

	// Show streaming sources
	streamingSources := sources.GetStreamingSources()
	if len(streamingSources) > 0 {
		fmt.Printf("  Streaming Services:\n")

		for _, source := range streamingSources {
			fmt.Printf("    • %s", source.GetDisplayName())

			if source.SourceAccount != "" {
				fmt.Printf(" (%s)", source.SourceAccount)
			}

			fmt.Println()
		}
	}

	return nil
}

// selectSource handles selecting an audio source
func selectSource(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	sourceName := strings.ToUpper(c.String("source"))
	sourceAccount := c.String("account")

	PrintDeviceHeader(fmt.Sprintf("Selecting source '%s'", sourceName), clientConfig.Host, clientConfig.Port)

	err = client.SelectSource(sourceName, sourceAccount)
	if err != nil {
		return fmt.Errorf("failed to select source: %w", err)
	}

	if sourceAccount != "" {
		PrintSuccess(fmt.Sprintf("Source '%s' with account '%s' selected", sourceName, sourceAccount))
	} else {
		PrintSuccess(fmt.Sprintf("Source '%s' selected", sourceName))
	}

	return nil
}

// selectSpotify handles selecting Spotify source
func selectSpotify(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Selecting Spotify source", clientConfig.Host, clientConfig.Port)

	err = client.SelectSpotify("")
	if err != nil {
		return fmt.Errorf("failed to select Spotify: %w", err)
	}

	PrintSuccess("Spotify source selected")

	return nil
}

// selectBluetooth handles selecting Bluetooth source
func selectBluetooth(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Selecting Bluetooth source", clientConfig.Host, clientConfig.Port)

	err = client.SelectBluetooth()
	if err != nil {
		return fmt.Errorf("failed to select Bluetooth: %w", err)
	}

	PrintSuccess("Bluetooth source selected")

	return nil
}

// selectAux handles selecting AUX input source
func selectAux(c *cli.Context) error {
	clientConfig := GetClientConfig(c)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		return err
	}

	PrintDeviceHeader("Selecting AUX input source", clientConfig.Host, clientConfig.Port)

	err = client.SelectAux()
	if err != nil {
		return fmt.Errorf("failed to select AUX: %w", err)
	}

	PrintSuccess("AUX input source selected")

	return nil
}
