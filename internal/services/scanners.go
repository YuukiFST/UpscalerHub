package services

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/sys/windows/registry"

	"UpscalerHub/internal/models"
)

// ScanSteam finds games installed via Steam.
func ScanSteam() []models.Game {
	installPath := steamInstallPath()
	if installPath == "" {
		return nil
	}

	var games []models.Game
	for _, lib := range steamLibraryFolders(installPath) {
		steamapps := filepath.Join(lib, "steamapps")
		entries, err := os.ReadDir(steamapps)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !strings.HasPrefix(e.Name(), "appmanifest_") || !strings.HasSuffix(e.Name(), ".acf") {
				continue
			}
			if g := parseSteamManifest(filepath.Join(steamapps, e.Name())); g != nil {
				if dirExists(g.InstallPath) {
					games = append(games, *g)
				}
			}
		}
	}
	return games
}

func steamInstallPath() string {
	// Try HKEY_LOCAL_MACHINE first (32-bit view)
	for _, root := range []registry.Key{registry.LOCAL_MACHINE, registry.CURRENT_USER} {
		for _, path := range []string{`SOFTWARE\Valve\Steam`, `SOFTWARE\WOW6432Node\Valve\Steam`} {
			k, err := registry.OpenKey(root, path, registry.READ)
			if err != nil {
				continue
			}
			val, _, err1 := k.GetStringValue("InstallPath")
			val2, _, err2 := k.GetStringValue("SteamPath")
			k.Close()
			if err1 == nil && val != "" && dirExists(val) {
				return val
			}
			if err2 == nil && val2 != "" && dirExists(val2) {
				return val2
			}
		}
	}
	// Fallback: check common default paths
	defaults := []string{
		`C:\Program Files (x86)\Steam`,
		`C:\Program Files\Steam`,
		`D:\Steam`,
		`D:\SteamLibrary`,
	}
	for _, d := range defaults {
		if dirExists(filepath.Join(d, "steamapps")) {
			return d
		}
	}
	return ""
}

var vdfPathRe = regexp.MustCompile(`"path"\s+"([^"]+)"`)

func steamLibraryFolders(steamPath string) []string {
	folders := []string{steamPath}
	data, err := os.ReadFile(filepath.Join(steamPath, "steamapps", "libraryfolders.vdf"))
	if err != nil {
		return folders
	}
	for _, m := range vdfPathRe.FindAllStringSubmatch(string(data), -1) {
		p := strings.ReplaceAll(m[1], `\\`, `\`)
		if !containsIgnoreCase(folders, p) {
			folders = append(folders, p)
		}
	}
	return folders
}

var (
	acfAppID      = regexp.MustCompile(`"appid"\s+"(\d+)"`)
	acfName       = regexp.MustCompile(`"name"\s+"([^"]+)"`)
	acfInstallDir = regexp.MustCompile(`"installdir"\s+"([^"]+)"`)
)

func parseSteamManifest(path string) *models.Game {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)

	appID := regexFirst(acfAppID, content)
	if appID == "" {
		appID = strings.TrimSuffix(strings.TrimPrefix(filepath.Base(path), "appmanifest_"), ".acf")
	}
	name := regexFirst(acfName, content)
	if name == "" {
		name = "Unknown Game"
	}
	installDir := regexFirst(acfInstallDir, content)
	if installDir == "" {
		return nil
	}

	steamapps := filepath.Dir(path)
	fullPath := filepath.Join(steamapps, "common", installDir)

	return &models.Game{
		AppID:       appID,
		Name:        name,
		InstallPath: fullPath,
		Platform:    models.PlatformSteam,
	}
}

func regexFirst(re *regexp.Regexp, s string) string {
	m := re.FindStringSubmatch(s)
	if len(m) > 1 {
		return m[1]
	}
	return ""
}

func containsIgnoreCase(list []string, s string) bool {
	for _, v := range list {
		if strings.EqualFold(v, s) {
			return true
		}
	}
	return false
}

func dirExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

// ScanEpic finds games installed via Epic Games Launcher.
func ScanEpic() []models.Game {
	programData := os.Getenv("ProgramData")
	if programData == "" {
		return nil
	}
	manifestsPath := filepath.Join(programData, "Epic", "EpicGamesLauncher", "Data", "Manifests")
	entries, err := os.ReadDir(manifestsPath)
	if err != nil {
		return nil
	}

	var games []models.Game
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".item") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(manifestsPath, e.Name()))
		if err != nil {
			continue
		}
		g := parseEpicManifest(data)
		if g != nil {
			games = append(games, *g)
		}
	}
	return games
}

func parseEpicManifest(data []byte) *models.Game {
	// Simple JSON field extraction without full unmarshal
	var m map[string]interface{}
	if err := jsonUnmarshal(data, &m); err != nil {
		return nil
	}

	// Check categories
	cats, ok := m["AppCategories"].([]interface{})
	if !ok {
		return nil
	}
	isGame := false
	for _, c := range cats {
		if s, _ := c.(string); s == "games" {
			isGame = true
			break
		}
	}
	if !isGame {
		return nil
	}

	appName, _ := m["AppName"].(string)
	mainApp, _ := m["MainGameAppName"].(string)
	if appName != mainApp {
		return nil // skip DLC
	}

	displayName, _ := m["DisplayName"].(string)
	installLoc, _ := m["InstallLocation"].(string)
	catalogID, _ := m["CatalogItemId"].(string)

	if displayName == "" || installLoc == "" || !dirExists(installLoc) {
		return nil
	}

	return &models.Game{
		Name:        displayName,
		InstallPath: installLoc,
		Platform:    models.PlatformEpic,
		AppID:       catalogID,
	}
}

