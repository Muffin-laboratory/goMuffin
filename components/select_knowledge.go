package components

import (
	"fmt"
	"strings"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var SelectKnowledgeComponent = &commands.Component{
	Parse: func(ctx *commands.ComponentContext) bool {
		return strings.HasPrefix(ctx.Inter.MessageComponentData().CustomID, utils.SelectKnowledge)
	},
	Run: func(ctx *commands.ComponentContext) error {
		var sections []*builders.Section
		var containers []*builders.Container

		i := ctx.Inter

		if err := i.DeferReply(&discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		}); err != nil {
			return err
		}

		command := utils.GetSelectKnowledgeCommand(i.MessageComponentData().CustomID)

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

		textDisplay := builders.TextDisplayBuilder(fmt.Sprintf("### %s에 대한 목록", command))
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
	},
}

func init() {
	commands.GetDiscommand().LoadComponent(SelectKnowledgeComponent)
}
