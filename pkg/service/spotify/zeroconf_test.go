package spotify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestRegisterSpeaker(t *testing.T) {
	svc := NewSpotifyService("cid", "csecret", "http://localhost/cb", t.TempDir())
	primer := NewZeroConfPrimer(svc)

	primer.RegisterSpeaker("acc1", "dev1", "192.168.1.100")

	primer.mu.RLock()
	defer primer.mu.RUnlock()

	speaker, ok := primer.speakers["dev1"]
	if !ok {
		t.Fatal("speaker dev1 not found in map")
	}
	if speaker.AccountID != "acc1" {
		t.Errorf("expected AccountID acc1, got %s", speaker.AccountID)
	}
	if speaker.DeviceID != "dev1" {
		t.Errorf("expected DeviceID dev1, got %s", speaker.DeviceID)
	}
	if speaker.IPAddress != "192.168.1.100" {
		t.Errorf("expected IPAddress 192.168.1.100, got %s", speaker.IPAddress)
	}
}

func TestRegisterSpeakerByIPWhenNoDeviceID(t *testing.T) {
	svc := NewSpotifyService("cid", "csecret", "http://localhost/cb", t.TempDir())
	primer := NewZeroConfPrimer(svc)

	primer.RegisterSpeaker("acc1", "", "192.168.1.100")

	primer.mu.RLock()
	defer primer.mu.RUnlock()

	_, ok := primer.speakers["192.168.1.100"]
	if !ok {
		t.Fatal("speaker should be keyed by IP when deviceID is empty")
	}
}

func TestRegisterSpeakerUpdatesExisting(t *testing.T) {
	svc := NewSpotifyService("cid", "csecret", "http://localhost/cb", t.TempDir())
	primer := NewZeroConfPrimer(svc)

	primer.RegisterSpeaker("acc1", "dev1", "192.168.1.100")
	primer.RegisterSpeaker("acc2", "dev1", "192.168.1.200")

	primer.mu.RLock()
	defer primer.mu.RUnlock()

	if len(primer.speakers) != 1 {
		t.Fatalf("expected 1 speaker, got %d", len(primer.speakers))
	}

	speaker := primer.speakers["dev1"]
	if speaker.AccountID != "acc2" {
		t.Errorf("expected updated AccountID acc2, got %s", speaker.AccountID)
	}
	if speaker.IPAddress != "192.168.1.200" {
		t.Errorf("expected updated IPAddress 192.168.1.200, got %s", speaker.IPAddress)
	}
}

func TestPrimeSpeakerAlreadyActive(t *testing.T) {
	svc := NewSpotifyService("cid", "csecret", "http://localhost/cb", t.TempDir())

	// Add an account so GetFreshToken would work (though it shouldn't be called)
	svc.mu.Lock()
	svc.accounts["user1"] = &SpotifyAccount{
		UserID:       "user1",
		AccessToken:  "token",
		RefreshToken: "refresh",
		ExpiresAt:    time.Now().Add(1 * time.Hour).Unix(),
	}
	svc.mu.Unlock()

	var addUserCalled atomic.Int32

	// Mock speaker that already has an active user
	speakerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("action") == "getInfo" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":         101,
				"statusString":   "OK",
				"activeUser":     "someone@spotify",
				"libraryVersion": "3.88.29",
			})
			return
		}

		// addUser should NOT be called
		if r.Method == http.MethodPost {
			addUserCalled.Add(1)
		}
	}))
	defer speakerServer.Close()

	primer := NewZeroConfPrimer(svc)
	primer.speakerURL = func(_ string) string { return speakerServer.URL }

	speaker := &TrackedSpeaker{
		AccountID: "acc1",
		DeviceID:  "dev1",
		IPAddress: "192.168.1.100",
	}

	err := primer.PrimeSpeaker(speaker)
	if err != nil {
		t.Fatalf("PrimeSpeaker failed: %v", err)
	}

	// Speaker was already active, so addUser should NOT have been called
	if addUserCalled.Load() != 0 {
		t.Errorf("addUser should not have been called, but was called %d times", addUserCalled.Load())
	}

	// LastPrimed should be set since speaker was already active
	if speaker.LastPrimed.IsZero() {
		t.Error("expected LastPrimed to be set for already-active speaker")
	}
}

