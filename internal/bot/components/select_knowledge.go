package components

import (
	"context"
	"strings"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var SelectKnowledgeComponent = &loader.Component{
	DeferredReply:    true,
	IsDeferEphemeral: true,
	Parse: func(ctx context.Context, inter *events.ComponentInteractionCreate) bool {
		return strings.HasPrefix(inter.Data.CustomID(), utils.SelectKnowledge)
	},
	Run: func(ctx context.Context, inter *events.ComponentInteractionCreate) error {
		var sections []discord.SectionComponent
		var containers []discord.ContainerComponent

		command := utils.GetSelectKnowledgeCommand(inter.Data.CustomID())

		filter := query.KnowledgeQueryBuilder().SetUserID(int64(inter.User().ID)).SetCommand(command)
		data, err := repository.GetDatabase().Knowledge.Find(ctx, filter)
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
				discord.NewSection(
					discord.NewTextDisplayf("**%s**\n", data.Result),
				).
					WithAccessory(
						discord.NewDangerButton("삭제", utils.MakeDeleteKnowledge(data.ID.Hex(), data.Result, inter.User().ID.String())),
					),
			)
		}

		textDisplay := discord.NewTextDisplayf("### %s에 대한 목록", command)
		container := discord.NewContainer().WithComponents(textDisplay)
		for i, section := range sections {
			container = container.AddComponents(section, discord.NewSeparator(discord.SeparatorSpacingSizeSmall))

			if (i+1)%10 == 0 {
				containers = append(containers, container)
				container = container.WithComponents(textDisplay)
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
