package knowledge

import (
	"fmt"

	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func Delete(i *utils.InteractionCreate, opts utils.CommandInteractionOptionsMap) error {
	var sections []discordgo.Section
	var containers []*discordgo.Container

	command := opts["단어"].StringValue()

	data, err := databases.GetDatabase().Knowledge.GetByFilter(databases.Knowledge{UserID: i.User.ID, Command: command})
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return utils.NewMessageSender(i).
			AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "해당 결과를 찾을 수 없어요."})).
			SetComponentsV2(true).
			SetEphemeral(true).
			Send()
	}

	for x, data := range data {
		sections = append(sections, discordgo.Section{
			Accessory: discordgo.Button{
				Label:    "삭제",
				Style:    discordgo.DangerButton,
				CustomID: utils.MakeDeleteKnowledge(data.ID.Hex(), x+1, i.User.ID),
			},
			Components: []discordgo.MessageComponent{
				discordgo.TextDisplay{
					Content: fmt.Sprintf("%d. %s\n", x+1, data.Result),
				},
			},
		})
	}

	textDisplay := discordgo.TextDisplay{Content: fmt.Sprintf("### %s 삭제", command)}
	container := &discordgo.Container{Components: []discordgo.MessageComponent{textDisplay}}
	for i, section := range sections {
		container.Components = append(container.Components, section, discordgo.Separator{})

		if (i+1)%10 == 0 {
			containers = append(containers, container)
			container = &discordgo.Container{Components: []discordgo.MessageComponent{textDisplay}}
			continue
		}
	}

	if len(container.Components) > 1 {
		containers = append(containers, container)
	}

	return utils.PaginationContainerBuilder(i).
		AddContainers(containers...).
		Start()
}
