package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"UpscalerHub/internal/models"
	"UpscalerHub/internal/nvapi"
	"UpscalerHub/internal/services"
)

// App is the main application struct bound to the Wails frontend.
type App struct {
	ctx        context.Context
	mu         sync.Mutex
	games      []models.Game
	persist    *services.Persistence
	components *services.ComponentManager
	metadata   *services.Metadata
	installer  *services.Installer
	updater    *services.AppUpdater
	swapper    *services.SwapperService
}

// NewApp creates the App instance.
func NewApp() *App {
	cm := services.NewComponentManager()
	return &App{
		persist:    services.NewPersistence(),
		components: cm,
		metadata:   services.NewMetadata(),
		installer:  &services.Installer{},
		updater:    services.NewAppUpdater(cm),
		swapper:    services.NewSwapperService(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.games = a.persist.Load()
}

// ── Game list ───────────────────────────────────────────────────────────

// GetGames returns the current game list.
func (a *App) GetGames() []models.Game {
	return a.games
}

// ScanGames triggers a full scan and returns the updated game list.
// Manually added games are preserved.
func (a *App) ScanGames() []models.Game {
	scanner := services.NewGameScanner(services.ConfigPath())
	scanned := scanner.ScanAll()

	// Preserve manually added games from the current list
	a.mu.Lock()
	var manual []models.Game
	for _, g := range a.games {
		if g.Platform == models.PlatformManual {
			manual = append(manual, g)
		}
	}
	a.mu.Unlock()

	// Merge: scanned games first, then manual games (avoiding duplicates)
	merged := make([]models.Game, 0, len(scanned)+len(manual))
	merged = append(merged, scanned...)
	for _, m := range manual {
		dupe := false
		for _, s := range scanned {
			if strings.EqualFold(s.InstallPath, m.InstallPath) {
				dupe = true
				break
			}
		}
		if !dupe {
			merged = append(merged, m)
		}
	}

	a.mu.Lock()
	// Restore preserved IsFavorite and CoverImageURL
	for i := range merged {
		for _, cg := range a.games {
			if strings.EqualFold(merged[i].InstallPath, cg.InstallPath) {
				merged[i].IsFavorite = cg.IsFavorite
				if cg.CoverImageURL != "" {
					merged[i].CoverImageURL = cg.CoverImageURL
				}
				break
			}
		}
	}

	// Set cover images immediately for Steam games (no API call needed)
	for i := range merged {
		if merged[i].CoverImageURL == "" {
			if merged[i].Platform == models.PlatformSteam && merged[i].AppID != "" {
				merged[i].CoverImageURL = a.metadata.FetchCoverURL(merged[i].Name, merged[i].AppID, string(merged[i].Platform))
			}
		}
	}

	a.games = merged
	a.persist.Save(a.games)
	a.mu.Unlock()

	// Fetch cover images for non-Steam games in background
	go func() {
		changed := false
		for i := range merged {
			if merged[i].CoverImageURL == "" {
				url := a.metadata.FetchCoverURL(merged[i].Name, merged[i].AppID, string(merged[i].Platform))
				if url != "" {
					merged[i].CoverImageURL = url
					changed = true
				}
			}
		}
		if changed {
			a.mu.Lock()
			a.games = merged
			a.persist.Save(a.games)
			a.mu.Unlock()
			runtime.EventsEmit(a.ctx, "games-updated", merged)
		}
	}()

	return merged
}

// AddGameManually opens a file dialog and adds a game from an executable path.
func (a *App) AddGameManually() (*models.Game, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Game Executable",
		Filters: []runtime.FileFilter{
			{DisplayName: "Executables", Pattern: "*.exe"},
		},
	})
	if err != nil || path == "" {
		return nil, err
	}

	// Check for duplicates
	for _, g := range a.games {
		if strings.EqualFold(g.InstallPath, filepath.Dir(path)) {
			return nil, nil
		}
	}

	game := models.Game{
		Name:           strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		InstallPath:    filepath.Dir(path),
		ExecutablePath: path,
		Platform:       models.PlatformManual,
	}

	analyzer := &services.GameAnalyzer{}
	analyzer.Analyze(&game)

	game.CoverImageURL = a.metadata.FetchCoverURL(game.Name, game.AppID, string(game.Platform))

	a.mu.Lock()
	a.games = append(a.games, game)
	a.persist.Save(a.games)
	a.mu.Unlock()

	return &game, nil
}

