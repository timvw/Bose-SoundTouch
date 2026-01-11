package models

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// BearerToken represents a bearer token response from the device
type BearerToken struct {
	XMLName xml.Name `xml:"bearertoken"`
	Value   string   `xml:"value,attr"`
}

// GetToken returns the bearer token value
func (b *BearerToken) GetToken() string {
	return b.Value
}

// GetTokenWithoutPrefix returns the bearer token value without the "Bearer " prefix
func (b *BearerToken) GetTokenWithoutPrefix() string {
	token := b.Value
	if strings.HasPrefix(token, "Bearer ") {
		return token[7:]
	}

	return token
}

// IsValid returns true if the bearer token has a valid value
func (b *BearerToken) IsValid() bool {
	return b.Value != "" && strings.HasPrefix(b.Value, "Bearer ")
}

// String returns a string representation of the bearer token
func (b *BearerToken) String() string {
	if !b.IsValid() {
		return "Invalid bearer token"
	}

	// Show only first and last 10 characters for security
	token := b.GetTokenWithoutPrefix()
	if len(token) > 20 {
		return fmt.Sprintf("Bearer %s...%s", token[:10], token[len(token)-10:])
	}

	return b.Value
}

// GetAuthHeader returns the token formatted for use in Authorization headers
func (b *BearerToken) GetAuthHeader() string {
	if !b.IsValid() {
		return ""
	}

	return b.Value
}

// NewBearerToken creates a new BearerToken with the given token value
func NewBearerToken(token string) *BearerToken {
	// Ensure token has "Bearer " prefix
	if !strings.HasPrefix(token, "Bearer ") {
		token = "Bearer " + token
	}

	return &BearerToken{
		Value: token,
	}
}
