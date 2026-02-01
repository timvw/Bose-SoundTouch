package models

import (
	"encoding/xml"
	"testing"
)

func TestServiceAvailability_UnmarshalXML(t *testing.T) {
	tests := []struct {
		name     string
		xmlData  string
		validate func(t *testing.T, sa *ServiceAvailability)
	}{
		{
			name: "complete service availability response",
			xmlData: `<?xml version="1.0" encoding="UTF-8" ?>
<serviceAvailability>
  <services>
    <service type="AIRPLAY" isAvailable="true" />
    <service type="ALEXA" isAvailable="false" />
    <service type="AMAZON" isAvailable="true" />
    <service type="BLUETOOTH" isAvailable="false" reason="INVALID_SOURCE_TYPE" />
    <service type="BMX" isAvailable="false" />
    <service type="DEEZER" isAvailable="true" />
    <service type="IHEART" isAvailable="true" />
    <service type="LOCAL_INTERNET_RADIO" isAvailable="true" />
    <service type="LOCAL_MUSIC" isAvailable="true" />
    <service type="NOTIFICATION" isAvailable="false" />
    <service type="PANDORA" isAvailable="true" />
    <service type="SPOTIFY" isAvailable="true" />
    <service type="TUNEIN" isAvailable="true" />
  </services>
</serviceAvailability>`,
			validate: func(t *testing.T, sa *ServiceAvailability) {
				t.Helper()
				if sa.Services == nil {
					t.Fatal("services should not be nil")
				}
				if len(sa.Services.Service) != 13 {
					t.Errorf("expected 13 services, got %d", len(sa.Services.Service))
				}

				// Check specific services
				spotifyService := sa.GetServiceByType(ServiceTypeSpotify)
				if spotifyService == nil {
					t.Fatal("spotify service should not be nil")
				}
				if !spotifyService.IsAvailable {
					t.Error("spotify service should be available")
				}
				if spotifyService.Reason != "" {
					t.Error("spotify service should not have a reason")
				}

				bluetoothService := sa.GetServiceByType(ServiceTypeBluetooth)
				if bluetoothService == nil {
					t.Fatal("bluetooth service should not be nil")
				}
				if bluetoothService.IsAvailable {
					t.Error("bluetooth service should not be available")
				}
				if bluetoothService.Reason != "INVALID_SOURCE_TYPE" {
					t.Errorf("bluetooth service reason should be 'INVALID_SOURCE_TYPE', got '%s'", bluetoothService.Reason)
				}
			},
		},
		{
			name: "empty services",
			xmlData: `<?xml version="1.0" encoding="UTF-8" ?>
<serviceAvailability>
  <services>
  </services>
</serviceAvailability>`,
			validate: func(t *testing.T, sa *ServiceAvailability) {
				t.Helper()
				if sa.Services == nil {
					t.Fatal("services should not be nil")
				}
				if len(sa.Services.Service) != 0 {
					t.Errorf("expected 0 services, got %d", len(sa.Services.Service))
				}
			},
		},
		{
			name: "minimal response",
			xmlData: `<serviceAvailability>
  <services>
    <service type="SPOTIFY" isAvailable="true" />
  </services>
</serviceAvailability>`,
			validate: func(t *testing.T, sa *ServiceAvailability) {
				t.Helper()
				if sa.Services == nil {
					t.Fatal("services should not be nil")
				}
				if len(sa.Services.Service) != 1 {
					t.Errorf("expected 1 service, got %d", len(sa.Services.Service))
				}
				if sa.Services.Service[0].Type != "SPOTIFY" {
					t.Errorf("expected SPOTIFY, got %s", sa.Services.Service[0].Type)
				}
				if !sa.Services.Service[0].IsAvailable {
					t.Error("service should be available")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sa ServiceAvailability
			err := xml.Unmarshal([]byte(tt.xmlData), &sa)
			if err != nil {
				t.Fatalf("failed to unmarshal XML: %v", err)
			}

			tt.validate(t, &sa)
		})
	}
}

func TestServiceAvailability_GetAvailableServices(t *testing.T) {
	sa := &ServiceAvailability{
		Services: &ServiceList{
			Service: []Service{
				{Type: "SPOTIFY", IsAvailable: true},
				{Type: "BLUETOOTH", IsAvailable: false, Reason: "UNAVAILABLE"},
				{Type: "AIRPLAY", IsAvailable: true},
				{Type: "ALEXA", IsAvailable: false},
			},
		},
	}

	available := sa.GetAvailableServices()
	if len(available) != 2 {
		t.Errorf("expected 2 available services, got %d", len(available))
	}
	if available[0].Type != "SPOTIFY" {
		t.Errorf("expected first service to be SPOTIFY, got %s", available[0].Type)
	}
	if available[1].Type != "AIRPLAY" {
		t.Errorf("expected second service to be AIRPLAY, got %s", available[1].Type)
	}
}

func TestServiceAvailability_GetUnavailableServices(t *testing.T) {
	sa := &ServiceAvailability{
		Services: &ServiceList{
			Service: []Service{
				{Type: "SPOTIFY", IsAvailable: true},
				{Type: "BLUETOOTH", IsAvailable: false, Reason: "UNAVAILABLE"},
				{Type: "AIRPLAY", IsAvailable: true},
				{Type: "ALEXA", IsAvailable: false},
			},
		},
	}

	unavailable := sa.GetUnavailableServices()
	if len(unavailable) != 2 {
		t.Errorf("expected 2 unavailable services, got %d", len(unavailable))
	}
	if unavailable[0].Type != "BLUETOOTH" {
		t.Errorf("expected first service to be BLUETOOTH, got %s", unavailable[0].Type)
	}
	if unavailable[1].Type != "ALEXA" {
		t.Errorf("expected second service to be ALEXA, got %s", unavailable[1].Type)
	}
}

func TestServiceAvailability_IsServiceAvailable(t *testing.T) {
	sa := &ServiceAvailability{
		Services: &ServiceList{
			Service: []Service{
				{Type: "SPOTIFY", IsAvailable: true},
				{Type: "BLUETOOTH", IsAvailable: false},
			},
		},
	}

	if !sa.IsServiceAvailable(ServiceTypeSpotify) {
		t.Error("Spotify should be available")
	}
	if sa.IsServiceAvailable(ServiceTypeBluetooth) {
		t.Error("Bluetooth should not be available")
	}
	if sa.IsServiceAvailable(ServiceTypeAlexa) {
		t.Error("Alexa should not be available (not in list)")
	}
}

func TestServiceAvailability_GetServiceByType(t *testing.T) {
	sa := &ServiceAvailability{
		Services: &ServiceList{
			Service: []Service{
				{Type: "SPOTIFY", IsAvailable: true},
				{Type: "BLUETOOTH", IsAvailable: false, Reason: "DEVICE_NOT_FOUND"},
			},
		},
	}

	spotifyService := sa.GetServiceByType(ServiceTypeSpotify)
	if spotifyService == nil {
		t.Fatal("spotify service should not be nil")
	}
	if spotifyService.Type != "SPOTIFY" {
		t.Errorf("expected SPOTIFY, got %s", spotifyService.Type)
	}
	if !spotifyService.IsAvailable {
		t.Error("spotify service should be available")
	}

	bluetoothService := sa.GetServiceByType(ServiceTypeBluetooth)
	if bluetoothService == nil {
		t.Fatal("bluetooth service should not be nil")
	}
	if bluetoothService.Type != "BLUETOOTH" {
		t.Errorf("expected BLUETOOTH, got %s", bluetoothService.Type)
	}
	if bluetoothService.IsAvailable {
		t.Error("bluetooth service should not be available")
	}
	if bluetoothService.Reason != "DEVICE_NOT_FOUND" {
		t.Errorf("expected DEVICE_NOT_FOUND, got %s", bluetoothService.Reason)
	}

	nonExistentService := sa.GetServiceByType(ServiceTypeAlexa)
	if nonExistentService != nil {
		t.Error("non-existent service should be nil")
	}
}

func TestServiceAvailability_ConvenienceMethods(t *testing.T) {
	sa := &ServiceAvailability{
		Services: &ServiceList{
			Service: []Service{
				{Type: "SPOTIFY", IsAvailable: true},
				{Type: "BLUETOOTH", IsAvailable: false},
				{Type: "AIRPLAY", IsAvailable: true},
				{Type: "ALEXA", IsAvailable: false},
				{Type: "TUNEIN", IsAvailable: true},
				{Type: "PANDORA", IsAvailable: true},
				{Type: "LOCAL_MUSIC", IsAvailable: true},
			},
		},
	}

	if !sa.HasSpotify() {
		t.Error("should have Spotify")
	}
	if sa.HasBluetooth() {
		t.Error("should not have Bluetooth")
	}
	if !sa.HasAirPlay() {
		t.Error("should have AirPlay")
	}
	if sa.HasAlexa() {
		t.Error("should not have Alexa")
	}
	if !sa.HasTuneIn() {
		t.Error("should have TuneIn")
	}
	if !sa.HasPandora() {
		t.Error("should have Pandora")
	}
	if !sa.HasLocalMusic() {
		t.Error("should have Local Music")
	}
}

func TestServiceAvailability_GetStreamingServices(t *testing.T) {
	sa := &ServiceAvailability{
		Services: &ServiceList{
			Service: []Service{
				{Type: "SPOTIFY", IsAvailable: true},
				{Type: "BLUETOOTH", IsAvailable: false},
				{Type: "PANDORA", IsAvailable: true},
				{Type: "TUNEIN", IsAvailable: false},
				{Type: "AMAZON", IsAvailable: true},
				{Type: "DEEZER", IsAvailable: false},
				{Type: "IHEART", IsAvailable: true},
				{Type: "LOCAL_INTERNET_RADIO", IsAvailable: true},
				{Type: "ALEXA", IsAvailable: false}, // Not a streaming service
			},
		},
	}

	streaming := sa.GetStreamingServices()
	if len(streaming) != 7 {
		t.Errorf("expected 7 streaming services, got %d", len(streaming))
	}

	streamingTypes := make(map[string]bool)
	for _, service := range streaming {
		streamingTypes[service.Type] = true
	}

	expectedStreaming := []string{"SPOTIFY", "PANDORA", "TUNEIN", "AMAZON", "DEEZER", "IHEART", "LOCAL_INTERNET_RADIO"}
	for _, expected := range expectedStreaming {
		if !streamingTypes[expected] {
			t.Errorf("expected to find streaming service %s", expected)
		}
	}

	notExpected := []string{"BLUETOOTH", "ALEXA"}
	for _, notExp := range notExpected {
		if streamingTypes[notExp] {
			t.Errorf("did not expect to find %s in streaming services", notExp)
		}
	}
}

func TestServiceAvailability_GetLocalServices(t *testing.T) {
	sa := &ServiceAvailability{
		Services: &ServiceList{
			Service: []Service{
				{Type: "SPOTIFY", IsAvailable: true},
				{Type: "BLUETOOTH", IsAvailable: false},
				{Type: "AIRPLAY", IsAvailable: true},
				{Type: "LOCAL_MUSIC", IsAvailable: true},
				{Type: "ALEXA", IsAvailable: false},
			},
		},
	}

	local := sa.GetLocalServices()
	if len(local) != 3 {
		t.Errorf("expected 3 local services, got %d", len(local))
	}

	localTypes := make(map[string]bool)
	for _, service := range local {
		localTypes[service.Type] = true
	}

	expectedLocal := []string{"BLUETOOTH", "AIRPLAY", "LOCAL_MUSIC"}
	for _, expected := range expectedLocal {
		if !localTypes[expected] {
			t.Errorf("expected to find local service %s", expected)
		}
	}

	notExpected := []string{"SPOTIFY", "ALEXA"}
	for _, notExp := range notExpected {
		if localTypes[notExp] {
			t.Errorf("did not expect to find %s in local services", notExp)
		}
	}
}

func TestServiceAvailability_CountMethods(t *testing.T) {
	sa := &ServiceAvailability{
		Services: &ServiceList{
			Service: []Service{
				{Type: "SPOTIFY", IsAvailable: true},
				{Type: "BLUETOOTH", IsAvailable: false},
				{Type: "AIRPLAY", IsAvailable: true},
				{Type: "ALEXA", IsAvailable: false},
			},
		},
	}

	if sa.GetServiceCount() != 4 {
		t.Errorf("expected 4 total services, got %d", sa.GetServiceCount())
	}
	if sa.GetAvailableServiceCount() != 2 {
		t.Errorf("expected 2 available services, got %d", sa.GetAvailableServiceCount())
	}
	if sa.GetUnavailableServiceCount() != 2 {
		t.Errorf("expected 2 unavailable services, got %d", sa.GetUnavailableServiceCount())
	}
}

func TestServiceAvailability_NilServicesHandling(t *testing.T) {
	sa := &ServiceAvailability{}

	if len(sa.GetAvailableServices()) != 0 {
		t.Error("available services should be empty")
	}
	if len(sa.GetUnavailableServices()) != 0 {
		t.Error("unavailable services should be empty")
	}
	if sa.IsServiceAvailable(ServiceTypeSpotify) {
		t.Error("Spotify should not be available")
	}
	if sa.GetServiceByType(ServiceTypeSpotify) != nil {
		t.Error("service should be nil")
	}
	if sa.HasSpotify() {
		t.Error("should not have Spotify")
	}
	if sa.HasBluetooth() {
		t.Error("should not have Bluetooth")
	}
	if len(sa.GetStreamingServices()) != 0 {
		t.Error("streaming services should be empty")
	}
	if len(sa.GetLocalServices()) != 0 {
		t.Error("local services should be empty")
	}
	if sa.GetServiceCount() != 0 {
		t.Error("service count should be 0")
	}
	if sa.GetAvailableServiceCount() != 0 {
		t.Error("available service count should be 0")
	}
	if sa.GetUnavailableServiceCount() != 0 {
		t.Error("unavailable service count should be 0")
	}
}

func TestService_Methods(t *testing.T) {
	t.Run("IsType", func(t *testing.T) {
		service := Service{Type: "SPOTIFY", IsAvailable: true}
		if !service.IsType(ServiceTypeSpotify) {
			t.Error("service should be of type Spotify")
		}
		if service.IsType(ServiceTypeBluetooth) {
			t.Error("service should not be of type Bluetooth")
		}
	})

	t.Run("GetReason", func(t *testing.T) {
		serviceWithReason := Service{
			Type:        "BLUETOOTH",
			IsAvailable: false,
			Reason:      "DEVICE_NOT_CONNECTED",
		}
		if serviceWithReason.GetReason() != "DEVICE_NOT_CONNECTED" {
			t.Errorf("expected DEVICE_NOT_CONNECTED, got %s", serviceWithReason.GetReason())
		}

		serviceWithoutReason := Service{Type: "SPOTIFY", IsAvailable: true}
		if serviceWithoutReason.GetReason() != "" {
			t.Errorf("expected empty reason, got %s", serviceWithoutReason.GetReason())
		}
	})
}

func TestServiceType_Constants(t *testing.T) {
	// Test that all service type constants are properly defined
	if ServiceTypeAirPlay != ServiceType("AIRPLAY") {
		t.Error("ServiceTypeAirPlay constant mismatch")
	}
	if ServiceTypeAlexa != ServiceType("ALEXA") {
		t.Error("ServiceTypeAlexa constant mismatch")
	}
	if ServiceTypeAmazon != ServiceType("AMAZON") {
		t.Error("ServiceTypeAmazon constant mismatch")
	}
	if ServiceTypeBluetooth != ServiceType("BLUETOOTH") {
		t.Error("ServiceTypeBluetooth constant mismatch")
	}
	if ServiceTypeBMX != ServiceType("BMX") {
		t.Error("ServiceTypeBMX constant mismatch")
	}
	if ServiceTypeDeezer != ServiceType("DEEZER") {
		t.Error("ServiceTypeDeezer constant mismatch")
	}
	if ServiceTypeIHeart != ServiceType("IHEART") {
		t.Error("ServiceTypeIHeart constant mismatch")
	}
	if ServiceTypeLocalInternetRadio != ServiceType("LOCAL_INTERNET_RADIO") {
		t.Error("ServiceTypeLocalInternetRadio constant mismatch")
	}
	if ServiceTypeLocalMusic != ServiceType("LOCAL_MUSIC") {
		t.Error("ServiceTypeLocalMusic constant mismatch")
	}
	if ServiceTypeNotification != ServiceType("NOTIFICATION") {
		t.Error("ServiceTypeNotification constant mismatch")
	}
	if ServiceTypePandora != ServiceType("PANDORA") {
		t.Error("ServiceTypePandora constant mismatch")
	}
	if ServiceTypeSpotify != ServiceType("SPOTIFY") {
		t.Error("ServiceTypeSpotify constant mismatch")
	}
	if ServiceTypeTuneIn != ServiceType("TUNEIN") {
		t.Error("ServiceTypeTuneIn constant mismatch")
	}
}

func TestServiceAvailability_MarshalXML(t *testing.T) {
	sa := &ServiceAvailability{
		Services: &ServiceList{
			Service: []Service{
				{Type: "SPOTIFY", IsAvailable: true},
				{Type: "BLUETOOTH", IsAvailable: false, Reason: "UNAVAILABLE"},
			},
		},
	}

	data, err := xml.Marshal(sa)
	if err != nil {
		t.Fatalf("failed to marshal XML: %v", err)
	}

	// Unmarshal back to verify roundtrip
	var unmarshaled ServiceAvailability
	err = xml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal XML: %v", err)
	}

	if sa.GetServiceCount() != unmarshaled.GetServiceCount() {
		t.Error("service count mismatch after roundtrip")
	}
	if sa.HasSpotify() != unmarshaled.HasSpotify() {
		t.Error("Spotify availability mismatch after roundtrip")
	}
	if sa.HasBluetooth() != unmarshaled.HasBluetooth() {
		t.Error("Bluetooth availability mismatch after roundtrip")
	}

	bluetoothService := unmarshaled.GetServiceByType(ServiceTypeBluetooth)
	if bluetoothService == nil {
		t.Fatal("bluetooth service should not be nil after roundtrip")
	}
	if bluetoothService.Reason != "UNAVAILABLE" {
		t.Errorf("expected UNAVAILABLE reason, got %s", bluetoothService.Reason)
	}
}

func BenchmarkServiceAvailability_GetAvailableServices(b *testing.B) {
	sa := &ServiceAvailability{
		Services: &ServiceList{
			Service: []Service{
				{Type: "SPOTIFY", IsAvailable: true},
				{Type: "BLUETOOTH", IsAvailable: false},
				{Type: "AIRPLAY", IsAvailable: true},
				{Type: "ALEXA", IsAvailable: false},
				{Type: "PANDORA", IsAvailable: true},
				{Type: "TUNEIN", IsAvailable: true},
				{Type: "AMAZON", IsAvailable: false},
				{Type: "DEEZER", IsAvailable: true},
			},
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sa.GetAvailableServices()
	}
}

func BenchmarkServiceAvailability_IsServiceAvailable(b *testing.B) {
	sa := &ServiceAvailability{
		Services: &ServiceList{
			Service: []Service{
				{Type: "SPOTIFY", IsAvailable: true},
				{Type: "BLUETOOTH", IsAvailable: false},
				{Type: "AIRPLAY", IsAvailable: true},
				{Type: "ALEXA", IsAvailable: false},
				{Type: "PANDORA", IsAvailable: true},
			},
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sa.IsServiceAvailable(ServiceTypeSpotify)
	}
}
