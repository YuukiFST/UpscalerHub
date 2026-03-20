package models

// RepositoryConfig holds GitHub repository coordinates.
type RepositoryConfig struct {
	RepoOwner string `json:"RepoOwner"`
	RepoName  string `json:"RepoName"`
}

// AppConfiguration is the root config loaded from config.json.
type AppConfiguration struct {
	App        RepositoryConfig `json:"App"`
	OptiScaler RepositoryConfig `json:"OptiScaler"`
	Fakenvapi  RepositoryConfig `json:"Fakenvapi"`
	NukemFG    RepositoryConfig `json:"NukemFG"`
	Language   string           `json:"Language,omitempty"`

	ScanExclusions []ScanExclusion `json:"ScanExclusions,omitempty"`
}

// DefaultAppConfiguration returns a sensible default config.
func DefaultAppConfiguration() AppConfiguration {
	return AppConfiguration{
		App:        RepositoryConfig{RepoOwner: "Agustinm28", RepoName: "UpscalerHub"},
		OptiScaler: RepositoryConfig{RepoOwner: "optiscaler", RepoName: "OptiScaler"},
		Fakenvapi:  RepositoryConfig{RepoOwner: "optiscaler", RepoName: "fakenvapi"},
		NukemFG:    RepositoryConfig{RepoOwner: "Nukem9", RepoName: "dlssg-to-fsr3"},
		Language:   "en",
		ScanExclusions: []ScanExclusion{
			{Name: "Wallpaper Engine", PathSegment: "wallpaper_engine"},
			{Name: "Steamworks Common Redistributables", PathSegment: "Steamworks Shared"},
			{Name: "GameSave", PathSegment: "GameSave"},
		},
	}
}

// ComponentVersions tracks locally-cached version tags.
type ComponentVersions struct {
	OptiScalerVersion string `json:"OptiScalerVersion,omitempty"`
	FakenvapiVersion  string `json:"FakenvapiVersion,omitempty"`
	NukemFGVersion    string `json:"NukemFGVersion,omitempty"`
}
