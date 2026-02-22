package information

import (
	"strings"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func InfoPatchLogs(e *handler.CommandEvent) error {
	var containers []discord.ContainerComponent
	var releases []repository.Release

	bot, _ := e.Client().Caches.SelfUser()

	ghConfig := configs.GetConfig().GitHub
	ghClient := repository.GetGHClient()

	ghReleases, _, err := ghClient.Repositories.ListReleases(e.Ctx, ghConfig.Owner, ghConfig.Repository, nil)
	if err != nil {
		return err
	}

	for _, ghRelease := range ghReleases {
		releases = append(releases, repository.Release{Version: *ghRelease.TagName, Body: *ghRelease.Body})
	}

	if ghConfig.OldRepository != "" {
		oldReleases, _, err := ghClient.Repositories.ListReleases(e.Ctx, ghConfig.Owner, ghConfig.OldRepository, nil)
		if err != nil {
			return err
		}

		for _, oldRelease := range oldReleases {
			releases = append(releases, repository.Release{Version: *oldRelease.TagName, Body: *oldRelease.Body})
		}

		releases = append(releases, repository.OldPatchLogs()...)
	}

	for _, release := range releases {
		containers = append(containers,
			discord.NewContainer(
				discord.NewSection(
					discord.NewTextDisplayf("# %s", release.Version),
					discord.NewTextDisplayf("%s", strings.ReplaceAll(release.Body, "#", "##")),
				).
					WithAccessory(discord.NewThumbnail(*bot.AvatarURL())),
			),
		)
	}

	return builders.NewPaginatedContainer(e, true).
		AddContainers(containers...).
		Start()
}
