package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"UpscalerHub/internal/archive"
	"UpscalerHub/internal/models"
)

// SwapperService handles reading the dlss-swapper manifest and managing the native DLLs cache
type SwapperService struct {
	client      *http.Client
	cacheDir    string
	manifest    *models.SwapperManifest
	manifestURL string
}

func NewSwapperService() *SwapperService {
	// Cache directory in APPDATA
	appData := os.Getenv("APPDATA")
	cacheDir := filepath.Join(appData, "UpscalerHub", "swapper")
	os.MkdirAll(cacheDir, 0755)

	return &SwapperService{
		client:      &http.Client{},
		cacheDir:    cacheDir,
		manifestURL: "https://beeradmoore.github.io/dlss-swapper/manifest.json",
	}
}

// FetchManifest retrieves the latest manifest from the dlss-swapper repository
func (s *SwapperService) FetchManifest() (*models.SwapperManifest, error) {
	resp, err := s.client.Get(s.manifestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch manifest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("manifest returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest body: %w", err)
	}

	var m models.SwapperManifest
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("failed to decode manifest JSON: %w", err)
	}

	s.manifest = &m

	// Add the "Type" string internally and check cache
	s.processRecords(m.DLSS, "dlss")
	s.processRecords(m.DLSSG, "dlssg")
	s.processRecords(m.DLSSD, "dlssd")
	s.processRecords(m.FSR31DX12, "fsr31dx12")
	s.processRecords(m.FSR31VK, "fsr31vk")
	s.processRecords(m.XeSS, "xess")
	s.processRecords(m.XeSSFG, "xessfg")

	return s.manifest, nil
}

// processRecords attaches the type and checks if the file exists locally
func (s *SwapperService) processRecords(records []models.DLLRecord, recordType string) {
	for i := range records {
		records[i].Type = recordType
		localPath := s.getDLLPath(&records[i])
		if _, err := os.Stat(localPath); err == nil {
			records[i].Downloaded = true
			records[i].LocalPath = localPath
		} else {
			records[i].Downloaded = false
		}
	}
}

// getDLLPath returns the expected local path for a downloaded DLL
func (s *SwapperService) getDLLPath(record *models.DLLRecord) string {
	// e.g. dlss_v3.7.10.0_887B23...
	folder := filepath.Join(s.cacheDir, record.Type)
	os.MkdirAll(folder, 0755)

	name := fmt.Sprintf("%s_v%s_%s.dll", record.Type, record.Version, record.MD5Hash)
	return filepath.Join(folder, name)
}

// DownloadDLL downloads and extracts the requested DLL
func (s *SwapperService) DownloadDLL(record *models.DLLRecord) error {
	if record.DownloadUrl == "" {
		return fmt.Errorf("no download URL found for this version")
	}

	// 1. Download zip to a temp file
	resp, err := s.client.Get(record.DownloadUrl)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	tempZip, err := os.CreateTemp(s.cacheDir, "*.zip")
	if err != nil {
		return err
	}
	tempZipPath := tempZip.Name()
	defer os.Remove(tempZipPath)

	_, err = io.Copy(tempZip, resp.Body)
	tempZip.Close()
	if err != nil {
		return fmt.Errorf("failed to save zip: %w", err)
	}

	// 2. Extract specific .dll file from zip
	destPath := s.getDLLPath(record)

	// Ensure the user hasn't already extracted it.
	if _, err := os.Stat(destPath); err == nil {
		record.Downloaded = true
		record.LocalPath = destPath
		return nil
	}

	// For extraction, we expect one .dll inside the zip that we care about.
	// E.g., nvngx_dlss.dll, libxess.dll, etc.
	// DLSS Swapper's zip files usually contain the DLL at the root. We can extract all to a temp dir,
	// find the heaviest .dll file, and rename it to destPath.

	tempExtractedDir, err := os.MkdirTemp(s.cacheDir, "extracted*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempExtractedDir)

	err = archive.Extract(tempZipPath, tempExtractedDir)
	if err != nil {
		return fmt.Errorf("failed to extract zip: %w", err)
	}

	// Find the .dll extension in the tempExtractedDir
	var dllFile string
	err = filepath.Walk(tempExtractedDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".dll") {
			dllFile = path
		}
		return nil
	})
	if err != nil {
		return err
	}

	if dllFile == "" {
		return fmt.Errorf("no DLL found inside the extracted zip")
	}

	// Move the DLL to cache
	err = os.Rename(dllFile, destPath)
	if err != nil {
		// fallback to copying if rename fails across drives.
		input, err := os.ReadFile(dllFile)
		if err != nil {
			return err
		}
		err = os.WriteFile(destPath, input, 0644)
		if err != nil {
			return err
		}
	}

	record.Downloaded = true
	record.LocalPath = destPath

	return nil
}

// targetDLLNames defines the expected file name for each upscaler type
var targetDLLNames = map[string]string{
	"dlss":      "nvngx_dlss.dll",
	"dlssg":     "nvngx_dlssg.dll",
	"dlssd":     "nvngx_dlssd.dll",
	"fsr31dx12": "ffx_fsr3_x12_x64.dll",
	"fsr31vk":   "ffx_fsr3_vk_x64.dll", // typically
	"xess":      "libxess.dll",
}

// swapBackupExtension is used to backup the original native DLL
const swapBackupExtension = ".orig.backup"

// SwapDLL takes a downloaded DLL record and places it into the target game path.
// It backs up the original DLL if not already backed up.
func (s *SwapperService) SwapDLL(record *models.DLLRecord, targetDir string) error {
	if !record.Downloaded || record.LocalPath == "" {
		return fmt.Errorf("DLL is not downloaded yet")
	}

	expectedName, ok := targetDLLNames[record.Type]
	if !ok {
		return fmt.Errorf("unknown DLL expected name for type %s", record.Type)
	}

	targetFile := filepath.Join(targetDir, expectedName)
	backupFile := targetFile + swapBackupExtension

	// Check if the original DLL is there and we don't have a backup yet
	if _, err := os.Stat(targetFile); err == nil {
		if _, bErr := os.Stat(backupFile); os.IsNotExist(bErr) {
			// Create a backup
			err = os.Rename(targetFile, backupFile)
			if err != nil {
				return fmt.Errorf("failed to backup original dll: %w", err)
			}
		} else {
			// Backup already exists, just remove the current targetFile to overwrite
			os.Remove(targetFile)
		}
	}

	// Copy the local cached DLL to the target file
	input, err := os.ReadFile(record.LocalPath)
	if err != nil {
		return fmt.Errorf("failed to read cached DLL: %w", err)
	}
	err = os.WriteFile(targetFile, input, 0644)
	if err != nil {
		return fmt.Errorf("failed to write DLL to game directory: %w", err)
	}

	return nil
}

// RestoreOriginalDLL restores the backed-up native DLL for a specific type
func (s *SwapperService) RestoreOriginalDLL(targetDir string, upscalerType string) error {
	expectedName, ok := targetDLLNames[upscalerType]
	if !ok {
		return fmt.Errorf("unknown DLL expected name for type %s", upscalerType)
	}

	targetFile := filepath.Join(targetDir, expectedName)
	backupFile := targetFile + swapBackupExtension

	// If backup doesn't exist, we assume original is already there or was never backed up
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		return fmt.Errorf("no backup found to restore")
	}

	// Remove whatever is currently there
	os.Remove(targetFile)

	// Rename backup back to target
	err := os.Rename(backupFile, targetFile)
	if err != nil {
		return fmt.Errorf("failed to restore original dll: %w", err)
	}

	return nil
}
