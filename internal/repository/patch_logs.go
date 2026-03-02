package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/configs"
)

type Release struct {
	Version string `json:"version"`
	Body    string `json:"body"`
}

type releaseCacheManager struct {
	releases []Release
	mu       *sync.RWMutex
}

var releaseManager = &releaseCacheManager{[]Release{}, &sync.RWMutex{}}

var (
	ErrConfigNotDefined = fmt.Errorf("config is not defined")
)

func Releases() ([]Release, error) {
	releaseManager.mu.RLock()
	if len(releaseManager.releases) != 0 {
		releasesCopy := make([]Release, len(releaseManager.releases))
		copy(releasesCopy, releaseManager.releases)
		releaseManager.mu.RUnlock()
		return releasesCopy, nil
	}

	releaseManager.mu.RUnlock()

	var patchedReleases []Release

	config := configs.Configs().GitHub
	client := GetGHClient()

	isOwnerEmpty := configs.Configs().GitHub.Owner == ""
	isRepoEmpty := configs.Configs().GitHub.Repository == ""
	if isOwnerEmpty && isRepoEmpty {
		return nil, ErrConfigNotDefined
	}

	goMuffinReleases, _, err := client.Repositories.ListReleases(context.Background(), config.Owner, config.Repository, nil)
	if err != nil {
		return nil, err
	}

	oldMuffinReleases, _, err := client.Repositories.ListReleases(context.Background(), config.Owner, config.OldRepository, nil)
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile("old_muffin_patch_logs.json")
	if err != nil {
		return nil, err
	}

	var oldMuffinSecondReleases []Release
	if err = json.Unmarshal(bytes, &oldMuffinSecondReleases); err != nil {
		return nil, err
	}

	for _, release := range goMuffinReleases {
		patchedReleases = append(patchedReleases, Release{*release.TagName, *release.Body})
	}

	for _, release := range oldMuffinReleases {
		patchedReleases = append(patchedReleases, Release{*release.TagName, *release.Body})
	}

	patchedReleases = append(patchedReleases, oldMuffinSecondReleases...)
	releaseManager.mu.Lock()
	copy(releaseManager.releases, patchedReleases)
	releaseManager.mu.Unlock()

	go func() {
		time.Sleep(1 * time.Hour)
		releaseManager.mu.Lock()
		releaseManager.releases = nil
		releaseManager.mu.Unlock()
	}()

	return patchedReleases, nil
}
