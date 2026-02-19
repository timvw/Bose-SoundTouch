package handlers

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// proxySpeakerGET forwards a GET request to the speaker and returns the raw response body.
func (s *Server) proxySpeakerGET(ip, path string) ([]byte, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(fmt.Sprintf("http://%s:8090%s", ip, path))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// proxySpeakerPOST forwards a POST request with XML body to the speaker and returns the raw response body.
func (s *Server) proxySpeakerPOST(ip, path string, xmlBody []byte) ([]byte, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Post(
		fmt.Sprintf("http://%s:8090%s", ip, path),
		"text/xml",
		bytes.NewReader(xmlBody),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// writeJSON is a convenience helper to marshal and write a JSON response.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeJSONError writes a JSON error response.
func writeJSONError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// HandleAPISpeakersList returns all discovered speakers.
func (s *Server) HandleAPISpeakersList(w http.ResponseWriter, _ *http.Request) {
	allDevices, err := s.ds.ListAllDevices()
	if err != nil {
		log.Printf("[SpeakerProxy] Failed to list devices: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "failed to list devices")

		return
	}

	type speakerInfo struct {
		IPAddress string `json:"ipAddress"`
		Name      string `json:"name"`
		DeviceID  string `json:"deviceId"`
		Type      string `json:"type"`
	}

	speakers := make([]speakerInfo, 0, len(allDevices))
	for _, d := range allDevices {
		speakers = append(speakers, speakerInfo{
			IPAddress: d.IPAddress,
			Name:      d.Name,
			DeviceID:  d.DeviceID,
			Type:      d.ProductCode,
		})
	}

	writeJSON(w, http.StatusOK, speakers)
}

// --- /api/speakers/{id}/info ---

type xmlInfo struct {
	XMLName          xml.Name `xml:"info"`
	DeviceID         string   `xml:"deviceID,attr"`
	Name             string   `xml:"name"`
	Type             string   `xml:"type"`
	MargeURL         string   `xml:"margeURL"`
	MargeAccountUUID string   `xml:"margeAccountUUID"`
}

// HandleAPISpeakerInfo proxies GET :8090/info and returns JSON.
func (s *Server) HandleAPISpeakerInfo(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "id")

	data, err := s.proxySpeakerGET(ip, "/info")
	if err != nil {
		log.Printf("[SpeakerProxy] info error for %s: %v", ip, err)
		writeJSONError(w, http.StatusBadGateway, "failed to reach speaker")

		return
	}

	var info xmlInfo
	if err := xml.Unmarshal(data, &info); err != nil {
		log.Printf("[SpeakerProxy] info XML parse error for %s: %v", ip, err)
		writeJSONError(w, http.StatusInternalServerError, "failed to parse speaker response")

		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"deviceID":    info.DeviceID,
		"name":        info.Name,
		"type":        info.Type,
		"margeURL":    info.MargeURL,
		"accountUUID": info.MargeAccountUUID,
	})
}

// --- /api/speakers/{id}/volume ---

type xmlVolume struct {
	XMLName      xml.Name `xml:"volume"`
	TargetVolume int      `xml:"targetvolume"`
	ActualVolume int      `xml:"actualvolume"`
	MuteEnabled  bool     `xml:"muteenabled"`
}

// HandleAPISpeakerVolume proxies GET :8090/volume and returns JSON.
func (s *Server) HandleAPISpeakerVolume(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "id")

	data, err := s.proxySpeakerGET(ip, "/volume")
	if err != nil {
		log.Printf("[SpeakerProxy] volume error for %s: %v", ip, err)
		writeJSONError(w, http.StatusBadGateway, "failed to reach speaker")

		return
	}

	var vol xmlVolume
	if err := xml.Unmarshal(data, &vol); err != nil {
		log.Printf("[SpeakerProxy] volume XML parse error for %s: %v", ip, err)
		writeJSONError(w, http.StatusInternalServerError, "failed to parse speaker response")

		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"targetVolume": vol.TargetVolume,
		"actualVolume": vol.ActualVolume,
		"muteEnabled":  vol.MuteEnabled,
	})
}

