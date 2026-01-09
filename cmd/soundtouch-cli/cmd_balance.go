// Package main provides the soundtouch-cli balance control commands.
package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// getBalance retrieves the current balance level from the device
func getBalance(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Getting balance level", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	balance, err := client.GetBalance()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get balance: %v", err))
		return err
	}

	fmt.Printf("Current balance level: %d\n", balance.ActualBalance)

	if balance.TargetBalance != balance.ActualBalance {
		fmt.Printf("Target balance level: %d\n", balance.TargetBalance)
	}

	// Display balance direction
	switch {
	case balance.ActualBalance > 0:
		fmt.Printf("Balance direction: Right (+%d)\n", balance.ActualBalance)
	case balance.ActualBalance < 0:
		fmt.Printf("Balance direction: Left (%d)\n", balance.ActualBalance)
	default:
		fmt.Println("Balance direction: Center (0)")
	}

	return nil
}

// setBalance sets the balance level on the device
func setBalance(c *cli.Context) error {
	level := c.Int("level")
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader(fmt.Sprintf("Setting balance level to %d", level), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.SetBalanceSafe(level)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to set balance: %v", err))
		return err
	}

	PrintSuccess(fmt.Sprintf("Balance level set to %d", level))
	return nil
}

// balanceLeft shifts balance to the left
func balanceLeft(c *cli.Context) error {
	amount := c.Int("amount")
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader(fmt.Sprintf("Shifting balance left by %d", amount), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	// Get current balance level first
	currentBalance, err := client.GetBalance()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get current balance: %v", err))
		return err
	}

	newLevel := currentBalance.ActualBalance - amount
	err = client.SetBalanceSafe(newLevel)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to shift balance left: %v", err))
		return err
	}

	PrintSuccess(fmt.Sprintf("Balance shifted from %d to %d (left)", currentBalance.ActualBalance, newLevel))
	return nil
}

// balanceRight shifts balance to the right
func balanceRight(c *cli.Context) error {
	amount := c.Int("amount")
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader(fmt.Sprintf("Shifting balance right by %d", amount), clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	// Get current balance level first
	currentBalance, err := client.GetBalance()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to get current balance: %v", err))
		return err
	}

	newLevel := currentBalance.ActualBalance + amount
	err = client.SetBalanceSafe(newLevel)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to shift balance right: %v", err))
		return err
	}

	PrintSuccess(fmt.Sprintf("Balance shifted from %d to %d (right)", currentBalance.ActualBalance, newLevel))
	return nil
}

// balanceCenter centers the balance (sets to 0)
func balanceCenter(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Centering balance", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	err = client.SetBalanceSafe(0)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to center balance: %v", err))
		return err
	}

	PrintSuccess("Balance centered")
	return nil
}
