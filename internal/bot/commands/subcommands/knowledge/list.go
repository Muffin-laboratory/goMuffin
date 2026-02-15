package knowledge

import (
	"slices"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func List(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	var command string
	var items []string
	var sections []discord.SectionComponent
	var containers []discord.ContainerComponent

	if value, ok := data.OptString("단어"); ok {
		command = value
	}

	filter := query.KnowledgeQueryBuilder().SetUserID(int64(e.User().ID))
	if command != "" {
		filter.SetCommandByRegex(command)
	} else {
		command = "전체"
	}

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
		if slices.Contains(items, data.Command) {
			continue
		}

		items = append(items, data.Command)
	}

	for _, item := range items {
		sections = append(sections,
			discord.NewSection(
				discord.NewTextDisplayf("- **%s**", item),
			).
				WithAccessory(
					discord.NewPrimaryButton("자세히 보기", utils.MakeSelectKnowledge(item)),
				),
		)
	}

	title := discord.NewSection(
		discord.NewTextDisplayf("### %s님에 대한 지식", *e.User().GlobalName),
		discord.NewTextDisplayf("- 총 `%d`개", len(items)),
		discord.NewTextDisplayf("> %s에 대한 검색 결과", command),
	).
		WithAccessory(discord.NewThumbnail(*e.User().AvatarURL()))
	container := discord.NewContainer(title, discord.NewSmallSeparator())
	for i, section := range sections {
		container = container.AddComponents(section, discord.NewSmallSeparator())

		if (i+1)%5 == 0 {
			containers = append(containers, container)
			container = container.WithComponents(title, discord.NewSmallSeparator())
			continue
		}
	}

	if len(container.Components) > 1 {
		containers = append(containers, container)
	}

	return builders.NewPaginatedContainer(e.User().ID, true).
		AddContainers(containers...).
		Start(e)
}
