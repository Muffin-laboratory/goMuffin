package knowledge

import (
	"slices"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/bwmarrin/discordgo"
)

func List(i *builders.InteractionCreate, opts *discordgo.ApplicationCommandInteractionDataOption) error {
	var command string
	var items []string
	var sections []*builders.Section
	var containers []*builders.Container

	if opt := opts.GetOption("단어"); opt != nil {
		command = opt.StringValue()
	}

	filter := query.KnowledgeQueryBuilder().SetUserID(i.User.ID)

	if command != "" {
		filter.SetCommandByRegex(command)
	} else {
		command = "전체"
	}

	data, err := repository.GetDatabase().Knowledge.Find(i.Ctx, filter)
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
				AddText("- **%s**", item),
		)
	}

	title := builders.SectionBuilder().
		SetAccessory(builders.ThumbnailBuilder(i.User.AvatarURL("512"))).
		AddText("### %s님에 대한 지식", i.User.GlobalName).
		AddText("- 총 `%d`개", len(items)).
		AddText("> %s에 대한 검색 결과", command)
	container := builders.ContainerBuilder().AddComponents(title, builders.SeparatorBuilder())

	for i, section := range sections {
		container.AddComponents(section, builders.SeparatorBuilder())

		if (i+1)%5 == 0 {
			containers = append(containers, container)
			container = builders.ContainerBuilder().AddComponents(title, builders.SeparatorBuilder())
			continue
		}
	}

	if container.GetComponentsLength() > 1 {
		containers = append(containers, container)
	}

	return builders.PaginationContainerBuilder(i).
		AddContainers(containers...).
		Start()
}