// HandleAPISpeakerSetVolume proxies POST :8090/volume with an XML volume element.
func (s *Server) HandleAPISpeakerSetVolume(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "id")

	var req struct {
		Volume int `json:"volume"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	xmlBody := []byte(fmt.Sprintf("<volume>%d</volume>", req.Volume))

	_, err := s.proxySpeakerPOST(ip, "/volume", xmlBody)
	if err != nil {
		log.Printf("[SpeakerProxy] set volume error for %s: %v", ip, err)
		writeJSONError(w, http.StatusBadGateway, "failed to reach speaker")

		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// --- /api/speakers/{id}/now-playing ---

type xmlNowPlaying struct {
	XMLName    xml.Name `xml:"nowPlaying"`
	Source     string   `xml:"source,attr"`
	Track      string   `xml:"track"`
	Artist     string   `xml:"artist"`
	Album      string   `xml:"album"`
	Art        string   `xml:"art"`
	PlayStatus string   `xml:"playStatus"`

	ContentItem struct {
		Source        string `xml:"source,attr"`
		Location      string `xml:"location,attr"`
		SourceAccount string `xml:"sourceAccount,attr"`
	} `xml:"ContentItem"`
}

// HandleAPISpeakerNowPlaying proxies GET :8090/nowPlaying and returns JSON.
func (s *Server) HandleAPISpeakerNowPlaying(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "id")

	data, err := s.proxySpeakerGET(ip, "/nowPlaying")
	if err != nil {
		log.Printf("[SpeakerProxy] nowPlaying error for %s: %v", ip, err)
		writeJSONError(w, http.StatusBadGateway, "failed to reach speaker")

		return
	}

	var np xmlNowPlaying
	if err := xml.Unmarshal(data, &np); err != nil {
		log.Printf("[SpeakerProxy] nowPlaying XML parse error for %s: %v", ip, err)
		writeJSONError(w, http.StatusInternalServerError, "failed to parse speaker response")

		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"source":     np.Source,
		"track":      np.Track,
		"artist":     np.Artist,
		"album":      np.Album,
		"art":        np.Art,
		"playStatus": np.PlayStatus,
		"contentItem": map[string]string{
			"source":        np.ContentItem.Source,
			"location":      np.ContentItem.Location,
			"sourceAccount": np.ContentItem.SourceAccount,
		},
	})
}

// --- /api/speakers/{id}/play-control ---

// HandleAPISpeakerPlayControl proxies POST :8090/userPlayControl with an XML PlayControl element.
func (s *Server) HandleAPISpeakerPlayControl(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "id")

	var req struct {
		Control string `json:"control"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Control == "" {
		writeJSONError(w, http.StatusBadRequest, "missing control field")
		return
	}

	xmlBody := []byte(fmt.Sprintf("<PlayControl>%s</PlayControl>", req.Control))

	_, err := s.proxySpeakerPOST(ip, "/userPlayControl", xmlBody)
	if err != nil {
		log.Printf("[SpeakerProxy] playControl error for %s: %v", ip, err)
		writeJSONError(w, http.StatusBadGateway, "failed to reach speaker")

		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// --- /api/speakers/{id}/presets ---

type xmlPresets struct {
	XMLName xml.Name    `xml:"presets"`
	Presets []xmlPreset `xml:"preset"`
}

type xmlPreset struct {
	ID          int `xml:"id,attr"`
	ContentItem struct {
		Source        string `xml:"source,attr"`
		Type          string `xml:"type,attr"`
		Location      string `xml:"location,attr"`
		SourceAccount string `xml:"sourceAccount,attr"`
		Name          string `xml:"itemName"`
		Image         string `xml:"containerArt"`
	} `xml:"ContentItem"`
}

// HandleAPISpeakerPresets proxies GET :8090/presets and returns JSON.
func (s *Server) HandleAPISpeakerPresets(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "id")

	data, err := s.proxySpeakerGET(ip, "/presets")
	if err != nil {
		log.Printf("[SpeakerProxy] presets error for %s: %v", ip, err)
		writeJSONError(w, http.StatusBadGateway, "failed to reach speaker")

		return
	}

	var presets xmlPresets
	if err := xml.Unmarshal(data, &presets); err != nil {
		log.Printf("[SpeakerProxy] presets XML parse error for %s: %v", ip, err)
		writeJSONError(w, http.StatusInternalServerError, "failed to parse speaker response")

		return
	}

	type presetJSON struct {
		ID            int    `json:"id"`
		Source        string `json:"source"`
		Type          string `json:"type"`
		Location      string `json:"location"`
		SourceAccount string `json:"sourceAccount"`
		Name          string `json:"name"`
		Image         string `json:"image"`
	}

	result := make([]presetJSON, 0, len(presets.Presets))
	for _, p := range presets.Presets {
		result = append(result, presetJSON{
			ID:            p.ID,
			Source:        p.ContentItem.Source,
			Type:          p.ContentItem.Type,
			Location:      p.ContentItem.Location,
			SourceAccount: p.ContentItem.SourceAccount,
			Name:          p.ContentItem.Name,
			Image:         p.ContentItem.Image,
		})
	}

	writeJSON(w, http.StatusOK, result)
}

// --- /api/speakers/{id}/store-preset ---

// HandleAPISpeakerStorePreset proxies POST :8090/storePreset with an XML ContentItem.
func (s *Server) HandleAPISpeakerStorePreset(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "id")

	var req struct {
		ID            int    `json:"id"`
		Source        string `json:"source"`
		Type          string `json:"type"`
		Location      string `json:"location"`
		SourceAccount string `json:"sourceAccount"`
		Name          string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	xmlBody := []byte(fmt.Sprintf(
		`<preset id="%d"><ContentItem source="%s" type="%s" location="%s" sourceAccount="%s"><itemName>%s</itemName></ContentItem></preset>`,
		req.ID, req.Source, req.Type, req.Location, req.SourceAccount, req.Name,
	))

	_, err := s.proxySpeakerPOST(ip, "/storePreset", xmlBody)
	if err != nil {
		log.Printf("[SpeakerProxy] storePreset error for %s: %v", ip, err)
		writeJSONError(w, http.StatusBadGateway, "failed to reach speaker")

		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// --- /api/speakers/{id}/remove-preset ---

// HandleAPISpeakerRemovePreset proxies POST :8090/removePreset with an XML preset element.
func (s *Server) HandleAPISpeakerRemovePreset(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "id")

	var req struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	xmlBody := []byte(fmt.Sprintf(`<preset id="%d"></preset>`, req.ID))

	_, err := s.proxySpeakerPOST(ip, "/removePreset", xmlBody)
	if err != nil {
		log.Printf("[SpeakerProxy] removePreset error for %s: %v", ip, err)
		writeJSONError(w, http.StatusBadGateway, "failed to reach speaker")

		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// --- /api/speakers/{id}/recents ---

type xmlRecents struct {
	XMLName xml.Name    `xml:"recents"`
	Recents []xmlRecent `xml:"recent"`
}

type xmlRecent struct {
	DeviceID    string `xml:"deviceID,attr"`
	ContentItem struct {
		Source        string `xml:"source,attr"`
		Type          string `xml:"type,attr"`
		Location      string `xml:"location,attr"`
		SourceAccount string `xml:"sourceAccount,attr"`
		Name          string `xml:"itemName"`
		Image         string `xml:"containerArt"`
	} `xml:"ContentItem"`
}

// HandleAPISpeakerRecents proxies GET :8090/recents and returns JSON.
func (s *Server) HandleAPISpeakerRecents(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "id")

	data, err := s.proxySpeakerGET(ip, "/recents")
	if err != nil {
		log.Printf("[SpeakerProxy] recents error for %s: %v", ip, err)
		writeJSONError(w, http.StatusBadGateway, "failed to reach speaker")

		return
	}

	var recents xmlRecents
	if err := xml.Unmarshal(data, &recents); err != nil {
		log.Printf("[SpeakerProxy] recents XML parse error for %s: %v", ip, err)
		writeJSONError(w, http.StatusInternalServerError, "failed to parse speaker response")

		return
	}

	type recentJSON struct {
		Source        string `json:"source"`
		Type          string `json:"type"`
		Location      string `json:"location"`
		SourceAccount string `json:"sourceAccount"`
		Name          string `json:"name"`
		Image         string `json:"image"`
	}

	result := make([]recentJSON, 0, len(recents.Recents))
	for _, rc := range recents.Recents {
		result = append(result, recentJSON{
			Source:        rc.ContentItem.Source,
			Type:          rc.ContentItem.Type,
			Location:      rc.ContentItem.Location,
			SourceAccount: rc.ContentItem.SourceAccount,
			Name:          rc.ContentItem.Name,
			Image:         rc.ContentItem.Image,
		})
	}

	writeJSON(w, http.StatusOK, result)
}

// --- /api/speakers/{id}/select ---

// HandleAPISpeakerSelect proxies POST :8090/select with an XML ContentItem.
func (s *Server) HandleAPISpeakerSelect(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "id")

	var req struct {
		Source        string `json:"source"`
		Type          string `json:"type"`
		Location      string `json:"location"`
		SourceAccount string `json:"sourceAccount"`
		Name          string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	xmlBody := []byte(fmt.Sprintf(
		`<ContentItem source="%s" type="%s" location="%s" sourceAccount="%s"><itemName>%s</itemName></ContentItem>`,
		req.Source, req.Type, req.Location, req.SourceAccount, req.Name,
	))

	_, err := s.proxySpeakerPOST(ip, "/select", xmlBody)
	if err != nil {
		log.Printf("[SpeakerProxy] select error for %s: %v", ip, err)
		writeJSONError(w, http.StatusBadGateway, "failed to reach speaker")

		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// --- /api/speakers/{id}/key ---

// HandleAPISpeakerKey proxies POST :8090/key with an XML key element.
func (s *Server) HandleAPISpeakerKey(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "id")

	var req struct {
		Key   string `json:"key"`
		State string `json:"state"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Key == "" {
		writeJSONError(w, http.StatusBadRequest, "missing key field")
		return
	}

	if req.State == "" {
		req.State = "press"
	}

	xmlBody := []byte(fmt.Sprintf(`<key state="%s" sender="Gabbo">%s</key>`, req.State, req.Key))

	_, err := s.proxySpeakerPOST(ip, "/key", xmlBody)
	if err != nil {
		log.Printf("[SpeakerProxy] key error for %s: %v", ip, err)
		writeJSONError(w, http.StatusBadGateway, "failed to reach speaker")

		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// --- /api/speakers/{id}/standby ---

// HandleAPISpeakerStandby proxies GET :8090/standby and returns 200 OK.
func (s *Server) HandleAPISpeakerStandby(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "id")

	_, err := s.proxySpeakerGET(ip, "/standby")
	if err != nil {
		log.Printf("[SpeakerProxy] standby error for %s: %v", ip, err)
		writeJSONError(w, http.StatusBadGateway, "failed to reach speaker")

		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// --- /api/speakers/{id}/name ---

// HandleAPISpeakerSetName proxies POST :8090/name with an XML name element.
func (s *Server) HandleAPISpeakerSetName(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "id")

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		writeJSONError(w, http.StatusBadRequest, "missing name field")
		return
	}

	xmlBody := []byte(fmt.Sprintf("<name>%s</name>", req.Name))

	_, err := s.proxySpeakerPOST(ip, "/name", xmlBody)
	if err != nil {
		log.Printf("[SpeakerProxy] setName error for %s: %v", ip, err)
		writeJSONError(w, http.StatusBadGateway, "failed to reach speaker")

		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// --- /api/speakers/{id}/zones ---

type xmlZone struct {
	XMLName xml.Name        `xml:"zone"`
	Master  string          `xml:"master,attr"`
	Members []xmlZoneMember `xml:"member"`
}

type xmlZoneMember struct {
	IPAddress string `xml:"ipaddress,attr"`
	DeviceID  string `xml:",chardata"`
}

// HandleAPISpeakerZones proxies GET :8090/getZone and returns JSON.
func (s *Server) HandleAPISpeakerZones(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "id")

	data, err := s.proxySpeakerGET(ip, "/getZone")
	if err != nil {
		log.Printf("[SpeakerProxy] getZone error for %s: %v", ip, err)
		writeJSONError(w, http.StatusBadGateway, "failed to reach speaker")

		return
	}

	var zone xmlZone
	if err := xml.Unmarshal(data, &zone); err != nil {
		log.Printf("[SpeakerProxy] getZone XML parse error for %s: %v", ip, err)
		writeJSONError(w, http.StatusInternalServerError, "failed to parse speaker response")

		return
	}

	type memberJSON struct {
		IPAddress string `json:"ipAddress"`
		DeviceID  string `json:"deviceId"`
	}

	members := make([]memberJSON, 0, len(zone.Members))
	for _, m := range zone.Members {
		members = append(members, memberJSON{
			IPAddress: m.IPAddress,
			DeviceID:  m.DeviceID,
		})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"master":  zone.Master,
		"members": members,
	})
}
