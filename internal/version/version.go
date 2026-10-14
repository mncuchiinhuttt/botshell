package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	BaseVersion = "v1.0.0"
	Author      = "@mncuchiinhuttt"
	RepoOwner   = "mncuchiinhuttt"
	RepoName    = "botshell"
)

// GitCommit will be set via ldflags or default to short commit
var GitCommit = "dev"

// CurrentVersion returns formatted version string
func CurrentVersion() string {
	if GitCommit != "" && GitCommit != "dev" {
		return fmt.Sprintf("%s (%s)", BaseVersion, GitCommit)
	}
	return BaseVersion
}

type githubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

// UpdateInfo holds release check details.
type UpdateInfo struct {
	HasUpdate     bool
	LatestVersion string
	ReleaseURL    string
}

// CheckLatest checks GitHub Releases API for the latest version tag.
func CheckLatest() (*UpdateInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", RepoOwner, RepoName)
	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "BotShell-Updater")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned status %d", resp.StatusCode)
	}

	var rel githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, err
	}

	latest := strings.TrimSpace(rel.TagName)
	hasUpdate := latest != "" && latest != BaseVersion && !strings.EqualFold(latest, BaseVersion)

	return &UpdateInfo{
		HasUpdate:     hasUpdate,
		LatestVersion: latest,
		ReleaseURL:    rel.HTMLURL,
	}, nil
}
