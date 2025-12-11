package update

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/creativeprojects/go-selfupdate"
)

const (
	repoOwner = "mr-kaynak"
	repoName  = "go-ssh"
)

type Checker struct {
	currentVersion string
}

func NewChecker(version string) *Checker {
	return &Checker{currentVersion: version}
}

func (c *Checker) CheckForUpdate(ctx context.Context) (*selfupdate.Release, bool, error) {
	if c.currentVersion == "dev" || c.currentVersion == "unknown" {
		return nil, false, fmt.Errorf("development build - updates disabled")
	}

	latest, found, err := selfupdate.DetectLatest(ctx, selfupdate.ParseSlug(repoOwner+"/"+repoName))
	if err != nil {
		return nil, false, err
	}

	if !found {
		return nil, false, fmt.Errorf("no releases found")
	}

	currentVer, err := semver.NewVersion(strings.TrimPrefix(c.currentVersion, "v"))
	if err != nil {
		return nil, false, fmt.Errorf("invalid current version: %w", err)
	}

	latestVer, err := semver.NewVersion(strings.TrimPrefix(latest.Version(), "v"))
	if err != nil {
		return nil, false, fmt.Errorf("invalid latest version: %w", err)
	}

	if !latestVer.GreaterThan(currentVer) {
		return nil, false, nil
	}

	return latest, true, nil
}

func (c *Checker) PerformUpdate(ctx context.Context, release *selfupdate.Release) error {
	exe, err := selfupdate.ExecutablePath()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	return selfupdate.UpdateTo(ctx, release.AssetURL, release.AssetName, exe)
}

func CheckAndNotify(currentVersion string) {
	state, err := LoadState()
	if err != nil || !state.ShouldCheck() {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	checker := NewChecker(currentVersion)
	release, found, err := checker.CheckForUpdate(ctx)
	if err != nil || !found {
		state.LastCheck = time.Now()
		state.Save()
		return
	}

	fmt.Printf("\nUpdate available: %s → %s\n", currentVersion, release.Version())
	fmt.Println("Run 'go-ssh --update' to update")

	state.LastCheck = time.Now()
	state.LastVersion = release.Version()
	state.Save()
}
