package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"UpscalerHub/internal/models"
)

// GameAnalyzer inspects a game directory for upscaling DLLs and OptiScaler state.
type GameAnalyzer struct{}

// Analyze populates the technology fields of a Game by scanning its directory tree.
func (a *GameAnalyzer) Analyze(game *models.Game) {
	if game.InstallPath == "" {
		return
	}

	// Walk the game directory (scan deep enough for UE/Unity games)
	_ = filepath.WalkDir(game.InstallPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip inaccessible
		}
		if d.IsDir() {
			rel, _ := filepath.Rel(game.InstallPath, path)
			if strings.Count(rel, string(os.PathSeparator)) > 8 {
				return filepath.SkipDir
			}
			lower := strings.ToLower(d.Name())
			// Skip directories that won't contain game DLLs
			switch lower {
			case "node_modules", ".git", "__pycache__", "saves", "logs",
				"screenshots", "crashreports", "shader_cache", "shadercache",
				"localization", "movies", "video", "audio", "music", "sounds":
				return filepath.SkipDir
			}
			return nil
		}

		name := strings.ToLower(d.Name())

		// Only process DLL files and specific config files
		if !strings.HasSuffix(name, ".dll") && !strings.HasSuffix(name, ".json") && !strings.HasSuffix(name, ".ini") && !strings.HasSuffix(name, ".exe") {
			return nil
		}

		// DLSS detection
		if game.DLSSVersion == "" {
			switch name {
			case "nvngx_dlss.dll", "nvngx_dlss.dll.dlss":
				game.DLSSVersion = getFileVersion(path)
				game.DLSSPath = path
			}
		}

		// DLSS Frame Generation
		if game.DLSSFrameGenVersion == "" {
			switch name {
			case "nvngx_dlssg.dll", "dlssg_to_fsr3.dll":
				game.DLSSFrameGenVersion = getFileVersion(path)
				game.DLSSFrameGenPath = path
			}
		}

		// FSR detection (many possible names)
		if game.FSRVersion == "" {
			switch name {
			case "amd_fidelityfx_vk.dll", "amd_fidelityfx_dx12.dll",
				"ffx_fsr2_api_x64.dll", "ffx_fsr3upscaler_x64.dll",
				"ffx_fsr2_api_dx12_x64.dll", "ffx_fsr2_api_vk_x64.dll",
				"fsr2steamvrdx11.dll", "ffx_fsr3_x64.dll",
				"amd_ags_x64.dll":
				game.FSRVersion = getFileVersion(path)
				game.FSRPath = path
			default:
				if strings.HasPrefix(name, "ffx_fsr") && strings.HasSuffix(name, ".dll") {
					game.FSRVersion = getFileVersion(path)
					game.FSRPath = path
				}
			}
		}

		// XeSS detection
		if game.XeSSVersion == "" {
			if name == "libxess.dll" || name == "libxess.so" {
				game.XeSSVersion = getFileVersion(path)
				game.XeSSPath = path
			}
		}

		// Detect OptiScaler installation via manifest
		if name == "optiscaler_manifest.json" {
			game.IsOptiScalerInstalled = true
			data, err := os.ReadFile(path)
			if err == nil {
				var manifest models.InstallationManifest
				if jsonUnmarshal(data, &manifest) == nil && manifest.OptiScalerVersion != "" {
					game.OptiScalerVersion = manifest.OptiScalerVersion
				}
			}
		}

		// Also detect OptiScaler by presence of OptiScaler.ini
		if name == "optiscaler.ini" && !game.IsOptiScalerInstalled {
			game.IsOptiScalerInstalled = true
		}

		// Detect executable
		if game.ExecutablePath == "" && strings.HasSuffix(name, ".exe") {
			dir := filepath.Dir(path)
			if dir == game.InstallPath || strings.Contains(strings.ToLower(dir), "bin") {
				game.ExecutablePath = path
			}
		}

		return nil
	})
}

func getFileVersion(path string) string {
	ver, err := readPEVersion(path)
	if err == nil && ver != "" {
		return ver
	}
	// Fallback to file size
	fi, err := os.Stat(path)
	if err != nil {
		return "detected"
	}
	mb := float64(fi.Size()) / (1024 * 1024)
	return fmt.Sprintf("detected (%.1f MB)", mb)
}
