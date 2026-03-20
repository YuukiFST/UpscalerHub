package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"UpscalerHub/internal/models"
)

// GameScanner orchestrates all platform scanners.
type GameScanner struct {
	exclusions []models.ScanExclusion
}

// NewGameScanner creates a scanner with exclusions loaded from the config path.
func NewGameScanner(configPath string) *GameScanner {
	return &GameScanner{exclusions: loadExclusions(configPath)}
}

// ScanAll runs all platform scanners and returns the combined, filtered list.
func (s *GameScanner) ScanAll() []models.Game {
	type scanFunc func() []models.Game
	scanners := []scanFunc{ScanSteam, ScanEpic, ScanGOG, ScanXbox, ScanEA, ScanBattleNet, ScanUbisoft}

	var all []models.Game
	analyzer := &GameAnalyzer{}

	for _, fn := range scanners {
		func() {
			defer func() { recover() }() // don't let one scanner crash everything
			for _, g := range fn() {
				if !s.isExcluded(&g) {
					analyzer.Analyze(&g)
					all = append(all, g)
				}
			}
		}()
	}

	sort.Slice(all, func(i, j int) bool {
		if all[i].Platform != all[j].Platform {
			return all[i].Platform < all[j].Platform
		}
		return strings.ToLower(all[i].Name) < strings.ToLower(all[j].Name)
	})
	return all
}

func (s *GameScanner) isExcluded(g *models.Game) bool {
	for _, rule := range s.exclusions {
		if rule.Name != "" && strings.EqualFold(g.Name, rule.Name) {
			return true
		}
		if rule.PathSegment != "" && strings.Contains(strings.ToLower(g.InstallPath), strings.ToLower(rule.PathSegment)) {
			return true
		}
	}
	return false
}

func loadExclusions(configPath string) []models.ScanExclusion {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return defaultExclusions()
	}
	var cfg struct {
		ScanExclusions []models.ScanExclusion `json:"ScanExclusions"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil || len(cfg.ScanExclusions) == 0 {
		return defaultExclusions()
	}
	return cfg.ScanExclusions
}

func defaultExclusions() []models.ScanExclusion {
	return []models.ScanExclusion{
		{Name: "Wallpaper Engine", PathSegment: "wallpaper_engine"},
		{Name: "Steamworks Common Redistributables", PathSegment: "Steamworks Shared"},
	}
}

// jsonUnmarshal is a package-level helper used by scanners.
func jsonUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// ConfigPath returns the path to config.json next to the executable.
func ConfigPath() string {
	exe, _ := os.Executable()
	return filepath.Join(filepath.Dir(exe), "config.json")
}
