package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	// ZeroConfPort is the port used by SoundTouch speakers for ZeroConf.
	ZeroConfPort = 8200
)

// TrackedSpeaker represents a speaker registered for ZeroConf priming.
type TrackedSpeaker struct {
	AccountID     string
	DeviceID      string
	IPAddress     string
	LastPrimed    time.Time
	PrimeFailures int
}

// ZeroConfPrimer manages Spotify token priming for SoundTouch speakers.
type ZeroConfPrimer struct {
	spotify  *SpotifyService
	speakers map[string]*TrackedSpeaker // key: deviceID or IP
	mu       sync.RWMutex
	ticker   *time.Ticker
	done     chan struct{}
	client   *http.Client

	// speakerURL builds the ZeroConf URL for a speaker IP.
	// Defaults to http://{ip}:8200. Overridable for testing.
	speakerURL func(ip string) string
}

// NewZeroConfPrimer creates a new ZeroConfPrimer.
func NewZeroConfPrimer(spotify *SpotifyService) *ZeroConfPrimer {
	return &ZeroConfPrimer{
		spotify:  spotify,
		speakers: make(map[string]*TrackedSpeaker),
		done:     make(chan struct{}),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		speakerURL: func(ip string) string {
			return fmt.Sprintf("http://%s:%d", ip, ZeroConfPort)
		},
	}
}

// RegisterSpeaker adds or updates a speaker in the tracking map.
func (p *ZeroConfPrimer) RegisterSpeaker(accountID, deviceID, ipAddress string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	key := deviceID
	if key == "" {
		key = ipAddress
	}
	if key == "" {
		return
	}

	if existing, ok := p.speakers[key]; ok {
		if accountID != "" {
			existing.AccountID = accountID
		}
		if ipAddress != "" {
			existing.IPAddress = ipAddress
		}
		return
	}

	p.speakers[key] = &TrackedSpeaker{
		AccountID: accountID,
		DeviceID:  deviceID,
		IPAddress: ipAddress,
	}
	log.Printf("[ZeroConf] Speaker registered: key=%s account=%s ip=%s", key, accountID, ipAddress)
}

// OnPowerOn is called when a speaker sends a power_on event.
// It registers the speaker and retries priming with backoff delays.
func (p *ZeroConfPrimer) OnPowerOn(accountID, deviceID, ipAddress string) {
	p.RegisterSpeaker(accountID, deviceID, ipAddress)

	go func() {
		delays := []time.Duration{5 * time.Second, 10 * time.Second, 20 * time.Second}

		p.mu.RLock()
		speakers := make([]*TrackedSpeaker, 0, len(p.speakers))
		for _, s := range p.speakers {
			if s.IPAddress != "" {
				speakers = append(speakers, s)
			}
		}
		p.mu.RUnlock()

		if len(speakers) == 0 {
			log.Printf("[ZeroConf] No speakers with IP addresses to prime")
			return
		}

		for _, d := range delays {
			log.Printf("[ZeroConf] Speaker booted — waiting %v before priming %d speaker(s)...", d, len(speakers))
			time.Sleep(d)

			allOK := true
			for _, speaker := range speakers {
				if err := p.PrimeSpeaker(speaker); err != nil {
					allOK = false
				}
			}

			if allOK {
				log.Printf("[ZeroConf] All speakers primed successfully")
				return
			}
		}

		log.Printf("[ZeroConf] Some speakers failed to prime after all retries")
	}()
}

// StartPeriodic starts a goroutine that primes all speakers at the given interval.
func (p *ZeroConfPrimer) StartPeriodic(interval time.Duration) {
	p.ticker = time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-p.ticker.C:
				p.primeAll()
			case <-p.done:
				return
			}
		}
	}()
}

// Stop stops the periodic primer.
func (p *ZeroConfPrimer) Stop() {
	if p.ticker != nil {
		p.ticker.Stop()
	}
	close(p.done)
}

// PrimeSpeaker checks whether a speaker needs priming and sends the addUser
// command if the speaker has no active Spotify user.
func (p *ZeroConfPrimer) PrimeSpeaker(speaker *TrackedSpeaker) error {
	if speaker.IPAddress == "" {
		return fmt.Errorf("no IP address for speaker %s", speaker.DeviceID)
	}

	// Check if speaker already has an active user
	activeUser, err := p.getActiveUser(speaker.IPAddress)
	if err != nil {
		log.Printf("[ZeroConf] Could not check activeUser for %s: %v", speaker.IPAddress, err)
		// Continue to try priming anyway
	} else if activeUser != "" {
		log.Printf("[ZeroConf] Speaker %s already primed (activeUser=%s)", speaker.IPAddress, activeUser)
		speaker.LastPrimed = time.Now()
		return nil
	}

	// Get fresh token
	accessToken, username, err := p.spotify.GetFreshToken()
	if err != nil {
		speaker.PrimeFailures++
		return fmt.Errorf("get token: %w", err)
	}

	// Send addUser
	if err := p.sendAddUser(speaker.IPAddress, username, accessToken); err != nil {
		speaker.PrimeFailures++
		return fmt.Errorf("addUser to %s: %w", speaker.IPAddress, err)
	}

	log.Printf("[ZeroConf] Speaker %s primed for Spotify (user=%s)", speaker.IPAddress, username)
	speaker.LastPrimed = time.Now()
	speaker.PrimeFailures = 0
	return nil
}

// primeAll iterates all registered speakers and primes each.
func (p *ZeroConfPrimer) primeAll() {
	p.mu.RLock()
	speakers := make([]*TrackedSpeaker, 0, len(p.speakers))
	for _, s := range p.speakers {
		if s.IPAddress != "" {
			speakers = append(speakers, s)
		}
	}
	p.mu.RUnlock()

	log.Printf("[ZeroConf] Periodic primer check: %d speaker(s)", len(speakers))

	for _, speaker := range speakers {
		if err := p.PrimeSpeaker(speaker); err != nil {
			log.Printf("[ZeroConf] Failed to prime %s: %v", speaker.IPAddress, err)
		}
	}
}

// getActiveUser queries the speaker's ZeroConf endpoint for the active Spotify user.
// The response is JSON: {"activeUser": "username", ...}
func (p *ZeroConfPrimer) getActiveUser(ip string) (string, error) {
	reqURL := p.speakerURL(ip) + "/zc?action=getInfo"

	resp, err := p.client.Get(reqURL)
	if err != nil {
		return "", fmt.Errorf("getInfo request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read getInfo response: %w", err)
	}

	var info map[string]interface{}
	if err := json.Unmarshal(body, &info); err != nil {
		return "", fmt.Errorf("parse getInfo response: %w", err)
	}

	activeUser, _ := info["activeUser"].(string)
	return activeUser, nil
}

// sendAddUser sends a Spotify token to the speaker's ZeroConf endpoint.
func (p *ZeroConfPrimer) sendAddUser(ip, username, token string) error {
	reqURL := p.speakerURL(ip) + "/zc"

	formData := url.Values{
		"action":    {"addUser"},
		"userName":  {username},
		"blob":      {token},
		"clientKey": {""},
		"tokenType": {"accesstoken"},
	}

	resp, err := p.client.Post(reqURL, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("addUser request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read addUser response: %w", err)
	}

	// Parse response to check status
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("parse addUser response: %w", err)
	}

	// Status 101 means success in the Spotify ZeroConf protocol
	if status, ok := result["status"].(float64); ok && int(status) != 101 {
		statusString, _ := result["statusString"].(string)
		return fmt.Errorf("addUser returned status %d: %s", int(status), statusString)
	}

	return nil
}
