package information

import (
	"context"
	"fmt"
	"strings"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/google/go-github/v74/github"
)

func InfoPatchLogs(i *utils.InteractionCreate) error {
	var containers []*discordgo.Container

	ghClient := github.NewClient(nil)
	ghConfig := &configs.GetConfig().GitHub

	if ghConfig.Owner == "" || ghConfig.Repository == "" {
		return returnErrMsg(i)
	}

	releases, _, err := ghClient.Repositories.ListReleases(context.TODO(), ghConfig.Owner, ghConfig.Repository, nil)
	if err != nil {
		return err
	}

	for _, release := range releases {
		containers = append(containers, &discordgo.Container{
			Components: []discordgo.MessageComponent{
				discordgo.Section{
					Accessory: discordgo.Thumbnail{Media: discordgo.UnfurledMediaItem{URL: i.Session.State.User.AvatarURL("512")}},
					Components: []discordgo.MessageComponent{
						discordgo.TextDisplay{
							Content: fmt.Sprintf("# %s", *release.TagName),
						},
						discordgo.TextDisplay{
							Content: strings.ReplaceAll(*release.Body, "#", "##"),
						},
					},
				},
			},
		})
	}

	return utils.PaginationContainerBuilder(i).
		AddContainers(containers...).
		Start()
}

func returnErrMsg(i *utils.InteractionCreate) error {
	return utils.NewMessageSender(i).
		AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{
			Content: "해당 봇은 패치내역을 제공하지 않아요.",
		})).
		SetComponentsV2(true).
		SetEphemeral(true).
		Send()
}