// RemoveGame removes a game by its install path (more reliable than index).
func (a *App) RemoveGame(installPath string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i, g := range a.games {
		if strings.EqualFold(g.InstallPath, installPath) {
			a.games = append(a.games[:i], a.games[i+1:]...)
			a.persist.Save(a.games)
			return
		}
	}
}

// GetGame returns a single game by index.
func (a *App) GetGame(index int) *models.Game {
	if index < 0 || index >= len(a.games) {
		return nil
	}
	return &a.games[index]
}

// ToggleFavorite toggles the favorite status of a game.
func (a *App) ToggleFavorite(index int) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.games) {
		return false
	}
	a.games[index].IsFavorite = !a.games[index].IsFavorite
	a.persist.Save(a.games)
	return a.games[index].IsFavorite
}

// LaunchGame launches the game's executable or via platform launcher.
func (a *App) LaunchGame(index int) error {
	a.mu.Lock()
	if index < 0 || index >= len(a.games) {
		a.mu.Unlock()
		return nil
	}
	g := a.games[index]
	a.mu.Unlock()

	var cmd *exec.Cmd
	switch g.Platform {
	case models.PlatformSteam:
		if g.AppID != "" {
			cmd = exec.Command("cmd", "/c", "start", "steam://rungameid/"+g.AppID)
		} else {
			cmd = exec.Command(g.ExecutablePath)
		}
	case models.PlatformEpic:
		if g.AppID != "" {
			cmd = exec.Command("cmd", "/c", "start", "com.epicgames.launcher://apps/"+g.AppID+"?action=launch&silent=true")
		} else {
			cmd = exec.Command(g.ExecutablePath)
		}
	default:
		cmd = exec.Command(g.ExecutablePath)
	}

	cmd.Dir = filepath.Dir(g.ExecutablePath)
	return cmd.Start()
}

// GetNVAPIPresets retrieves the NVAPI presets for DLSS/DLSS-RR.
func (a *App) GetNVAPIPresets(exePath string) (*nvapi.PresetResult, error) {
	exeName := filepath.Base(exePath)
	return nvapi.GetPresets(exeName)
}

// SetNVAPIPreset updates the NVAPI preset for a given tech.
func (a *App) SetNVAPIPreset(exePath, tech string, preset uint) error {
	exeName := filepath.Base(exePath)
	return nvapi.SetPreset(exeName, tech, preset)
}

// ── GPU ─────────────────────────────────────────────────────────────────

// GpuResult holds GPU info for the frontend.
type GpuResult struct {
	Display string            `json:"display"`
	Info    *services.GpuInfo `json:"info"`
}

// GetGPU detects the system GPU.
func (a *App) GetGPU() GpuResult {
	info := services.DetectGPU()
	return GpuResult{
		Display: services.FormatGPU(info),
		Info:    info,
	}
}

// ── Components ──────────────────────────────────────────────────────────

// ComponentStatus holds the current version state for the frontend.
type ComponentStatus struct {
	Local              models.ComponentVersions `json:"local"`
	Remote             models.ComponentVersions `json:"remote"`
	OptiScalerVersions []string                 `json:"optiScalerVersions"`
	UpdateAvailable    bool                     `json:"updateAvailable"`
}

// GetComponentStatus returns current component version info.
func (a *App) GetComponentStatus() ComponentStatus {
	a.components.CheckForUpdates()
	return ComponentStatus{
		Local:              a.components.LocalVersions,
		Remote:             a.components.RemoteVersions,
		OptiScalerVersions: a.components.OptiScalerVersions(),
		UpdateAvailable:    a.components.IsOptiScalerUpdateAvailable(),
	}
}

// GetDownloadedVersions returns locally cached OptiScaler versions.
func (a *App) GetDownloadedVersions() []string {
	return a.components.GetDownloadedVersions()
}

// DeleteCachedVersion removes a cached OptiScaler version.
func (a *App) DeleteCachedVersion(version string) error {
	return a.components.DeleteOptiScalerCache(version)
}

// ── Installation ────────────────────────────────────────────────────────

