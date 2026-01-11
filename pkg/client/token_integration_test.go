package client

import (
	"os"
	"strings"
	"testing"
)

// TestRequestToken_Integration tests the RequestToken method against a real device
// This test only runs when SOUNDTOUCH_TEST_HOST environment variable is set
func TestRequestToken_Integration(t *testing.T) {
	host := os.Getenv("SOUNDTOUCH_TEST_HOST")
	if host == "" {
		t.Skip("Skipping integration test: SOUNDTOUCH_TEST_HOST not set")
	}

	// Create client for real device
	config := &Config{
		Host: host,
		Port: 8090,
	}
	client := NewClient(config)

	// Test RequestToken with real device
	token, err := client.RequestToken()
	if err != nil {
		t.Fatalf("RequestToken() failed with real device: %v", err)
	}

	if token == nil {
		t.Fatal("RequestToken() returned nil token from real device")
	}

	// Validate token properties without exposing actual values
	tokenValue := token.GetToken()

	// Token should have Bearer prefix
	if !strings.HasPrefix(tokenValue, "Bearer ") {
		t.Error("Real device token should have 'Bearer ' prefix")
	}

	// Token should be valid according to our validation
	if !token.IsValid() {
		t.Error("Real device token should be valid")
	}

	// Raw token should not include Bearer prefix
	rawToken := token.GetTokenWithoutPrefix()
	if strings.HasPrefix(rawToken, "Bearer ") {
		t.Error("Raw token should not include Bearer prefix")
	}

	// Token should be reasonably long (real bearer tokens are substantial)
	if len(rawToken) < 80 {
		t.Errorf("Real device token seems too short: %d characters", len(rawToken))
	}

	// Token should only contain base64-like characters plus common token chars
	validChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
	for _, char := range rawToken {
		if !strings.ContainsRune(validChars, char) {
			t.Errorf("Token contains unexpected character: %c", char)
		}
	}

	// Auth header should match full token value
	if token.GetAuthHeader() != tokenValue {
		t.Error("Auth header should match full token value")
	}

	// String representation should be truncated for security
	stringRepr := token.String()
	if len(stringRepr) >= len(tokenValue) {
		t.Error("String representation should be shorter than full token for security")
	}

	// String representation should contain "..." for long tokens
	if !strings.Contains(stringRepr, "...") {
		t.Error("String representation should contain '...' for long tokens")
	}

	// Multiple calls should generate different tokens (if the device supports it)
	token2, err := client.RequestToken()
	if err != nil {
		t.Fatalf("Second RequestToken() call failed: %v", err)
	}

	// Note: Some devices may return the same token, so we don't enforce uniqueness
	// but we do verify the second token is also valid
	if !token2.IsValid() {
		t.Error("Second token should also be valid")
	}

	t.Logf("Successfully validated real device token properties (length: %d chars)", len(rawToken))
}
