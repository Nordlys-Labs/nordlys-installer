package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/nordlys-labs/nordlys-installer/internal/constants"
)

const (
	GitHubReleasesURL = "https://api.github.com/repos/nordlys-labs/nordlys-installer/releases/latest"
)

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type Updater struct {
	client      HTTPClient
	currentVer  string
	releasesURL string
}

func NewUpdater() *Updater {
	return &Updater{
		client:      &http.Client{Timeout: 10 * time.Second},
		currentVer:  constants.Version,
		releasesURL: GitHubReleasesURL,
	}
}

func NewUpdaterWithClient(client HTTPClient, currentVer, releasesURL string) *Updater {
	return &Updater{
		client:      client,
		currentVer:  currentVer,
		releasesURL: releasesURL,
	}
}

func (u *Updater) CheckForUpdate() (string, bool, error) {
	resp, err := u.client.Get(u.releasesURL)
	if err != nil {
		return "", false, fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", false, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", false, fmt.Errorf("failed to parse release info: %w", err)
	}

	latestVersion := release.TagName
	if latestVersion == "" {
		return "", false, nil
	}

	if latestVersion[0] == 'v' {
		latestVersion = latestVersion[1:]
	}

	if latestVersion != u.currentVer {
		return latestVersion, true, nil
	}

	return latestVersion, false, nil
}

func CheckForUpdate() (string, bool, error) {
	u := NewUpdater()
	return u.CheckForUpdate()
}

func (u *Updater) SelfUpdate() error {
	latestVersion, needsUpdate, err := u.CheckForUpdate()
	if err != nil {
		return err
	}

	if !needsUpdate {
		return fmt.Errorf("already at the latest version (%s)", u.currentVer)
	}

	resp, err := u.client.Get(u.releasesURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return err
	}

	assetName := GetAssetName()
	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no release asset found for %s", assetName)
	}

	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	resp, err = u.client.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tmpFile := execPath + ".new"
	f, err := os.Create(tmpFile)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, resp.Body)
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(tmpFile)
		return err
	}

	if err := os.Chmod(tmpFile, 0o755); err != nil {
		_ = os.Remove(tmpFile)
		return err
	}

	backupPath := execPath + ".old"
	if err := os.Rename(execPath, backupPath); err != nil {
		_ = os.Remove(tmpFile)
		return err
	}

	if err := os.Rename(tmpFile, execPath); err != nil {
		_ = os.Rename(backupPath, execPath)
		return err
	}

	_ = os.Remove(backupPath)

	fmt.Printf("Updated to version %s\n", latestVersion)
	return nil
}

func SelfUpdate() error {
	u := &Updater{
		client:      &http.Client{Timeout: 60 * time.Second},
		currentVer:  constants.Version,
		releasesURL: GitHubReleasesURL,
	}
	return u.SelfUpdate()
}

func GetAssetName() string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	ext := ""
	if goos == "windows" {
		ext = ".exe"
	}

	return fmt.Sprintf("nordlys-installer-%s-%s%s", goos, goarch, ext)
}

func GetUpdateCheckCachePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cache", "nordlys-installer", "update-check")
}
