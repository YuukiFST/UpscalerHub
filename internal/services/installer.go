package services

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"UpscalerHub/internal/models"
)

// Installer handles installing and uninstalling OptiScaler for games.
type Installer struct{}

const (
	backupFolder = "OptiScalerBackup"
	manifestFile = "optiscaler_manifest.json"
)

// Install copies OptiScaler files into a game directory.
func (inst *Installer) Install(game *models.Game, cachePath, injectionMethod, optiVersion string, includeFakenvapi, includeNukem bool, fakenvapiPath, nukemPath string) error {
	gameDir := inst.determineInstallDir(game)
	if gameDir == "" {
		return fmt.Errorf("could not determine install directory")
	}

	backup := filepath.Join(gameDir, backupFolder)
	os.MkdirAll(backup, 0755)

	manifest := models.InstallationManifest{
		ManifestVersion:   1,
		OptiScalerVersion: optiVersion,
		InjectionMethod:   injectionMethod,
		InstallDate:       time.Now().Format(time.RFC3339),
		InstalledGameDir:  gameDir,
	}

	// Backup existing file if it exists, then copy injection DLL
	injDLL := injectionMethod // e.g. "dxgi.dll", "version.dll", "winmm.dll"
	existingDLL := filepath.Join(gameDir, injDLL)
	if fileExists(existingDLL) {
		backupDst := filepath.Join(backup, injDLL)
		copyFile(existingDLL, backupDst)
		manifest.BackedUpFiles = append(manifest.BackedUpFiles, injDLL)
	}

	// Find the OptiScaler proxy DLL in cache
	srcDLL := findOptiScalerDLL(cachePath)
	if srcDLL == "" {
		return fmt.Errorf("OptiScaler DLL not found in cache")
	}
	copyFile(srcDLL, filepath.Join(gameDir, injDLL))
	manifest.InstalledFiles = append(manifest.InstalledFiles, injDLL)

	// Copy OptiScaler.ini if present
	srcINI := filepath.Join(filepath.Dir(srcDLL), "OptiScaler.ini")
	if fileExists(srcINI) {
		copyFile(srcINI, filepath.Join(gameDir, "OptiScaler.ini"))
		manifest.InstalledFiles = append(manifest.InstalledFiles, "OptiScaler.ini")
	}

	// Copy optional Fakenvapi
	if includeFakenvapi && fakenvapiPath != "" {
		if err := inst.installComponent(fakenvapiPath, gameDir, "nvapi64.dll", &manifest); err == nil {
			// Also look for additional fakenvapi files
			for _, name := range []string{"fakenvapi.ini"} {
				src := filepath.Join(fakenvapiPath, name)
				if fileExists(src) {
					copyFile(src, filepath.Join(gameDir, name))
					manifest.InstalledFiles = append(manifest.InstalledFiles, name)
				}
			}
		}
	}

	// Copy optional NukemFG DLL
	if includeNukem && nukemPath != "" {
		nukemDLL := filepath.Join(nukemPath, "dlssg_to_fsr3_amd_is_better.dll")
		if fileExists(nukemDLL) {
			copyFile(nukemDLL, filepath.Join(gameDir, "dlssg_to_fsr3_amd_is_better.dll"))
			manifest.InstalledFiles = append(manifest.InstalledFiles, "dlssg_to_fsr3_amd_is_better.dll")
		}
	}

	// Write manifest
	manifestData, _ := json.MarshalIndent(manifest, "", "  ")
	os.WriteFile(filepath.Join(gameDir, manifestFile), manifestData, 0644)
	manifest.InstalledFiles = append(manifest.InstalledFiles, manifestFile)

	// Update game state
	game.IsOptiScalerInstalled = true
	game.OptiScalerVersion = optiVersion

	return nil
}

