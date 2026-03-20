package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Metadata fetches cover images from the Steam Store API.
type Metadata struct {
	client *http.Client
}

// NewMetadata creates a new metadata service.
func NewMetadata() *Metadata {
	return &Metadata{client: &http.Client{}}
}

// FetchCoverURL returns a cover image URL.
// For Steam games, it uses the AppID directly. Otherwise searches the API.
func (m *Metadata) FetchCoverURL(gameName, appID, platform string) string {
	// Steam games: build URL directly from AppID
	if platform == "Steam" && appID != "" {
		return fmt.Sprintf("https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/%s/library_600x900_2x.jpg", appID)
	}

	// Other platforms: search Steam for cover art
	return m.searchSteamCover(gameName)
}

func (m *Metadata) searchSteamCover(gameName string) string {
	apiURL := fmt.Sprintf("https://store.steampowered.com/api/storesearch/?term=%s&l=english&cc=US",
		url.QueryEscape(gameName))

	resp, err := m.client.Get(apiURL)
	if err != nil || resp.StatusCode != 200 {
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	var result struct {
		Total int `json:"total"`
		Items []struct {
			ID int `json:"id"`
		} `json:"items"`
	}
	if json.Unmarshal(body, &result) != nil || result.Total == 0 || len(result.Items) == 0 {
		return ""
	}

	return fmt.Sprintf("https://shared.akamai.steamstatic.com/store_item_assets/steam/apps/%d/library_600x900_2x.jpg",
		result.Items[0].ID)
}