// InstallResult holds the result of an install/uninstall operation.
type InstallResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// InstallOptiScaler installs OptiScaler for a game.
func (a *App) InstallOptiScaler(gameIndex int, version, injection string, fakenvapi, nukem bool) InstallResult {
	if gameIndex < 0 || gameIndex >= len(a.games) {
		return InstallResult{false, "Invalid game index"}
	}

	// Download if not cached
	cachePath, err := a.components.DownloadOptiScaler(version, func(pct float64) {
		runtime.EventsEmit(a.ctx, "download-progress", pct)
	})
	if err != nil {
		return InstallResult{false, "Download failed: " + err.Error()}
	}

	var fakenvapiPath, nukemPath string
	if fakenvapi {
		fakenvapiPath = a.components.FakenvapiCachePath()
		if entries, _ := os.ReadDir(fakenvapiPath); len(entries) == 0 {
			if err := a.components.DownloadFakenvapi(); err != nil {
				return InstallResult{false, "Fakenvapi download failed: " + err.Error()}
			}
		}
	}
	if nukem {
		nukemPath = a.components.NukemFGCachePath()
	}

	err = a.installer.Install(&a.games[gameIndex], cachePath, injection, version, fakenvapi, nukem, fakenvapiPath, nukemPath)
	if err != nil {
		return InstallResult{false, err.Error()}
	}

	a.persist.Save(a.games)
	return InstallResult{true, version + " installed successfully!"}
}

// UninstallOptiScaler uninstalls OptiScaler from a game.
func (a *App) UninstallOptiScaler(gameIndex int) InstallResult {
	if gameIndex < 0 || gameIndex >= len(a.games) {
		return InstallResult{false, "Invalid game index"}
	}

	err := a.installer.Uninstall(&a.games[gameIndex])
	if err != nil {
		return InstallResult{false, err.Error()}
	}

	a.persist.Save(a.games)
	return InstallResult{true, "OptiScaler uninstalled successfully"}
}

// ── Folders ─────────────────────────────────────────────────────────────

// OpenFolder opens a file explorer at the given path.
func (a *App) OpenFolder(path string) {
	if path == "" {
		return
	}
	exec.Command("explorer.exe", path).Start()
}

// ── Settings ────────────────────────────────────────────────────────────

// GetConfig returns the current app config.
func (a *App) GetConfig() models.AppConfiguration {
	return a.components.Config
}

// SaveConfig persists config changes.
func (a *App) SaveConfig(cfg models.AppConfiguration) {
	a.components.Config = cfg
	a.components.SaveConfig()
}

// ── App Update ──────────────────────────────────────────────────────────

// AppUpdateInfo holds update check result.
type AppUpdateInfo struct {
	Available bool   `json:"available"`
	Version   string `json:"version"`
	Notes     string `json:"notes"`
	URL       string `json:"url"`
}

// CheckAppUpdate checks for a newer version of UpscalerHub.
func (a *App) CheckAppUpdate() AppUpdateInfo {
	avail, _ := a.updater.CheckForUpdate("0.1.0")
	return AppUpdateInfo{
		Available: avail,
		Version:   a.updater.Version,
		Notes:     a.updater.Notes,
		URL:       a.updater.URL,
	}
}

// ── Swapper (Native DLLs) ───────────────────────────────────────────────

// FetchSwapperManifest fetches the latest DLSS Swapper manifest.
func (a *App) FetchSwapperManifest() *models.SwapperManifest {
	m, err := a.swapper.FetchManifest()
	if err != nil {
		return nil
	}
	return m
}

// DownloadOfficialDLL downloads and extracts a specific official DLL.
func (a *App) DownloadOfficialDLL(record models.DLLRecord) error {
	return a.swapper.DownloadDLL(&record)
}

// SwapOfficialDLL backups the original and copies the cached official DLL into the game folder.
func (a *App) SwapOfficialDLL(record models.DLLRecord, targetDir string) error {
	return a.swapper.SwapDLL(&record, targetDir)
}

// RestoreOfficialDLL restores the native backup DLL.
func (a *App) RestoreOfficialDLL(targetDir, upscalerType string) error {
	return a.swapper.RestoreOriginalDLL(targetDir, upscalerType)
}
