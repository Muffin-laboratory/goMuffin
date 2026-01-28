package information

import (
	"context"
	"fmt"
	"strings"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
)

func InfoPatchLogs(i *builders.InteractionCreate) error {
	var containers []*builders.Container

	ghConfig := &configs.GetConfig().GitHub
	if ghConfig.Owner == "" || ghConfig.Repository == "" {
		return returnErrMsg(i)
	}

	ghClient := repository.GetGHClient()

	releases, _, err := ghClient.Repositories.ListReleases(context.TODO(), ghConfig.Owner, ghConfig.Repository, nil)
	if err != nil {
		return err
	}

	for _, release := range releases {
		containers = append(containers,
			builders.ContainerBuilder().
				AddComponents(
					builders.SectionBuilder().
						SetAccessory(builders.ThumbnailBuilder(i.Session.State.User.AvatarURL("512"))).
						AddText(fmt.Sprintf("# %s", *release.TagName)).
						AddText(strings.ReplaceAll(*release.Body, "#", "##")),
				),
		)
	}

	return builders.PaginationContainerBuilder(i).
		AddContainers(containers...).
		Start()
}

func returnErrMsg(i *builders.InteractionCreate) error {
	return builders.NewMessageSender(i).
		AddComponents(builders.MakeErrorContainer("해당 봇은 패치내역을 제공하지 않아요.")).
		SetComponentsV2(true).
		SetEphemeral(true).
		Send()
}
