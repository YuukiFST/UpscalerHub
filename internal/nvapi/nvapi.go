package nvapi

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

//go:embed bin/nvapi-cli.exe bin/NvAPIWrapper.dll
var binFiles embed.FS

type PresetResult struct {
	Success      bool   `json:"success"`
	Error        string `json:"error,omitempty"`
	DlssPreset   uint   `json:"dlssPreset,omitempty"`
	DlssdPreset  uint   `json:"dlssdPreset,omitempty"`
	FoundProfile bool   `json:"foundProfile,omitempty"`
}

type SetResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

var extractedPath string

func extractBinaries() error {
	if extractedPath != "" {
		return nil
	}

	appData, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(appData, "UpscalerHub", "tools")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	extractedPath = filepath.Join(dir, "nvapi-cli.exe")

	// If it already exists, assume it's valid. This avoids writing exactly at launch every time.
	if _, err := os.Stat(extractedPath); err == nil {
		return nil
	}

	cliData, err := binFiles.ReadFile("bin/nvapi-cli.exe")
	if err != nil {
		return err
	}
	dllData, err := binFiles.ReadFile("bin/NvAPIWrapper.dll")
	if err != nil {
		return err
	}

	if err := os.WriteFile(extractedPath, cliData, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "NvAPIWrapper.dll"), dllData, 0644); err != nil {
		return err
	}

	return nil
}

func GetPresets(exeName string) (*PresetResult, error) {
	if err := extractBinaries(); err != nil {
		return nil, err
	}

	cmd := exec.Command(extractedPath, "get", exeName)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("nvapi-cli Get execution failed: %v, output: %s", err, out.String())
	}

	var res PresetResult
	if err := json.Unmarshal(out.Bytes(), &res); err != nil {
		return nil, fmt.Errorf("failed to parse nvapi-cli output: %v (output: %s)", err, out.String())
	}

	return &res, nil
}

func SetPreset(exeName, tech string, preset uint) error {
	if err := extractBinaries(); err != nil {
		return err
	}

	cmd := exec.Command(extractedPath, "set", exeName, tech, fmt.Sprintf("%d", preset))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("nvapi-cli Set execution failed: %v, output: %s", err, out.String())
	}

	var res SetResult
	if err := json.Unmarshal(out.Bytes(), &res); err != nil {
		return fmt.Errorf("failed to parse nvapi-cli output: %v (output: %s)", err, out.String())
	}

	if !res.Success {
		return fmt.Errorf("nvapi-cli set error: %s", res.Error)
	}

	return nil
}
