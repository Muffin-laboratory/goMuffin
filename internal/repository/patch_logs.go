package repository

import (
	"encoding/json"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/Muffin-laboratory/goMuffin/internal/configs"
)

type Release struct {
	Version string `json:"version"`
	Body    string `json:"body"`
}

var oldReleases []Release
var oldReleasesOnce sync.Once

func OldPatchLogs() []Release {
	oldReleasesOnce.Do(func() {
		if configs.GetConfig().GitHub.OldRepository != "" {
			bytes, err := os.ReadFile("old_muffin_patch_logs.json")
			if err != nil {
				return
			}

			if err = json.Unmarshal(bytes, &oldReleases); err != nil {
				return
			}
		}
	})

	oldReleasesCopy := make([]Release, len(oldReleases))
	copy(oldReleasesCopy, oldReleases)
	return oldReleasesCopy
}

var tags []string
var tagsOnce sync.Once

func Tags() []string {
	tagsOnce.Do(func() {
		var out strings.Builder

		cmd := exec.Command("git", "--no-pager", "tag", "--sort=-creatordate")
		cmd.Stdout = &out

		if err := cmd.Run(); err != nil {
			slog.Error("[Fatal] error while creating patch log choices", "error", err)
			os.Exit(1)
		}

		list := strings.Split(out.String(), "\n")

		if len(list) > 25 {
			list = list[:25]
		}

		for _, v := range list {
			if v == "" {
				continue
			}

			tags = append(tags, v)
		}

		tags = append(tags, "4.1.1-Pudding", "4.1.0-Pudding", "4.0.0-Pudding")
		tags = append(tags, "3.2.1-Cake", "3.2.0-Cake", "3.1.0-Cake", "3.0.2-Cake", "3.0.1-Cake", "3.0.0-Cake")

		for _, v := range OldPatchLogs() {
			tags = append(tags, v.Version)
		}
	})

	tagsCopy := make([]string, len(tags))
	copy(tagsCopy, tags)
	return tagsCopy
}
