package components

import (
	"context"
	"fmt"
	"strings"

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
		var data []*databases.Learn
		var sections []*discordgo.Section
		var containers []*discordgo.Container

		i := ctx.Inter

		if err := i.DeferReply(&discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		}); err != nil {
			return err
		}

		command := utils.GetSelectKnowledgeCommand(i.MessageComponentData().CustomID)

		cur, err := databases.GetDatabase().Learns.Find(context.TODO(), databases.Learn{UserID: i.User.ID, Command: command})
		if err != nil {
			return err
		}

		defer cur.Close(context.TODO())

		if err = cur.All(context.TODO(), &data); err != nil {
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
			sections = append(sections, &discordgo.Section{
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

		textDisplay := discordgo.TextDisplay{Content: fmt.Sprintf("### %s에 대한 목록", command)}
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
	},
}

func init() {
	commands.GetDiscommand().LoadComponent(SelectKnowledgeComponent)
}
