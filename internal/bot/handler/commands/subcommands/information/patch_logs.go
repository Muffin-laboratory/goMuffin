package information

import (
	"strings"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func InfoPatchLogs(e *handler.CommandEvent) error {
	var containers []discord.ContainerComponent

	startPage := e.SlashCommandInteractionData().Int("버전")
	bot, _ := e.Client().Caches.SelfUser()

	releases, err := repository.Releases()
	if err != nil {
		return err
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
		SetStartPage(startPage).
		Start()
}
