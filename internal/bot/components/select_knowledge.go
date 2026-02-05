package components

import (
	"fmt"
	"strings"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/bwmarrin/discordgo"
)

var SelectKnowledgeComponent = &loader.Component{
	DeferredReply: true,
	DeferReplyOptions: &discordgo.InteractionResponseData{
		Flags: discordgo.MessageFlagsEphemeral,
	},
	Parse: func(inter *builders.InteractionCreate) bool {
		return strings.HasPrefix(inter.MessageComponentData().CustomID, utils.SelectKnowledge)
	},
	Run: func(inter *builders.InteractionCreate) error {
		var sections []*builders.Section
		var containers []*builders.Container

		command := utils.GetSelectKnowledgeCommand(inter.MessageComponentData().CustomID)

		filter := query.KnowledgeQueryBuilder().SetUserID(inter.User.ID).SetCommand(command)
		data, err := repository.GetDatabase().Knowledge.Find(inter.Ctx, filter)
		if err != nil {
			return err
		}

		if len(data) == 0 {
			return builders.NewMessageSender(inter).
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
							SetCustomID(utils.MakeDeleteKnowledge(data.ID.Hex(), data.Result, inter.User.ID)),
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

		return builders.PaginationContainerBuilder(inter).
			AddContainers(containers...).
			Start()
	},
}

func init() {
	loader.GetDiscommand().LoadComponent(SelectKnowledgeComponent)
}
