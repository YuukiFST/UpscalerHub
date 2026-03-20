package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"UpscalerHub/internal/archive"
	"UpscalerHub/internal/models"
)

// ComponentManager handles OptiScaler, Fakenvapi, and NukemFG lifecycle.
type ComponentManager struct {
	baseDir     string
	cacheDir    string
	versionFile string
	configFile  string
	client      *http.Client

	Config         models.AppConfiguration
	LocalVersions  models.ComponentVersions
	RemoteVersions models.ComponentVersions

	cachedOptiVersions []string
	lastAPICheck       time.Time
}

// NewComponentManager creates the service and loads config + local versions.
func NewComponentManager() *ComponentManager {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		appData = "."
	}
	base := filepath.Join(appData, "UpscalerHub")
	os.MkdirAll(filepath.Join(base, "Cache"), 0755)

	cm := &ComponentManager{
		baseDir:     base,
		cacheDir:    filepath.Join(base, "Cache"),
		versionFile: filepath.Join(base, "versions.json"),
		configFile:  filepath.Join(base, "config.json"),
		client:      &http.Client{Timeout: 30 * time.Second},
	}
	cm.loadConfig()
	cm.loadVersions()
	return cm
}

func (cm *ComponentManager) loadConfig() {
	// Try local config first (next to exe)
	exe, _ := os.Executable()
	localCfg := filepath.Join(filepath.Dir(exe), "config.json")
	if data, err := os.ReadFile(localCfg); err == nil {
		if json.Unmarshal(data, &cm.Config) == nil {
			return
		}
	}
	// Try AppData config
	if data, err := os.ReadFile(cm.configFile); err == nil {
		json.Unmarshal(data, &cm.Config)
		return
	}
	// Use defaults
	cm.Config = models.DefaultAppConfiguration()
	cm.SaveConfig()
}

// SaveConfig persists the current configuration.
func (cm *ComponentManager) SaveConfig() {
	data, _ := json.MarshalIndent(cm.Config, "", "  ")
	// Save next to exe if that file exists, otherwise to AppData
	exe, _ := os.Executable()
	localCfg := filepath.Join(filepath.Dir(exe), "config.json")
	if _, err := os.Stat(localCfg); err == nil {
		os.WriteFile(localCfg, data, 0644)
	} else {
		os.WriteFile(cm.configFile, data, 0644)
	}
}

func (cm *ComponentManager) loadVersions() {
	data, err := os.ReadFile(cm.versionFile)
	if err != nil {
		return
	}
	json.Unmarshal(data, &cm.LocalVersions)
}

func (cm *ComponentManager) saveVersions() {
	data, _ := json.MarshalIndent(cm.LocalVersions, "", "  ")
	os.WriteFile(cm.versionFile, data, 0644)
}

// CheckForUpdates fetches latest version info from GitHub for all components.
func (cm *ComponentManager) CheckForUpdates() error {
	if len(cm.cachedOptiVersions) > 0 && time.Since(cm.lastAPICheck) < 15*time.Minute {
		return nil
	}

	// Fetch OptiScaler versions
	versions, err := cm.fetchAllVersions(cm.Config.OptiScaler)
	if err == nil && len(versions) > 0 {
		cm.cachedOptiVersions = versions
	}

	// Fetch Fakenvapi latest
	if v := cm.fetchLatestTag(cm.Config.Fakenvapi); v != "" {
		cm.RemoteVersions.FakenvapiVersion = v
	}
	// Fetch NukemFG latest
	if v := cm.fetchLatestTag(cm.Config.NukemFG); v != "" {
		cm.RemoteVersions.NukemFGVersion = v
	}

	// Set OptiScaler remote version (prefer stable, non-nightly)
	if len(cm.cachedOptiVersions) > 0 {
		for _, v := range cm.cachedOptiVersions {
			if !strings.Contains(strings.ToLower(v), "nightly") {
				cm.RemoteVersions.OptiScalerVersion = v
				break
			}
		}
		if cm.RemoteVersions.OptiScalerVersion == "" {
			cm.RemoteVersions.OptiScalerVersion = cm.cachedOptiVersions[0]
		}
	}

	cm.lastAPICheck = time.Now()
	return nil
}

// OptiScalerVersions returns available versions (from API + locally cached).
func (cm *ComponentManager) OptiScalerVersions() []string {
	if len(cm.cachedOptiVersions) > 0 {
		return cm.cachedOptiVersions
	}
	return cm.GetDownloadedVersions()
}

// IsOptiScalerUpdateAvailable checks if a newer OptiScaler version exists.
func (cm *ComponentManager) IsOptiScalerUpdateAvailable() bool {
	return cm.RemoteVersions.OptiScalerVersion != "" &&
		cm.LocalVersions.OptiScalerVersion != cm.RemoteVersions.OptiScalerVersion
}

func (cm *ComponentManager) fetchLatestTag(repo models.RepositoryConfig) string {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repo.RepoOwner, repo.RepoName)
	data, err := cm.apiGet(url)
	if err != nil {
		return ""
	}
	var release struct {
		TagName string `json:"tag_name"`
	}
	if json.Unmarshal(data, &release) != nil {
		return ""
	}
	return strings.TrimPrefix(release.TagName, "v")
}

func (cm *ComponentManager) fetchAllVersions(repo models.RepositoryConfig) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases?per_page=30", repo.RepoOwner, repo.RepoName)
	data, err := cm.apiGet(url)
	if err != nil {
		return nil, err
	}
	var releases []struct {
		TagName string `json:"tag_name"`
	}
	if json.Unmarshal(data, &releases) != nil {
		return nil, fmt.Errorf("failed to parse releases")
	}
	var versions []string
	for _, r := range releases {
		versions = append(versions, strings.TrimPrefix(r.TagName, "v"))
	}
	return versions, nil
}

