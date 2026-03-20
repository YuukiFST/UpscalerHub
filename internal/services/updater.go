package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// AppUpdater checks for application updates on GitHub.
type AppUpdater struct {
	client  *http.Client
	config  *ComponentManager
	Version string
	Notes   string
	URL     string
}

// NewAppUpdater creates an update checker.
func NewAppUpdater(cm *ComponentManager) *AppUpdater {
	return &AppUpdater{
		client: &http.Client{Timeout: 15 * time.Second},
		config: cm,
	}
}

// CheckForUpdate returns true if a newer version is available.
func (u *AppUpdater) CheckForUpdate(currentVersion string) (bool, error) {
	repo := u.config.Config.App
	if repo.RepoOwner == "" || repo.RepoName == "" {
		return false, nil
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repo.RepoOwner, repo.RepoName)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "UpscalerHub")

	resp, err := u.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return false, nil
	}

	var release struct {
		TagName string `json:"tag_name"`
		Body    string `json:"body"`
		Assets  []struct {
			BrowserDownloadURL string `json:"browser_download_url"`
			Name               string `json:"name"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return false, err
	}

	u.Version = strings.TrimPrefix(release.TagName, "v")
	u.Notes = release.Body
	for _, a := range release.Assets {
		if strings.HasSuffix(strings.ToLower(a.Name), ".zip") {
			u.URL = a.BrowserDownloadURL
			break
		}
	}

	if currentVersion == "" || u.Version == "" {
		return false, nil
	}
	return u.Version != currentVersion && u.Version > currentVersion, nil
}
