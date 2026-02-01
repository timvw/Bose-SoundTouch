package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchTuneInMetadata(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
<!doctype html>
<html>
<head>
    <meta property="og:title" content="WDR 2 Rheinland, 100.4 FM, Köln | Free Internet Radio | TuneIn"/>
    <meta property="og:image" content="https://cdn-radiotime-logos.tunein.com/s213886g.png"/>
</head>
<body></body>
</html>
`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer ts.Close()

	// Temporarily override httpClient to use test server
	oldClient := httpClient
	httpClient = ts.Client()
	defer func() { httpClient = oldClient }()

	metadata, err := fetchTuneInMetadata("https://tunein.com/radio/WDR-2-Rheinland-1004-s213886/")
	if err != nil {
		t.Fatalf("fetchTuneInMetadata() error = %v", err)
	}

	if metadata == nil {
		t.Fatal("fetchTuneInMetadata() returned nil metadata")
	}

	expectedName := "WDR 2 Rheinland"
	if metadata.Name != expectedName {
		t.Errorf("metadata.Name = %v, want %v", metadata.Name, expectedName)
	}

	expectedArtwork := "https://cdn-radiotime-logos.tunein.com/s213886g.png"
	if metadata.Artwork != expectedArtwork {
		t.Errorf("metadata.Artwork = %v, want %v", metadata.Artwork, expectedArtwork)
	}
}

func TestResolveLocation(t *testing.T) {
	tests := []struct {
		name             string
		source           string
		location         string
		expectedSource   string
		expectedLocation string
	}{
		{
			name:             "Plain location",
			source:           "TUNEIN",
			location:         "/v1/playback/station/s213886",
			expectedSource:   "TUNEIN",
			expectedLocation: "/v1/playback/station/s213886",
		},
		{
			name:             "TuneIn URL",
			source:           "",
			location:         "https://tunein.com/radio/WDR-2-Rheinland-1004-s213886/",
			expectedSource:   "TUNEIN",
			expectedLocation: "/v1/playback/station/s213886",
		},
		{
			name:             "TuneIn URL with source",
			source:           "SOMETHING",
			location:         "https://tunein.com/radio/WDR-2-Rheinland-1004-s213886/",
			expectedSource:   "TUNEIN",
			expectedLocation: "/v1/playback/station/s213886",
		},
		{
			name:             "TuneIn URL without trailing slash",
			source:           "",
			location:         "https://tunein.com/radio/WDR-2-Rheinland-1004-s213886",
			expectedSource:   "TUNEIN",
			expectedLocation: "/v1/playback/station/s213886",
		},
		{
			name:             "Non-TuneIn URL",
			source:           "OTHER",
			location:         "https://example.com/radio/s123",
			expectedSource:   "OTHER",
			expectedLocation: "https://example.com/radio/s123",
		},
		{
			name:             "TuneIn URL short form",
			source:           "",
			location:         "https://tunein.com/radio/s213886/",
			expectedSource:   "TUNEIN",
			expectedLocation: "/v1/playback/station/s213886",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSource, gotLocation := resolveLocation(tt.source, tt.location)
			if gotSource != tt.expectedSource {
				t.Errorf("resolveLocation() gotSource = %v, want %v", gotSource, tt.expectedSource)
			}
			if gotLocation != tt.expectedLocation {
				t.Errorf("resolveLocation() gotLocation = %v, want %v", gotLocation, tt.expectedLocation)
			}
		})
	}
}
