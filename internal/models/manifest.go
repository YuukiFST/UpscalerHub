package models

// InstallationManifest tracks files installed by OptiScaler for a game,
// enabling clean uninstallation.
type InstallationManifest struct {
	ManifestVersion      int      `json:"ManifestVersion"`
	OptiScalerVersion    string   `json:"OptiscalerVersion,omitempty"`
	InjectionMethod      string   `json:"InjectionMethod"`
	InstallDate          string   `json:"InstallDate"`
	InstalledGameDir     string   `json:"InstalledGameDirectory,omitempty"`
	InstalledFiles       []string `json:"InstalledFiles"`
	BackedUpFiles        []string `json:"BackedUpFiles"`
	InstalledDirectories []string `json:"InstalledDirectories"`
}
