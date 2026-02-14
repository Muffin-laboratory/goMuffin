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

	bot, _ := e.Client().Caches.SelfUser()

	ghConfig := &configs.GetConfig().GitHub
	if ghConfig.Owner == "" || ghConfig.Repository == "" {
		_, err := e.UpdateInteractionResponse(
			discord.NewMessageUpdateV2([]discord.LayoutComponent{
				builders.MakeErrorContainer("해당 봇은 패치내역을 제공하지 않아요."),
			}),
		)
		return err
	}

	ghClient := repository.GetGHClient()

	releases, _, err := ghClient.Repositories.ListReleases(e.Ctx, ghConfig.Owner, ghConfig.Repository, nil)
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

	return builders.PaginationContainerBuilder(e, true).
		AddContainers(containers...).
		Start()
}
