package knowledge

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders/customid"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func Delete(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	var sections []discord.SectionComponent
	var containers []discord.ContainerComponent

	command := data.String("단어")

	filter := query.KnowledgeQueryBuilder().SetUserID(int64(e.User().ID)).SetCommand(command)
	knowledge, err := repository.GetDatabase().Knowledge.Find(e.Ctx, filter)
	if err != nil {
		return err
	}

	if len(knowledge) == 0 {
		_, err := e.UpdateInteractionResponse(
			discord.NewMessageUpdateV2([]discord.LayoutComponent{
				builders.MakeErrorContainer("해당 결과를 찾을 수 없어요."),
			}),
		)
		return err
	}

	for _, data := range knowledge {
		sections = append(sections,
			discord.NewSection(
				discord.NewTextDisplayf("**%s**\n", data.Result),
			).
				WithAccessory(
					discord.NewDangerButton("삭제", customid.MakeDeleteKnowledge(data.ID.Hex(), e.User().ID.String())),
				),
		)
	}

	textDisplay := discord.NewTextDisplayf("### %s 삭제", command)
	container := discord.NewContainer(textDisplay)
	for i, section := range sections {
		container = container.AddComponents(section, discord.NewSmallSeparator())

		if (i+1)%10 == 0 {
			containers = append(containers, container)
			container = container.WithComponents(textDisplay)
			continue
		}
	}

	if len(container.Components) > 1 {
		containers = append(containers, container)
	}

	return builders.NewPaginatedContainer(e, true).
		AddContainers(containers...).
		Start()
}