func (cm *ComponentManager) apiGet(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "UpscalerHub")
	resp, err := cm.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// DownloadOptiScaler downloads and extracts a specific OptiScaler version.
// Returns the path to the extracted cache directory.
func (cm *ComponentManager) DownloadOptiScaler(version string, progressCb func(float64)) (string, error) {
	extractPath := cm.OptiScalerCachePath(version)
	if entries, _ := os.ReadDir(extractPath); len(entries) > 0 {
		return extractPath, nil // already cached
	}

	repo := cm.Config.OptiScaler
	// Try with "v" prefix first, then without
	var downloadURL string
	for _, tag := range []string{"v" + version, version} {
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", repo.RepoOwner, repo.RepoName, tag)
		data, err := cm.apiGet(url)
		if err != nil {
			continue
		}
		downloadURL = findAssetURL(data)
		if downloadURL != "" {
			break
		}
	}
	if downloadURL == "" {
		return "", fmt.Errorf("no downloadable asset found for OptiScaler %s", version)
	}

	os.MkdirAll(extractPath, 0755)

	tmpFile, err := cm.downloadFile(downloadURL, progressCb)
	if err != nil {
		os.RemoveAll(extractPath)
		return "", err
	}
	defer os.Remove(tmpFile)

	if err := archive.Extract(tmpFile, extractPath); err != nil {
		os.RemoveAll(extractPath)
		return "", fmt.Errorf("extraction failed: %w", err)
	}

	cm.LocalVersions.OptiScalerVersion = version
	cm.saveVersions()
	return extractPath, nil
}

// DownloadFakenvapi downloads the latest Fakenvapi release.
func (cm *ComponentManager) DownloadFakenvapi() error {
	if cm.RemoteVersions.FakenvapiVersion == "" {
		return fmt.Errorf("no remote version available")
	}
	repo := cm.Config.Fakenvapi
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repo.RepoOwner, repo.RepoName)
	data, err := cm.apiGet(url)
	if err != nil {
		return err
	}
	dlURL := findAssetURL(data)
	if dlURL == "" {
		return fmt.Errorf("no downloadable asset found")
	}

	extractPath := cm.FakenvapiCachePath()
	os.RemoveAll(extractPath)
	os.MkdirAll(extractPath, 0755)

	tmpFile, err := cm.downloadFile(dlURL, nil)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	if err := archive.Extract(tmpFile, extractPath); err != nil {
		return err
	}

	cm.LocalVersions.FakenvapiVersion = cm.RemoteVersions.FakenvapiVersion
	cm.saveVersions()
	return nil
}

func (cm *ComponentManager) downloadFile(url string, progressCb func(float64)) (string, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "UpscalerHub")
	resp, err := cm.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpFile, err := os.CreateTemp("", "UpscalerHub-*.zip")
	if err != nil {
		return "", err
	}

	totalBytes := resp.ContentLength
	if totalBytes <= 0 {
		totalBytes = 10 * 1024 * 1024
	}

	buf := make([]byte, 8192)
	var written int64
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			tmpFile.Write(buf[:n])
			written += int64(n)
			if progressCb != nil {
				progressCb(float64(written) / float64(totalBytes) * 100)
			}
		}
		if readErr != nil {
			break
		}
	}
	tmpFile.Close()
	return tmpFile.Name(), nil
}

func findAssetURL(releaseJSON []byte) string {
	var release struct {
		Assets []struct {
			BrowserDownloadURL string `json:"browser_download_url"`
			Name               string `json:"name"`
		} `json:"assets"`
	}
	if json.Unmarshal(releaseJSON, &release) != nil {
		return ""
	}
	for _, a := range release.Assets {
		lower := strings.ToLower(a.Name)
		if strings.HasSuffix(lower, ".zip") || strings.HasSuffix(lower, ".7z") {
			return a.BrowserDownloadURL
		}
	}
	if len(release.Assets) > 0 {
		return release.Assets[0].BrowserDownloadURL
	}
	return ""
}

// OptiScalerCachePath returns the cache directory for a given version.
func (cm *ComponentManager) OptiScalerCachePath(version string) string {
	return filepath.Join(cm.cacheDir, "OptiScaler", version)
}

// FakenvapiCachePath returns the Fakenvapi cache directory.
func (cm *ComponentManager) FakenvapiCachePath() string {
	return filepath.Join(cm.cacheDir, "Fakenvapi")
}

// NukemFGCachePath returns the NukemFG cache directory.
func (cm *ComponentManager) NukemFGCachePath() string {
	return filepath.Join(cm.cacheDir, "NukemFG")
}

// GetDownloadedVersions lists locally cached OptiScaler versions.
func (cm *ComponentManager) GetDownloadedVersions() []string {
	dir := filepath.Join(cm.cacheDir, "OptiScaler")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	skip := map[string]bool{"d3d12_optiscaler": true, "dlssoverrides": true, "licenses": true}
	var versions []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if skip[strings.ToLower(e.Name())] {
			continue
		}
		versions = append(versions, e.Name())
	}
	sort.Slice(versions, func(i, j int) bool {
		return versions[i] > versions[j]
	})
	return versions
}

// DeleteOptiScalerCache removes a specific cached version.
func (cm *ComponentManager) DeleteOptiScalerCache(version string) error {
	path := cm.OptiScalerCachePath(version)
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	if cm.LocalVersions.OptiScalerVersion == version {
		versions := cm.GetDownloadedVersions()
		if len(versions) > 0 {
			cm.LocalVersions.OptiScalerVersion = versions[0]
		} else {
			cm.LocalVersions.OptiScalerVersion = ""
		}
		cm.saveVersions()
	}
	return nil
}
