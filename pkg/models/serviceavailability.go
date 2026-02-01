package models

import "encoding/xml"

// ServiceAvailability represents the response from /serviceAvailability endpoint
type ServiceAvailability struct {
	XMLName  xml.Name     `xml:"serviceAvailability"`
	Services *ServiceList `xml:"services"`
}

// ServiceList contains the list of available services
type ServiceList struct {
	Service []Service `xml:"service"`
}

// Service represents an individual service availability status
type Service struct {
	Type        string `xml:"type,attr"`
	IsAvailable bool   `xml:"isAvailable,attr"`
	Reason      string `xml:"reason,attr,omitempty"`
}

// ServiceType represents known service types
type ServiceType string

const (
	// ServiceTypeAirPlay represents Apple AirPlay streaming service
	ServiceTypeAirPlay ServiceType = "AIRPLAY"
	// ServiceTypeAlexa represents Amazon Alexa voice assistant integration
	ServiceTypeAlexa ServiceType = "ALEXA"
	// ServiceTypeAmazon represents Amazon Music streaming service
	ServiceTypeAmazon ServiceType = "AMAZON"
	// ServiceTypeBluetooth represents Bluetooth audio connectivity
	ServiceTypeBluetooth ServiceType = "BLUETOOTH"
	// ServiceTypeBMX represents BMX streaming service
	ServiceTypeBMX ServiceType = "BMX"
	// ServiceTypeDeezer represents Deezer music streaming service
	ServiceTypeDeezer ServiceType = "DEEZER"
	// ServiceTypeIHeart represents iHeartRadio streaming service
	ServiceTypeIHeart ServiceType = "IHEART"
	// ServiceTypeLocalInternetRadio represents local internet radio stations
	ServiceTypeLocalInternetRadio ServiceType = "LOCAL_INTERNET_RADIO"
	// ServiceTypeLocalMusic represents local music library
	ServiceTypeLocalMusic ServiceType = "LOCAL_MUSIC"
	// ServiceTypeNotification represents system notifications
	ServiceTypeNotification ServiceType = "NOTIFICATION"
	// ServiceTypePandora represents Pandora music streaming service
	ServiceTypePandora ServiceType = "PANDORA"
	// ServiceTypeSpotify represents Spotify music streaming service
	ServiceTypeSpotify ServiceType = "SPOTIFY"
	// ServiceTypeTuneIn represents TuneIn internet radio service
	ServiceTypeTuneIn ServiceType = "TUNEIN"
)

// GetReason returns the reason why a service is unavailable (if any)
func (s *Service) GetReason() string {
	return s.Reason
}

// IsType checks if the service is of a specific type
func (s *Service) IsType(serviceType ServiceType) bool {
	return s.Type == string(serviceType)
}

// GetAvailableServices returns only services that are available
func (sa *ServiceAvailability) GetAvailableServices() []Service {
	if sa.Services == nil {
		return []Service{}
	}

	var available []Service
	for _, service := range sa.Services.Service {
		if service.IsAvailable {
			available = append(available, service)
		}
	}
	return available
}

// GetUnavailableServices returns only services that are unavailable
func (sa *ServiceAvailability) GetUnavailableServices() []Service {
	if sa.Services == nil {
		return []Service{}
	}

	var unavailable []Service
	for _, service := range sa.Services.Service {
		if !service.IsAvailable {
			unavailable = append(unavailable, service)
		}
	}
	return unavailable
}

// IsServiceAvailable checks if a specific service type is available
func (sa *ServiceAvailability) IsServiceAvailable(serviceType ServiceType) bool {
	if sa.Services == nil {
		return false
	}

	for _, service := range sa.Services.Service {
		if service.Type == string(serviceType) && service.IsAvailable {
			return true
		}
	}
	return false
}

// GetServiceByType returns the service information for a specific type
func (sa *ServiceAvailability) GetServiceByType(serviceType ServiceType) *Service {
	if sa.Services == nil {
		return nil
	}

	for _, service := range sa.Services.Service {
		if service.Type == string(serviceType) {
			return &service
		}
	}
	return nil
}

// HasSpotify returns true if Spotify service is available
func (sa *ServiceAvailability) HasSpotify() bool {
	return sa.IsServiceAvailable(ServiceTypeSpotify)
}

// HasAlexa returns true if Alexa service is available
func (sa *ServiceAvailability) HasAlexa() bool {
	return sa.IsServiceAvailable(ServiceTypeAlexa)
}

// HasBluetooth returns true if Bluetooth service is available
func (sa *ServiceAvailability) HasBluetooth() bool {
	return sa.IsServiceAvailable(ServiceTypeBluetooth)
}

// HasAirPlay returns true if AirPlay service is available
func (sa *ServiceAvailability) HasAirPlay() bool {
	return sa.IsServiceAvailable(ServiceTypeAirPlay)
}

// HasTuneIn returns true if TuneIn service is available
func (sa *ServiceAvailability) HasTuneIn() bool {
	return sa.IsServiceAvailable(ServiceTypeTuneIn)
}

// HasPandora returns true if Pandora service is available
func (sa *ServiceAvailability) HasPandora() bool {
	return sa.IsServiceAvailable(ServiceTypePandora)
}

// HasLocalMusic returns true if Local Music service is available
func (sa *ServiceAvailability) HasLocalMusic() bool {
	return sa.IsServiceAvailable(ServiceTypeLocalMusic)
}

// GetStreamingServices returns all streaming service types
func (sa *ServiceAvailability) GetStreamingServices() []Service {
	if sa.Services == nil {
		return []Service{}
	}

	streamingTypes := []ServiceType{
		ServiceTypeSpotify,
		ServiceTypePandora,
		ServiceTypeTuneIn,
		ServiceTypeAmazon,
		ServiceTypeDeezer,
		ServiceTypeIHeart,
		ServiceTypeLocalInternetRadio,
	}

	var streaming []Service
	for _, service := range sa.Services.Service {
		for _, streamingType := range streamingTypes {
			if service.Type == string(streamingType) {
				streaming = append(streaming, service)
				break
			}
		}
	}
	return streaming
}

// GetLocalServices returns all local service types
func (sa *ServiceAvailability) GetLocalServices() []Service {
	if sa.Services == nil {
		return []Service{}
	}

	localTypes := []ServiceType{
		ServiceTypeBluetooth,
		ServiceTypeAirPlay,
		ServiceTypeLocalMusic,
	}

	var local []Service
	for _, service := range sa.Services.Service {
		for _, localType := range localTypes {
			if service.Type == string(localType) {
				local = append(local, service)
				break
			}
		}
	}
	return local
}

// GetServiceCount returns the total number of services
func (sa *ServiceAvailability) GetServiceCount() int {
	if sa.Services == nil {
		return 0
	}
	return len(sa.Services.Service)
}

// GetAvailableServiceCount returns the number of available services
func (sa *ServiceAvailability) GetAvailableServiceCount() int {
	return len(sa.GetAvailableServices())
}

// GetUnavailableServiceCount returns the number of unavailable services
func (sa *ServiceAvailability) GetUnavailableServiceCount() int {
	return len(sa.GetUnavailableServices())
}
