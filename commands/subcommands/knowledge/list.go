package knowledge

import (
	"fmt"
	"slices"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/repository"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func List(i *builders.InteractionCreate, opts builders.CommandInteractionOptionsMap) error {
	var command string
	var items []string
	var sections []*builders.Section
	var containers []*builders.Container

	if opt, ok := opts["단어"]; ok {
		command = opt.StringValue()
	}

	filter := bson.D{{Key: "user_id", Value: i.User.ID}}

	if command != "" {
		filter = append(filter, bson.E{Key: "command", Value: bson.M{"$regex": command}})
	} else {
		command = "전체"
	}

	data, err := repository.GetDatabase().Knowledge.GetByFilter(filter)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return builders.NewMessageSender(i).
			AddComponents(builders.MakeErrorContainer("해당 결과를 찾을 수 없어요.")).
			SetComponentsV2(true).
			SetEphemeral(true).
			Send()
	}

	for _, data := range data {
		if slices.Contains(items, data.Command) {
			continue
		}

		items = append(items, data.Command)
	}

	for _, item := range items {
		sections = append(sections,
			builders.SectionBuilder().
				SetAccessory(
					builders.ButtonBuilder().
						SetStyle(discordgo.PrimaryButton).
						SetLabel("자세히 보기").
						SetCustomID(utils.MakeSelectKnowledge(item)),
				).
				AddText(fmt.Sprintf("- **%s**", item)),
		)
	}

	title := builders.SectionBuilder().
		SetAccessory(builders.ThumbnailBuilder(i.User.AvatarURL("512"))).
		AddText(fmt.Sprintf("### %s님에 대한 지식", i.User.GlobalName)).
		AddText(fmt.Sprintf("- 총 `%d`개", len(items))).
		AddText(fmt.Sprintf("> %s에 대한 검색 결과", command))
	container := builders.ContainerBuilder().AddComponents(title, builders.SeparatorBuilder())

	for i, section := range sections {
		container.Components = append(container.Components, section, discordgo.Separator{})

		if (i+1)%5 == 0 {
			containers = append(containers, container)
			container = builders.ContainerBuilder().AddComponents(title, builders.SeparatorBuilder())
			continue
		}
	}

	if len(container.Components) > 1 {
		containers = append(containers, container)
	}

	return builders.PaginationContainerBuilder(i).
		AddContainers(containers...).
		Start()
}
