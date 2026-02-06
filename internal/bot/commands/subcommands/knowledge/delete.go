package knowledge

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/bwmarrin/discordgo"
)

func Delete(i *builders.InteractionCreate, opts *discordgo.ApplicationCommandInteractionDataOption) error {
	var sections []*builders.Section
	var containers []*builders.Container

	command := opts.GetOption("단어").StringValue()

	filter := query.KnowledgeQueryBuilder().SetUserID(i.User.ID).SetCommand(command)
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
		sections = append(sections,
			builders.SectionBuilder().
				SetAccessory(
					builders.ButtonBuilder().
						SetStyle(discordgo.DangerButton).
						SetLabel("삭제").
						SetCustomID(utils.MakeDeleteKnowledge(data.ID.Hex(), data.Result, i.User.ID)),
				).
				AddText("**%s**\n", data.Result),
		)
	}

	textDisplay := builders.TextDisplayBuilder("### %s 삭제", command)
	container := builders.ContainerBuilder().AddComponents(textDisplay)
	for i, section := range sections {
		container.AddComponents(section, builders.SeparatorBuilder())

		if (i+1)%10 == 0 {
			containers = append(containers, container)
			container = builders.ContainerBuilder().AddComponents(textDisplay)
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
