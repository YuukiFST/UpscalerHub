package services

import (
	"encoding/json"
	"os"
	"path/filepath"

	"UpscalerHub/internal/models"
)

// Persistence handles saving and loading the game list from disk.
type Persistence struct {
	filePath string
}

// NewPersistence creates a persistence service that stores games in %APPDATA%/UpscalerHub/.
func NewPersistence() *Persistence {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		appData = "."
	}
	dir := filepath.Join(appData, "UpscalerHub")
	os.MkdirAll(dir, 0755)
	return &Persistence{filePath: filepath.Join(dir, "games.json")}
}

// Save writes games to disk as JSON.
func (p *Persistence) Save(games []models.Game) error {
	data, err := json.MarshalIndent(games, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p.filePath, data, 0644)
}

// Load reads the game list from disk.
func (p *Persistence) Load() []models.Game {
	data, err := os.ReadFile(p.filePath)
	if err != nil {
		return nil
	}
	var games []models.Game
	if json.Unmarshal(data, &games) != nil {
		return nil
	}
	return games
}
