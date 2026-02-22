package repository

import (
	"encoding/json"
	"os"

	"github.com/Muffin-laboratory/goMuffin/internal/configs"
)

type Release struct {
	Version string `json:"version"`
	Body    string `json:"body"`
}

var oldReleases []Release

func init() {
	if configs.GetConfig().GitHub.OldRepository != "" {
		bytes, err := os.ReadFile("old_muffin_patch_logs.json")
		if err != nil {
			return
		}

		if err = json.Unmarshal(bytes, &oldReleases); err != nil {
			return
		}
	}
}

func OldPatchLogs() []Release {
	oldReleasesCopy := make([]Release, len(oldReleases))
	copy(oldReleasesCopy, oldReleases)
	return oldReleasesCopy
}