// ScanGOG finds games installed via GOG Galaxy.
func ScanGOG() []models.Game {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\GOG.com\Games`, registry.READ|registry.WOW64_32KEY)
	if err != nil {
		return nil
	}
	defer k.Close()

	subkeys, err := k.ReadSubKeyNames(-1)
	if err != nil {
		return nil
	}

	var games []models.Game
	for _, name := range subkeys {
		sk, err := registry.OpenKey(k, name, registry.READ)
		if err != nil {
			continue
		}
		gameName, _, _ := sk.GetStringValue("gameName")
		gamePath, _, _ := sk.GetStringValue("path")
		gameID, _, _ := sk.GetStringValue("gameID")
		sk.Close()

		if gameID == "" {
			gameID = name
		}
		if gameName != "" && gamePath != "" && dirExists(gamePath) {
			games = append(games, models.Game{
				AppID:       gameID,
				Name:        gameName,
				InstallPath: gamePath,
				Platform:    models.PlatformGOG,
			})
		}
	}
	return games
}

// ScanXbox finds games installed via Xbox / Game Pass.
func ScanXbox() []models.Game {
	var games []models.Game
	for _, letter := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		root := fmt.Sprintf("%c:\\XboxGames", letter)
		entries, err := os.ReadDir(root)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			full := filepath.Join(root, e.Name())
			children, _ := os.ReadDir(full)
			if len(children) > 0 {
				games = append(games, models.Game{
					AppID:       e.Name(),
					Name:        e.Name(),
					InstallPath: full,
					Platform:    models.PlatformXbox,
				})
			}
		}
	}
	return games
}

// ScanEA finds games installed via EA / Origin.
func ScanEA() []models.Game {
	paths := []string{
		`SOFTWARE\WOW6432Node\Electronic Arts\EA Games`,
		`SOFTWARE\Electronic Arts\EA Games`,
	}
	var games []models.Game
	for _, p := range paths {
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, p, registry.READ)
		if err != nil {
			continue
		}
		subkeys, _ := k.ReadSubKeyNames(-1)
		for _, name := range subkeys {
			sk, err := registry.OpenKey(k, name, registry.READ)
			if err != nil {
				continue
			}
			gameName, _, _ := sk.GetStringValue("DisplayName")
			if gameName == "" {
				gameName = name
			}
			gamePath, _, _ := sk.GetStringValue("Install Dir")
			sk.Close()

			if gameName != "" && gamePath != "" && dirExists(gamePath) {
				games = append(games, models.Game{
					AppID:       name,
					Name:        gameName,
					InstallPath: gamePath,
					Platform:    models.PlatformEA,
				})
			}
		}
		k.Close()
	}
	return games
}

// ScanBattleNet finds games installed via Battle.net.
func ScanBattleNet() []models.Game {
	paths := []string{
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,
		`SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`,
	}
	var games []models.Game
	for _, p := range paths {
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, p, registry.READ)
		if err != nil {
			continue
		}
		subkeys, _ := k.ReadSubKeyNames(-1)
		for _, name := range subkeys {
			sk, err := registry.OpenKey(k, name, registry.READ)
			if err != nil {
				continue
			}
			publisher, _, _ := sk.GetStringValue("Publisher")
			if !strings.Contains(strings.ToLower(publisher), "blizzard entertainment") {
				sk.Close()
				continue
			}
			gameName, _, _ := sk.GetStringValue("DisplayName")
			gamePath, _, _ := sk.GetStringValue("InstallLocation")
			sk.Close()

			if gameName == "Battle.net" || gameName == "Blizzard Battle.net App" {
				continue
			}
			gameName = strings.ReplaceAll(gameName, " (PTR)", "")

			if gameName != "" && gamePath != "" && dirExists(gamePath) {
				games = append(games, models.Game{
					AppID:       name,
					Name:        gameName,
					InstallPath: gamePath,
					Platform:    models.PlatformBattleNet,
				})
			}
		}
		k.Close()
	}
	return games
}

// ScanUbisoft finds games installed via Ubisoft Connect / Uplay.
func ScanUbisoft() []models.Game {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`, registry.READ)
	if err != nil {
		return nil
	}
	defer k.Close()

	subkeys, _ := k.ReadSubKeyNames(-1)
	var games []models.Game
	for _, name := range subkeys {
		if !strings.HasPrefix(strings.ToLower(name), "uplay install ") {
			continue
		}
		sk, err := registry.OpenKey(k, name, registry.READ)
		if err != nil {
			continue
		}
		gameName, _, _ := sk.GetStringValue("DisplayName")
		gamePath, _, _ := sk.GetStringValue("InstallLocation")
		sk.Close()

		if gameName != "" && gamePath != "" && dirExists(gamePath) {
			appID := strings.TrimSpace(strings.TrimPrefix(name, "Uplay Install "))
			games = append(games, models.Game{
				AppID:       appID,
				Name:        gameName,
				InstallPath: gamePath,
				Platform:    models.PlatformUbisoft,
			})
		}
	}
	return games
}