func TestPrimeSpeakerNeedsActivation(t *testing.T) {
	svc := NewSpotifyService("cid", "csecret", "http://localhost/cb", t.TempDir())

	// Add an account with a valid token
	svc.mu.Lock()
	svc.accounts["user1"] = &SpotifyAccount{
		UserID:       "user1",
		AccessToken:  "valid-token",
		RefreshToken: "refresh",
		ExpiresAt:    time.Now().Add(1 * time.Hour).Unix(),
	}
	svc.mu.Unlock()

	var receivedBody string
	var addUserCalled atomic.Int32

	// Mock speaker with no active user
	speakerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("action") == "getInfo" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":       101,
				"statusString": "OK",
				"activeUser":   "", // empty = needs priming
			})
			return
		}

		if r.Method == http.MethodPost {
			addUserCalled.Add(1)
			body, _ := io.ReadAll(r.Body)
			receivedBody = string(body)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":       101,
				"statusString": "OK",
			})
		}
	}))
	defer speakerServer.Close()

	primer := NewZeroConfPrimer(svc)
	primer.speakerURL = func(_ string) string { return speakerServer.URL }

	speaker := &TrackedSpeaker{
		AccountID: "acc1",
		DeviceID:  "dev1",
		IPAddress: "192.168.1.100",
	}

	err := primer.PrimeSpeaker(speaker)
	if err != nil {
		t.Fatalf("PrimeSpeaker failed: %v", err)
	}

	if addUserCalled.Load() != 1 {
		t.Errorf("expected addUser to be called once, got %d", addUserCalled.Load())
	}

	// Verify the form body
	if !strings.Contains(receivedBody, "action=addUser") {
		t.Errorf("body should contain action=addUser, got: %s", receivedBody)
	}
	if !strings.Contains(receivedBody, "userName=user1") {
		t.Errorf("body should contain userName=user1, got: %s", receivedBody)
	}
	if !strings.Contains(receivedBody, "blob=valid-token") {
		t.Errorf("body should contain blob=valid-token, got: %s", receivedBody)
	}
	if !strings.Contains(receivedBody, "tokenType=accesstoken") {
		t.Errorf("body should contain tokenType=accesstoken, got: %s", receivedBody)
	}

	// Verify speaker state updated
	if speaker.PrimeFailures != 0 {
		t.Errorf("expected 0 failures, got %d", speaker.PrimeFailures)
	}
	if speaker.LastPrimed.IsZero() {
		t.Error("expected LastPrimed to be set")
	}
}

func TestPrimeSpeakerFailure(t *testing.T) {
	svc := NewSpotifyService("cid", "csecret", "http://localhost/cb", t.TempDir())

	svc.mu.Lock()
	svc.accounts["user1"] = &SpotifyAccount{
		UserID:       "user1",
		AccessToken:  "valid-token",
		RefreshToken: "refresh",
		ExpiresAt:    time.Now().Add(1 * time.Hour).Unix(),
	}
	svc.mu.Unlock()

	// Mock speaker that returns error status on addUser
	speakerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("action") == "getInfo" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":     101,
				"activeUser": "",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":       303,
			"statusString": "ERROR-SPOTIFY-NOT-READY",
		})
	}))
	defer speakerServer.Close()

	primer := NewZeroConfPrimer(svc)
	primer.speakerURL = func(_ string) string { return speakerServer.URL }

	speaker := &TrackedSpeaker{
		AccountID: "acc1",
		DeviceID:  "dev1",
		IPAddress: "192.168.1.100",
	}

	err := primer.PrimeSpeaker(speaker)
	if err == nil {
		t.Fatal("expected error from PrimeSpeaker")
	}

	if speaker.PrimeFailures != 1 {
		t.Errorf("expected 1 failure, got %d", speaker.PrimeFailures)
	}
}

func TestGetActiveUser(t *testing.T) {
	speakerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":         101,
			"statusString":   "OK",
			"activeUser":     "testuser",
			"libraryVersion": "3.88.29",
		})
	}))
	defer speakerServer.Close()

	svc := NewSpotifyService("cid", "csecret", "http://localhost/cb", t.TempDir())
	primer := NewZeroConfPrimer(svc)
	primer.speakerURL = func(_ string) string { return speakerServer.URL }

	user, err := primer.getActiveUser("192.168.1.100")
	if err != nil {
		t.Fatalf("getActiveUser failed: %v", err)
	}
	if user != "testuser" {
		t.Errorf("expected testuser, got %s", user)
	}
}

func TestGetActiveUserEmpty(t *testing.T) {
	speakerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     101,
			"activeUser": "",
		})
	}))
	defer speakerServer.Close()

	svc := NewSpotifyService("cid", "csecret", "http://localhost/cb", t.TempDir())
	primer := NewZeroConfPrimer(svc)
	primer.speakerURL = func(_ string) string { return speakerServer.URL }

	user, err := primer.getActiveUser("192.168.1.100")
	if err != nil {
		t.Fatalf("getActiveUser failed: %v", err)
	}
	if user != "" {
		t.Errorf("expected empty activeUser, got %s", user)
	}
}

func TestSendAddUser(t *testing.T) {
	var receivedContentType string
	var receivedBody string

	speakerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedContentType = r.Header.Get("Content-Type")
		body, _ := io.ReadAll(r.Body)
		receivedBody = string(body)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":       101,
			"statusString": "OK",
		})
	}))
	defer speakerServer.Close()

	svc := NewSpotifyService("cid", "csecret", "http://localhost/cb", t.TempDir())
	primer := NewZeroConfPrimer(svc)
	primer.speakerURL = func(_ string) string { return speakerServer.URL }

	err := primer.sendAddUser("192.168.1.100", "testuser", "testtoken")
	if err != nil {
		t.Fatalf("sendAddUser failed: %v", err)
	}

	if receivedContentType != "application/x-www-form-urlencoded" {
		t.Errorf("expected form content type, got: %s", receivedContentType)
	}

	if !strings.Contains(receivedBody, "action=addUser") {
		t.Errorf("body should contain action=addUser, got: %s", receivedBody)
	}
	if !strings.Contains(receivedBody, "userName=testuser") {
		t.Errorf("body should contain userName=testuser, got: %s", receivedBody)
	}
	if !strings.Contains(receivedBody, "blob=testtoken") {
		t.Errorf("body should contain blob=testtoken, got: %s", receivedBody)
	}
	if !strings.Contains(receivedBody, "tokenType=accesstoken") {
		t.Errorf("body should contain tokenType=accesstoken, got: %s", receivedBody)
	}
}
