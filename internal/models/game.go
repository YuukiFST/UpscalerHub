package models

// GamePlatform identifies where a game was installed from.
type GamePlatform string

const (
	PlatformSteam    GamePlatform = "Steam"
	PlatformEpic     GamePlatform = "Epic"
	PlatformGOG      GamePlatform = "GOG"
	PlatformXbox     GamePlatform = "Xbox"
	PlatformEA       GamePlatform = "EA"
	PlatformBattleNet GamePlatform = "BattleNet"
	PlatformUbisoft  GamePlatform = "Ubisoft"
	PlatformManual   GamePlatform = "Manual"
)

// Game represents a detected or manually added game.
type Game struct {
	Name           string       `json:"name"`
	InstallPath    string       `json:"installPath"`
	Platform       GamePlatform `json:"platform"`
	AppID          string       `json:"appId"`
	ExecutablePath string       `json:"executablePath"`
	CoverImageURL  string       `json:"coverImageUrl,omitempty"`
	IsFavorite     bool         `json:"isFavorite"`

	// Detected upscaling technologies
	DLSSVersion        string `json:"dlssVersion,omitempty"`
	DLSSPath           string `json:"dlssPath,omitempty"`
	DLSSFrameGenVersion string `json:"dlssFrameGenVersion,omitempty"`
	DLSSFrameGenPath   string `json:"dlssFrameGenPath,omitempty"`
	FSRVersion         string `json:"fsrVersion,omitempty"`
	FSRPath            string `json:"fsrPath,omitempty"`
	XeSSVersion        string `json:"xessVersion,omitempty"`
	XeSSPath           string `json:"xessPath,omitempty"`

	// OptiScaler installation state
	IsOptiScalerInstalled bool   `json:"isOptiScalerInstalled"`
	OptiScalerVersion     string `json:"optiScalerVersion,omitempty"`
}

// IsManual returns true if the game was manually added.
func (g *Game) IsManual() bool {
	return g.Platform == PlatformManual
}
