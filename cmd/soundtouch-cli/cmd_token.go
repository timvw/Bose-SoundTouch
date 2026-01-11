package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// requestToken requests a new bearer token from the device
func requestToken(c *cli.Context) error {
	clientConfig := GetClientConfig(c)
	PrintDeviceHeader("Requesting bearer token", clientConfig.Host, clientConfig.Port)

	client, err := CreateSoundTouchClient(clientConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to create client: %v", err))
		return err
	}

	token, err := client.RequestToken()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to request token: %v", err))
		return err
	}

	fmt.Println("Bearer Token Information:")

	if token.IsValid() {
		fmt.Printf("  Status: Valid\n")
		fmt.Printf("  Token: %s\n", token.String())
		fmt.Printf("  Full value: %s\n", token.GetToken())
		fmt.Printf("  Authorization header: %s\n", token.GetAuthHeader())

		// Display token without Bearer prefix for API usage
		fmt.Println("\nFor API Usage:")
		fmt.Printf("  Raw token: %s\n", token.GetTokenWithoutPrefix())

		// Usage instructions
		fmt.Println("\nUsage Instructions:")
		fmt.Println("  • Use the 'Authorization header' value in HTTP Authorization headers")
		fmt.Println("  • Use the 'Raw token' value when an API requires token without 'Bearer ' prefix")
		fmt.Println("  • Tokens are generated per request and may have expiration times")

		// Security notice
		fmt.Println("\nSecurity Notice:")
		fmt.Println("  • Store tokens securely and avoid logging them in plain text")
		fmt.Println("  • Tokens provide authentication - treat them as passwords")
		fmt.Println("  • Request new tokens when needed rather than reusing old ones")
	} else {
		fmt.Printf("  Status: Invalid\n")
		fmt.Printf("  Raw response: %s\n", token.GetToken())
		PrintError("Received invalid bearer token from device")
	}

	return nil
}
