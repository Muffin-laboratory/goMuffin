package knowledge

import (
	"fmt"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func Delete(i *builders.InteractionCreate, opts builders.CommandInteractionOptionsMap) error {
	var sections []*builders.Section
	var containers []*builders.Container

	command := opts["단어"].StringValue()

	data, err := databases.GetDatabase().Knowledge.GetByFilter(databases.Knowledge{UserID: i.User.ID, Command: command})
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
		// sections = append(sections, discordgo.Section{
		// 	Accessory: discordgo.Button{
		// 		Label:    "삭제",
		// 		Style:    discordgo.DangerButton,
		// 		CustomID: utils.MakeDeleteKnowledge(data.ID.Hex(), data.Result, i.User.ID),
		// 	},
		// 	Components: []discordgo.MessageComponent{
		// 		discordgo.TextDisplay{
		// 			Content: fmt.Sprintf("**%s**\n", data.Result),
		// 		},
		// 	},
		// })
		sections = append(sections,
			builders.SectionBuilder().
				SetAccessory(
					builders.ButtonBuilder().
						SetStyle(discordgo.DangerButton).
						SetLabel("삭제").
						SetCustomID(utils.MakeDeleteKnowledge(data.ID.Hex(), data.Result, i.User.ID)),
				).
				AddText(fmt.Sprintf("**%s**\n", data.Result)),
		)
	}

	textDisplay := builders.TextDisplayBuilder(fmt.Sprintf("### %s 삭제", command))
	container := builders.ContainerBuilder().AddComponents(textDisplay)
	for i, section := range sections {
		container.Components = append(container.Components, section, discordgo.Separator{})

		if (i+1)%10 == 0 {
			containers = append(containers, container)
			container = builders.ContainerBuilder().AddComponents(textDisplay)
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
