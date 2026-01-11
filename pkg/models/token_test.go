package models

import (
	"encoding/xml"
	"testing"
)

func TestBearerToken_GetToken(t *testing.T) {
	tests := []struct {
		name     string
		token    BearerToken
		expected string
	}{
		{
			name:     "Valid token",
			token:    BearerToken{Value: "Bearer abc123def456"},
			expected: "Bearer abc123def456",
		},
		{
			name:     "Empty token",
			token:    BearerToken{Value: ""},
			expected: "",
		},
		{
			name:     "Token without Bearer prefix",
			token:    BearerToken{Value: "abc123def456"},
			expected: "abc123def456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.token.GetToken()
			if result != tt.expected {
				t.Errorf("GetToken() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBearerToken_GetTokenWithoutPrefix(t *testing.T) {
	tests := []struct {
		name     string
		token    BearerToken
		expected string
	}{
		{
			name:     "Valid token with Bearer prefix",
			token:    BearerToken{Value: "Bearer abc123def456"},
			expected: "abc123def456",
		},
		{
			name:     "Token without Bearer prefix",
			token:    BearerToken{Value: "abc123def456"},
			expected: "abc123def456",
		},
		{
			name:     "Token with Bearer and extra spaces",
			token:    BearerToken{Value: "Bearer   abc123def456   "},
			expected: "  abc123def456   ",
		},
		{
			name:     "Empty token",
			token:    BearerToken{Value: ""},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.token.GetTokenWithoutPrefix()
			if result != tt.expected {
				t.Errorf("GetTokenWithoutPrefix() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBearerToken_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		token    BearerToken
		expected bool
	}{
		{
			name:     "Valid token",
			token:    BearerToken{Value: "Bearer abc123def456"},
			expected: true,
		},
		{
			name:     "Empty token",
			token:    BearerToken{Value: ""},
			expected: false,
		},
		{
			name:     "Token without Bearer prefix",
			token:    BearerToken{Value: "abc123def456"},
			expected: false,
		},
		{
			name:     "Just Bearer",
			token:    BearerToken{Value: "Bearer"},
			expected: false,
		},
		{
			name:     "Bearer with space but no token",
			token:    BearerToken{Value: "Bearer "},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.token.IsValid()
			if result != tt.expected {
				t.Errorf("IsValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBearerToken_String(t *testing.T) {
	tests := []struct {
		name     string
		token    BearerToken
		expected string
	}{
		{
			name:     "Valid long token",
			token:    BearerToken{Value: "Bearer abcdefghij1234567890klmnopqrstuvwxyz"},
			expected: "Bearer abcdefghij...qrstuvwxyz",
		},
		{
			name:     "Valid short token",
			token:    BearerToken{Value: "Bearer abc123"},
			expected: "Bearer abc123",
		},
		{
			name:     "Invalid token",
			token:    BearerToken{Value: "abc123"},
			expected: "Invalid bearer token",
		},
		{
			name:     "Empty token",
			token:    BearerToken{Value: ""},
			expected: "Invalid bearer token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.token.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBearerToken_GetAuthHeader(t *testing.T) {
	tests := []struct {
		name     string
		token    BearerToken
		expected string
	}{
		{
			name:     "Valid token",
			token:    BearerToken{Value: "Bearer abc123def456"},
			expected: "Bearer abc123def456",
		},
		{
			name:     "Invalid token",
			token:    BearerToken{Value: "abc123"},
			expected: "",
		},
		{
			name:     "Empty token",
			token:    BearerToken{Value: ""},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.token.GetAuthHeader()
			if result != tt.expected {
				t.Errorf("GetAuthHeader() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewBearerToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected BearerToken
	}{
		{
			name:     "Token with Bearer prefix",
			input:    "Bearer abc123def456",
			expected: BearerToken{Value: "Bearer abc123def456"},
		},
		{
			name:     "Token without Bearer prefix",
			input:    "abc123def456",
			expected: BearerToken{Value: "Bearer abc123def456"},
		},
		{
			name:     "Empty token",
			input:    "",
			expected: BearerToken{Value: "Bearer "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewBearerToken(tt.input)
			if result.Value != tt.expected.Value {
				t.Errorf("NewBearerToken() = %v, want %v", result.Value, tt.expected.Value)
			}
		})
	}
}

func TestBearerToken_XMLMarshaling(t *testing.T) {
	// Test unmarshaling from XML (device response)
	xmlData := `<?xml version="1.0" encoding="UTF-8" ?><bearertoken value="Bearer abc123def456xyz789" />`

	var token BearerToken

	err := xml.Unmarshal([]byte(xmlData), &token)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	if token.Value != "Bearer abc123def456xyz789" {
		t.Errorf("Unmarshaled token value = %v, want %v", token.Value, "Bearer abc123def456xyz789")
	}

	if !token.IsValid() {
		t.Error("Unmarshaled token should be valid")
	}

	// Test marshaling to XML
	token2 := BearerToken{Value: "Bearer test123token456"}

	data, err := xml.Marshal(token2)
	if err != nil {
		t.Fatalf("Failed to marshal to XML: %v", err)
	}

	expected := `<bearertoken value="Bearer test123token456"></bearertoken>`
	if string(data) != expected {
		t.Errorf("Marshaled XML = %v, want %v", string(data), expected)
	}
}

func TestBearerToken_RealDeviceExample(t *testing.T) {
	// Test with example device response format (generic token)
	xmlData := `<?xml version="1.0" encoding="UTF-8" ?><bearertoken value="Bearer vUApzBVT6Lh0nw1xVu/plr1UDRNdMYMEpe0cStm4wCH5mWSjrrtORnGGirMn3pspkJ8mNR1MFh/J4OcsbEikMplcDGJVeuZOnDPAskQALvDBCF0PW74qXRms2k1AfLJ/" />`

	var token BearerToken

	err := xml.Unmarshal([]byte(xmlData), &token)
	if err != nil {
		t.Fatalf("Failed to unmarshal device XML: %v", err)
	}

	// Verify it's valid
	if !token.IsValid() {
		t.Error("Device token should be valid")
	}

	// Verify auth header
	authHeader := token.GetAuthHeader()
	if authHeader != token.Value {
		t.Errorf("Auth header = %v, want %v", authHeader, token.Value)
	}

	// Verify string representation shows truncated version for long tokens
	stringRepr := token.String()
	if stringRepr == token.Value {
		t.Error("String representation should be truncated for long tokens")
	}

	// Should contain "Bearer" and "..." for long tokens
	if len(stringRepr) >= len(token.Value) {
		t.Error("String representation should be shorter than full token")
	}

	// Verify raw token extraction
	rawToken := token.GetTokenWithoutPrefix()

	expectedRawToken := "vUApzBVT6Lh0nw1xVu/plr1UDRNdMYMEpe0cStm4wCH5mWSjrrtORnGGirMn3pspkJ8mNR1MFh/J4OcsbEikMplcDGJVeuZOnDPAskQALvDBCF0PW74qXRms2k1AfLJ/"
	if rawToken != expectedRawToken {
		t.Errorf("Raw token = %v, want %v", rawToken, expectedRawToken)
	}

	// Verify token length is reasonable (typical bearer tokens are long)
	if len(rawToken) < 50 {
		t.Errorf("Token seems too short: %d characters", len(rawToken))
	}
}
