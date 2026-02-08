package information

import (
	"context"
	"strings"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
)

func InfoPatchLogs(ctx context.Context, i *builders.CommandCreate) error {
	var containers []discord.ContainerComponent

	bot, _ := i.Client().Caches.SelfUser()

	ghConfig := &configs.GetConfig().GitHub
	if ghConfig.Owner == "" || ghConfig.Repository == "" {
		return returnErrMsg(i)
	}

	ghClient := repository.GetGHClient()

	releases, _, err := ghClient.Repositories.ListReleases(ctx, ghConfig.Owner, ghConfig.Repository, nil)
	if err != nil {
		return err
	}

	for _, release := range releases {
		containers = append(containers,
			discord.NewContainer(
				discord.NewSection(
					discord.NewTextDisplayf("# %s", *release.TagName),
					discord.NewTextDisplayf("%s", strings.ReplaceAll(*release.Body, "#", "##")),
				).
					WithAccessory(discord.NewThumbnail(*bot.AvatarURL())),
			),
		)
	}

	return builders.PaginationContainerBuilder(i).
		AddContainers(containers...).
		Start()
}

func returnErrMsg(i *builders.CommandCreate) error {
	return builders.NewMessageSender(i).
		AddComponents(builders.MakeErrorContainer("해당 봇은 패치내역을 제공하지 않아요.")).
		SetComponentsV2(true).
		SetEphemeral(true).
		Send()
}
