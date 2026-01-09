package models

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestBalanceMarshalXML(t *testing.T) {
	balance := Balance{
		DeviceID:      "1234567890AB",
		TargetBalance: -25,
		ActualBalance: -25,
	}

	var buf strings.Builder
	encoder := xml.NewEncoder(&buf)
	err := balance.MarshalXML(encoder, xml.StartElement{Name: xml.Name{Local: "balance"}})
	if err != nil {
		t.Fatalf("MarshalXML failed: %v", err)
	}
	encoder.Flush()

	// Convert to string for easier testing
	xmlStr := buf.String()

	// Check that XML contains expected elements
	expectedElements := []string{
		`deviceID="1234567890AB"`,
		`<targetbalance>-25</targetbalance>`,
		`<actualbalance>-25</actualbalance>`,
	}

	for _, expected := range expectedElements {
		if !strings.Contains(xmlStr, expected) {
			t.Errorf("MarshalXML result %q does not contain expected element %q", xmlStr, expected)
		}
	}
}

func TestBalanceMarshalXML_PositiveValue(t *testing.T) {
	balance := Balance{
		DeviceID:      "ABCDEF123456",
		TargetBalance: 30,
		ActualBalance: 30,
	}

	var buf strings.Builder
	encoder := xml.NewEncoder(&buf)
	err := balance.MarshalXML(encoder, xml.StartElement{Name: xml.Name{Local: "balance"}})
	if err != nil {
		t.Fatalf("MarshalXML failed: %v", err)
	}
	encoder.Flush()

	xmlStr := buf.String()

	expectedElements := []string{
		`deviceID="ABCDEF123456"`,
		`<targetbalance>30</targetbalance>`,
		`<actualbalance>30</actualbalance>`,
	}

	for _, expected := range expectedElements {
		if !strings.Contains(xmlStr, expected) {
			t.Errorf("MarshalXML result %q does not contain expected element %q", xmlStr, expected)
		}
	}
}

func TestBalanceMarshalXML_ZeroValue(t *testing.T) {
	balance := Balance{
		DeviceID:      "ZERO0000TEST",
		TargetBalance: 0,
		ActualBalance: 0,
	}

	var buf strings.Builder
	encoder := xml.NewEncoder(&buf)
	err := balance.MarshalXML(encoder, xml.StartElement{Name: xml.Name{Local: "balance"}})
	if err != nil {
		t.Fatalf("MarshalXML failed: %v", err)
	}
	encoder.Flush()

	xmlStr := buf.String()

	expectedElements := []string{
		`deviceID="ZERO0000TEST"`,
		`<targetbalance>0</targetbalance>`,
		`<actualbalance>0</actualbalance>`,
	}

	for _, expected := range expectedElements {
		if !strings.Contains(xmlStr, expected) {
			t.Errorf("MarshalXML result %q does not contain expected element %q", xmlStr, expected)
		}
	}
}

func TestBalanceMarshalXML_ExtremeValues(t *testing.T) {
	balance := Balance{
		DeviceID:      "EXTREME_TEST",
		TargetBalance: -50, // Min value
		ActualBalance: 50,  // Max value
	}

	var buf strings.Builder
	encoder := xml.NewEncoder(&buf)
	err := balance.MarshalXML(encoder, xml.StartElement{Name: xml.Name{Local: "balance"}})
	if err != nil {
		t.Fatalf("MarshalXML failed: %v", err)
	}
	encoder.Flush()

	xmlStr := buf.String()

	expectedElements := []string{
		`deviceID="EXTREME_TEST"`,
		`<targetbalance>-50</targetbalance>`,
		`<actualbalance>50</actualbalance>`,
	}

	for _, expected := range expectedElements {
		if !strings.Contains(xmlStr, expected) {
			t.Errorf("MarshalXML result %q does not contain expected element %q", xmlStr, expected)
		}
	}
}