// Uninstall removes OptiScaler from a game using the manifest.
func (inst *Installer) Uninstall(game *models.Game) error {
	gameDir := game.InstallPath
	if game.ExecutablePath != "" {
		gameDir = filepath.Dir(game.ExecutablePath)
	}

	// Try to find manifest in possible directories
	var manifest models.InstallationManifest
	var manifestPath string
	for _, dir := range []string{gameDir, game.InstallPath} {
		mp := filepath.Join(dir, manifestFile)
		if data, err := os.ReadFile(mp); err == nil {
			if json.Unmarshal(data, &manifest) == nil {
				manifestPath = mp
				if manifest.InstalledGameDir != "" {
					gameDir = manifest.InstalledGameDir
				}
				break
			}
		}
	}

	if manifestPath == "" {
		// Fallback: try to find manifest anywhere in game dir
		filepath.WalkDir(game.InstallPath, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			if strings.ToLower(d.Name()) == manifestFile {
				if data, err := os.ReadFile(path); err == nil {
					if json.Unmarshal(data, &manifest) == nil {
						manifestPath = path
						gameDir = filepath.Dir(path)
						if manifest.InstalledGameDir != "" {
							gameDir = manifest.InstalledGameDir
						}
						return filepath.SkipAll
					}
				}
			}
			return nil
		})
	}

	// Remove installed files
	for _, f := range manifest.InstalledFiles {
		os.Remove(filepath.Join(gameDir, f))
	}

	// Restore backed-up files
	backupDir := filepath.Join(gameDir, backupFolder)
	for _, f := range manifest.BackedUpFiles {
		src := filepath.Join(backupDir, f)
		dst := filepath.Join(gameDir, f)
		if fileExists(src) {
			copyFile(src, dst)
			os.Remove(src)
		}
	}

	// Clean up backup dir if empty
	if entries, _ := os.ReadDir(backupDir); len(entries) == 0 {
		os.Remove(backupDir)
	}

	// Remove empty installed directories
	for _, d := range manifest.InstalledDirectories {
		full := filepath.Join(gameDir, d)
		if entries, _ := os.ReadDir(full); len(entries) == 0 {
			os.Remove(full)
		}
	}

	// Remove manifest
	if manifestPath != "" {
		os.Remove(manifestPath)
	}

	game.IsOptiScalerInstalled = false
	game.OptiScalerVersion = ""
	return nil
}

func (inst *Installer) determineInstallDir(game *models.Game) string {
	if game.ExecutablePath != "" {
		return filepath.Dir(game.ExecutablePath)
	}
	// Look for common game executable patterns
	for _, pattern := range []string{"*.exe"} {
		matches, _ := filepath.Glob(filepath.Join(game.InstallPath, pattern))
		if len(matches) > 0 {
			return game.InstallPath
		}
	}
	// Check Binaries/Win64 (UE games)
	ue := filepath.Join(game.InstallPath, "Binaries", "Win64")
	if dirExists(ue) {
		return ue
	}
	// Check common subdirs
	for _, sub := range []string{"Bin", "bin", "x64", "Game"} {
		d := filepath.Join(game.InstallPath, sub)
		if dirExists(d) {
			return d
		}
	}
	return game.InstallPath
}

func (inst *Installer) installComponent(srcDir, gameDir, dllName string, manifest *models.InstallationManifest) error {
	// Find the DLL recursively
	var srcPath string
	filepath.WalkDir(srcDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if strings.EqualFold(d.Name(), dllName) {
			srcPath = path
			return filepath.SkipAll
		}
		return nil
	})
	if srcPath == "" {
		return fmt.Errorf("%s not found in %s", dllName, srcDir)
	}
	copyFile(srcPath, filepath.Join(gameDir, dllName))
	manifest.InstalledFiles = append(manifest.InstalledFiles, dllName)
	return nil
}

func findOptiScalerDLL(cachePath string) string {
	// Look for OptiScaler.dll or the proxy DLL in the cache
	names := []string{"OptiScaler.dll", "optiscaler.dll"}
	var result string
	filepath.WalkDir(cachePath, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		for _, n := range names {
			if strings.EqualFold(d.Name(), n) {
				result = path
				return filepath.SkipAll
			}
		}
		return nil
	})
	return result
}

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir()
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
